package mbserver

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"

	"github.com/hootrhino/goodbusserver/handler"
	"github.com/hootrhino/goodbusserver/protocol"
	"github.com/hootrhino/goodbusserver/store"
)

type Server struct {
	listener       net.Listener
	ctx            context.Context
	cancel         context.CancelFunc
	wg             sync.WaitGroup
	errorHandler   func(error)
	logger         *log.Logger
	store          store.Store
	handlers       map[byte]handler.Handler
	customHandler  func(Request)
	customHandlers map[byte]func(Request, store.Store) ([]byte, error)
	connSem        chan struct{}
	activeConns    int64
}

type Request struct {
	Frame        []byte
	SlaveID      byte
	FuncCode     byte
	StartAddress uint16
	Quantity     uint16
}

func NewServer(ctx context.Context, Store store.Store, maxConns int) *Server {
	ctx, cancel := context.WithCancel(ctx)
	server := &Server{
		ctx:            ctx,
		cancel:         cancel,
		store:          Store,
		handlers:       make(map[byte]handler.Handler),
		customHandlers: make(map[byte]func(Request, store.Store) ([]byte, error)),
		connSem:        make(chan struct{}, maxConns),
	}

	// Register built-in handlers
	server.handlers[protocol.FuncCodeReadCoils] = &handler.CoilsHandler{}
	server.handlers[protocol.FuncCodeReadDiscreteInputs] = &handler.DiscreteInputsHandler{}
	server.handlers[protocol.FuncCodeReadHoldingRegisters] = &handler.HoldingRegistersHandler{}
	server.handlers[protocol.FuncCodeReadInputRegisters] = &handler.InputRegistersHandler{}
	server.handlers[protocol.FuncCodeWriteSingleCoil] = &handler.SingleCoilHandler{}
	server.handlers[protocol.FuncCodeWriteSingleRegister] = &handler.SingleRegisterHandler{}
	server.handlers[protocol.FuncCodeWriteMultipleCoils] = &handler.MultipleCoilsHandler{}
	server.handlers[protocol.FuncCodeWriteMultipleRegisters] = &handler.MultipleRegistersHandler{}

	return server
}

func (s *Server) SetErrorHandler(h func(error)) {
	s.errorHandler = h
}

func (s *Server) SetLogger(w io.Writer) {
	s.logger = log.New(w, "[MODBUS SERVER] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func (s *Server) SetCoils(values []byte) error          { return s.store.SetCoils(values) }
func (s *Server) SetDiscreteInputs(values []byte) error { return s.store.SetDiscreteInputs(values) }
func (s *Server) SetHoldingRegisters(values []uint16) error {
	return s.store.SetHoldingRegisters(values)
}
func (s *Server) SetInputRegisters(values []uint16) error { return s.store.SetInputRegisters(values) }

func (s *Server) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		s.handleError(nil, "failed to start listener", err)
		return err
	}
	s.listener = listener

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case <-s.ctx.Done():
				return
			default:
			}

			s.connSem <- struct{}{}
			conn, err := listener.Accept()
			if err != nil {
				<-s.connSem
				s.handleError(nil, "accept failed", err)
				continue
			}

			atomic.AddInt64(&s.activeConns, 1)
			s.wg.Add(1)
			go s.handleConnection(conn)
		}
	}()

	return nil
}

func (s *Server) Stop() {
	s.cancel()
	if s.listener != nil {
		s.listener.Close()
	}
	s.wg.Wait()
}

func (s *Server) OnCustomRequest(h func(Request)) {
	s.customHandler = h
}

