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

func TestCoilsHandler_Handle(t *testing.T) {
	handler := &CoilsHandler{}
	memStore := store.NewInMemoryStore().(*store.InMemoryStore)
	values := []byte{1, 0, 1}
	memStore.SetCoils(values)

	request := Request{
		Frame:        []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x06, 0x01, 0x01, 0x00, 0x00, 0x00, 0x03},
		SlaveID:      0x01,
		FuncCode:     protocol.FuncCodeReadCoils,
		StartAddress: 0,
		Quantity:     3,
	}

	response, err := handler.Handle(request, memStore)
	if err != nil {
		t.Fatalf("Failed to handle request: %v", err)
	}

	if response[7] != request.FuncCode {
		t.Errorf("Response function code mismatch: got %d, want %d", response[7], request.FuncCode)
	}
}
