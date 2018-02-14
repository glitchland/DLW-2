package asm 

import (
  "fmt"
)

/*
 * General syntax errors
 */
type syntaxError struct {
	msg    string // description of error
	code   string // the parsed code
	Offset uint8  // error occurred after reading Offset bytes
}

func (e *syntaxError) Error() string {
	return fmt.Sprintf("Syntax error on line %d - [%s] - [%s]", e.Offset, e.code, e.msg)
}