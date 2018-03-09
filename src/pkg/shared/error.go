package shared

import (
	"fmt"
	"os"
)

/*
 * General syntax errors
 */
type SyntaxError struct {
	Msg    string // description of error
	Code   string // the parsed code
	Offset uint8  // error occurred after reading Offset bytes
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("Syntax error on line %d - [%s] - [%s]", e.Offset, e.Code, e.Msg)
}

type ParserError struct {
	Msg string // description of error
}

func (e *ParserError) Error() string {
	return fmt.Sprintf("Parser error %s\n", e.Msg)
}

type OpcodeError struct {
	Msg string // description of error
}

func (e *OpcodeError) Error() string {
	return fmt.Sprintf("Opcode generation error %s\n", e.Msg)
}

/*
 * General emulation errors
 */
type EmulatorError struct {
	Msg string // description of error
}

func (e *EmulatorError) Error() string {
	return fmt.Sprintf("Emulation error %s\n", e.Msg)
}

////
type MemoryError struct {
	Msg  string // description of error
	Addr uint8
}

func (e *MemoryError) Error() string {
	return fmt.Sprintf("Memory error [%s] accessing address [%d].\n", e.Msg, e.Addr)
}

func ChkFatalError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
