package mbserver

import (
	"context"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hootrhino/goodbusserver/store"
)

type mockStore struct{ store.Store }

func (m *mockStore) SetCoils(values []byte) error              { return nil }
func (m *mockStore) SetDiscreteInputs(values []byte) error     { return nil }
func (m *mockStore) SetHoldingRegisters(values []uint16) error { return nil }
func (m *mockStore) SetInputRegisters(values []uint16) error   { return nil }

// fakeConn implements net.Conn for testing
// It simulates a client connection with predefined input/output

type fakeConn struct {
	inBuf    []byte
	outBuf   []byte
	readOnce bool
	closed   bool
}

func (c *fakeConn) Read(b []byte) (int, error) {
	if c.readOnce {
		return 0, net.ErrClosed
	}
	copy(b, c.inBuf)
	c.readOnce = true
	return len(c.inBuf), nil
}

func (c *fakeConn) Write(b []byte) (int, error) {
	c.outBuf = append(c.outBuf, b...)
	return len(b), nil
}

func (c *fakeConn) Close() error                       { c.closed = true; return nil }
func (c *fakeConn) LocalAddr() net.Addr                { return &net.IPAddr{} }
func (c *fakeConn) RemoteAddr() net.Addr               { return &net.IPAddr{} }
func (c *fakeConn) SetDeadline(t time.Time) error      { return nil }
func (c *fakeConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *fakeConn) SetWriteDeadline(t time.Time) error { return nil }

func TestParseRequestSafe_InvalidLength(t *testing.T) {
	s := NewServer(context.Background(), &mockStore{}, 1)
	_, err := s.parseRequestSafe([]byte{0x01, 0x02})
	if err == nil {
		t.Fatal("expected error for short frame, got nil")
	}
}

func TestDispatchRequest_NoHandler(t *testing.T) {
	s := NewServer(context.Background(), &mockStore{}, 1)
	req := Request{FuncCode: 0x99}
	_, err := s.dispatchRequest(req)
	if err == nil {
		t.Fatal("expected error for unknown func code, got nil")
	}
}

func TestDispatchRequest_CustomHandler(t *testing.T) {
	s := NewServer(context.Background(), &mockStore{}, 1)
	called := false
	s.RegisterCustomHandler(0x64, func(r Request, st store.Store) ([]byte, error) {
		called = true
		return []byte{0x01, 0x02}, nil
	})

	req := Request{FuncCode: 0x64}
	resp, err := s.dispatchRequest(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("custom handler was not called")
	}
	if string(resp) != string([]byte{0x01, 0x02}) {
		t.Fatalf("unexpected response: %v", resp)
	}
}

func TestHandleConnection_ClosesProperly(t *testing.T) {
	s := NewServer(context.Background(), &mockStore{}, 1)
	atomic.StoreInt64(&s.activeConns, 1)
	s.wg.Add(1)

	frame := make([]byte, 12)
	frame[6] = 0x01 // slave id
	frame[7] = 0x64 // func code (no handler)
	c := &fakeConn{inBuf: frame}

	go s.handleConnection(c)
	time.Sleep(50 * time.Millisecond)

	if !c.closed {
		t.Fatal("connection was not closed")
	}
}

func TestServer_StartStop(t *testing.T) {
	s := NewServer(context.Background(), &mockStore{}, 1)

	addr := "127.0.0.1:0"
	if err := s.Start(addr); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}

	go func() {
		conn, _ := net.Dial("tcp", s.listener.Addr().String())
		if conn != nil {
			conn.Close()
		}
	}()

	time.Sleep(50 * time.Millisecond)
	s.Stop()
}
