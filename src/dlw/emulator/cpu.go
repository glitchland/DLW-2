package emu

import (
        "fmt"
        )

///
/// This represents CPU
///
type Cpu struct {
	A       uint8
	B       uint8
	C       uint8
	D       uint8
	PC      uint8
	Instruction uint16	
	mem     Memory
	halted  bool
}

// TODO, set the top of the rom
// also set the end insruction in the code after loading
// halt the CPU if the PC increments past the end instruction

func (c *Cpu) Init(romData []uint16) {
	c.mem.Init()
	c.mem.LoadRom(romData)
	c.LoadInstruction(0)
	c.halted = false
}

func (c *Cpu) NextInstruction() {
	c.IncPc()
	c.Instruction = c.mem.GetRomWordAt(c.PC)
	if (c.PC >= c.mem.RomTop) {
		c.HaltCpu()
	}	
}

func (c *Cpu) CurrentInstruction() uint16 {
	return c.Instruction
}

func (c *Cpu) LoadInstruction(addr uint8) {
	c.Instruction = c.mem.GetRomWordAt(addr)
}

func (c *Cpu) IncPc() {
	c.PC++
}

func (c *Cpu) SetPc(v uint8) {
	c.PC = v 
}

func (c *Cpu) SetA(v uint8) {
	c.A = v
}

func (c *Cpu) SetB(v uint8) {
	c.A = v
}

func (c *Cpu) SetC(v uint8) {
	c.A = v
}

func (c *Cpu) SetD(v uint8) {
	c.A = v
}

func (c *Cpu) IsHalted() bool {
	return c.halted
}

func (c *Cpu) HaltCpu () {
	c.halted = true
}

func (c *Cpu) PrintState() {
	i := 0
	b := uint16(0)
	s := ""
	s += fmt.Sprintf("Registers: \n")
	s += fmt.Sprintf("A: %v B: %v C: %v D: %v \n", c.A, c.B, c.C, c.D)
	s += fmt.Sprintf("PC: %v \n", c.PC)
	s += fmt.Sprintf("Memory:\n")
	fmt.Println(len(c.mem.Ram))
	for _, b = range c.mem.Ram {
		i++	
        s += fmt.Sprintf("%04x ", b)
		if (i % 8 == 0) {
			s += fmt.Sprintf("\n")
		}
	}
	fmt.Println(s)
}