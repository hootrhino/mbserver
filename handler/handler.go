package handler

import "github.com/hootrhino/mbserver/store"

type Handler interface {
	Handle(request Request, store store.Store) ([]byte, error)
}

type Request struct {
	Frame        []byte
	SlaveID      byte
	FuncCode     byte
	StartAddress uint16
	Quantity     uint16
}
