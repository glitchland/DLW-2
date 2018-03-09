package emu

import (
	"fmt"
	s "pkg/shared"
)

///
/// This represents CPU
///
type Cpu struct {
	regs        Registers
	PC          uint8
	PCBranchMod bool // has the PC been modified via branch
	Instruction uint16
	alu         Alu
	rom         RomMem
	ram         RamMem
	halted      bool
}

func (c *Cpu) Init(romData []uint16) {
	c.rom.Init(romData)
	c.ram.Init()
	c.regs.Init()
	c.LoadInstruction(0)
	c.halted = false
}

///////////////////////////////////////////////////
// Store instructions
///////////////////////////////////////////////////

func (c *Cpu) StoreAtRegReference(srcr uint8, dstr uint8) {
	addr, err := c.regs.Read(dstr)
	s.ChkFatalError(err)

	v, err := c.regs.Read(srcr)
	s.ChkFatalError(err)

	c.ram.Write(addr, v)
}

func (c *Cpu) StoreAtRelative(srcr uint8, baser uint8, offset uint8) {
	base, err := c.regs.Read(baser)
	s.ChkFatalError(err)

	addr := base + offset

	v, err := c.regs.Read(srcr)
	s.ChkFatalError(err)

	c.ram.Write(addr, v)
}

func (c *Cpu) StoreAtAddr(addr uint8, rv uint8) {
	v, err := c.regs.Read(rv)
	s.ChkFatalError(err)

	c.ram.Write(addr, v)
}

///////////////////////////////////////////////////
// Load instructions
///////////////////////////////////////////////////

func (c *Cpu) LoadFromRegReference(srcr uint8, dstr uint8) {
	addr, err := c.regs.Read(srcr)
	s.ChkFatalError(err)

	v := c.ram.Read(addr)

	err = c.regs.Write(dstr, v)
	s.ChkFatalError(err)
}

func (c *Cpu) LoadFromRelative(baser uint8, offset uint8, dstr uint8) {
	base, err := c.regs.Read(baser)
	s.ChkFatalError(err)

	addr := base + offset
	v := c.ram.Read(addr)

	err = c.regs.Write(dstr, v)
	s.ChkFatalError(err)
}

func (c *Cpu) LoadFromAddr(addr uint8, dstr uint8) {
	v := c.ram.Read(addr)
	err := c.regs.Write(dstr, v)
	s.ChkFatalError(err)
}

///////////////////////////////////////////////////
// Add instructions
///////////////////////////////////////////////////
func (c *Cpu) AddImmediate(src1 uint8, imm uint8, dest uint8) {
	x, err := c.regs.Read(src1)
	s.ChkFatalError(err)
	y := imm
	v := c.alu.Add(x, y)
	err = c.regs.Write(dest, v)
	s.ChkFatalError(err)
}

func (c *Cpu) AddRegister(src1 uint8, src2 uint8, dest uint8) {
	x, err := c.regs.Read(src1)
	s.ChkFatalError(err)
	y, err := c.regs.Read(src2)
	s.ChkFatalError(err)
	v := c.alu.Add(x, y)
	err = c.regs.Write(dest, v)
	s.ChkFatalError(err)
}

///////////////////////////////////////////////////
// Sub instructions
///////////////////////////////////////////////////
func (c *Cpu) SubImmediate(src1 uint8, imm uint8, dest uint8) {
	x, err := c.regs.Read(src1)
	s.ChkFatalError(err)
	y := imm
	v := c.alu.Sub(x, y)
	err = c.regs.Write(dest, v)
	s.ChkFatalError(err)
}

func (c *Cpu) SubRegister(src1 uint8, src2 uint8, dest uint8) {
	x, err := c.regs.Read(src1)
	s.ChkFatalError(err)
	y, err := c.regs.Read(src2)
	s.ChkFatalError(err)
	v := c.alu.Sub(x, y)
	err = c.regs.Write(dest, v)
	s.ChkFatalError(err)
}

func (c *Cpu) Cycle() {
	if c.PCBranchMod { // go directly to the location of the branch
		c.LoadInstruction(c.PC)
		c.unsetBranchedPC()
	} else {
		c.IncPc()
		c.LoadInstruction(c.PC)
	}
}

