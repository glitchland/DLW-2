package emu

import (
	s "pkg/shared"
	"testing"
)

func TestInit(t *testing.T) {
	var r Registers
	var ex uint8 = 0

	r.Init()
	if r.A != ex {
		t.Errorf("r.A(%x) want %x", r.A, ex)
	}
	if r.B != ex {
		t.Errorf("r.B(%x) want %x", r.B, ex)
	}
	if r.C != ex {
		t.Errorf("r.C(%x) want %x", r.C, ex)
	}
	if r.D != ex {
		t.Errorf("r.D(%x) want %x", r.D, ex)
	}
}

var testSet = []struct {
	v  uint8
	ex uint8
}{
	{255, 255},
	{0, 0},
	{0x80, 0x80},
	{47, 47},
}

func TestSetA(t *testing.T) {
	var r Registers
	r.Init()

	for _, tst := range testSet {
		r.SetA(tst.v)
		if r.A != tst.ex {
			t.Errorf("[Direct]SetA(%x) got %x, want %x", tst.v, r.A, tst.ex)
		}
		rv, _ := r.Read(s.A)
		if rv != tst.ex {
			t.Errorf("[Read]SetA(%x) got %x, want %x", tst.v, rv, tst.ex)
		}
	}
}

func TestSetB(t *testing.T) {
	var r Registers
	r.Init()

	for _, tst := range testSet {
		r.SetB(tst.v)
		if r.B != tst.ex {
			t.Errorf("[Direct]SetB(%x) got %x, want %x", tst.v, r.B, tst.ex)
		}
		rv, _ := r.Read(s.B)
		if rv != tst.ex {
			t.Errorf("[Read]SetB(%x) got %x, want %x", tst.v, rv, tst.ex)
		}
	}
}

func TestSetC(t *testing.T) {
	var r Registers
	r.Init()

	for _, tst := range testSet {
		r.SetC(tst.v)
		if r.C != tst.ex {
			t.Errorf("[Direct]SetC(%x) got %x, want %x", tst.v, r.C, tst.ex)
		}
		rv, _ := r.Read(s.C)
		if rv != tst.ex {
			t.Errorf("[Read]SetC(%x) got %x, want %x", tst.v, rv, tst.ex)
		}
	}
}

func TestSetD(t *testing.T) {
	var r Registers
	r.Init()

	for _, tst := range testSet {
		r.SetD(tst.v)
		if r.D != tst.ex {
			t.Errorf("[Direct]SetD(%x) got %x, want %x", tst.v, r.D, tst.ex)
		}
		rv, _ := r.Read(s.D)
		if rv != tst.ex {
			t.Errorf("[Read]SetD(%x) got %x, want %x", tst.v, rv, tst.ex)
		}
	}
}

// This tests Read && Write
func TestWrite(t *testing.T) {
	var regList = []uint8{s.A, s.B, s.C, s.D}
	var r Registers

	r.Init()
	for _, reg := range regList {
		for _, tst := range testSet {
			r.Write(reg, tst.v)
			rv, _ := r.Read(reg)
			if rv != tst.ex {
				t.Errorf("[TestRead|Write]Reg(%x)(%x) got %x, want %x", reg, tst.v, rv, tst.ex)
			}
		}
	}

	// test error
	err := r.Write(255, 1)
	if err == nil {
		t.Errorf("Expected error when writing to non-existent register")
	}
	_, err = r.Read(255)
	if err == nil {
		t.Errorf("Expected error when writing to non-existent register")
	}
}
