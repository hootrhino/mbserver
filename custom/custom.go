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

package custom

import (
	"modbus_server"
)

type CustomRequestHandler struct {
	handler func(modbus_server.Request)
}

func NewCustomRequestHandler(h func(modbus_server.Request)) *CustomRequestHandler {
	return &CustomRequestHandler{handler: h}
}

func (h *CustomRequestHandler) Handle(request modbus_server.Request) {
	if h.handler != nil {
		h.handler(request)
	}
}

