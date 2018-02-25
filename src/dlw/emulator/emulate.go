package emu

import (
        "fmt"
         s "dlw/shared" 
         "time"
        )


// check the opcode bits and pass to the correct handler
// arithmetic, memory, 

func isModeSet(opcode uint16) bool {
	if ( s.IsBitSet(opcode, 0) ) {
		return true
	}
	return false
}

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

// handlers 
func whichReg(opcode uint16, bitIdxl uint8, bitIdxr uint8) uint8 {
	switch {
	case !s.IsBitSet(opcode, bitIdxl) && !s.IsBitSet(opcode, bitIdxr): //A 00
		return s.A
	case !s.IsBitSet(opcode, bitIdxl) && s.IsBitSet(opcode, bitIdxr):  //B 01 
		return s.B
 	case s.IsBitSet(opcode, bitIdxl) && !s.IsBitSet(opcode, bitIdxr):  //C 10
 		return s.C
 	case s.IsBitSet(opcode, bitIdxl) &&  s.IsBitSet(opcode, bitIdxr):  //D 11
 		return s.D
 	default:
 		return s.X
 	}
}

func getImmediate(opcode uint16) uint8 {
	return uint8(opcode &  0x7f) 
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
	src := whichReg(opcode, 4, 5)

	if (isModeSet(opcode)) {
		// immediate type store
		addr  := getImmediate(opcode)
		c.StoreImmediate(addr, src)
	} else {
		// register type 

	}
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
	
	if (isModeSet(opcode)) {
		// immediate type add
		// immediate is not allowed for destination
		// only the second argument can be immediate
		src1 := whichReg(opcode, 4, 5)
		imm  := getImmediate(opcode)
		dest := whichReg(opcode, 8, 9)
		c.AddImmediate(src1, imm, dest)
	} else {
		// register type add
		src1 := whichReg(opcode, 4, 5)
		src2 := whichReg(opcode, 6, 7)
		dest := whichReg(opcode, 8, 9) 
		c.AddRegister(src1, src2, dest)
	}

	//fmt.Printf("[Source %s] ", s.Ritos(whichReg(opcode, 4, 5)))
	//fmt.Printf(" [Dest %s]\n", s.Ritos(whichReg(opcode, 6, 7)))
}

func handleJump(c *Cpu) {
	opcode := c.CurrentInstruction()
	imm  := getImmediate(opcode)

	if (isModeSet(opcode)) {

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
		baseReg := whichReg(opcode, 8, 9)
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

			//fmt.Printf("[OPCODE] %04x -> [PC] %02x\n", cpu.Instruction, cpu.PC)
		    if (cpu.IsHalted()) {
		    	fmt.Println("CPU Halted")
		    	break
		    }
		    cpu.PrintState()

	}

}