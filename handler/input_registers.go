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

type InputRegistersHandler struct{}

func (h *InputRegistersHandler) Handle(request Request, store store.Store) ([]byte, error) {
	values, err := store.GetInputRegisters(request.StartAddress, request.Quantity)
	if err != nil {
		return nil, protocol.ErrIllegalDataAddress
	}

	// 验证数据长度
	if len(values) < int(request.Quantity) {
		return nil, protocol.ErrIllegalDataAddress
	}

	// 计算实际需要返回的字节数
	byteCount := int(request.Quantity) * 2

	// Extract transaction ID from the request frame
	transactionID := protocol.ExtractTransactionID(request.Frame)

	// Construct response PDU
	pdu := make([]byte, 0, 2+byteCount)
	pdu = append(pdu, request.FuncCode)
	pdu = append(pdu, byte(byteCount))
	
	// 将uint16数组转换为字节流
	for i := uint16(0); i < request.Quantity; i++ {
		value := values[i]
		pdu = append(pdu, byte(value>>8))
		pdu = append(pdu, byte(value))
	}

	// Build MBAP header
	header := protocol.BuildResponseHeader(transactionID, 0, uint16(len(pdu)+1), request.SlaveID)

	// Combine header and PDU
	response := append(header, pdu...)
	return response, nil
}
