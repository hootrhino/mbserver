// Copyright (C) 2025 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package protocol

const (
	FuncCodeReadCoils = 0x01
	FuncCodeReadDiscreteInputs = 0x02
	FuncCodeReadHoldingRegisters = 0x03
	FuncCodeReadInputRegisters = 0x04
	FuncCodeWriteSingleCoil = 0x05
	FuncCodeWriteSingleRegister = 0x06
	FuncCodeWriteMultipleCoils = 0x0F
	FuncCodeWriteMultipleRegisters = 0x10 // Add this line
	// Add other standard function codes
)

func IsCustomFuncCode(code byte) bool {
	return code >= 0x80
}

// ExtractTransactionID extracts the transaction ID from the Modbus TCP frame
func ExtractTransactionID(frame []byte) uint16 {
    if len(frame) < 2 {
        return 0
    }
    return uint16(frame[0])<<8 | uint16(frame[1])
}

// BuildResponseHeader builds the MBAP header for the response
func BuildResponseHeader(transactionID uint16, protocolID uint16, length uint16, unitID byte) []byte {
    header := make([]byte, 7)
    header[0] = byte(transactionID >> 8)
    header[1] = byte(transactionID)
    header[2] = byte(protocolID >> 8)
    header[3] = byte(protocolID)
    header[4] = byte(length >> 8)
    header[5] = byte(length)
    header[6] = unitID
    return header
}

var ErrIllegalDataValue = &ProtocolError{Code: "ILLEGAL_DATA_VALUE", Message: "Illegal data value"}

type ProtocolError struct {
    Code    string
    Message string
}

func (e *ProtocolError) Error() string {
    return e.Message
}

// ... existing code ...

