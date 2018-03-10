package emu

import (
	s "pkg/shared"
	"strconv"
	"testing"
)

// Add
var isaddt = []struct {
	bs string
	ex bool
}{
	{"1000111111111111", true},
	{"1001111111111111", false},
	{"1011111111111111", false},
	{"1101111111111111", false},
	{"1010111111111111", false},
	{"1100111111111111", false},
	{"1111111111111111", false},
}

func TestIsAdd(t *testing.T) {
	for _, tst := range isaddt {
		o := binStrToUint16(t, tst.bs)
		r := isAdd(o)
		if r != tst.ex {
			t.Errorf("isAdd(%s) => %v, want %v", tst.bs, r, tst.ex)
		}
	}
}

// Sub
var issubt = []struct {
	bs string
	ex bool
}{
	{"1000111111111111", false},
	{"1001111111111111", true},
	{"1011111111111111", false},
	{"1101111111111111", false},
	{"1010111111111111", false},
	{"1100111111111111", false},
	{"1111111111111111", false},
}

func TestIsSub(t *testing.T) {
	for _, tst := range issubt {
		o := binStrToUint16(t, tst.bs)
		r := isSub(o)
		if r != tst.ex {
			t.Errorf("isSub(%s) => %v, want %v", tst.bs, r, tst.ex)
		}
	}
}

// Load
var isloadt = []struct {
	bs string
	ex bool
}{
	{"1000111111111111", false},
	{"1001111111111111", false},
	{"1011111111111111", false},
	{"1101111111111111", false},
	{"1010111111111111", true},
	{"1100111111111111", false},
	{"1111111111111111", false},
}

func TestIsLoad(t *testing.T) {
	for _, tst := range isloadt {
		o := binStrToUint16(t, tst.bs)
		r := isLoad(o)
		if r != tst.ex {
			t.Errorf("isLoad(%s) => %v, want %v", tst.bs, r, tst.ex)
		}
	}
}

// Store
var isstoret = []struct {
	bs string
	ex bool
}{
	{"1000111111111111", false},
	{"1001111111111111", false},
	{"1011111111111111", true},
	{"1101111111111111", false},
	{"1010111111111111", false},
	{"1100111111111111", false},
	{"1111111111111111", false},
}

func TestIsStore(t *testing.T) {
	for _, tst := range isstoret {
		o := binStrToUint16(t, tst.bs)
		r := isStore(o)
		if r != tst.ex {
			t.Errorf("isStore(%s) => %v, want %v", tst.bs, r, tst.ex)
		}
	}
}

// Jump
var isjumpt = []struct {
	bs string
	ex bool
}{
	{"1000111111111111", false},
	{"1001111111111111", false},
	{"1011111111111111", false},
	{"1101111111111111", false},
	{"1010111111111111", false},
	{"1100111111111111", true},
	{"1111111111111111", false},
}

func TestIsJump(t *testing.T) {
	for _, tst := range isjumpt {
		o := binStrToUint16(t, tst.bs)
		r := isJump(o)
		if r != tst.ex {
			t.Errorf("isJump(%s) => %v, want %v", tst.bs, r, tst.ex)
		}
	}
}

// Jumpz
var isjumpzt = []struct {
	bs string
	ex bool
}{
	{"1000111111111111", false},
	{"1001111111111111", false},
	{"1011111111111111", false},
	{"1101111111111111", true},
	{"1010111111111111", false},
	{"1100111111111111", false},
	{"1111111111111111", false},
}

func TestIsJumpz(t *testing.T) {
	for _, tst := range isjumpzt {
		o := binStrToUint16(t, tst.bs)
		r := isJumpz(o)
		if r != tst.ex {
			t.Errorf("isJumpz(%s) => %v, want %v", tst.bs, r, tst.ex)
		}
	}
}

//handleArithmetic
//handleAdd
//handleSub
var arithtsts = []struct {
	opType uint64
	desc   string
	bts    string
	ex     uint8
}{
	{s.SUB, "sub A, B, C", "0001000110000000", 18},
	{s.SUB, "sub A, 15, C", "1001001000001111", 5},
	{s.ADD, "add A, B, C", "0000000110000000", 22},
	{s.ADD, "add A, 15, C", "1000001000001111", 35},
}

func TestHandleArithmetic(t *testing.T) {
	cpu := new(Cpu)

	for _, tst := range arithtsts {

		cpu = new(Cpu) // initialize new cpu
		cpu.regs.Write(s.A, 20)
		cpu.regs.Write(s.B, 2)
		cpu.regs.Write(s.C, 0)
		cpu.regs.Write(s.D, 0) // set up the entry state

		o := binStrToUint16(t, tst.bts)
		cpu.Instruction = o

		switch tst.opType {
		case s.ADD:
			handleAdd(cpu)
		case s.SUB:
			handleSub(cpu)
		}

		r, _ := cpu.regs.Read(s.C)

		if r != tst.ex {
			t.Errorf("instruction (%s)[%s] => %v, want %v", tst.desc, tst.bts, r, tst.ex)
		}

	}
}

//handleStore
//handleLoad

//handleJump
//handleJumpz

/* test data
sub A, B, C --> 0001000110000000
sub C, B, A --> 0001100100000000
add A, B, C --> 0000000110000000
add C, B, A --> 0000100100000000
sub D, D, D --> 0001111111000000
add A, 15, C --> 1000001000001111
sub A, 15, D --> 1001001100001111
store C, #(D + 16) --> 1011101100010000
store C, #1 --> 1011100000000001
store B, #D --> 1011011100000000
load #13, B --> 1010000100001101
load #(C + 2), C --> 1010101000000010
load #D, D --> 1010111100000000
jumpz LBL1 --> 1100000000001010
jumpz #C --> 1100001000000000
jumpz #(D + 16) --> 1100001100010000
jumpz 1 --> 1100100000000001
jumpz -1       --> 1100100011111111
jump LBL1 --> 1100000000000101
jump #C --> 1100001000000000
jump #(D + 16) --> 1100001100010000
jump 1 --> 1100100000000001
jump -1 --> 1100100011111111
LBL1: add A, B, B --> 0000000101000000
*/
func binStrToUint16(t *testing.T, bs string) uint16 {
	i, err := strconv.ParseInt(bs, 2, 17)
	if err != nil {
		t.Error(err)
	}
	return uint16(i)
}
