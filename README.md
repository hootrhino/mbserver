# Modbus Server

## Overview
This is a Modbus server library implemented in Go. It supports both in-memory and SQLite storage modes, and provides a set of callback-based APIs to handle Modbus requests. This library aims to help developers quickly build Modbus servers, supporting standard function codes and custom function code processing.

## Features
- **Callback-based Request Handling**: Developers can flexibly handle different Modbus requests through callback functions.
- **Support for Standard Modbus Function Codes**: Supports common Modbus function codes, such as reading coils, reading discrete inputs, and reading holding registers.
- **Custom Request Handling**: Allows users to register custom function code handlers to meet special business requirements.
- **Data Persistence**: Supports data persistence using SQLite to ensure data is not lost after the server restarts.
- **Configurable Logging**: You can customize the log output location and format, facilitating debugging and monitoring.
- **Middleware Support**: Provides a middleware mechanism, allowing developers to add custom logic before and after request processing.

## Installation
Use the following command to install the library:
```bash
go get github.com/yourusername/modbus_server
```

## Usage Examples

### Basic Usage
The following example demonstrates how to start a Modbus server and set some sample data:
```go:e:\workspace\modbus_server\examples\main.go
package main

import (
	"log"
	"modbus_server"
	"modbus_server/protocol"
	"modbus_server/store"
	"os"
)

func main() {
	// Create an in-memory store instance
	memStore := store.NewInMemoryStore().(*store.InMemoryStore)

	// Set sample coil data
	defaultCoilsSize := 110 // Original 100 + 10
	sampleCoils := make([]byte, defaultCoilsSize)
	for i := 0; i < defaultCoilsSize; i++ {
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
	defaultHoldingRegistersSize := 10 // Original 0 + 10
	memStore.SetHoldingRegisters(make([]uint16, defaultHoldingRegistersSize))

	// Initialize a Modbus server
	server := modbus_server.NewServer(memStore)

	// Set an error handler
	server.SetErrorHandler(func(err error) {
		log.Printf("Modbus server error: %v", err)
	})

	// Set up logger
	server.SetLogger(os.Stdout)

	// Set more sample holding register data
	sampleHoldingRegisters := make([]uint16, 12) // Original 2 + 10
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
	select {}
}
```

### Custom Function Code Handling
In the above example, we registered a custom function code `0x81` handler. When the server receives a request with this function code, it will execute the custom processing logic and return a response.

## Contributing
If you want to contribute to this project, please follow these steps:
1. Fork this repository.
2. Create a new branch (`git checkout -b feature/your-feature`).
3. Commit your changes (`git commit -am 'Add some feature'`).
4. Push to the branch (`git push origin feature/your-feature`).
5. Create a new Pull Request.

## License
```
Copyright (C) 2025 wwhai

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
```