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
	"sync"
)

type InMemoryStore struct {
	coils            []byte
	discreteInputs   []byte
	holdingRegisters []uint16
	inputRegisters   []uint16
	mu               sync.RWMutex
}

// GetHoldingRegisters implements Store.
func (s *InMemoryStore) GetHoldingRegisters(start uint16, quantity uint16) ([]uint16, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	startIdx := int(start)
	endIdx := startIdx + int(quantity)

	if startIdx < 0 || endIdx > len(s.holdingRegisters) {
		return nil, ErrInvalidAddress
	}

	return s.holdingRegisters[startIdx:endIdx], nil
}

// GetInputRegisters implements Store.
func (s *InMemoryStore) GetInputRegisters(start uint16, quantity uint16) ([]uint16, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    startIdx := int(start)
    endIdx := startIdx + int(quantity)

    if startIdx < 0 || endIdx > len(s.inputRegisters) {
        return nil, ErrInvalidAddress
    }

    return s.inputRegisters[startIdx:endIdx], nil
}

// SetDiscreteInputs implements Store.
func (s *InMemoryStore) SetDiscreteInputs(values []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.discreteInputs = values
	return nil
}

// SetHoldingRegisters implements Store.
func (s *InMemoryStore) SetHoldingRegisters(values []uint16) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.holdingRegisters = values
	return nil
}

// SetInputRegisters implements Store.
func (s *InMemoryStore) SetInputRegisters(values []uint16) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.inputRegisters = values
	return nil
}

func (s *InMemoryStore) SetHoldingRegistersAt(start uint16, values []uint16) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	startIdx := int(start)
	endIdx := startIdx + len(values)

	if startIdx < 0 || endIdx > len(s.holdingRegisters) {
		return ErrInvalidAddress
	}

	copy(s.holdingRegisters[startIdx:endIdx], values)
	return nil
}

func NewInMemoryStore() Store {
	defaultDiscreteInputsSize := 100 // You can adjust this value
	defaultCoilsSize := 100          // You can adjust this value
	return &InMemoryStore{
		coils:            make([]byte, defaultCoilsSize),
		discreteInputs:   make([]byte, defaultDiscreteInputsSize),
		holdingRegisters: make([]uint16, 0),
		inputRegisters:   make([]uint16, 0),
	}
}

// The following method implementation is added to make *InMemoryStore implement the Store interface
// GetDiscreteInputs method implementation
func (s *InMemoryStore) GetDiscreteInputs(start, quantity uint16) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	startIdx := int(start)
	endIdx := startIdx + int(quantity)

	if startIdx < 0 || startIdx >= len(s.discreteInputs) || endIdx > len(s.discreteInputs) {
		return nil, ErrInvalidAddress
	}

	return s.discreteInputs[startIdx:endIdx], nil
}
func (s *InMemoryStore) GetCoils(start, quantity uint16) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	end := int(start) + int(quantity)
	if end > len(s.coils) {
		return nil, ErrInvalidAddress
	}
	return s.coils[start:end], nil
}

func (s *InMemoryStore) SetCoils(values []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.coils = values
	return nil
}

func (s *InMemoryStore) SetCoilsAt(start uint16, values []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	startIdx := int(start)
	endIdx := startIdx + len(values)

	if startIdx < 0 || endIdx > len(s.coils) {
		return ErrInvalidAddress
	}

	copy(s.coils[startIdx:endIdx], values)
	return nil
}

var ErrInvalidAddress = &StoreError{Code: "INVALID_ADDRESS", Message: "Invalid address"}

type StoreError struct {
	Code    string
	Message string
}

func (e *StoreError) Error() string {
	return e.Message
}
