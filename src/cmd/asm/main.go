package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"pkg/assembler"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: line_parser <path>")
		return
	}

	lines, err := asm.ParseLines(os.Args[1], func(s string) (string, bool) {
		return s, true
	})

	if err != nil {
		fmt.Println("Error while parsing file", err)
		return
	}

	romdata := make([]byte, len(lines)*2)
	i := 0
	for _, l := range lines {
		binary.LittleEndian.PutUint16(romdata[i:], l)
		i += 2
	}

	// XXX check errors
	binfile, _ := os.Create("rom.bin")

	// XXX check errors
	binary.Write(binfile, binary.LittleEndian, romdata)

	fmt.Println(romdata)
}
