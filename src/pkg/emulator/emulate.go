package emu

import (
	"fmt"
	s "pkg/shared"
	"time"
)

// check the opcode bits and pass to the correct handler
// arithmetic, memory,

//000
func isAdd(opcode uint16) bool {
	if !s.IsBitSet(opcode, 1) && !s.IsBitSet(opcode, 2) && !s.IsBitSet(opcode, 3) {
		return true
	}
	return false
}

//001
func isSub(opcode uint16) bool {
	if !s.IsBitSet(opcode, 1) && !s.IsBitSet(opcode, 2) && s.IsBitSet(opcode, 3) {
		return true
	}
	return false
}

//010
func isLoad(opcode uint16) bool {
	if !s.IsBitSet(opcode, 1) && s.IsBitSet(opcode, 2) && !s.IsBitSet(opcode, 3) {
		return true
	}
	return false
}

//011
func isStore(opcode uint16) bool {
	if !s.IsBitSet(opcode, 1) && s.IsBitSet(opcode, 2) && s.IsBitSet(opcode, 3) {
		return true
	}
	return false
}

//100
func isJump(opcode uint16) bool {
	if s.IsBitSet(opcode, 1) && !s.IsBitSet(opcode, 2) && !s.IsBitSet(opcode, 3) {
		return true
	}
	return false
}

//101
func isJumpz(opcode uint16) bool {
	if s.IsBitSet(opcode, 1) && !s.IsBitSet(opcode, 2) && s.IsBitSet(opcode, 3) {
		return true
	}
	return false
}

/*
   A 00
   B 01
   C 10
   D 11
   +----------------------------------------------------------------------+
   | 0    | 1,2,3  | 4,5      | 6,7      | 8,9,10,11,12,13,14,15          |
   |----------------------------------------------------------------------|
   | mode | opcode | source   | dest     | 8-bit dest address             |
   +----------------------------------------------------------------------+

   STORE REG, (#REG || #(REG + OFFSET || #Memory)
*/
func handleStore(c *Cpu) {
	// is it immediate?
	//src = constant register -- check bits 4, 5
	opcode := c.CurrentInstruction()
	src := s.WhichMemOpSrcReg(opcode)

	// is this an immediate type store?
	if s.IsModeSet(opcode) {

		// mode = 1 AND dest REGISTER == 00 (#ADDRESS form)
		if s.IsMemOpDstZero(opcode) {
			addr := s.GetImmediate(opcode)
			c.StoreAtAddr(addr, src)
		} else {
			// mode = 1 AND dest REGISTER != 00 #(REGISTER + OFFSET) form
			offset := s.GetImmediate(opcode)
			base := s.WhichMemOpDstReg(opcode)
			c.StoreAtRelative(src, base, offset)
		}

	} else {
		dst := s.WhichMemOpDstReg(opcode)
		c.StoreAtRegReference(src, dst)
	}
}

/*
   LOAD (#REG || #(REG + OFFSET) || #MEMORY), REG)
   A 00
   B 01
   C 10
   D 11
   +----------------------------------------------------------------------+
   |0     | 1,2,3  | 4,5      | 6,7      | 8,9,10,11,12,13,14,15          |
   |----------------------------------------------------------------------|
   | mode | opcode | source   | dest     | 8-bit source address           |
   +----------------------------------------------------------------------+
*/
func handleLoad(c *Cpu) {
	//is it immediate?
	//src = constant register -- check bits 4, 5
	opcode := c.CurrentInstruction()
	src := s.WhichMemOpSrcReg(opcode)
	dst := s.WhichMemOpDstReg(opcode)

	// is this an immediate type store?
	if s.IsModeSet(opcode) {

		// mode = 1 AND src REGISTER == 00 (#ADDRESS form)
		if s.IsMemOpSrcZero(opcode) {
			addr := s.GetImmediate(opcode)
			c.LoadFromAddr(addr, dst) // load from addr into dest
		} else {
			// mode = 1 AND src REGISTER != 00 #(REGISTER + OFFSET) form
			offset := s.GetImmediate(opcode)
			base := s.WhichMemOpDstReg(opcode)
			c.LoadFromRelative(base, offset, dst)
		}

	} else {
		c.LoadFromRegReference(src, dst) // load from src reg into dest
	}
}

