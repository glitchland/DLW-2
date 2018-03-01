package shared 

const (
	MEMOP_SRC_BIT_IDX = 4 
	MEMOP_DST_BIT_IDX = 6
	ARITHOP_SRC1_BIT_IDX = 4
	ARITHOP_SRC2_BIT_IDX = 6
	ARITHOP_DST_BIT_IDX = 8
	JMPOP_DST_BIT_IDX = 8
)

// bit manipulation functions, using the leftmost bit as the zeroth bit
func FlipBit(input *uint16, offset uint8) {
	*input ^= (1 << (15 - offset))
}

func SetBit(input *uint16, offset uint8) {
	*input |= (1 << (15 - offset))
}

func UnsetBit(input *uint16, offset uint8) {
	*input &^= (1 << (15 - offset))
}

func IsBitSet(n uint16, pos uint8) bool {
    val := n & (1 << (15 - pos))
    return (val > 0)
}

func IsModeSet(opcode uint16) bool {
	if (IsBitSet(opcode, 0)) {
		return true
	}
	return false
}

func IsMemOpSrcZero(opcode uint16) bool {
	if (IsBitSet(opcode, MEMOP_SRC_BIT_IDX) || IsBitSet(opcode, MEMOP_SRC_BIT_IDX + 1)) {
		return false
	} else {
		return true
	}
}

func IsMemOpDstZero(opcode uint16) bool {
	if (IsBitSet(opcode, MEMOP_DST_BIT_IDX) || IsBitSet(opcode, MEMOP_DST_BIT_IDX + 1)) {
		return false
	} else {
		return true
	}
}

func whichReg(opcode uint16, bitIdxl uint8, bitIdxr uint8) uint8 {
	switch {
	case !IsBitSet(opcode, bitIdxl) && !IsBitSet(opcode, bitIdxr): //A 00
		return A
	case !IsBitSet(opcode, bitIdxl) && IsBitSet(opcode, bitIdxr):  //B 01 
		return B
 	case IsBitSet(opcode, bitIdxl) && !IsBitSet(opcode, bitIdxr):  //C 10
 		return C
 	case IsBitSet(opcode, bitIdxl) &&  IsBitSet(opcode, bitIdxr):  //D 11
 		return D
 	default:
 		return X
 	}
}

func WhichMemOpSrcReg(opcode uint16) uint8 {
	return whichReg(opcode, MEMOP_SRC_BIT_IDX, MEMOP_SRC_BIT_IDX + 1)
}

func WhichMemOpDstReg(opcode uint16) uint8 {
	return whichReg(opcode, MEMOP_DST_BIT_IDX, MEMOP_DST_BIT_IDX + 1)	
}

func WhichArithOpSrc1Reg(opcode uint16) uint8 {
	return whichReg(opcode, ARITHOP_SRC1_BIT_IDX, ARITHOP_SRC1_BIT_IDX + 1)
}

func WhichArithOpSrc2Reg(opcode uint16) uint8 {
	return whichReg(opcode, ARITHOP_SRC2_BIT_IDX, ARITHOP_SRC2_BIT_IDX + 1)
}

func WhichArithOpDstReg(opcode uint16) uint8 {
	return whichReg(opcode, ARITHOP_DST_BIT_IDX, ARITHOP_DST_BIT_IDX + 1)	
}

func WhichJmpOpDstReg(opcode uint16) uint8 {
	return whichReg(opcode, JMPOP_DST_BIT_IDX, JMPOP_DST_BIT_IDX + 1)	
}

func GetImmediate(opcode uint16) uint8 {
	return uint8(opcode & 0xFF) 
}