func (s *Server) RegisterCustomHandler(code byte, handler func(Request, store.Store) ([]byte, error)) {
	s.customHandlers[code] = handler
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		atomic.AddInt64(&s.activeConns, -1)
		<-s.connSem
		s.wg.Done()
	}()

	if s.logger != nil {
		s.logger.Printf("New connection from %s. Active connections: %d", conn.RemoteAddr(), atomic.LoadInt64(&s.activeConns))
	}

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
		}

		// 为每个请求分配独立缓冲区，避免并发竞争
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if err != net.ErrClosed {
				s.handleError(conn, "read failed", err)
			}
			return
		}

		// 复制数据避免后续处理中的竞争
		frame := make([]byte, n)
		copy(frame, buf[:n])

		req, err := s.parseRequestSafe(frame)
		if err != nil {
			s.handleError(conn, "parse failed", err)
			continue
		}

		resp, err := s.dispatchRequest(req)
		if err != nil {
			s.handleError(conn, "dispatch failed", err)
			continue
		}

		if err := writeResponse(conn, resp); err != nil {
			s.handleError(conn, "write failed", err)
			return
		}
	}
}

func (s *Server) dispatchRequest(req Request) ([]byte, error) {
	if s.logger != nil {
		s.logger.Printf("Dispatching request: SlaveID=%d, FuncCode=0x%x, StartAddress=%d, Quantity=%d",
			req.SlaveID, req.FuncCode, req.StartAddress, req.Quantity)
	}

	if h, ok := s.customHandlers[req.FuncCode]; ok {
		resp, err := h(req, s.store)
		if s.logger != nil {
			if err != nil {
				s.logger.Printf("Custom handler for FuncCode=0x%x failed: %v", req.FuncCode, err)
			} else {
				s.logger.Printf("Custom handler for FuncCode=0x%x succeeded, response length=%d", req.FuncCode, len(resp))
			}
		}
		return resp, err
	}

	if h, ok := s.handlers[req.FuncCode]; ok {
		resp, err := h.Handle(convertToHandlerRequest(req), s.store)
		if s.logger != nil {
			if err != nil {
				s.logger.Printf("Built-in handler for FuncCode=0x%x failed: %v", req.FuncCode, err)
			} else {
				s.logger.Printf("Built-in handler for FuncCode=0x%x succeeded, response length=%d", req.FuncCode, len(resp))
			}
		}
		return resp, err
	}

	if s.customHandler != nil {
		s.customHandler(req)
		if s.logger != nil {
			s.logger.Printf("Fallback custom handler invoked for FuncCode=0x%x", req.FuncCode)
		}
	}

	err := fmt.Errorf("no handler for func code %x", req.FuncCode)
	s.handleError(nil, "dispatchRequest failed", err)
	return nil, err
}

func (s *Server) handleError(conn net.Conn, msg string, err error) {
	if s.errorHandler != nil {
		s.errorHandler(err)
	}
	if s.logger != nil {
		if conn != nil {
			s.logger.Printf("%s (%s): %v", msg, conn.RemoteAddr(), err)
		} else {
			s.logger.Printf("%s: %v", msg, err)
		}
	}
}

func writeResponse(conn net.Conn, response []byte) error {
	_, err := conn.Write(response)
	return err
}

func convertToHandlerRequest(req Request) handler.Request {
	return handler.Request{
		Frame:        req.Frame,
		SlaveID:      req.SlaveID,
		FuncCode:     req.FuncCode,
		StartAddress: req.StartAddress,
		Quantity:     req.Quantity,
	}
}

