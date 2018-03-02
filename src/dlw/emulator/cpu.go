package emu

import (
        "fmt"
  		s "dlw/shared"         
        )

const (
	OverflowFlagBitIndex = 0
	SignFlagBitIndex     = 1
	ZeroFlagBitIndex     = 2
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
	PCBranchMod bool // has the PC been modified via branch
	PSW     uint16
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

/*
  These Processor Status Word methods should be private.
 */

//set sign flag
func (c *Cpu) setSignFlag() {
	s.SetBit(&c.PSW, SignFlagBitIndex)
}

func (c *Cpu) unsetSignFlag() {
	s.UnsetBit(&c.PSW, SignFlagBitIndex)
}

func (c *Cpu) getSignFlag() bool {
	return s.GetBit(&c.PSW, SignFlagBitIndex)
}

//set overflow flag
func (c *Cpu) setOverflowFlag() {
	s.SetBit(&c.PSW, OverflowFlagBitIndex)
}

func (c *Cpu) unsetOverflowFlag() {
	s.UnsetBit(&c.PSW, OverflowFlagBitIndex)
}

func (c *Cpu) getOverflowFlag() bool {
	return s.GetBit(&c.PSW, OverflowFlagBitIndex)
}

//set zero flag
func (c *Cpu) setZeroFlag() {
	s.SetBit(&c.PSW, ZeroFlagBitIndex)
}

func (c *Cpu) unsetZeroFlag() {
	s.UnsetBit(&c.PSW, ZeroFlagBitIndex)
}

func (c *Cpu) getZeroFlag() bool {
	return s.GetBit(&c.PSW, ZeroFlagBitIndex)
}
/////////////////////////////////////////  

func (c *Cpu) readReg(r uint8) uint8 {
	switch {
	case r == s.A:
		return c.A
	case r == s.B:
		return c.B
	case r == s.C:
		return c.C
	case r == s.D:
		return c.D				
 	default:
 		panic("Unknown register") // XXX Handle Errors
 	}
}

func (c *Cpu) writeReg(r uint8, v uint8) {
	switch {
	case r == s.A:
		c.A = v
	case r == s.B:
		c.B = v
	case r == s.C:
		c.C = v
	case r == s.D:
		c.D = v			
 	default:
 		panic("Unknown register") // XXX Handle Errors
 	}
}

///////////////////////////////////////////////////
// Store instructions
///////////////////////////////////////////////////

func (c *Cpu) StoreAtRegReference(srcr uint8, dstr uint8) {
	addr := c.readReg(dstr)
	v := c.readReg(srcr)
	c.mem.SetRamByteAt(addr, v)
}

func (c *Cpu) StoreAtRelative(srcr uint8, baser uint8, offset uint8) {
	base := c.readReg(baser)
	addr := base + offset
	v := c.readReg(srcr)
	c.mem.SetRamByteAt(addr, v) 
}

func (c *Cpu) StoreAtAddr(addr uint8, rv uint8) {
	v := c.readReg(rv)
	c.mem.SetRamByteAt(addr, v)
}

///////////////////////////////////////////////////
// Load instructions
///////////////////////////////////////////////////

func (c *Cpu) LoadFromRegReference(srcr uint8, dstr uint8) {
	addr := c.readReg(srcr)
	v := c.mem.GetRamByteAt(addr)
	c.writeReg(dstr, v)
}

func (c *Cpu) LoadFromRelative(baser uint8, offset uint8, dstr uint8) {
	base := c.readReg(baser)
	addr := base + offset
	v := c.mem.GetRamByteAt(addr)	
	c.writeReg(dstr, v)	
}

func (c *Cpu) LoadFromAddr(addr uint8, dstr uint8) {
	v := c.mem.GetRamByteAt(addr)
	c.writeReg(dstr, v)	
}

///////////////////////////////////////////////////
// Add instructions
///////////////////////////////////////////////////
func (c *Cpu) AddImmediate (src1 uint8, imm uint8, dest uint8) {
	v := c.readReg(src1) + imm 
	c.writeReg(dest, v)
}

func (c *Cpu) AddRegister (src1 uint8, src2 uint8, dest uint8) {
	v := c.readReg(src1) + c.readReg(src2)
	c.writeReg(dest, v)
}

///////////////////////////////////////////////////
// Sub instructions
///////////////////////////////////////////////////
func (c *Cpu) SubImmediate (src1 uint8, imm uint8, dest uint8) {
	v := c.readReg(src1) - imm 
	c.writeReg(dest, v)
}

func (c *Cpu) SubRegister (src1 uint8, src2 uint8, dest uint8) {
	v := c.readReg(src1) - c.readReg(src2)
	c.writeReg(dest, v)
}

func (c *Cpu) Cycle() {
	if (c.PCBranchMod) { // go directly to the location of the branch
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
	c.Instruction = c.mem.GetRomWordAt(addr)
}


func (c *Cpu) setBranchedPC () {
	c.PCBranchMod = true
}

func (c *Cpu) unsetBranchedPC () {
	c.PCBranchMod = false
}

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

func (c *Cpu) IncPc() {
	c.PC++
	if (c.PC >= c.mem.RomTop) {
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
	twsCmplmnt := ( ( ^v ) + 1 ) 

	if ( (twsCmplmnt & 0x80) > 0 ) { // negative
		twsCmplmnt := twsCmplmnt & 0x7f
		c.PC -= twsCmplmnt	
	} else { // positive
		c.PC += v 
	}
}

func (c *Cpu) AdjustPC(baseRegister uint8, offset uint8) {
	c.PC = c.readReg(baseRegister) + offset
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
	fmt.Printf("\033[0;0H")
	fmt.Printf("                                                                                                         \n")
	fmt.Printf("                                                                                                         \n")		
	fmt.Printf("CPU State:                                                    \n")
	fmt.Printf("-----------------------\n")	
	fmt.Printf("A:\t%02X B:\t%02X \nC:\t%02X D:\t%02X\n", c.A, c.B, c.C, c.D)
	fmt.Printf("PC:\t%02X ZF:%s OF:%s SF:%s\n", c.PC, s.BoolToIntStr(c.getZeroFlag()), 
													  s.BoolToIntStr(c.getOverflowFlag()), 
													  s.BoolToIntStr(c.getSignFlag()) )
	fmt.Printf("Ins:[%b]\n", c.CurrentInstruction())
	fmt.Printf("-----------------------\n")
	c.mem.PrintRam()
}