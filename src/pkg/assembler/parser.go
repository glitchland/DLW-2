package asm

import (
	"fmt"
	s "pkg/shared"
	"regexp"
	"strconv"
	"strings"
)

func isValidRegister(register string) bool {

	var regs = regexp.MustCompile(`(?i)\b[ABCD]\b`)

	if regs.MatchString(register) {
		return true
	} else {
		return false
	}

}

func isValidAddress(address string) bool {
	return false
}

// remove spaces and tabs
func normalizeString(input string) string {
	replacer := strings.NewReplacer(" ", "", "\t", "")
	return replacer.Replace(input)
}

func standardizeSpaces(input string) string {
	return strings.Join(strings.Fields(input), " ")
}

func isReg(r string, t string) bool {
	if r == strings.ToLower(t) || r == strings.ToUpper(t) {
		return true
	} else {
		return false
	}
}

func isRegA(register string) bool {
	return isReg(register, "A")
}

func isRegB(register string) bool {
	return isReg(register, "B")
}

func isRegC(register string) bool {
	return isReg(register, "C")
}

func isRegD(register string) bool {
	return isReg(register, "D")
}

func whichReg(register string) uint8 {
	switch {
	case isRegA(register):
		return s.A
	case isRegB(register):
		return s.B
	case isRegC(register):
		return s.C
	case isRegD(register):
		return s.D
	default:
		fmt.Printf("register is unknown %s.", register) // raise error
		return s.X
	}
}

////

func removeDerefChars(input string) string {
	replacer := strings.NewReplacer("#", "", "(", "", ")", "")
	return replacer.Replace(input)
}

func getBaseAndOffset(a string, t string) (uint8, uint8, error) {

	if t != "-" && t != "+" {
		return 0, 0, &s.ParserError{"Only '-' and '+' are supported."}
	}

	wsRemoved := normalizeString(a)
	baseAndOffset := strings.Split(wsRemoved, t)
	baseRegister := whichReg(baseAndOffset[0])
	offset, err := strconv.Atoi(baseAndOffset[1])
	if err != nil {
		return 0, 0, &s.ParserError{"Failed to convert base intege."}
	}
	return baseRegister, uint8(offset), nil
}

func regArg(a *Argument, argStr string) {
	r := whichReg(argStr)
	a.MakeRegister(r)
}

func derefRegOrAddrArg(a *Argument, argStr string) {
	r := removeDerefChars(argStr)

	// check if this is an integer memory address
	if addr, err := strconv.Atoi(r); err == nil {
		a.MakeAddress(uint8(addr))
	} else {
		// this is a register
		baseRegister := whichReg(r)
		a.MakeDereference(baseRegister, 0)
	}
}

func derefRegAndOffsetArg(a *Argument, argStr string) {
	var plus = regexp.MustCompile(`\+`)
	var minus = regexp.MustCompile(`\-`)

	baseAndOffsetStr := removeDerefChars(argStr)
	if plus.MatchString(baseAndOffsetStr) {
		// this is a positive offset
		baseRegister, offset, err := getBaseAndOffset(baseAndOffsetStr, "+")
		s.ChkFatalError(err)
		a.MakeDereference(baseRegister, offset)
	}
	if minus.MatchString(baseAndOffsetStr) {
		// this is a negative ofset
		baseRegister, offset, err := getBaseAndOffset(baseAndOffsetStr, "-")
		s.ChkFatalError(err)
		a.MakeDereference(baseRegister, offset)
	}
}

func getArgument(a string) *Argument {

	arg := new(Argument)
	arg.Init()

	var deref = regexp.MustCompile(`^#`)
	var parenthesis = regexp.MustCompile(`\(`)

	// this is a dereference style argument #C / #(D+16)
	if deref.MatchString(a) {

		if !parenthesis.MatchString(a) {
			// this is of the form #C or #123
			derefRegOrAddrArg(arg, a)
			return arg
		} else {
			// this is of the form #(D+16) extract the base + offset
			derefRegAndOffsetArg(arg, a)
			return arg
		}

	}

	// this is a register
	if len(a) == 1 && isValidRegister(a) {

		regArg(arg, a)
		return arg

	} else {

		// this is a label or an immediate integer
		// check if this is an integer
		if ai, err := strconv.Atoi(a); err == nil {
			arg.MakeImmediate(uint8(ai))
			return arg
		} else {
			// this is a label
			arg.MakeLabel(a)
			return arg
		}

	}

	fmt.Printf("argument type is not known %s.", a) // raise error
	return arg
}

