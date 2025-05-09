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
	"mbserver/protocol"
	"mbserver/store"
)

type MultipleCoilsHandler struct{}

func (h *MultipleCoilsHandler) Handle(request Request, store store.Store) ([]byte, error) {
	byteCount := request.Frame[12]
	values := request.Frame[13 : 13+byteCount]

	err := store.SetCoilsAt(request.StartAddress, values)
	if err != nil {
		return nil, protocol.ErrIllegalDataAddress
	}

	// Construct the response PDU
	pdu := []byte{
		request.FuncCode,
		byte(request.StartAddress >> 8),
		byte(request.StartAddress),
		byte(request.Quantity >> 8),
		byte(request.Quantity),
	}

	// Extract transaction ID from the request frame
	transactionID := protocol.ExtractTransactionID(request.Frame)

	// Build MBAP header
	header := protocol.BuildResponseHeader(transactionID, 0, uint16(len(pdu)+1), request.SlaveID)

	// Combine header and PDU
	response := append(header, pdu...)
	return response, nil
}
