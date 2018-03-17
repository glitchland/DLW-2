package emu

import (
	"fmt"
)

const RamAddrLimits = 0x100

type RamMem struct {
	ram Mem
	Top uint8
}

func (r *RamMem) Init() {
	r.Top = uint8(RamAddrLimits - 1)
	r.ram.Init8(RamAddrLimits)
}

func (r *RamMem) ToStr() string {
	s := ""
	for i := 1; i <= RamAddrLimits; i++ {
		v, e := r.ram.GetByteAt(uint8(i - 1))
		if e != nil {
			panic("Unable to read memory")
		}
		s += fmt.Sprintf("%02X ", v)
		if i%16 == 0 {
			s += fmt.Sprintf("\n")
		}
	}
	return s
}

//XXX check errors
func (r *RamMem) Write(addr uint8, v uint8) {
	_ = r.ram.SetByteAt(uint8(addr), v)
}

//XXX check errors
func (r *RamMem) Read(addr uint8) uint8 {
	v, _ := r.ram.GetByteAt(uint8(addr))
	return v
}
