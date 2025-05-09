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

import (
	"testing"
)

func TestInMemoryStore_SetGetCoils(t *testing.T) {
	store := NewInMemoryStore().(*InMemoryStore)
	values := []byte{1, 0, 1}

	// 设置线圈值
	err := store.SetCoils(values)
	if err != nil {
		t.Fatalf("Failed to set coils: %v", err)
	}

	// 获取线圈值
	result, err := store.GetCoils(0, uint16(len(values)))
	if err != nil {
		t.Fatalf("Failed to get coils: %v", err)
	}

	for i := range values {
		if result[i] != values[i] {
			t.Errorf("Coil value at index %d mismatch: got %d, want %d", i, result[i], values[i])
		}
	}
}

func TestInMemoryStore_SetGetHoldingRegisters(t *testing.T) {
	store := NewInMemoryStore().(*InMemoryStore)
	values := []uint16{0x1234, 0x5678}

	// 设置保持寄存器值
	err := store.SetHoldingRegisters(values)
	if err != nil {
		t.Fatalf("Failed to set holding registers: %v", err)
	}

	// 获取保持寄存器值
	result, err := store.GetHoldingRegisters(0, uint16(len(values)))
	if err != nil {
		t.Fatalf("Failed to get holding registers: %v", err)
	}

	for i := range values {
		if result[i] != values[i] {
			t.Errorf("Holding register value at index %d mismatch: got %d, want %d", i, result[i], values[i])
		}
	}
}

