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

package handler

import (
	"testing"

	"github.com/hootrhino/mbserver/protocol"
	"github.com/hootrhino/mbserver/store"
)

// TestHoldingRegistersHandler_Handle tests the normal handling of holding registers requests
func TestHoldingRegistersHandler_Handle(t *testing.T) {
	handler := &HoldingRegistersHandler{}
	memStore := store.NewInMemoryStore().(*store.InMemoryStore)
	values := []uint16{0x1234, 0x5678}
	memStore.SetHoldingRegisters(values)

	request := Request{
		Frame:        []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x06, 0x01, 0x03, 0x00, 0x00, 0x00, 0x02},
		SlaveID:      0x01,
		FuncCode:     protocol.FuncCodeReadHoldingRegisters,
		StartAddress: 0,
		Quantity:     2,
	}

	response, err := handler.Handle(request, memStore)
	if err != nil {
		t.Fatalf("Failed to handle request: %v", err)
	}

	if response[7] != request.FuncCode {
		t.Errorf("Response function code mismatch: got %d, want %d", response[7], request.FuncCode)
	}
}

// TestHoldingRegistersHandler_Handle_Error tests error handling when getting holding registers fails
func TestHoldingRegistersHandler_Handle_Error(t *testing.T) {
	handler := &HoldingRegistersHandler{}
	memStore := &store.InMemoryStore{}
	request := Request{
		Frame:        []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x06, 0x01, 0x03, 0x00, 0x00, 0x00, 0x02},
		SlaveID:      0x01,
		FuncCode:     protocol.FuncCodeReadHoldingRegisters,
		StartAddress: 0,
		Quantity:     2,
	}

	response, err := handler.Handle(request, memStore)
	if err == nil {
		t.Fatalf("Expected an error, but got nil")
	}

	if response != nil {
		t.Errorf("Expected response to be nil, but got %v", response)
	}
}
