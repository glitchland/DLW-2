package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"pkg/emulator"
)

func main() {
	var opcodes []uint16

	romFile := flag.String("rom", "", "the name of the rom to load and emulate")
	flag.Parse()

	fmt.Println(romFile)
	if *romFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	fp, err := os.Open(*romFile)

	if err != nil {
		fmt.Printf("Could not find the '%s' rom file\n", *romFile)
		os.Exit(1)
	}

	defer fp.Close()

	fi, err := fp.Stat()
	if err != nil {
		fmt.Printf("Problem stating rom file.\n")
		os.Exit(1)
	}

	romLength := fi.Size()
	romBuffer := make([]byte, romLength)

	// read the ROM contents
	_, err = fp.Read(romBuffer)

	if err == io.EOF {
		fmt.Printf("Problem readiing rom file.\n")
		os.Exit(1)
	}

	i := 0
	for _, _ = range romBuffer {
		// 0, 2, 4, 8
		if i%2 == 0 {
			wVal := uint16(binary.LittleEndian.Uint16(romBuffer[i:])) //same as: int16(uint32(b[0]) | uint32(b[1])<<8)
			opcodes = append(opcodes, wVal)
		}
		i += 1
	}

	emu.Emulate(opcodes)
}
