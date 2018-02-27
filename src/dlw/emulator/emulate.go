package emu

import (
        "fmt"
         s "dlw/shared" 
         "time"
        )


// check the opcode bits and pass to the correct handler
// arithmetic, memory, 

func isAdd(opcode uint16) bool {
	if (!s.IsBitSet(opcode, 1) && !s.IsBitSet(opcode, 2) && !s.IsBitSet(opcode, 3) ) {
		return true
	}
	return false
}

func isSub(opcode uint16) bool {
	if (!s.IsBitSet(opcode, 1) && !s.IsBitSet(opcode, 2) && s.IsBitSet(opcode, 3) ) {
		return true
	}
	return false
}

func isLoad(opcode uint16) bool {
	if (!s.IsBitSet(opcode, 1) && s.IsBitSet(opcode, 2) && !s.IsBitSet(opcode, 3) ) {
		return true
	}
	return false
}

func isStore(opcode uint16) bool {
	if (!s.IsBitSet(opcode, 1) && s.IsBitSet(opcode, 2) && s.IsBitSet(opcode, 3) ) {
		return true
	}
	return false
}

// 100
func isJump(opcode uint16) bool {
	if (s.IsBitSet(opcode, 1) && !s.IsBitSet(opcode, 2) && !s.IsBitSet(opcode, 3) ) {
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
	if (s.IsModeSet(opcode)) {

		// mode = 1 AND dest REGISTER == 00 (#ADDRESS form)
		if (s.IsMemOpSrcZero(opcode)) {
			addr  := s.GetImmediate(opcode)
			c.StoreAtAddr(addr, src)			
		} else {
			// mode = 1 AND dest REGISTER != 00 #(REGISTER + OFFSET) form
			offset := s.GetImmediate(opcode)
			base   := s.WhichMemOpDstReg(opcode)
			c.StoreAtRelative(base, offset, src)
		}

	} else {
		// register type 
		reg := s.WhichMemOpDstReg(opcode)
		c.StoreAtRegReference(reg, src)
	}

	// dest is #(REGISTER + OFFSET)
	// mode = 1 AND dest REGISTER != 00 indicates this form

	// dest is #ADDRESS
	// mode = 1 AND dest REGISTER == 00 indicates this form

	//fmt.Printf("[Source %s] ", s.Ritos(whichReg(opcode, 4, 5)))
	//fmt.Printf(" [Dest %s]\n", s.Ritos(whichReg(opcode, 6, 7)))
}

/*
   ADD SRC1, SRC2, DESTINATION
   ADD 000
   A 00
   B 01
   C 10
   D 11
   +----------------------------------------------------------------------+
   |0     |1,2,3   |4,5       |6,7       |8,9          | 10,11,12,13,14,15|
   |----------------------------------------------------------------------|
   | mode | opcode | source 1 | source 2 | destination | padding          |
   +----------------------------------------------------------------------+
*/
func handleAdd(c *Cpu) {
	// is it immediate?
	//src = constant register -- check bits 4, 5
	opcode := c.CurrentInstruction()
	
	if (s.IsModeSet(opcode)) {
		// immediate type add
		// immediate is not allowed for destination
		// only the second argument can be immediate
		src1 := s.WhichArithOpSrc1Reg(opcode)
		imm  := s.GetImmediate(opcode)
		dest := s.WhichArithOpDstReg(opcode)
		c.AddImmediate(src1, imm, dest)
	} else {
		// register type add
		src1 := s.WhichArithOpSrc1Reg(opcode)
		src2 := s.WhichArithOpSrc2Reg(opcode)
		dest := s.WhichArithOpDstReg(opcode)
		c.AddRegister(src1, src2, dest)
	}
}

func handleJump(c *Cpu) {
	opcode := c.CurrentInstruction()
	imm  := s.GetImmediate(opcode)

	if (s.IsModeSet(opcode)) {

		// if the top bit on the source register is set
		// then this is an #ADDRESS type.
		if (s.IsBitSet(opcode, 4)) {
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

func Emulate(romData []uint16) {

	cpu := new(Cpu)
   	cpu.Init(romData)

	for {
			// CPU tick time
			ticker := time.Tick(time.Second/2)
			<-ticker
			// bits 1, 2, 3
			// 000 add
			// 001 sub
			// 010 load
			// 011 store
			// 100 jump
			switch {
			case isAdd(cpu.Instruction):
				//fmt.Printf("Add instruction [%04x]\n", cpu.Instruction)
				handleAdd(cpu)
			case isSub(cpu.Instruction):
				//fmt.Printf("Sub instruction [%04x]\n", cpu.Instruction)
			case isLoad(cpu.Instruction):
				//fmt.Printf("Load instruction [%04x]\n", cpu.Instruction)
			case isStore(cpu.Instruction):
				//fmt.Printf("Store instruction [%04x]\n", cpu.Instruction)
				handleStore(cpu)
			case isJump(cpu.Instruction):
				//fmt.Printf("Jump instruction [%04x]\n", cpu.Instruction)
				handleJump(cpu)
			default:
				fmt.Printf("instruction is unknown %04x.", cpu.Instruction) // raise error
			}

			cpu.Cycle()
		    cpu.PrintState()

		    if (cpu.IsHalted()) {
		    	fmt.Println("CPU Halted")
		    	break
		    }
	}

}