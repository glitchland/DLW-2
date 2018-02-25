package shared 

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