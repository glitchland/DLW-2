package main

import (
	"pkg/assembler"
	"encoding/binary"	
	"fmt"
	"os"
)


func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: line_parser <path>")
		return
	}

	lines, err := asm.ParseLines(os.Args[1], func(s string) (string, bool) {
		fmt.Println(s)
		return s, true
	})

	if err != nil {
		fmt.Println("Error while parsing file", err)
		return
	}

	// pack the words into a byte string
	// https://golang.org/pkg/encoding/binary/#example_Read
	romdata := make([]byte, len(lines) * 2)
	i := 0
	for _, l := range lines {
		binary.LittleEndian.PutUint16(romdata[i:], l)
		i += 2
	}

	// check errors
	binfile, _ := os.Create("rom.bin")
	
	// check errors
	binary.Write(binfile, binary.LittleEndian, romdata)

	fmt.Println(romdata)
}
