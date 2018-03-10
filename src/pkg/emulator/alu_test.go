package emu

import "testing"

var addtests = []struct {
	x  uint8
	y  uint8
	r  uint8
	zf bool
	of bool
	sf bool
}{
	{1, 2, 3, false, false, false},
	{2, 1, 3, false, false, false},
	{255, 255, 254, false, true, true},
	{255, 1, 0, true, true, false},
	{1, 255, 0, true, true, false},
	{254, 1, 255, false, false, true},
	{0, 0, 0, true, false, false},
}

var subtests = []struct {
	x  uint8
	y  uint8
	r  uint8
	zf bool
	of bool
	sf bool
}{
	{1, 2, 255, false, true, true},
	{2, 1, 1, false, false, false},
	{255, 255, 0, true, false, false},
	{255, 1, 254, false, false, true},
	{1, 255, 2, false, true, false},
	{0, 0, 0, true, false, false},
}

func TestAdd(t *testing.T) {
	var alu Alu
	for _, tsts := range addtests {
		r := alu.Add(tsts.x, tsts.y)
		if r != tsts.r {
			t.Errorf("alu.Add(%d, %d) => %d, want %d", tsts.x, tsts.y, r, tsts.r)
		}
		if alu.ZeroFlag() != tsts.zf {
			t.Errorf("alu.Add(%d, %d) => %d, ZF (%v) want (%v)", tsts.x, tsts.y, r, alu.ZeroFlag(), tsts.zf)
		}
		if alu.OverflowFlag() != tsts.of {
			t.Errorf("alu.Add(%d, %d) => %d, OF (%v) want (%v)", tsts.x, tsts.y, r, alu.OverflowFlag(), tsts.of)
		}
		if alu.SignFlag() != tsts.sf {
			t.Errorf("alu.Add(%d, %d) => %d, SF (%v) want (%v)", tsts.x, tsts.y, r, alu.SignFlag(), tsts.sf)
		}
	}
}

func TestSub(t *testing.T) {
	var alu Alu
	for _, tsts := range subtests {
		r := alu.Sub(tsts.x, tsts.y)
		if r != tsts.r {
			t.Errorf("alu.Sub(%d, %d) => %d, want %d", tsts.x, tsts.y, r, tsts.r)
		}
		if alu.ZeroFlag() != tsts.zf {
			t.Errorf("alu.Sub(%d, %d) => %d, ZF (%v) want (%v)", tsts.x, tsts.y, r, alu.ZeroFlag(), tsts.zf)
		}
		if alu.OverflowFlag() != tsts.of {
			t.Errorf("alu.Sub(%d, %d) => %d, OF (%v) want (%v)", tsts.x, tsts.y, r, alu.OverflowFlag(), tsts.of)
		}
		if alu.SignFlag() != tsts.sf {
			t.Errorf("alu.Sub(%d, %d) => %d, SF (%v) want (%v)", tsts.x, tsts.y, r, alu.SignFlag(), tsts.sf)
		}
	}
}
