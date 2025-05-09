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

package protocol

import (
	"encoding/binary"
	"testing"
)

func TestEncodeUint16(t *testing.T) {
	value := uint16(0x1234)
	expected := make([]byte, 2)
	binary.BigEndian.PutUint16(expected, value)

	result := EncodeUint16(value)
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("EncodeUint16(%d) byte %d = %d; want %d", value, i, result[i], expected[i])
		}
	}
}

