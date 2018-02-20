package asm 

import (
	"bufio"
	s "dlw/shared"	
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func isValidRegister(register string) bool {

	var regs = regexp.MustCompile(`(?i)^A|B|C|D`)

	if regs.MatchString(register) {
		return true
	} else {
		return false
	}

}

func isValidAddress(address string) bool {
	return false
}

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

func getBaseAndOffset(a string, t string) (uint8,uint8) {
	baseAndOffset := strings.Split(a, "+")
	baseRegister := whichReg(baseAndOffset[0])
	offset, err := strconv.Atoi(baseAndOffset[1])
	if err != nil {
		// handle error
	}
	return baseRegister, uint8(offset)
}

func regArg(a *Argument, argStr string) {
	r := whichReg(argStr)
	a.MakeRegister(r)
}

func derefRegOrAddrArg(a *Argument, argStr string) {
	fmt.Println("derefRegOrAddrArg")
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

   	s := removeDerefChars(argStr)
	if plus.MatchString(s) {
		// this is a positive offset
		baseRegister, offset := getBaseAndOffset(s, "+")
		a.MakeDereference(baseRegister, offset)
	}
	if minus.MatchString(s) {
		// this is a negative ofset
		baseRegister, offset := getBaseAndOffset(s, "-")
		a.MakeDereference(baseRegister, offset)
	}
}

func getArgument(a string) *Argument {

   	arg := new(Argument)
   	arg.Init()

   	var deref = regexp.MustCompile(`^#`)
   	var parenthesis = regexp.MustCompile(`\(`)

   	fmt.Println(a)

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
	if (len(a) == 1 && isValidRegister(a)) {
		
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

	var bytecode uint16 = 0

    fmt.Println(asm)

	// collect the arguments from the line
	arguments := parseAsmLine(asm)

	// process the arguments
	src1 := getArgument(arguments[0])
	src2 := getArgument(arguments[1])
	dest := getArgument(arguments[2])

	// add does not handle deref arguments 
	if ( (src1.IsDereference || src1.IsLabel) || 
		 (src2.IsDereference || src2.IsLabel) ||
		 (dest.IsDereference || dest.IsLabel) ) {
      return 0, &s.SyntaxError{"Arithmetic instructions can operate on register or immediate", asm, lineNumber}  	
    }

    // generate the bytecode based on the types of arguments
    if (src1.IsRegister && src2.IsRegister && dest.IsRegister) {
		registerArithmeticByteCode(src1, src2, dest, &bytecode)    	
    } else {
    	immediateArithmeticByteCode(src1, src2, dest, &bytecode)
    }

	// set the instruction bits last
	setInstructionOpcodeBits(instructionType, &bytecode)

	// print for debugging
	fmt.Println(strconv.FormatInt(int64(bytecode), 2))

    // set the instruction code 
	return bytecode, nil
}

// load addr, dest_reg
// store src_reg, addr deref reg or offset
func HandleMemoryOperation(instructionType uint64, lineNumber uint8, asm string) (uint16, error) {

	var bytecode uint16 = 0

    fmt.Println(asm)

	// collect the arguments from the line
	arguments := parseAsmLine(asm)

	// process the arguments
	arg1 := getArgument(arguments[0])
	arg2 := getArgument(arguments[1])

	if ( instructionType == s.LOAD ) {
		e := getLoadByteCode(arg1, arg2, &bytecode, asm, lineNumber)
		return bytecode, e 
	}

	if ( instructionType == s.STORE ) {
		e := getStoreByteCode(arg1, arg2, &bytecode, asm, lineNumber)
		return bytecode, e
	}
  
	return 0, nil
}

// jump/jumpz #reg, #(reg + offset), label
func HandleBranchOperation(branchType uint64, labelOffsets map[string]uint8, currentLineNumber uint8, asm string) (uint16, error) {

	var bytecode uint16 = 0

    fmt.Println(asm)

	// collect the arguments from the line
	arguments := parseAsmLine(asm)

	// process the argument, there is only one for jump*
	arg := getArgument(arguments[0])

	// if this is a label, then seek the offset and populate the argument with it
	if(arg.IsLabel) {
		if labelLineNumber, ok := labelOffsets[arg.Label]; ok {
    		arg.SetLabelRelativeOffset(labelLineNumber, currentLineNumber)
    		fmt.Println(arg.ToString())
		} else {
			eS := fmt.Sprintf("the label %s does not exist in assembly", arg.Label)
			return 0, &s.SyntaxError{eS, asm, currentLineNumber}
		}
	}

	if ( branchType == s.JUMP ) {
		e := getJumpByteCode(arg, &bytecode, asm, currentLineNumber)
		return bytecode, e 
	}
  
	return 0, nil
}

func ParseLines(filePath string, parse func(string) (string, bool)) ([]uint16, error) {

	inputFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer inputFile.Close()

	/*
	 * Regular expressions for initial parser.
	 */
	var add = regexp.MustCompile(`(?i)add`)
	var sub = regexp.MustCompile(`(?i)sub`)
	var load  = regexp.MustCompile(`(?i)load`)
	var store = regexp.MustCompile(`(?i)store`)
	var jump = regexp.MustCompile(`(?i)jump`)
	//var jmpz  =
	var comment = regexp.MustCompile(`;`)
	var label = regexp.MustCompile(`(?i)\w+:`)

	var results []uint16
	var labelOffsets map[string]uint8
	labelOffsets = make(map[string]uint8)
	var lineNumber uint8 = 0

	// first run a pass to detect tags, make a map of tags -> line number
	// this will become index into array of words
	firstPassScanner := bufio.NewScanner(inputFile)
	fmt.Println("First scan pass")
	for firstPassScanner.Scan() {
		if output, _e := parse(firstPassScanner.Text()); _e {
			if label.MatchString(output) {
				tag := strings.Split(output, ":")[0]
				labelOffsets[tag] = lineNumber
				fmt.Println(labelOffsets)
			}
			lineNumber++
		}
	}

	// now run a second pass to build the set of instructions
	inputFile.Seek(0, 0) // rewind the file pointer
	secondPassScanner := bufio.NewScanner(inputFile)
	fmt.Println("Second scan pass")
	lineNumber = 0 // reuse this counter
	for secondPassScanner.Scan() {

		/*
		 * for each line, detect what it is
		 * call the parser for that line
		 * raise syntax error, with descriptive message and abort
		 * or return a valid 2-byte opcode
		 * handle invalid lines, abort with desriptive message
		 */

		if output, _e := parse(secondPassScanner.Text()); _e {

			// first check if this is a label

			//results = append(results, output)
			switch {

			// is this an add instruction
			case add.MatchString(output):
				if opcode, e := HandleArithmetic(s.ADD, lineNumber, output); e != nil {
					fmt.Println("Handle add failed:", e)
				} else {
					fmt.Println("Handle add worked:", opcode)
					results = append(results, opcode)
				}

			case sub.MatchString(output):
				if opcode, e := HandleArithmetic(s.SUB, lineNumber, output); e != nil {
					fmt.Println("HandleSub failed:", e)
				} else {
					fmt.Println("Handle sub worked:", opcode)
					results = append(results, opcode)
				}	

			case load.MatchString(output):
				if opcode, e := HandleMemoryOperation(s.LOAD, lineNumber, output); e != nil {
					fmt.Println("Handle load failed:", e)
				} else {
					fmt.Println("Handle load worked:", opcode)
					results = append(results, opcode)
				}

			case store.MatchString(output):
				if opcode, e := HandleMemoryOperation(s.STORE, lineNumber, output); e != nil {
					fmt.Println("Handle store failed:", e)
				} else {
					fmt.Println("Handle store worked:", opcode)
					results = append(results, opcode)	
				}

			case jump.MatchString(output):
				if opcode, e := HandleBranchOperation(s.JUMP, labelOffsets, lineNumber, output); e != nil {
					fmt.Println("Handle jump failed:", e)
				} else {
					fmt.Println("Handle jump worked:", opcode)
					results = append(results, opcode)		
				}
			case comment.MatchString(output):
				// skip these lines
			default:
				fmt.Println(output, "is not recognized")
			}

			lineNumber++
		}
	}

	// XXX refactor this
	if err := secondPassScanner.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
