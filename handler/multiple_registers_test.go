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

	"mbserver/protocol"
	"mbserver/store"
)

func TestMultipleRegistersHandler_Handle(t *testing.T) {
	handler := &MultipleRegistersHandler{}
	memStore := store.NewInMemoryStore().(*store.InMemoryStore)
	request := Request{
		Frame:        []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x06, 0x01, 0x10, 0x00, 0x00, 0x00, 0x02},
		SlaveID:      0x01,
		FuncCode:     protocol.FuncCodeWriteMultipleRegisters,
		StartAddress: 0,
	}

	response, err := handler.Handle(request, memStore)
	if err != nil {
		t.Fatalf("Failed to handle request: %v", err)
	}

	if response[7] != request.FuncCode {
		t.Errorf("Response function code mismatch: got %d, want %d", response[7], request.FuncCode)
	}
}

func TestMultipleRegistersHandler_Handle_Error(t *testing.T) {
	handler := &MultipleRegistersHandler{}
	memStore := &store.InMemoryStore{}
	request := Request{
		Frame:        []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x06, 0x01, 0x10, 0x00, 0x00, 0x00, 0x02},
		SlaveID:      0x01,
		FuncCode:     protocol.FuncCodeWriteMultipleRegisters,
		StartAddress: 0,
	}

	response, err := handler.Handle(request, memStore)
	if err == nil {
		t.Fatalf("Expected an error, but got nil")
	}

	if response != nil {
		t.Errorf("Expected response to be nil, but got %v", response)
	}
}
