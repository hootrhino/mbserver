package main

import (
	"context"
	modbus_server "github.com/hootrhino/mbserver"
	"github.com/hootrhino/mbserver/protocol"
	"github.com/hootrhino/mbserver/store"
	"log"
	"os"
)

func main() {
	// Create an in-memory store instance
	memStore := store.NewInMemoryStore().(*store.InMemoryStore)

	// Set sample coil data
	defaultCoilsSize := 110
	sampleCoils := make([]byte, defaultCoilsSize)
	for i := range defaultCoilsSize {
		if i%2 == 0 {
			sampleCoils[i] = 1
		} else {
			sampleCoils[i] = 0
		}
	}
	if err := memStore.SetCoils(sampleCoils); err != nil {
		log.Fatalf("Failed to set coils: %v", err)
	}

	// Set sample holding register data
	defaultHoldingRegistersSize := 10
	memStore.SetHoldingRegisters(make([]uint16, defaultHoldingRegistersSize))

	// Set maximum concurrent connections
	maxConns := 100
	// Initialize a Modbus server
	server := modbus_server.NewServer(memStore, maxConns)

	// Set an error handler
	server.SetErrorHandler(func(err error) {
		log.Printf("Modbus server error: %v", err)
	})

	// Set up logger
	server.SetLogger(os.Stdout)

	// Set more sample holding register data
	sampleHoldingRegisters := make([]uint16, 12)
	sampleHoldingRegisters[0] = 0x1234
	sampleHoldingRegisters[1] = 0x5678
	for i := 2; i < 12; i++ {
		sampleHoldingRegisters[i] = uint16(i) * 100
	}
	if err := server.SetHoldingRegisters(sampleHoldingRegisters); err != nil {
		log.Fatalf("Failed to set holding registers: %v", err)
	}

	// Set sample input register data
	defaultInputRegistersSize := 10
	sampleInputRegisters := make([]uint16, defaultInputRegistersSize)
	for i := 0; i < defaultInputRegistersSize; i++ {
		sampleInputRegisters[i] = uint16(i) * 10
	}
	if err := server.SetInputRegisters(sampleInputRegisters); err != nil {
		log.Fatalf("Failed to set input registers: %v", err)
	}

	// Register a custom function code handler
	customCode := byte(0x81)
	server.RegisterCustomHandler(customCode, func(request modbus_server.Request, store store.Store) ([]byte, error) {
		transactionID := protocol.ExtractTransactionID(request.Frame)
		pdu := []byte{customCode, 0x01, 0x02}
		header := protocol.BuildResponseHeader(transactionID, 0, uint16(len(pdu)+1), request.SlaveID)
		response := append(header, pdu...)
		return response, nil
	})

	// Start the Modbus server
	log.Println("Starting Modbus server on :502")
	if err := server.Start(":502"); err != nil {
		log.Fatalf("Failed to start Modbus server: %v", err)
	}
	defer server.Stop()

	// Keep the main goroutine alive
	<-context.Background().Done()
	log.Println("Server stopped")
}
