package asm

import "testing"

// Function is supposed to ensure a valid register name is passed
func TestIsValidRegister(t *testing.T) {

	// Test all valid registers
	registers := []string{"A","B","C","D","a","b","c","d"}

	for _, r := range registers {
	 	result := isValidRegister(r)
	    if result != true {
	       t.Errorf("Failed register %s, got: %v, want: %v.", r, result, true)
	    }
    }

	// Test a number of invalid registers
	registers = []string{"X","BAAA","%","///","ABCD","AA","aa","BB","bb","CC","bD"}

	for _, r := range registers {
	 	result := isValidRegister(r)
	    if result != false {
	       t.Errorf("Failed register %s, got: %v, want: %v.", r, result, false)
	    }
    }
}

// Function is supposed to strip whitespace from strings
func TestNormalizeString(t *testing.T) {
	test := "   Foo Bar			Tabular                 AA X                   "
	expected := "FooBarTabularAAX"
	result := normalizeString(test)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)		
	}
}

// Function is supposed to make all spaces a single space
func TestStandardizeSpaces(t *testing.T) {
	test := "    Foo Bar	Tabular                 AA X                   "
	expected := "Foo Bar Tabular AA X"
	result := standardizeSpaces(test)	
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)		
	}
}

// These functions test if the register string is a register
func TestIsRegA(t *testing.T) {

	// Expect to pass
	registers := []string{"A","a"}

	for _, r := range registers {
	 	result := isRegA(r)
	    if result != true {
	       t.Errorf("Failed %s, got: %v, want: %v.", r, result, true)
	    }
    }

	// Expect to fail
	registers = []string{"X","BAAA","%","///","ABCD","AA","aa","BB","bb","CC","bD"}

	for _, r := range registers {
	 	result := isRegA(r)
	    if result != false {
	       t.Errorf("Failed %s, got: %v, want: %v.", r, result, false)
	    }
    }

}

func TestIsRegB(t *testing.T) {

	// Expect to pass
	registers := []string{"B","b"}

	for _, r := range registers {
	 	result := isRegB(r)
	    if result != true {
	       t.Errorf("Failed %s, got: %v, want: %v.", r, result, true)
	    }
    }

	// Expect to fail
	registers = []string{"X","BAAA","%","///","ABCD","AA","aa","BB","bb","CC","bD"}

	for _, r := range registers {
	 	result := isRegB(r)
	    if result != false {
	       t.Errorf("Failed %s, got: %v, want: %v.", r, result, false)
	    }
    }
    
}

func TestIsRegC(t *testing.T) {

	// Expect to pass
	registers := []string{"C","c"}

	for _, r := range registers {
	 	result := isRegC(r)
	    if result != true {
	       t.Errorf("Failed %s, got: %v, want: %v.", r, result, true)
	    }
    }

	// Expect to fail
	registers = []string{"X","BAAA","%","///","ABCD","AA","aa","BB","bb","CC","bD"}

	for _, r := range registers {
	 	result := isRegC(r)
	    if result != false {
	       t.Errorf("Failed %s, got: %v, want: %v.", r, result, false)
	    }
    }
    
}

func TestIsRegD(t *testing.T) {

	// Expect to pass
	registers := []string{"D","d"}

	for _, r := range registers {
	 	result := isRegD(r)
	    if result != true {
	       t.Errorf("Failed %s, got: %v, want: %v.", r, result, true)
	    }
    }

	// Expect to fail
	registers = []string{"X","BAAA","%","///","ABCD","AA","aa","BB","bb","CC","bD"}

	for _, r := range registers {
	 	result := isRegD(r)
	    if result != false {
	       t.Errorf("Failed %s, got: %v, want: %v.", r, result, false)
	    }
    }
    
}

func TestWhichReg(t *testing.T) {
	
	registers := []string{"A","a","B","b","C","c","D","d"}
	values := []uint8{0,0,1,1,2,2,3,3}

	for i, r := range registers {
		expected := values[i]
		result := whichReg(r)
		if result != expected {
			t.Errorf("Failed %s, got: %v, want: %v.", r, result, expected)
		}
	}

}

func TestRemoveDerefChars(t *testing.T) {

	tests := []string { "#(REG + TEST)",
					    "#REG",
						"##(REG    + TEST   )",
					    "	#((((((((REG+TEST))))))",
	}
	expected := []string {  "REG + TEST",
							"REG",
							"REG    + TEST   ",
							"	REG+TEST",
	}

	for i, tst := range tests {
		expect := expected[i]
		result := removeDerefChars(tst)
		if result != expect {
			t.Errorf("Failed |%s|, got: %v, want: %v.", tst, result, expected)
		}
	}

}

// XXX : Expand this
func TestGetBaseAndOffset(t * testing.T) {

	// Expect these to all pass
	tests := [][]string{ {"A+2"  , "+"},
						 {"B+11" , "+"},
						 {"C+5"  , "+"},
						 {"D+90" , "+"},
						 {"A-12" , "-"},
						 {"B-13" , "-"},
						 {"C-14" , "-"},
						 {"D-15" , "-"},						 						 						 						 						 						 
	}
	expected := [][]uint8{ {0, 2 }, 
						   {1, 11},
						   {2, 5 },
						   {3, 90},
						   {0, 12},
						   {1, 13},
						   {2, 14},
						   {3, 15},
	}

	for i, tst := range tests {
		expect := expected[i]
		rbase, roffset := getBaseAndOffset(tst[0], tst[1])
		ebase := expect[0]
		eoffset := expect[1]
		if rbase != ebase || roffset != eoffset {
			t.Errorf("Failed |%s|, got: (%v,%v), want: (%v,%v).", tst, rbase, roffset, ebase, eoffset)
		}
	}

}

//regArg
//derefRegOrAddrArg
//derefRegAndOffsetArg
//getArgument
//parseAsmLine
//HandleArithmetic
//HandleBranchOperation
//ParseLines
