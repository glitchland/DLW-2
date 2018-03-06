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
	r.ram.Init(RamAddrLimits)
	r.ram.Zero8()
}

func (r *RamMem) Print() {
	s := "RAM:\n"
	for i := 1; i <= RamAddrLimits; i++ {
		s += fmt.Sprintf("%02X ", r.ram.GetByteAt(uint8(i-1)))
		if i%16 == 0 {
			s += fmt.Sprintf("\n")
		}
	}
	fmt.Println(s)
}

func (r *RamMem) Write(addr uint8, v uint8) {
	r.ram.SetByteAt(uint8(addr), v)
}

func (r *RamMem) Read(addr uint8) uint8 {
	return r.ram.GetByteAt(uint8(addr))
}
