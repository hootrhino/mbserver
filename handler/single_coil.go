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
	"github.com/hootrhino/goodbusserver/protocol"
	"github.com/hootrhino/goodbusserver/store"
)

type SingleCoilHandler struct{}

func (h *SingleCoilHandler) Handle(request Request, store store.Store) ([]byte, error) {
	// Extract the coil value from the request frame
	coilValue := uint16(request.Frame[10])<<8 | uint16(request.Frame[11])
	var value byte
	switch coilValue {
	case 0xFF00:
		value = 1
	case 0x0000:
		value = 0
	default:
		return nil, protocol.ErrIllegalDataValue
	}

	// Write the coil value to the store
	err := store.SetCoilsAt(request.StartAddress, []byte{value})
	if err != nil {
		return nil, protocol.ErrIllegalDataAddress
	}

	// Construct the response PDU
	pdu := []byte{
		request.FuncCode,
		byte(request.StartAddress >> 8),
		byte(request.StartAddress),
		byte(coilValue >> 8),
		byte(coilValue),
	}

	// Extract transaction ID from the request frame
	transactionID := protocol.ExtractTransactionID(request.Frame)

	// Build MBAP header
	header := protocol.BuildResponseHeader(transactionID, 0, uint16(len(pdu)+1), request.SlaveID)

	// Combine header and PDU
	response := append(header, pdu...)
	return response, nil
}
