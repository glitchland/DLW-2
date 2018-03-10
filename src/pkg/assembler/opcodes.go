package asm

import (
	"fmt"
	s "pkg/shared"
)

// There can only be one immediate
// these instructions cannot have a dereference type, it has to be register
// or immediate

func setInstructionOpcodeBits(instruction uint64, opcode uint16) uint16 {

	switch instruction {
	case s.ADD:
		// ADD is 000, so nothing to do
		return opcode
	case s.SUB:
		// SUB is 001, bit 3 index from 0
		return s.SetBit(opcode, 3)
	case s.LOAD:
		// LOAD is 010, bits index from 0
		return s.SetBit(opcode, 2)
	case s.STORE:
		// STORE is 011
		opcode = s.SetBit(opcode, 2)
		return s.SetBit(opcode, 3)
	case s.JUMP:
		// JUMP is 100
		return s.SetBit(opcode, 1)
	case s.JUMPZ:
		// JUMPZ is 101
		opcode = s.SetBit(opcode, 1)
		return s.SetBit(opcode, 3)
	default:
		fmt.Printf("unknown instruction") // raise error
	}

	return opcode
}

// ADD, SUB, MUL, DIV
func setRegisterOpcodeBits(reg uint8, opcode uint16, offset uint8) uint16 {

	switch reg {
	case s.A:
		// A is 0,0 no action required
		return opcode
	case s.B:
		// B is 0,1
		return s.SetBit(opcode, offset+1)
	case s.C:
		// C is 1,0
		return s.SetBit(opcode, offset)
	case s.D:
		// D is 1,1
		opcode = s.SetBit(opcode, offset)
		return s.SetBit(opcode, offset+1)
	default:
		fmt.Printf("unknown register") // raise error
	}

	return opcode
}

// Bits |8|9|10|11|12|13|16|15| are the 8 bit immediate
func setImmediateOpcodeBits(immediateValue uint8, opcode uint16) uint16 {
	// extend the immediate value to 16 bits
	ei := uint16(immediateValue)

	// or in the low 8 bits, which will be safe because the upper bits are all 0
	opcode |= ei

	// set the top bit to indicate immediate
	return s.SetBit(opcode, 0)
}

// registerType
//|0   |1|2|3 |4|5   |6|7     |8|9 |10|11|12|13|16|15|
//|mode|opcode|source1|source2|dest|                 |
func registerArithmeticopcode(src1 *Argument, src2 *Argument, dest *Argument, opcode uint16) uint16 {

	// offset 4 :: bit 4, 5 (src1)
	opcode = setRegisterOpcodeBits(src1.Register, opcode, 4)
	// offset 6 :: bit 6,7 (src2)
	opcode = setRegisterOpcodeBits(src2.Register, opcode, 6)
	// offset 8 :: bit 8,9 (dest)
	opcode = setRegisterOpcodeBits(dest.Register, opcode, 8)

	return opcode
}

// imediateType
//|0   |1|2|3 |4|5   |6|7 |8|9|10|11|12|13|16|15|
//|mode|opcode|source|dest|8-bit immediate      |
func immediateArithmeticopcode(src1 *Argument, src2 *Argument, dest *Argument, opcode uint16) (uint16, error) {

	// Only the second argument can be immediate
	if src1.IsImmediate {
		return opcode, &s.OpcodeError{"Error, only the second argument can be immediate"}
	} else {
		// offset 4 :: bit 4, 5 (src1)
		opcode = setRegisterOpcodeBits(src1.Register, opcode, 4)
	}

	if src2.IsImmediate {
		opcode = setImmediateOpcodeBits(src2.ImmediateInt, opcode)
	} else {
		return opcode, &s.OpcodeError{"Immediate must be set as the second argument in this form"}
	}

	if dest.IsImmediate {
		return opcode, &s.OpcodeError{"Error, only the second argument can be immediate"}
	} else {
		// offset 6 :: bit 6,7 (dest)
		opcode = setRegisterOpcodeBits(dest.Register, opcode, 6)
	}

	return opcode, nil
}

