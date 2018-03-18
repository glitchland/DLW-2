package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"pkg/assembler"
)

func main() {

	var asmLines []string

	asmFileName := flag.String("asm", "", "the name of the file containing the assembly")
	outFileName := flag.String("out", "", "the name of the file to write the assembly too")
	flag.Parse()

	if *asmFileName == "" || *outFileName == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// open the file containing the assembly
	asmFilePtr, err := os.Open(*asmFileName)

	if err != nil {
		fmt.Printf("Could not find the '%s' file to assemble\n", *asmFileName)
		os.Exit(1)
	}
	defer asmFilePtr.Close()

	_, err = asmFilePtr.Stat()
	if err != nil {
		fmt.Printf("Problem stating assembly file.\n")
		os.Exit(1)
	}

	// open the file we will write the rom to
	outFilePtr, err := os.Create(*outFileName)
	if err != nil {
		fmt.Printf("Could not open the '%s' rom file to write\n", *outFileName)
		os.Exit(1)
	}
	defer outFilePtr.Close()

	// read the file into a slice of strings for parsing
	scanner := bufio.NewScanner(asmFilePtr)
	for scanner.Scan() {
		asmLines = append(asmLines, scanner.Text())
	}

	opcodeArray, err := asm.ParseLines(asmLines)
	if err != nil {
		fmt.Printf("Failed to get opcodes\n")
		os.Exit(1)
	}

	romdata := make([]byte, len(opcodeArray)*2)
	i := 0
	for _, l := range opcodeArray {
		binary.LittleEndian.PutUint16(romdata[i:], l)
		i += 2
	}

	err = binary.Write(outFilePtr, binary.LittleEndian, romdata)
	if err != nil {
		fmt.Printf("Failed to write the rom file '%s'\n", *outFileName)
		os.Exit(1)
	}
}
