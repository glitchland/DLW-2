package shared 

import(
	"fmt"
)

// instructions enum 
const (
	ADD = iota 
	SUB
	LOAD
	STORE
	JUMP
)

// registers enum
const (
	A = iota
	B
	C
	D
	X
)

func Ritos(register uint8) string {
	switch {
	case register == A:
		return "A"
	case register == B:
		return "B"
	case register == C:
		return "C"
	case register == D:
		return "D"						
	default:
		return fmt.Sprintf("unknown register %d", register)
	}
}