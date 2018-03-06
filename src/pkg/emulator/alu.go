package emu

// TODO: Change bit flipping functions to pass by reference and return 

import (
  		s "pkg/shared"    // refactor this     
        )

const (
	OverflowFlagBitIndex = 0
	SignFlagBitIndex     = 1
	ZeroFlagBitIndex     = 2
)

type Alu struct {
	PSW     uint16
}

func (a *Alu) Init() {
	a.PSW = 0
}

func (a *Alu) Add(x uint8, y uint8) uint8 {
	v := x + y
	return v
}

func (a *Alu) Sub(x uint8, y uint8) uint8 {
	v := x - y
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
