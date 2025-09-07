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

type CoilsHandler struct{}

func (h *CoilsHandler) Handle(request Request, store store.Store) ([]byte, error) {
	values, err := store.GetCoils(request.StartAddress, request.Quantity)
	if err != nil {
		return nil, protocol.ErrIllegalDataAddress
	}

	// 验证数据长度
	expectedLength := (int(request.Quantity) + 7) / 8
	if len(values) < expectedLength {
		return nil, protocol.ErrIllegalDataAddress
	}

	// 计算实际需要返回的字节数
	byteCount := expectedLength

	// Extract transaction ID from the request frame
	transactionID := protocol.ExtractTransactionID(request.Frame)

	// Construct response PDU
	pdu := make([]byte, 0, 2+byteCount)
	pdu = append(pdu, request.FuncCode)
	pdu = append(pdu, byte(byteCount))
	pdu = append(pdu, values[:byteCount]...)

	// Build MBAP header
	header := protocol.BuildResponseHeader(transactionID, 0, uint16(len(pdu)+1), request.SlaveID)

	// Combine header and PDU
	response := append(header, pdu...)
	return response, nil
}
