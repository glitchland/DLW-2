package emu

import (
	s "pkg/shared"
)

type Registers struct {
	A uint8
	B uint8
	C uint8
	D uint8
}

func (r *Registers) Init() {
	r.A = uint8(0)
	r.B = uint8(0)
	r.C = uint8(0)
	r.D = uint8(0)
}

func (r *Registers) Read(reg uint8) (uint8, error) {
	switch {
	case reg == s.A:
		return r.A, nil
	case reg == s.B:
		return r.B, nil
	case reg == s.C:
		return r.C, nil
	case reg == s.D:
		return r.D, nil
	default:
		return 0, &s.EmulatorError{"Fatal read attempt from unknown register"}
	}
}

func (r *Registers) Write(reg uint8, v uint8) error {
	switch {
	case reg == s.A:
		r.A = v
	case reg == s.B:
		r.B = v
	case reg == s.C:
		r.C = v
	case reg == s.D:
		r.D = v
	default:
		return &s.EmulatorError{"Fatal write attempt to unknown register"}
	}
	return nil
}

func (r *Registers) SetA(v uint8) {
	r.A = v
}

func (r *Registers) SetB(v uint8) {
	r.B = v
}

func (r *Registers) SetC(v uint8) {
	r.C = v
}

func (r *Registers) SetD(v uint8) {
	r.D = v
}
