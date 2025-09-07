# Modbus Server

A comprehensive Modbus server library implemented in Go, supporting both TCP and RTU protocols with flexible storage backends.

## Features

- **Multiple Protocol Support**: Modbus TCP and RTU protocols
- **Flexible Storage**: In-memory and SQLite storage backends
- **Standard Function Codes**: Complete support for standard Modbus function codes
- **Custom Handlers**: Extensible callback system for custom function codes
- **Concurrent Safe**: Thread-safe operations with proper locking
- **Configurable**: Flexible configuration options
- **Production Ready**: Comprehensive test coverage and error handling

## Supported Function Codes

### Read Operations
- `01` - Read Coils
- `02` - Read Discrete Inputs
- `03` - Read Holding Registers
- `04` - Read Input Registers

### Write Operations
- `05` - Write Single Coil
- `06` - Write Single Register
- `15` - Write Multiple Coils
- `16` - Write Multiple Registers

## Installation

```bash
go get github.com/hootrhino/goodbusserver
```

## Quick Start

### Basic Server Setup

```go
package main

import (
	"log"
	"github.com/hootrhino/goodbusserver"
	"github.com/hootrhino/goodbusserver/store"
)

func main() {
	// Create storage backend
	store := store.NewInMemoryStore()
	
	// Create and configure server
	server := mbserver.NewServer(store)
	
	// Start server
	if err := server.Start(":502"); err != nil {
		log.Fatal(err)
	}
	defer server.Stop()
	
	log.Println("Modbus server started on :502")
	select {}
}
```

### Setting Data

```go
// Set coil data
coils := []byte{0x01, 0x00, 0x01, 0x01}
store.SetCoils(coils)

// Set holding registers
registers := []uint16{1000, 2000, 3000, 4000}
store.SetHoldingRegisters(registers)

// Set coils at specific address
store.SetCoilsAt(100, []byte{0xFF, 0x00})

// Set holding registers at specific address
store.SetHoldingRegistersAt(200, []uint16{1234, 5678})
```

### Custom Function Handlers

```go
// Register custom function code handler
server.RegisterCustomHandler(0x81, func(request mbserver.Request, store store.Store) ([]byte, error) {
	// Custom processing logic
	response := []byte{request.FuncCode, 0x01, 0x02}
	return response, nil
})
```

### SQLite Storage

```go
// Use SQLite for persistence
store, err := store.NewSqliteStore("modbus.db")
if err != nil {
	log.Fatal(err)
}
defer store.Close()
```

## Configuration

### Server Options

```go
server := mbserver.NewServer(store)

// Set timeout
server.SetTimeout(5 * time.Second)

// Set logger
server.SetLogger(log.New(os.Stdout, "MODBUS: ", log.LstdFlags))

// Set error handler
server.SetErrorHandler(func(err error) {
	log.Printf("Server error: %v", err)
})
```

## API Reference

### Store Interface

```go
type Store interface {
	GetCoils(start, quantity uint16) ([]byte, error)
	GetDiscreteInputs(start, quantity uint16) ([]byte, error)
	GetHoldingRegisters(start, quantity uint16) ([]uint16, error)
	GetInputRegisters(start, quantity uint16) ([]uint16, error)
	
	SetCoils(values []byte) error
	SetDiscreteInputs(values []byte) error
	SetHoldingRegisters(values []uint16) error
	SetInputRegisters(values []uint16) error
	
	SetCoilsAt(start uint16, values []byte) error
	SetHoldingRegistersAt(start uint16, values []uint16) error
}
```

### Error Handling

The library provides comprehensive error handling:

- `ErrIllegalFunction` - Invalid function code
- `ErrIllegalDataAddress` - Invalid data address
- `ErrIllegalDataValue` - Invalid data value
- `ErrServerDeviceFailure` - Server device failure

## Testing

Run the test suite:

```bash
go test -v ./...
```

## Examples

See the `examples/` directory for complete working examples:
- Basic TCP server
- RTU server
- SQLite persistence
- Custom handlers

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the GNU Affero General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## Changelog

### v1.0.0
- Initial release
- Complete Modbus TCP support
- In-memory and SQLite storage
- Comprehensive test coverage
- Custom handler support