func (c *Cpu) CurrentInstruction() uint16 {
	return c.Instruction
}

func (c *Cpu) LoadInstruction(addr uint8) {
	c.Instruction = c.rom.Read(addr)
}

func (c *Cpu) setBranchedPC() {
	c.PCBranchMod = true
}

func (c *Cpu) unsetBranchedPC() {
	c.PCBranchMod = false
}

///////////////////////////////////////////////////
// Branch unconditional
///////////////////////////////////////////////////
func (c *Cpu) BranchAddress(a uint8) {
	c.setBranchedPC()
	c.SetPc(a)
}

func (c *Cpu) BranchRelative(v uint8) {
	c.setBranchedPC()
	c.AddPc(v)
}

func (c *Cpu) BranchArithmetic(baseRegister uint8, offset uint8) {
	c.AdjustPC(baseRegister, offset)
}

///////////////////////////////////////////////////
// Branch zero flag
///////////////////////////////////////////////////
func (c *Cpu) BranchAddressCondZ(a uint8) {
	if c.alu.ZeroFlag() {
		c.setBranchedPC()
		c.SetPc(a)
	}
}

func (c *Cpu) BranchRelativeCondZ(v uint8) {
	if c.alu.ZeroFlag() {
		c.setBranchedPC()
		c.AddPc(v)
	}
}

func (c *Cpu) BranchArithmeticCondZ(baseRegister uint8, offset uint8) {
	if c.alu.ZeroFlag() {
		c.AdjustPC(baseRegister, offset)
	}
}

func (c *Cpu) IncPc() {
	c.PC++
	if c.PC >= c.rom.Top {
		c.HaltCpu()
	}
}

func (c *Cpu) DecPc() {
	c.PC++
}

func (c *Cpu) SetPc(v uint8) {
	c.PC = v
}

// this is signed so that we can use negative offsets
func (c *Cpu) AddPc(v uint8) {
	twsCmplmnt := ((^v) + 1)

	if (twsCmplmnt & 0x80) > 0 { // negative
		twsCmplmnt := twsCmplmnt & 0x7f
		c.PC -= twsCmplmnt
	} else { // positive
		c.PC += v
	}
}

func (c *Cpu) AdjustPC(baseRegister uint8, offset uint8) {
	v, err := c.regs.Read(baseRegister)
	s.ChkFatalError(err)
	c.PC = v + offset
}

func (c *Cpu) IsHalted() bool {
	return c.halted
}

func (c *Cpu) HaltCpu() {
	c.halted = true
}

func (c *Cpu) ClearLine() {
	for i := 1; i <= 200; i++ {
		fmt.Printf("  ")
	}
	fmt.Println()
}

func (c *Cpu) ClearScreen() {
	for i := 1; i <= 20; i++ {
		c.ClearLine()
	}
}

func (c *Cpu) PrintState() {
	fmt.Printf("\033[0;0H")
	fmt.Printf("+----------------------------------------------+\n")
	fmt.Printf("|CPU STATE                                     |\n")
	fmt.Printf("+----------------------------------------------+\n")
	fmt.Printf("|A:\t%02X B:\t%02X                             |\n", c.regs.A, c.regs.B)
	fmt.Printf("|C:\t%02X D:\t%02X                             |\n", c.regs.C, c.regs.D)
	fmt.Printf("+----------------------------------------------+\n")
	fmt.Printf("|Flags: ZF:%s OF:%s SF:%s                         |\n",
		s.BoolToIntStr(c.alu.ZeroFlag()),
		s.BoolToIntStr(c.alu.OverflowFlag()),
		s.BoolToIntStr(c.alu.SignFlag()))
	fmt.Printf("+----------------------------------------------+\n")
	fmt.Printf("|PC:\t[%02X] -> [%b]             |\n", c.PC, c.CurrentInstruction())
	fmt.Printf("+----------------------------------------------+\n")
	fmt.Printf("RAM:\n%s\n", c.ram.ToStr())
	fmt.Printf("ROM:\n%s\n", c.rom.ToStr())
}
