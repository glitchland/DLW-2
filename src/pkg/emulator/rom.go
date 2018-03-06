package emu

import (
	"fmt"
)

const RomAddrLimits = 0x100

type RomMem struct {
	rom Mem
	Top uint8
}

func (r *RomMem) Init(opcodes []uint16) {
	r.Top = (RomAddrLimits - 1)
	r.rom.Init(RomAddrLimits)
	r.rom.Zero16()
	r.Load(opcodes)
}

// ROM access, ROM is a series of 16 bit word opcodes
func (r *RomMem) Load(opcodes []uint16) {
	for i, opcode := range opcodes {
		r.Write(uint8(i), opcode)
	}
	r.Top = uint8(len(opcodes))
	r.rom.WriteLock()
}

func (r *RomMem) Print() {
	s := "ROM:\n"
	for i := 1; i <= RomAddrLimits; i++ {
		s += fmt.Sprintf("%04X ", r.rom.GetWordAt(uint8(i-1)))
		if i%16 == 0 {
			s += fmt.Sprintf("\n")
		}
	}
	fmt.Println(s)
}

func (r *RomMem) Write(addr uint8, v uint16) {
	r.rom.SetWordAt(uint8(addr), v)
}

func (r *RomMem) Read(addr uint8) uint16 {
	return r.rom.GetWordAt(uint8(addr))
}
