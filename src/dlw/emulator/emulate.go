package emu

import (
        "fmt"
         s "dlw/shared" 
        )


// check the opcode bits and pass to the correct handler
// arithmetic, memory, 

func hasBit(n uint16, pos uint8) bool {
    val := n & (1 << (15 - pos))
    return (val > 0)
}

func isAdd(opcode uint16) bool {
	if (!hasBit(opcode, 1) && !hasBit(opcode, 2) && !hasBit(opcode, 3) ) {
		return true
	}
	return false
}

func isSub(opcode uint16) bool {
	if (!hasBit(opcode, 1) && !hasBit(opcode, 2) && hasBit(opcode, 3) ) {
		return true
	}
	return false
}

func isLoad(opcode uint16) bool {
	if (!hasBit(opcode, 1) && hasBit(opcode, 2) && !hasBit(opcode, 3) ) {
		return true
	}
	return false
}

func isStore(opcode uint16) bool {
	if (!hasBit(opcode, 1) && hasBit(opcode, 2) && hasBit(opcode, 3) ) {
		return true
	}
	return false
}

// 100
func isJump(opcode uint16) bool {
	if (hasBit(opcode, 1) && !hasBit(opcode, 2) && !hasBit(opcode, 3) ) {
		return true
	}
	return false
}

// handlers 
func whichReg(opcode uint16, bitIdxl uint8, bitIdxr uint8 ) uint8 {
	switch {
	case !hasBit(opcode, bitIdxl) && !hasBit(opcode, bitIdxr): //A 00
		return s.A
	case !hasBit(opcode, bitIdxl) && hasBit(opcode, bitIdxr):  //B 01 
		return s.B
 	case hasBit(opcode, bitIdxl) && !hasBit(opcode, bitIdxr):  //C 10
 		return s.C
 	case hasBit(opcode, bitIdxl) &&  hasBit(opcode, bitIdxr):  //D 11
 		return s.D
 	default:
 		return s.X
 	}
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
*/
func handleStore(c *Cpu) {
	// is it immediate?
	//src = constant register -- check bits 4, 5
	opcode := c.CurrentInstruction()
	fmt.Printf("[Source %s] ", s.Ritos(whichReg(opcode, 4, 5)))
	fmt.Printf(" [Dest %s]\n", s.Ritos(whichReg(opcode, 6, 7)))
}

func Emulate(romData []uint16) {

	cpu := new(Cpu)
   	cpu.Init(romData)

	for {
			// bits 1, 2, 3
			// 000 add
			// 001 sub
			// 010 load
			// 011 store
			// 100 jump

			switch {
			case isAdd(cpu.Instruction):
				fmt.Printf("Add instruction [%04x]\n", cpu.Instruction)
			case isSub(cpu.Instruction):
				fmt.Printf("Sub instruction [%04x]\n", cpu.Instruction)
			case isLoad(cpu.Instruction):
				fmt.Printf("Load instruction [%04x]\n", cpu.Instruction)
			case isStore(cpu.Instruction):
				fmt.Printf("Store instruction [%04x]\n", cpu.Instruction)
				handleStore(cpu)
			case isJump(cpu.Instruction):
				fmt.Printf("Jump instruction [%04x]\n", cpu.Instruction)
			default:
				fmt.Printf("instruction is unknown %04x.", cpu.Instruction) // raise error
			}

			fmt.Printf("[OPCODE] %04x -> [PC] %02x\n", cpu.Instruction, cpu.PC)
		    cpu.NextInstruction()
		    if (cpu.IsHalted()) {
		    	fmt.Println("CPU Halted")
		    	break
		    }
		    cpu.PrintState()

	}

}