/*
registerType
|0   |1|2|3 |4|5    |6|7    |8|9 |10|11|12|13|16|15|
|mode|opcode|source1|source2|dest|                 |

imediateType
|1   |1|2|3 |4|5   |6|7     |8|9|10|11|12|13|16|15|
|mode|opcode|source|dest    |   8-bit immediate   |
*/
func handleAdd(c *Cpu) {
	err := handleArithmetic(c, s.ADD)
	s.ChkFatalError(err)
}

func handleSub(c *Cpu) {
	err := handleArithmetic(c, s.SUB)
	s.ChkFatalError(err)
}

func handleArithmetic(c *Cpu, opType uint64) error {
	opcode := c.CurrentInstruction()

	if s.IsModeSet(opcode) {

		// immediate type add
		src1 := s.WhichArithOpSrc1Reg(opcode)
		dest := s.WhichArithOpSrc2Reg(opcode) // src2 is dest in this form
		imm := s.GetImmediate(opcode)

		switch opType {
		case s.ADD:
			c.AddImmediate(src1, imm, dest)
		case s.SUB:
			c.SubImmediate(src1, imm, dest)
		default:
			return &s.EmulatorError{"Unknown arithmetic instruction type"}
		}

	} else {

		// register type add
		src1 := s.WhichArithOpSrc1Reg(opcode)
		src2 := s.WhichArithOpSrc2Reg(opcode)
		dest := s.WhichArithOpDstReg(opcode)

		switch opType {
		case s.ADD:
			c.AddRegister(src1, src2, dest)
		case s.SUB:
			c.SubRegister(src1, src2, dest)
		default:
			return &s.EmulatorError{"Unknown arithmetic instruction type"}
		}
	}

	return nil
}

func handleJump(c *Cpu) {
	opcode := c.CurrentInstruction()
	imm := s.GetImmediate(opcode)

	if s.IsModeSet(opcode) {
		// if the top bit on the source register is set
		// then this is an #ADDRESS type.
		if s.IsBitSet(opcode, 4) {
			c.BranchAddress(imm)
		} else {
			// otherwise its a relative jump
			c.BranchRelative(imm)
		}
	} else {
		// this is a #REGISTER + IMMEDIATE type
		// bits 6,7 is register portion
		baseReg := s.WhichJmpOpDstReg(opcode)
		c.BranchArithmetic(baseReg, imm)
	}
}

func handleJumpz(c *Cpu) {
	opcode := c.CurrentInstruction()
	imm := s.GetImmediate(opcode)

	if s.IsModeSet(opcode) {
		// if the top bit on the source register is set
		// then this is an #ADDRESS type.
		if s.IsBitSet(opcode, 4) {
			c.BranchAddressCondZ(imm)
		} else {
			// otherwise its a relative jump
			c.BranchRelativeCondZ(imm)
		}
	} else {
		// this is a #REGISTER + IMMEDIATE type
		// bits 6,7 is register portion
		baseReg := s.WhichJmpOpDstReg(opcode)
		c.BranchArithmeticCondZ(baseReg, imm)
	}
}

func Emulate(romData []uint16) {

	cpu := new(Cpu)
	cpu.Init(romData)
	cpu.ClearScreen()
	for {
		// CPU tick time
		timer := time.Tick(time.Second / 50)
		<-timer
		switch {
		case isAdd(cpu.Instruction):
			handleAdd(cpu)
		case isSub(cpu.Instruction):
			handleSub(cpu)
		case isLoad(cpu.Instruction):
			handleLoad(cpu)
		case isStore(cpu.Instruction):
			handleStore(cpu)
		case isJump(cpu.Instruction):
			handleJump(cpu)
		default:
			fmt.Printf("instruction is unknown %04x.", cpu.Instruction) // raise error
		}

		cpu.Cycle()
		cpu.PrintState()

		if cpu.IsHalted() {
			fmt.Println("CPU Halted")
			break
		}
	}

}