func (s *Server) parseRequestSafe(frame []byte) (Request, error) {
	if len(frame) < 12 {
		err := fmt.Errorf("invalid frame length: %d", len(frame))
		s.handleError(nil, "parseRequestSafe failed", err)
		return Request{}, err
	}

	// 验证MBAP头部长度字段
	length := uint16(frame[4])<<8 | uint16(frame[5])
	if int(length)+6 > len(frame) {
		err := fmt.Errorf("invalid length field: declared %d, actual %d", length, len(frame)-6)
		s.handleError(nil, "parseRequestSafe failed", err)
		return Request{}, err
	}

	// 验证协议标识符（必须为0）
	protocolID := uint16(frame[2])<<8 | uint16(frame[3])
	if protocolID != 0 {
		err := fmt.Errorf("invalid protocol ID: %d", protocolID)
		s.handleError(nil, "parseRequestSafe failed", err)
		return Request{}, err
	}

	req := Request{
		Frame:        frame,
		SlaveID:      frame[6],
		FuncCode:     frame[7],
		StartAddress: uint16(frame[8])<<8 | uint16(frame[9]),
		Quantity:     uint16(frame[10])<<8 | uint16(frame[11]),
	}

	// 验证功能码特定的要求
	switch req.FuncCode {
	case 0x01, 0x02, 0x03, 0x04: // 读操作
		if req.Quantity == 0 || req.Quantity > 2000 {
			err := fmt.Errorf("invalid quantity for read: %d (must be 1-2000)", req.Quantity)
			s.handleError(nil, "parseRequestSafe failed", err)
			return Request{}, err
		}
	case 0x05: // 写单个线圈
		if len(frame) < 12 {
			err := fmt.Errorf("frame too short for write single coil")
			s.handleError(nil, "parseRequestSafe failed", err)
			return Request{}, err
		}
		value := uint16(frame[10])<<8 | uint16(frame[11])
		if value != 0x0000 && value != 0xFF00 {
			err := fmt.Errorf("invalid coil value: 0x%04X (must be 0x0000 or 0xFF00)", value)
			s.handleError(nil, "parseRequestSafe failed", err)
			return Request{}, err
		}
	case 0x06: // 写单个寄存器
		if len(frame) < 12 {
			err := fmt.Errorf("frame too short for write single register")
			s.handleError(nil, "parseRequestSafe failed", err)
			return Request{}, err
		}
	case 0x0F: // 写多个线圈
		if len(frame) < 14 {
			err := fmt.Errorf("frame too short for write multiple coils")
			s.handleError(nil, "parseRequestSafe failed", err)
			return Request{}, err
		}
		byteCount := int(frame[12])
		expectedBytes := (int(req.Quantity) + 7) / 8
		if byteCount != expectedBytes {
			err := fmt.Errorf("invalid byte count: %d, expected %d", byteCount, expectedBytes)
			s.handleError(nil, "parseRequestSafe failed", err)
			return Request{}, err
		}
		if len(frame) < 13+byteCount {
			err := fmt.Errorf("frame too short for coil data")
			s.handleError(nil, "parseRequestSafe failed", err)
			return Request{}, err
		}
	case 0x10: // 写多个寄存器
		if len(frame) < 14 {
			err := fmt.Errorf("frame too short for write multiple registers")
			s.handleError(nil, "parseRequestSafe failed", err)
			return Request{}, err
		}
		byteCount := int(frame[12])
		expectedBytes := int(req.Quantity) * 2
		if byteCount != expectedBytes {
			err := fmt.Errorf("invalid byte count: %d, expected %d", byteCount, expectedBytes)
			s.handleError(nil, "parseRequestSafe failed", err)
			return Request{}, err
		}
		if len(frame) < 13+byteCount {
			err := fmt.Errorf("frame too short for register data")
			s.handleError(nil, "parseRequestSafe failed", err)
			return Request{}, err
		}
	}

	// 验证地址范围
	maxAddress := uint32(req.StartAddress) + uint32(req.Quantity)
	if maxAddress > 0xFFFF {
		err := fmt.Errorf("address overflow: start=%d, quantity=%d", req.StartAddress, req.Quantity)
		s.handleError(nil, "parseRequestSafe failed", err)
		return Request{}, err
	}

	if s.logger != nil {
		s.logger.Printf("Parsed request: SlaveID=%d, FuncCode=0x%x, StartAddress=%d, Quantity=%d",
			req.SlaveID, req.FuncCode, req.StartAddress, req.Quantity)
	}

	return req, nil
}
