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
	r.rom.Init16(RomAddrLimits)
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

func (r *RomMem) ToStr() string {
	s := ""
	for i := 1; i <= RomAddrLimits; i++ {
		v, e := r.rom.GetWordAt(uint8(i - 1))
		if e != nil {
			panic("Unable to read memory")
		}
		s += fmt.Sprintf("%04X ", v)
		if i%16 == 0 {
			s += fmt.Sprintf("\n")
		}
	}
	return s
}

//XXX check errors
func (r *RomMem) Write(addr uint8, v uint16) {
	_ = r.rom.SetWordAt(uint8(addr), v)
}

//XXX check errors
func (r *RomMem) Read(addr uint8) uint16 {
	v, _ := r.rom.GetWordAt(uint8(addr))
	return v
}
