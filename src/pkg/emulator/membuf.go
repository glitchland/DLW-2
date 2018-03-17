package emu

import (
	s "pkg/shared"
)

type Mem struct {
	buf         []interface{}
	WriteLocked bool
}

func (m *Mem) Init8(size uint16) {
	m.buf = make([]interface{}, size)
	m.zero8()
}

func (m *Mem) Init16(size uint16) {
	m.buf = make([]interface{}, size)
	m.zero16()
}

func (m *Mem) GetWordAt(addr uint8) (uint16, error) {
	e := m.boundsCheck(addr)
	if e != nil {
		return 0, e
	}

	return m.buf[addr].(uint16), nil
}

func (m *Mem) SetWordAt(addr uint8, v uint16) error {
	e := m.boundsCheck(addr)
	if e != nil {
		return e
	}

	if !m.WriteLocked {
		m.buf[addr] = v
	}
	return nil
}

func (m *Mem) GetByteAt(addr uint8) (uint8, error) {
	e := m.boundsCheck(addr)
	if e != nil {
		return 0, e
	}

	return m.buf[addr].(uint8), nil
}

func (m *Mem) SetByteAt(addr uint8, v uint8) error {
	e := m.boundsCheck(addr)
	if e != nil {
		return e
	}

	if !m.WriteLocked {
		m.buf[addr] = v
	}
	return nil
}

func (m *Mem) WriteLock() {
	m.WriteLocked = true
}

func (m *Mem) WriteUnlock() {
	m.WriteLocked = false
}

func (m *Mem) IsWriteLocked() bool {
	return m.WriteLocked
}

func (m *Mem) zero8() {
	for i, _ := range m.buf {
		m.buf[i] = uint8(0)
	}
}

func (m *Mem) zero16() {
	for i, _ := range m.buf {
		m.buf[i] = uint16(0)
	}
}

func (m *Mem) boundsCheck(addr uint8) error {
	if uint16(addr) > uint16(len(m.buf)) {
		return &s.MemoryError{"Out Of Range", addr}
	}
	return nil
}
