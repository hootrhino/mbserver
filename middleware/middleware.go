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

package middleware

import (
	modbus_server "mbserver"
)

// Since modbus_server.Handler is undefined, we'll create a placeholder type for it.
// This should be replaced with the actual definition once it's available.
type MockHandler func(request interface{}) ([]byte, error)

type Middleware func(next MockHandler) MockHandler

type Logger interface {
	LogRequest(request modbus_server.Request)
	LogResponse(response []byte, err error)
}
