package main

import (
        "os"
        "encoding/binary"
		//"dlw/common"        
       	"dlw/emulator"
        "io"
        )

func main() {
	var opcodes []uint16
    fp, err := os.Open("rom.bin")

    if err != nil {
        panic(err)
    }

    defer fp.Close()

	fi, err := fp.Stat()
	if err != nil {
	  // Could not obtain stat, handle error
	}

	romLength := fi.Size()
    romBuffer := make([]byte, romLength) 

    // read the ROM contents
    _, err = fp.Read(romBuffer)

    if err == io.EOF{
        
    }

    i := 0
	for _, _ = range romBuffer {
        // 0, 2, 4, 8
		if (i % 2 == 0) {
			wVal := uint16(binary.LittleEndian.Uint16(romBuffer[i:])) //same as: int16(uint32(b[0]) | uint32(b[1])<<8)
			opcodes = append(opcodes, wVal)
		}
		i += 1 
	}

    emu.Emulate(opcodes)
}