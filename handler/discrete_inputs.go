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

type DiscreteInputsHandler struct{}

func (h *DiscreteInputsHandler) Handle(request Request, store store.Store) ([]byte, error) {
	values, err := store.GetDiscreteInputs(request.StartAddress, request.Quantity)
	if err != nil {
		return nil, protocol.ErrIllegalDataAddress
	}

	// Calculate the actual byte count for discrete inputs
	// Similar to coils, discrete inputs are packed into bytes, each byte contains 8 inputs
	byteCount := (int(request.Quantity) + 7) / 8
	if byteCount > len(values) {
		byteCount = len(values)
	}

	// Extract transaction ID from the request frame
	transactionID := protocol.ExtractTransactionID(request.Frame)

	// Construct response PDU
	pdu := []byte{request.FuncCode, byte(byteCount)}
	pdu = append(pdu, values[:byteCount]...)

	// Build MBAP header
	header := protocol.BuildResponseHeader(transactionID, 0, uint16(len(pdu)+1), request.SlaveID)

	// Combine header and PDU
	response := append(header, pdu...)
	return response, nil
}
