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
	rows, err := s.db.Query("SELECT value FROM coils WHERE address BETWEEN ? AND ?", start, start+quantity-1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []byte
	for rows.Next() {
		var val int
		if err := rows.Scan(&val); err != nil {
			return nil, err
		}
		values = append(values, byte(val))
	}
	return values, nil
}

// Implement other methods similarly...

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

