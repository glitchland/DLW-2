package emu

import(
//	"dlw/common"
	"fmt"
)

const AddrLimits = ^uint8(0)
const MaxRom = 1024
type Memory struct {
	Rom	[]uint16
	RomWriteLocked bool
	Ram	[]uint8
	MemTop uint8
	RomTop uint8	
}

func (m *Memory) Init() {
	m.Rom = make([]uint16, MaxRom) 
	m.Ram = make([]uint8, AddrLimits) 
	m.RomWriteLocked = false
	m.MemTop = AddrLimits
}

// RAM access - RAM is a series of 8 bit bytes
func (m *Memory) GetRamByteAt(addr uint8) uint8 {
	return m.Ram[addr]
}

func (m *Memory) SetRamByteAt(addr uint8, v uint8) {
	m.Ram[addr] = v
}

// ROM access, ROM is a series of 16 bit word opcodes
func (m *Memory) LoadRom(opcodes []uint16) {
	for i, opcode := range opcodes {
		m.SetRomWordAt(uint8(i), opcode)
	}
	m.RomTop = uint8(len(opcodes))
	m.RomLock()
}

func (m *Memory) GetRomWordAt(addr uint8) uint16 {
	return m.Rom[addr]
}

func (m *Memory) SetRomWordAt(addr uint8, v uint16) {//error {
	if(!m.RomWriteLocked) {
		m.Rom[addr] = v
		//return nil
	} //else {
	//	return &memoryError{"Attempted to write to read-only memory", addr} 
	//}
}

func (m *Memory) RomLock() {
	m.RomWriteLocked = true
}

func (m *Memory) RomUnlock() {
	m.RomWriteLocked = false
}

func (m *Memory) IsRomLocked() bool {
	return m.RomWriteLocked
}

func (m *Memory) PrintRam() {
	i := 0
	b := uint8(0)
	s := "RAM:\n"
	for _, b = range m.Ram {
		i++	
        s += fmt.Sprintf("%02x ", b)
		if (i % 16 == 0) {
			s += fmt.Sprintf("\n")
		}
	}
	fmt.Println(s)
}

