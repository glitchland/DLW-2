package asm

import (
	"fmt"
	s "pkg/shared"
)

///
/// Argument type
///
type Argument struct {
	IsImmediate   bool
	IsLabel       bool
	IsRegister    bool
	IsAddress     bool
	IsDereference bool
	Register      uint8
	BaseRegister  uint8
	Address       uint8
	Offset        uint8
	Label         string
	LabelOffset   int8
	ImmediateInt  uint8
}

func (a *Argument) Init() {
	a.IsImmediate = false
	a.IsLabel = false
	a.IsRegister = false
	a.IsAddress = false
	a.IsDereference = false
	a.Register = 255
	a.BaseRegister = 255
	a.Address = 0
	a.Offset = 0
	a.ImmediateInt = 0
	a.Label = ""
	a.LabelOffset = 0
}

func (a *Argument) MakeLabel(label string) {
	a.IsLabel = true
	a.Label = label
}

func (a *Argument) SetLabelRelativeOffset(labelLineNumber uint8, currentLineNumber uint8) {
	a.LabelOffset = int8(labelLineNumber) - int8(currentLineNumber)
}

func (a *Argument) MakeRegister(register uint8) {
	a.IsRegister = true
	a.Register = register
}

func (a *Argument) MakeAddress(address uint8) {
	a.IsAddress = true
	a.Address = address
}

func (a *Argument) MakeDereference(baseRegister uint8, offset uint8) {
	a.IsDereference = true
	a.BaseRegister = baseRegister
	a.Offset = offset
}

func (a *Argument) MakeImmediate(immediate uint8) {
	a.IsImmediate = true
	a.ImmediateInt = immediate
}

func (a *Argument) RegIntToStr(register uint8) string {
	switch {
	case register == s.A:
		return "A"
	case register == s.B:
		return "B"
	case register == s.C:
		return "C"
	case register == s.D:
		return "D"
	default:
		return fmt.Sprintf("unknown register %d", register)
	}
}

func (a *Argument) ToString() string {
	s := ""
	s += fmt.Sprintf("It is a label.............: %v\n", a.IsLabel)
	s += fmt.Sprintf("It is a register..........: %v\n", a.IsRegister)
	s += fmt.Sprintf("It is a address...........: %v\n", a.IsAddress)
	s += fmt.Sprintf("It is a dereference.......: %v\n", a.IsDereference)
	s += fmt.Sprintf("It is a immediate.........: %v\n", a.IsImmediate)
	s += fmt.Sprintf("literal register value....: %s\n", a.RegIntToStr(a.Register))
	s += fmt.Sprintf("base register value.......: %s\n", a.RegIntToStr(a.BaseRegister))
	s += fmt.Sprintf("address value.............: %d\n", a.Address)
	s += fmt.Sprintf("offset value..............: %d\n", a.Offset)
	s += fmt.Sprintf("label value...............: '%s'\n", a.Label)
	if a.IsLabel {
		s += fmt.Sprintf("label offset..............: %d\n", a.LabelOffset)
	}
	s += fmt.Sprintf("immediate int value.......: %d\n", a.ImmediateInt)
	return s
}
