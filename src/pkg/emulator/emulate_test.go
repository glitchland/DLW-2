package emu

import (
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

func binStrToUint16(t *testing.T, bs string) uint16 {
	i, err := strconv.ParseInt(bs, 2, 17)
	if err != nil {
		t.Error(err)
	}
	return uint16(i)
}
