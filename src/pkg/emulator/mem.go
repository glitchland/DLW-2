package emu

type Mem struct {
	buf         []interface{}
	WriteLocked bool
}

func (m *Mem) Init(size uint16) {
	m.buf = make([]interface{}, size)
}

func (m *Mem) Zero16() {
	for i, _ := range m.buf {
		m.buf[i] = uint16(0)
	}
}

func (m *Mem) Zero8() {
	for i, _ := range m.buf {
		m.buf[i] = uint8(0)
	}
}

func (m *Mem) GetWordAt(addr uint8) uint16 {
	return m.buf[addr].(uint16)
}

func (m *Mem) SetWordAt(addr uint8, v uint16) {
	if !m.WriteLocked {
		m.buf[addr] = v
	}
}

func (m *Mem) GetByteAt(addr uint8) uint8 {
	return m.buf[addr].(uint8)
}

func (m *Mem) SetByteAt(addr uint8, v uint8) {
	m.buf[addr] = v
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
