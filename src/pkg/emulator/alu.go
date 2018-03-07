package emu

// TODO: Change bit flipping functions to pass by reference and return

import (
	s "pkg/shared" // refactor this
)

const (
	OverflowFlagBitIndex = 0
	SignFlagBitIndex     = 1
	ZeroFlagBitIndex     = 2
	SignBitIndex         = 0
	MaxUint8             = uint16(^uint8(0))
)

type Alu struct {
	PSW uint16
}

func (a *Alu) Init() {
	a.PSW = 0
}

func (a *Alu) Add(x uint8, y uint8) uint8 {
	a.resetFlags()
	v := x + y

	a.checkAndSetAddOF(x, y)
	a.checkAndSetZF(v)
	a.checkAndSetSF(v)

	return v
}

func (a *Alu) Sub(x uint8, y uint8) uint8 {
	a.resetFlags()
	v := x - y

	a.checkAndSetAddOF(x, y)
	a.checkAndSetZF(v)
	a.checkAndSetSF(v)

	return v
}

func (a *Alu) SignFlag() bool {
	return s.GetBit(a.PSW, SignFlagBitIndex)
}

func (a *Alu) ZeroFlag() bool {
	return s.GetBit(a.PSW, ZeroFlagBitIndex)
}

func (a *Alu) OverflowFlag() bool {
	return s.GetBit(a.PSW, OverflowFlagBitIndex)
}

func (a *Alu) resetFlags() {
	a.unsetOverflowFlag()
	a.unsetZeroFlag()
	a.unsetSignFlag()
}

func (a *Alu) checkAndSetAddOF(x uint8, y uint8) {
	if uint16(x)+uint16(y) > MaxUint8 { //0x100
		a.setOverflowFlag()
	}
}

func (a *Alu) checkAndSetSubOF(x uint8, y uint8) {
	if uint16(x)-uint16(y) > MaxUint8 { //0xffff
		a.setOverflowFlag()
	}
}

func (a *Alu) checkAndSetSF(v uint8) {
	if (v & (1 << 7)) > 0 {
		a.setSignFlag()
	}
}

func (a *Alu) checkAndSetZF(v uint8) {
	if v == 0 {
		a.setZeroFlag()
	}
}

//set sign flag
func (a *Alu) setSignFlag() {
	s.SetBit(&a.PSW, SignFlagBitIndex)
}

func (a *Alu) unsetSignFlag() {
	s.UnsetBit(&a.PSW, SignFlagBitIndex)
}

//set overflow flag
func (a *Alu) setOverflowFlag() {
	s.SetBit(&a.PSW, OverflowFlagBitIndex)
}

func (a *Alu) unsetOverflowFlag() {
	s.UnsetBit(&a.PSW, OverflowFlagBitIndex)
}

//set zero flag
func (a *Alu) setZeroFlag() {
	s.SetBit(&a.PSW, ZeroFlagBitIndex)
}

func (a *Alu) unsetZeroFlag() {
	s.UnsetBit(&a.PSW, ZeroFlagBitIndex)
}