/* load needs a memory/deref for first arg, reg for second
   ------------------------------------------------------------------------
   LOAD (#REG || #(REG + OFFSET || #Memory), REG
   A 00
   B 01
   C 10
   D 11
   +----------------------------------------------------------------------+
   | 0    | 1,2,3  | 4,5      | 6,7      | 8,9,10,11,12,13,14,15          |
   |----------------------------------------------------------------------|
   | mode | opcode | source   | dest     | 8-bit source address           |
   +----------------------------------------------------------------------+
*/
func getLoadopcode(src *Argument, dest *Argument, asm string, lineNumber uint8) (uint16, error) {

	var opcode uint16

	if !(src.IsDereference || src.IsAddress) && !dest.IsRegister {
		return opcode, &s.SyntaxError{"load must have the form: load #deref/mem, register", asm, lineNumber}
	}

	// set the instruction bits
	// bits 1, 2, 3
	opcode = setInstructionOpcodeBits(s.LOAD, opcode)

	// source is #ADDRESS
	// mode = 1 AND source REGISTER = 00 indicates this form
	if src.IsAddress {
		// just write the address into the last 8-bit
		// bits, leave the destination register 00
		opcode = setImmediateOpcodeBits(src.Address, opcode)
	}

	// source is #(REGISTER + OFFSET)
	// mode = 1 and source REGISTER != 00 indicates this form
	if src.IsDereference {
		if src.BaseRegister == s.A {
			return opcode, &s.SyntaxError{"A register is not legal for load #(REGISTER + OFFSET) form", asm, lineNumber}
		}
		opcode = setImmediateOpcodeBits(src.Offset, opcode)
		// set the register parts
		// bits 4,5
		opcode = setRegisterOpcodeBits(src.BaseRegister, opcode, 4)
	}

	// set the destination register
	// bits 6,7
	opcode = setRegisterOpcodeBits(dest.Register, opcode, 6)

	return opcode, nil
}

/* store needs a reg for first arg, deref for second

   ------------------------------------------------------------------------
   STORE REG, (#REG || #(REG + OFFSET || #Memory)
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
func getStoreopcode(src *Argument, dest *Argument, asm string, lineNumber uint8) (uint16, error) {

	var opcode uint16

	if !src.IsRegister && !(dest.IsDereference || dest.IsAddress) {
		return opcode, &s.SyntaxError{"store must have the form: store register, #deref/mem", asm, lineNumber}
	}

	// set the instruction bits
	// bits 1, 2, 3
	opcode = setInstructionOpcodeBits(s.STORE, opcode)

	// set the source register
	// bits 4, 5
	opcode = setRegisterOpcodeBits(src.Register, opcode, 4)

	// dest is #ADDRESS
	// mode = 1 AND dest REGISTER == 00 indicates this form
	if dest.IsAddress {
		// just write the address into the last 8-bit
		// bits, leave the destination register 00
		opcode = setImmediateOpcodeBits(dest.Address, opcode)
	}

	// dest is #(REGISTER + OFFSET)
	// mode = 1 AND dest REGISTER != 00 indicates this form
	if dest.IsDereference {
		if dest.BaseRegister == s.A {
			return opcode, &s.SyntaxError{"A register is not legal for store #(REGISTER + OFFSET) form", asm, lineNumber}
		}
		// set the address portion of the destination
		opcode = setImmediateOpcodeBits(dest.Offset, opcode)
		// set the register portion of the destination
		// bits 6,7
		opcode = setRegisterOpcodeBits(dest.BaseRegister, opcode, 6)
	}
	return opcode, nil
}

/* jump has a single argument
   ------------------------------------------------------------------------
   JUMP (#REG || #(REG + OFFSET) || LABEL)
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
func getJumpopcode(jmpType uint64, dest *Argument, asm string, lineNumber uint8) (uint16, error) {

	var opcode uint16

	if !(dest.IsDereference || dest.IsLabel || dest.IsImmediate) {
		eS := "jump must have the form: jump #reg, or jump #(reg + offset), or jump label, or jump #immediate"
		return opcode, &s.SyntaxError{eS, asm, lineNumber}
	}

	// set the instruction bits
	// bits 1, 2, 3
	switch jmpType {
	case s.JUMP:
		opcode = setInstructionOpcodeBits(s.JUMP, opcode)
	case s.JUMPZ:
		opcode = setInstructionOpcodeBits(s.JUMPZ, opcode)
	default:
		return opcode, &s.SyntaxError{"Unknown jump type passed", asm, lineNumber}
	}

	// mode 0
	// if this is a #reg, use dest field + 0 as 8-bit offset
	// if this is a #reg + offset, use dest field + 8-bit offset address
	if dest.IsDereference {
		// set the address portion of the destination
		opcode = setImmediateOpcodeBits(dest.Offset, opcode)
		// set the register portion of the destination
		// bits 6,7
		opcode = setRegisterOpcodeBits(dest.BaseRegister, opcode, 6)
	}

	// mode 1
	// if this is a label type store the offset in the 8-bit dest address
	if dest.IsLabel {
		opcode = setImmediateOpcodeBits(uint8(dest.LabelOffset), opcode)
	}

	// mode 1
	// if this is an immediate type store the address in the 8-bit dest
	if dest.IsImmediate {
		opcode = setImmediateOpcodeBits(dest.ImmediateInt, opcode)
		opcode = s.SetBit(opcode, 4) // set the top bit on the source register
		// to identify immediate
	}

	return opcode, nil
}
