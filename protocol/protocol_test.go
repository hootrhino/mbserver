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

import (
	"testing"
)

func TestIsCustomFuncCode(t *testing.T) {
	tests := []struct {
		name     string
		code     byte
		expected bool
	}{
		{
			name:     "Standard function code",
			code:     FuncCodeReadCoils,
			expected: false,
		},
		{
			name:     "Custom function code",
			code:     0x80,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsCustomFuncCode(tt.code)
			if result != tt.expected {
				t.Errorf("IsCustomFuncCode(%d) = %v; want %v", tt.code, result, tt.expected)
			}
		})
	}
}

func TestExtractTransactionID(t *testing.T) {
	frame := []byte{0x12, 0x34, 0x00, 0x00, 0x00, 0x06, 0x01}
	expected := uint16(0x1234)
	result := ExtractTransactionID(frame)
	if result != expected {
		t.Errorf("ExtractTransactionID() = %d; want %d", result, expected)
	}
}

func TestBuildResponseHeader(t *testing.T) {
	transactionID := uint16(0x1234)
	protocolID := uint16(0)
	length := uint16(5)
	unitID := byte(0x01)

	expected := []byte{0x12, 0x34, 0x00, 0x00, 0x00, 0x05, 0x01}
	result := BuildResponseHeader(transactionID, protocolID, length, unitID)

	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("BuildResponseHeader() byte %d = %d; want %d", i, result[i], expected[i])
		}
	}
}

func TestProtocolError_Error(t *testing.T) {
	err := &ProtocolError{
		Code:    "TEST_CODE",
		Message: "Test error message",
	}
	expected := "Test error message"
	result := err.Error()
	if result != expected {
		t.Errorf("Error() = %s; want %s", result, expected)
	}
}

