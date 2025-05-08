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

package store

type Store interface {
	GetCoils(start, quantity uint16) ([]byte, error)
	GetDiscreteInputs(start, quantity uint16) ([]byte, error)
	GetHoldingRegisters(start, quantity uint16) ([]uint16, error)
	GetInputRegisters(start, quantity uint16) ([]uint16, error)
	SetCoils(values []byte) error
	SetDiscreteInputs(values []byte) error
	SetHoldingRegisters(values []uint16) error
	SetInputRegisters(values []uint16) error
	SetCoilsAt(start uint16, values []byte) error
	SetHoldingRegistersAt(start uint16, values []uint16) error // Add this line
}