func parseAsmLine(input string) []string {
	// XXX error check
	var labelRemoved string = ""

	labelString := strings.Split(input, ":") // remove label
	if len(labelString) == 2 {
		labelRemoved = labelString[1]
	} else {
		labelRemoved = labelString[0]
	}

	// make multiple spaces single spaces
	standardizedLine := standardizeSpaces(labelRemoved)

	// split the line, and then remove the instruction part
	splitLine := strings.Split(standardizedLine, " ")
	argsString := strings.Join(append(splitLine[:0], splitLine[1:]...), "")

	// remove all of the whitespace from the arguments
	// and then split them into parts
	normalizedArgs := normalizeString(argsString)
	parsedLine := strings.Split(normalizedArgs, ",")

	return parsedLine
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
func HandleArithmetic(instructionType uint64, lineNumber uint8, asm string) (uint16, error) {

	var opcode uint16
	var err error

	// collect the arguments from the line
	arguments := parseAsmLine(asm)

	// process the arguments
	src1 := getArgument(arguments[0])
	src2 := getArgument(arguments[1])
	dest := getArgument(arguments[2])

	// add does not handle deref arguments
	if (src1.IsDereference || src1.IsLabel) ||
		(src2.IsDereference || src2.IsLabel) ||
		(dest.IsDereference || dest.IsLabel) {
		return opcode, &s.SyntaxError{"Arithmetic instructions can operate on register or immediate", asm, lineNumber}
	}

	// generate the opcode based on the types of arguments
	if src1.IsRegister && src2.IsRegister && dest.IsRegister {
		opcode = registerArithmeticopcode(src1, src2, dest, opcode)
	} else {
		opcode, err = immediateArithmeticopcode(src1, src2, dest, opcode)
		if err != nil {
			// return the error up the stack
			return opcode, err
		}
	}

	// set the instruction bits last
	opcode = setInstructionOpcodeBits(instructionType, opcode)

	// set the instruction code
	return opcode, nil
}

// load addr, dest_reg
// store src_reg, addr deref reg or offset
func HandleMemoryOperation(instructionType uint64, lineNumber uint8, asm string) (uint16, error) {

	var opcode uint16
	var err error

	// collect the arguments from the line
	arguments := parseAsmLine(asm)

	// process the arguments
	arg1 := getArgument(arguments[0])
	arg2 := getArgument(arguments[1])

	if instructionType == s.LOAD {
		opcode, err = getLoadopcode(arg1, arg2, asm, lineNumber)
		return opcode, err
	}

	if instructionType == s.STORE {
		opcode, err = getStoreopcode(arg1, arg2, asm, lineNumber)
		return opcode, err
	}

	return opcode, nil
}

// jump/jumpz #reg, #(reg + offset), label
func HandleBranchOperation(branchType uint64, labelOffsets map[string]uint8, currentLineNumber uint8, asm string) (uint16, error) {

	var opcode uint16
	var err error

	// collect the arguments from the line
	arguments := parseAsmLine(asm)

	// process the argument, there is only one for jump*
	arg := getArgument(arguments[0])

	// if this is a label, then seek the offset and populate the argument with it
	if arg.IsLabel {
		if labelLineNumber, ok := labelOffsets[arg.Label]; ok {
			arg.SetLabelRelativeOffset(labelLineNumber, currentLineNumber)
		} else {
			eS := fmt.Sprintf("the label %s does not exist in assembly", arg.Label)
			return opcode, &s.SyntaxError{eS, asm, currentLineNumber}
		}
	}

	opcode, err = getJumpopcode(branchType, arg, asm, currentLineNumber)

	return opcode, err
}

func ParseLines(asmLines []string) ([]uint16, error) {

	var add = regexp.MustCompile(`(?i)add`)
	var sub = regexp.MustCompile(`(?i)sub`)
	var load = regexp.MustCompile(`(?i)load`)
	var store = regexp.MustCompile(`(?i)store`)
	var jump = regexp.MustCompile(`(?i)jump`)
	var jumpz = regexp.MustCompile(`(?i)jumpz`)
	var comment = regexp.MustCompile(`;`)
	var label = regexp.MustCompile(`(?i)\w+:`)

	var results []uint16
	var labelOffsets map[string]uint8
	labelOffsets = make(map[string]uint8)
	var lineNumber uint8 = 0

	// first run a pass to detect tags, make a map of tags -> line number
	// this will become index into array of words
	for _, l := range asmLines {
		if label.MatchString(l) {
			tag := strings.Split(l, ":")[0]
			labelOffsets[tag] = lineNumber
		}
		lineNumber++
	}

	// now run a second pass to build the set of instructions
	lineNumber = 0 // reuse this counter
	for _, line := range asmLines {

		/*
		 * for each line, detect what it is
		 * call the parser for that line
		 * raise syntax error, with descriptive message and abort
		 * or return a valid 2-byte opcode
		 * handle invalid lines, abort with desriptive message
		 */

		// first check if this is a label
		switch {

		// is this an add instruction
		case add.MatchString(line):
			if opcode, e := HandleArithmetic(s.ADD, lineNumber, line); e != nil {
				s.ChkFatalError(e)
			} else {
				fmt.Printf("Add [%s] -> |%016b|\n", line, opcode)
				results = append(results, opcode)
			}

		case sub.MatchString(line):
			if opcode, e := HandleArithmetic(s.SUB, lineNumber, line); e != nil {
				s.ChkFatalError(e)
			} else {
				fmt.Printf("Sub [%s] -> |%016b|\n", line, opcode)
				results = append(results, opcode)
			}

		case load.MatchString(line):
			if opcode, e := HandleMemoryOperation(s.LOAD, lineNumber, line); e != nil {
				s.ChkFatalError(e)
			} else {
				fmt.Printf("Load [%s] -> |%016b|\n", line, opcode)
				results = append(results, opcode)
			}

		case store.MatchString(line):
			if opcode, e := HandleMemoryOperation(s.STORE, lineNumber, line); e != nil {
				s.ChkFatalError(e)
			} else {
				fmt.Printf("Store [%s] -> |%016b|\n", line, opcode)
				results = append(results, opcode)
			}

		case jump.MatchString(line):
			if opcode, e := HandleBranchOperation(s.JUMP, labelOffsets, lineNumber, line); e != nil {
				s.ChkFatalError(e)
			} else {
				fmt.Printf("Jump [%s] -> |%016b|\n", line, opcode)
				results = append(results, opcode)
			}
		case jumpz.MatchString(line):
			if opcode, e := HandleBranchOperation(s.JUMPZ, labelOffsets, lineNumber, line); e != nil {
				s.ChkFatalError(e)
			} else {
				fmt.Printf("Jumpz [%s] -> |%016b|\n", line, opcode)
				results = append(results, opcode)
			}
		case comment.MatchString(line):
			// skip these lines
		default:
			fmt.Println(line, "is not recognized")
		}

		lineNumber++
	}
	// XXX return errors
	return results, nil
}
