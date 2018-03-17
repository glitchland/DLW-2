package emu

import (
	"testing"
)

const AddrLimits = 0x100

var wordWidthSuccessTests = []struct {
	v    uint16
	addr uint8
}{
	{0xffff, 0},
	{0x8000, 0},
	{0x0001, (AddrLimits - 1)},
}

func TestSuccessSetGetWord(t *testing.T) {
	var m Mem
	var e error
	var r uint16
	m.Init16(AddrLimits)

	for _, tst := range wordWidthSuccessTests {

		e = m.SetWordAt(tst.addr, tst.v)
		if e != nil {
			t.Errorf("TestSetGetWord failed to set %x at address %x", tst.v, tst.addr)
		}

		r, e = m.GetWordAt(tst.addr)

		if e != nil {
			t.Errorf("TestSetGetWord failed to get %x at address %x", tst.v, tst.addr)
		}

		if r != tst.v {
			t.Errorf("TestSetGetWord got %x, want %x", r, tst.v)
		}
	}
}

var byteWidthSuccessTests = []struct {
	v    uint8
	addr uint8
}{
	{0xff, 0},
	{0x80, 0},
	{0x01, (AddrLimits - 1)},
}

func TestSuccessSetGetByte(t *testing.T) {
	var m Mem
	var e error
	var r uint8
	m.Init8(AddrLimits)

	for _, tst := range byteWidthSuccessTests {

		e = m.SetByteAt(tst.addr, tst.v)
		if e != nil {
			t.Errorf("TestSetGetByte failed to set %x at address %x", tst.v, tst.addr)
		}

		r, e = m.GetByteAt(tst.addr)

		if e != nil {
			t.Errorf("TestSetGetByte failed to get %x at address %x", tst.v, tst.addr)
		}

		if r != tst.v {
			t.Errorf("TestSetGetByte got %x, want %x", r, tst.v)
		}
	}
}

func TestWriteLock(t *testing.T) {
	var m8 Mem
	var m16 Mem

	m8.Init8(AddrLimits)
	m16.Init16(AddrLimits)

	m8.WriteLock()
	v := m8.IsWriteLocked()
	if !v {
		t.Errorf("TestWriteLock(8) - locking - got %v, expected %v", v, true)
	}

	m8.WriteUnlock()
	v = m8.IsWriteLocked()
	if v {
		t.Errorf("TestWriteLock(8) - unlocking - got %v, expected %v", v, true)
	}

	m16.WriteLock()
	v = m16.IsWriteLocked()
	if !v {
		t.Errorf("TestWriteLock(16) - locking - got %v, expected %v", v, true)
	}

	m16.WriteUnlock()
	v = m16.IsWriteLocked()
	if v {
		t.Errorf("TestWriteLock(16) - unlocking - got %v, expected %v", v, true)
	}
}
