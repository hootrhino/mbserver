package modbus_server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"

	"modbus_server/handler"
	"modbus_server/protocol"
	"modbus_server/store"
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
	connSem        chan struct{} // 用于限制并发连接数
	activeConns    int64         // 记录当前活跃连接数
}

type Request struct {
	Frame        []byte
	SlaveID      byte
	FuncCode     byte
	StartAddress uint16
	Quantity     uint16
}

func NewServer(Store store.Store, maxConns int) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	server := &Server{
		ctx:            ctx,
		cancel:         cancel,
		store:          Store,
		handlers:       make(map[byte]handler.Handler),
		customHandlers: make(map[byte]func(Request, store.Store) ([]byte, error)),
		connSem:        make(chan struct{}, maxConns),
	}

	// Register handlers
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

func (s *Server) SetCoils(values []byte) error {
	return s.store.SetCoils(values)
}

func (s *Server) SetDiscreteInputs(values []byte) error {
	return s.store.SetDiscreteInputs(values)
}

func (s *Server) SetHoldingRegisters(values []uint16) error {
	return s.store.SetHoldingRegisters(values)
}

func (s *Server) SetInputRegisters(values []uint16) error {
	return s.store.SetInputRegisters(values)
}

func (s *Server) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		if s.errorHandler != nil {
			s.errorHandler(err)
		}
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
				s.connSem <- struct{}{}
				conn, err := listener.Accept()
				if err != nil {
					<-s.connSem
					if s.errorHandler != nil {
						s.errorHandler(err)
					}
					continue
				}
				atomic.AddInt64(&s.activeConns, 1)
				s.wg.Add(1)
				go func(c net.Conn) {
					defer func() {
						<-s.connSem
						atomic.AddInt64(&s.activeConns, -1)
						s.wg.Done()
					}()
					s.handleConnection(c)
				}(conn)
			}
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

// RegisterCustomHandler allows users to register custom function code handlers
func (s *Server) RegisterCustomHandler(code byte, handler func(Request, store.Store) ([]byte, error)) {
	s.customHandlers[code] = handler
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	if s.logger != nil {
		s.logger.Printf("New connection from %s. Active connections: %d", conn.RemoteAddr(), atomic.LoadInt64(&s.activeConns))
	}

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if s.errorHandler != nil {
				s.errorHandler(err)
			}
			if s.logger != nil {
				s.logger.Printf("Error reading from %s: %v", conn.RemoteAddr(), err)
			}
			return
		}

		if s.logger != nil {
			s.logger.Printf("Received %d bytes from %s: %s", n, conn.RemoteAddr(), BeautifulByte(buf[:n]))
		}

		request := parseRequest(buf[:n])
		if protocol.IsCustomFuncCode(request.FuncCode) {
			if handler, ok := s.customHandlers[request.FuncCode]; ok {
				response, err := handler(request, s.store)
				if err != nil {
					if s.errorHandler != nil {
						s.errorHandler(err)
					}
					if s.logger != nil {
						s.logger.Printf("Error handling custom request from %s: %v", conn.RemoteAddr(), err)
					}
					continue
				}
				if s.logger != nil {
					s.logger.Printf("Successfully generated response for custom request from %s, preparing to send", conn.RemoteAddr())
				}
				s.logger.Printf("Generated response for custom request from %s: %s", conn.RemoteAddr(), BeautifulByte(response))
				_, err = conn.Write(response)
				if err != nil {
					if s.errorHandler != nil {
						s.errorHandler(err)
					}
					if s.logger != nil {
						s.logger.Printf("Error writing response to %s: %v", conn.RemoteAddr(), err)
					}
					return
				}
				if s.logger != nil {
					s.logger.Printf("Sent response to %s: %s", conn.RemoteAddr(), BeautifulByte(response))
				}
			} else if s.customHandler != nil {
				s.customHandler(request)
			}
			continue
		}

		if handler, ok := s.handlers[request.FuncCode]; ok {
			handlerRequest := convertToHandlerRequest(request)
			if s.logger != nil {
				s.logger.Printf("Starting to handle request from %s with function code %x", conn.RemoteAddr(), request.FuncCode)
			}
			response, err := handler.Handle(handlerRequest, s.store)
			if err != nil {
				if s.errorHandler != nil {
					s.errorHandler(err)
				}
				if s.logger != nil {
					s.logger.Printf("Error handling request from %s: %v", conn.RemoteAddr(), err)
				}
				continue
			}
			if s.logger != nil {
				s.logger.Printf("Successfully generated response for %s, preparing to send", conn.RemoteAddr())
			}
			s.logger.Printf("Generated response for %s: %s", conn.RemoteAddr(), BeautifulByte(response))
			_, err = conn.Write(response)
			if err != nil {
				if s.errorHandler != nil {
					s.errorHandler(err)
				}
				if s.logger != nil {
					s.logger.Printf("Error writing response to %s: %v", conn.RemoteAddr(), err)
				}
				return
			}
			if s.logger != nil {
				s.logger.Printf("Sent response to %s: %s", conn.RemoteAddr(), BeautifulByte(response))
			}
		} else {
			if s.logger != nil {
				s.logger.Printf("No handler found for function code %x from %s", request.FuncCode, conn.RemoteAddr())
			}
		}
	}
}

func BeautifulByte(bytes []byte) string {
	var result string
	for i, b := range bytes {
		if i > 0 && i%16 == 0 {
			result += "\n"
		}
		result += fmt.Sprintf("%02x ", b)
	}
	return result
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

func parseRequest(frame []byte) Request {
	return Request{
		Frame:        frame,
		SlaveID:      frame[6],
		FuncCode:     frame[7],
		StartAddress: uint16(frame[8])<<8 | uint16(frame[9]),
		Quantity:     uint16(frame[10])<<8 | uint16(frame[11]),
	}
}
