package shared 

import (
  "fmt"
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

/*
 * General emulation errors
 */
type EmulatorError struct {
	Msg    string // description of error
}

func (e *EmulatorError) Error() string {
	return fmt.Sprintf("Emulation error %s\n", e.Msg)
}

////
type MemoryError struct {
	Msg    string // description of error
	Addr   uint8
}

func (e *MemoryError) Error() string {
	return fmt.Sprintf("Memory error [%s] accessing address [%d].\n", e.Msg, e.Addr)
}