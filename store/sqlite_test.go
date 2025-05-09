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
	"os"
	"testing"
)

func TestNewSqliteStore(t *testing.T) {
	dsn := "test.db"
	defer os.Remove(dsn)

	store, err := NewSqliteStore(dsn)
	if err != nil {
		t.Fatalf("NewSqliteStore() error = %v", err)
	}
	defer store.Close()

	if store.db == nil {
		t.Error("Expected database to be initialized, but got nil")
	}
}

func TestGetCoils(t *testing.T) {
	dsn := "test.db"
	defer os.Remove(dsn)

	store, err := NewSqliteStore(dsn)
	if err != nil {
		t.Fatalf("NewSqliteStore() error = %v", err)
	}
	defer store.Close()

	_, err = store.db.Exec("INSERT INTO coils (address, value) VALUES (0, 1), (1, 0), (2, 1)")
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	start := uint16(0)
	quantity := uint16(3)
	values, err := store.GetCoils(start, quantity)
	if err != nil {
		t.Fatalf("GetCoils() error = %v", err)
	}

	expected := []byte{1, 0, 1}
	for i := range expected {
		if values[i] != expected[i] {
			t.Errorf("GetCoils() got %v, want %v", values, expected)
			break
		}
	}
}
