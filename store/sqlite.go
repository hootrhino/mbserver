package store

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteStore struct {
	db *sql.DB
}

func NewSqliteStore(dsn string) (*SqliteStore, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS coils (
			address INTEGER PRIMARY KEY,
			value INTEGER
		);
		CREATE TABLE IF NOT EXISTS discrete_inputs (
			address INTEGER PRIMARY KEY,
			value INTEGER
		);
		CREATE TABLE IF NOT EXISTS holding_registers (
			address INTEGER PRIMARY KEY,
			value INTEGER
		);
		CREATE TABLE IF NOT EXISTS input_registers (
			address INTEGER PRIMARY KEY,
			value INTEGER
		);
	`)
	if err != nil {
		return nil, err
	}

	return &SqliteStore{db: db}, nil
}

func (s *SqliteStore) GetCoils(start, quantity uint16) ([]byte, error) {
	if quantity == 0 {
		return []byte{}, nil
	}
	
	rows, err := s.db.Query("SELECT value FROM coils WHERE address BETWEEN ? AND ? ORDER BY address", start, start+quantity-1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := make([]byte, 0, quantity)
	for rows.Next() {
		var val int
		if err := rows.Scan(&val); err != nil {
			return nil, err
		}
		values = append(values, byte(val))
	}
	
	// 填充缺失的地址
	if len(values) < int(quantity) {
		missing := int(quantity) - len(values)
		values = append(values, make([]byte, missing)...)
	}
	
	return values, nil
}

func (s *SqliteStore) GetDiscreteInputs(start, quantity uint16) ([]byte, error) {
	if quantity == 0 {
		return []byte{}, nil
	}
	
	rows, err := s.db.Query("SELECT value FROM discrete_inputs WHERE address BETWEEN ? AND ? ORDER BY address", start, start+quantity-1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := make([]byte, 0, quantity)
	for rows.Next() {
		var val int
		if err := rows.Scan(&val); err != nil {
			return nil, err
		}
		values = append(values, byte(val))
	}
	
	if len(values) < int(quantity) {
		missing := int(quantity) - len(values)
		values = append(values, make([]byte, missing)...)
	}
	
	return values, nil
}

func (s *SqliteStore) GetHoldingRegisters(start, quantity uint16) ([]uint16, error) {
	if quantity == 0 {
		return []uint16{}, nil
	}
	
	rows, err := s.db.Query("SELECT value FROM holding_registers WHERE address BETWEEN ? AND ? ORDER BY address", start, start+quantity-1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := make([]uint16, 0, quantity)
	for rows.Next() {
		var val int
		if err := rows.Scan(&val); err != nil {
			return nil, err
		}
		values = append(values, uint16(val))
	}
	
	if len(values) < int(quantity) {
		missing := int(quantity) - len(values)
		values = append(values, make([]uint16, missing)...)
	}
	
	return values, nil
}

func (s *SqliteStore) GetInputRegisters(start, quantity uint16) ([]uint16, error) {
	if quantity == 0 {
		return []uint16{}, nil
	}
	
	rows, err := s.db.Query("SELECT value FROM input_registers WHERE address BETWEEN ? AND ? ORDER BY address", start, start+quantity-1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := make([]uint16, 0, quantity)
	for rows.Next() {
		var val int
		if err := rows.Scan(&val); err != nil {
			return nil, err
		}
		values = append(values, uint16(val))
	}
	
	if len(values) < int(quantity) {
		missing := int(quantity) - len(values)
		values = append(values, make([]uint16, missing)...)
	}
	
	return values, nil
}

func (s *SqliteStore) SetCoils(values []byte) error {
	if len(values) == 0 {
		return nil
	}
	
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO coils(address, value) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i, val := range values {
		if _, err := stmt.Exec(i, val); err != nil {
			return err
		}
	}
	
	return tx.Commit()
}

func (s *SqliteStore) SetDiscreteInputs(values []byte) error {
	if len(values) == 0 {
		return nil
	}
	
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO discrete_inputs(address, value) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i, val := range values {
		if _, err := stmt.Exec(i, val); err != nil {
			return err
		}
	}
	
	return tx.Commit()
}

func (s *SqliteStore) SetHoldingRegisters(values []uint16) error {
	if len(values) == 0 {
		return nil
	}
	
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO holding_registers(address, value) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i, val := range values {
		if _, err := stmt.Exec(i, val); err != nil {
			return err
		}
	}
	
	return tx.Commit()
}

func (s *SqliteStore) SetInputRegisters(values []uint16) error {
	if len(values) == 0 {
		return nil
	}
	
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO input_registers(address, value) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i, val := range values {
		if _, err := stmt.Exec(i, val); err != nil {
			return err
		}
	}
	
	return tx.Commit()
}

func (s *SqliteStore) SetCoilsAt(start uint16, values []byte) error {
	if len(values) == 0 {
		return nil
	}
	
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO coils(address, value) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i, val := range values {
		if _, err := stmt.Exec(start+uint16(i), val); err != nil {
			return err
		}
	}
	
	return tx.Commit()
}

func (s *SqliteStore) SetHoldingRegistersAt(start uint16, values []uint16) error {
	if len(values) == 0 {
		return nil
	}
	
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO holding_registers(address, value) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i, val := range values {
		if _, err := stmt.Exec(start+uint16(i), val); err != nil {
			return err
		}
	}
	
	return tx.Commit()
}

func (s *SqliteStore) Close() error {
	return s.db.Close()
}

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

