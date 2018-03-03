
``` 
$ export GOPATH=$(pwd)
$ go run src/cmd/asm/main.go src/examples/fill_ram_demo.asm 
$ go run src/cmd/emu/main.go 
```

```
8-bit addressing, 16-bit instructions, 16-bit memory and register width

- 1 bit mode            (2^1 = 2)
- 3 bits opcode         (2^3 = 8)
  LOAD 
  STORE
  ADD
  SUB
   - MUL
   - DIV
  JMP
  JMPZ
- 2 bits register type  (2^2 = 4)
 A
 B
 C
 D
- 8 bit address         (2^8 = 256)
```

# TODO   

- Check if this LOAD and STORE works properly in the parser and emulator
- Revise the rules around the JUMP instruction and the top bit of the source
  register (and in general clean this up)
- Implement the status register
- Write a JUMPZ instruction
- Handle code comments start and end of line, and relative jumps when there are commented lines

# Processor Status Word Register (PSW)

- sign flag (S) (bit 0)
 Indicates whether the result of the last mathematical operation resulted 
 in a value in which the most significant bit was set.
 (ref)[https://en.wikipedia.org/wiki/Negative_flag]

- overflow flag (V) (bit 1)
 Used to indicate when an arithmetic overflow has occurred in an operation, 
 indicating that the signed two's-complement result would not fit in the number 
 of bits used for the operation.
 (ref) [https://en.wikipedia.org/wiki/Overflow_flag]

- zero flag Z (bit 2)
 It is set if an arithmetic result is zero, and reset otherwise.
 (ref) [https://en.wikipedia.org/wiki/Zero_flag]

# Supported Instructions 

### LOAD 

The load instruction takes a memory location for the first argument, and a 
register for the second argument. 

The value in the memory location described by the first argument is written 
into the register specified as the second argument.

The memory location of the first argument can have 3 forms: 

- #REGISTER

If the mode flag is set to 0, then the value stored in the memory address 
that the value stored in the first argument points to is stored in the 
register that the second argument describes.

- #(REGISTER + IMMEDIATE_OFFSET)

If the mode flag is set to 1 and the _source_ register is not 00 then the 
value in the memory address that the source register + immediate value points 
to is written into the register that the second argument describes.

- #MEMORY_ADDRESS

If the mode flag is set to 1 and the _source_ register is 00, then the 
value in the in the memory address specified in the immediate value portion 
of the instruction is stored in the register that the second argument 
describes.

_Important_ : Due to this design, the A register is not legal as the _source_ 
register of this instruction.

```
   ------------------------------------------------------------------------
   LOAD (#REG || #(REG + OFFSET || #Memory), REG
   A 00
   B 01
   C 10
   D 11
```

### STORE

The store instruction takes a register as the first argument and a memory 
location for the second argument. 

The value in the first argument is written to the memory location described 
by the second argument.

The memory location of the second argument can have 3 forms: 

- #REGISTER

If the mode flag is set to 0, then the value in the first argument is stored
in the memory address that the value stored in the second argument points to.

- #(REGISTER + IMMEDIATE_OFFSET)

If the mode flag is set to 1 and the _destination_ register is not 00 then the 
value in the first argument (which is always a register) is stored in the
memory address that the destination register + immediate value points to.

- #MEMORY_ADDRESS

If the mode flag is set to 1 and the _destination_ register is 00, then the 
value in the first argument is stored in the memory address specified in the 
immediate value portion of the instruction.

_Important_ : Due to this design, the A register is not legal as the _destination_ register 
of this instruction.

```   
   ------------------------------------------------------------------------
   STORE REG, (#REG || #(REG + OFFSET || #Memory)
   A 00
   B 01
   C 10
   D 11
```

### ADD

The add instruction takes three arguments, it adds the first two, and writes the result
into the third.

It has two forms, the first form has a register for every argument. The second form has 
an immediate value for the second value.

If the mode flag is set to 1, then it is the immediate form. Take the immediate value
and add it to the contents of the first argument -- store the results in the third.

If the mode flag is set to 0, then it is the register form. Take the values from the 
registers specified in arg 1, and 2 and add them -- store the results in the register 
that is specified in arg 3.

This instruction should update the status register

```
can be all registers
 mode is 0, all values src1, src2, dest are registers

or can contain immediate, but only as src2 
A, 1, A

registerType
|0   |1|2|3 |4|5    |6|7    |8|9 |10|11|12|13|16|15|
|mode|opcode|source1|source2|dest|                 |

imediateType
|1   |1|2|3 |4|5   |6|7     |8|9|10|11|12|13|16|15|
|mode|opcode|source|dest    |   8-bit immediate   |
```

### SUB

The sub instruction takes three arguments, it subtracts the second argument from
the first, and writes the result into the third.

It has two forms, the first form has a register for every argument. The second form has 
an immediate argument for the second value.

If the mode flag is set to 1, then it is the immediate form. Take the immediate value
and subtract it from the contents of the first argument -- store the results in the 
third.

If the mode flag is set to 0, then it is the register form. Take the values from the 
registers specified in arg 1, and 2 and subtract them -- store the results in the 
register that is specified in arg 3.

This instruction should update the status register


```
can be all registers
 mode is 0, all values src1, src2, dest are registers

or can contain immediate, but only as src2 
A, 1, A

(mode will be set to 1, immediate value will be in bits |8|9|10|11|12|13|16|15|)

   ADD SRC1, SRC2, DESTINATION
   A 00
   B 01
   C 10
   D 11
```


### JMP
this is the unconditional form of jump, this instruction needs to be revised

- Mode should be reversed to be consistent.
- Rules for immediate / relative need to be refined.

- #REGISTER

Mode is 0. It advances the program counter by the value specified in the source
register portion of the instruction.

- #(REGISTER + IMMEDIATE_OFFSET)

Mode is 1. All bits of source register are set. It advances te program counter by
the value specified in the immediate portion of the instruction.

- LABEL

Mode is 1. This is a relative offset from the current position. It advances the program
counter by the value stored in the immediate portion of the instruction.


```
   ------------------------------------------------------------------------
   JUMP (#REG || #(REG + OFFSET) || LABEL)
   A 00
   B 01
   C 10
   D 11
   +----------------------------------------------------------------------+
   | 0    | 1,2,3  | 4,5      | 6,7      | 8,9,10,11,12,13,14,15          |
   |----------------------------------------------------------------------|
   | mode | opcode | source   | dest     | 8-bit dest address             |
   +----------------------------------------------------------------------+
```

-----------------------------------------------
-----------------------------------------------

# Parsing Load Instruction  

```    
   LOAD (#REG || #(REG + OFFSET) || #MEMORY), REG)
   A 00
   B 01
   C 10
   D 11
   +----------------------------------------------------------------------+
   |0     | 1,2,3  | 4,5      | 6,7      | 8,9,10,11,12,13,14,15          |
   |----------------------------------------------------------------------|
   | mode | opcode | source   | dest     | 8-bit source address           |
   +----------------------------------------------------------------------+
```

The source can be a '#' prefixed register (which is a dereference) or an immediate
type. 

The register type uses a mode flag value of 0. The contents of the register are 
used as an address.

The immediate type is a memory address literal -- in the form #ADDRESS, or in the 
form #(REGISTER + OFFSET) where the contents of REGISTER are treated as the base.

A mode flag value of 1 is used to identify this type. However, we also need to 
differentiate the #ADDRESS and #(REGISTER + OFFSET) sub-types. Due to the binary 
encoding format chosen in the 'Inside the Machine' book we need to zero the source 
register flags, and use the following parsing rules:

If the mode is immediate (1)
  If the source register is 00 then:
    It is the #ADDRESS form
  else:
    It is the #(REGISTER + OFFSET) form

This collides with the value for the A register, thus, the A register is invalid 
for all LOAD instructions.

# Parsing Jump Instructions     

Jump instructions can have the syntax variants:

```
jump #20
jump #A
jump #(C + 30)
jump LABEL
```    

There is always a single argument, and it is a dereferenced register, or register + offset,
or a label, or an immediate address.

The instruction will be in the form:    

```
   JUMP (#REG || #(REG + OFFSET) || LABEL)
   A 00
   B 01
   C 10
   D 11
   +----------------------------------------------------------------------+
   | 0    | 1,2,3  | 4,5      | 6,7      | 8,9,10,11,12,13,14,15          |
   |----------------------------------------------------------------------|
   | mode | opcode | source   | dest     | 8-bit dest offset              |
   +----------------------------------------------------------------------+
```

If the mode bit is unset, and bit 4 is unset, then a signed relative offset 
is present in bits 8-15. If bit 4 is set, then the 8-bit dest is interpreted 
as an immediate address.

If the mode bit is set, then bits 6,7 contain the base register encoding 
and bits 8-15 store an unsigned relative offset that is added to the base.

# References 

project org:

https://talks.golang.org/2014/organizeio.slide#9
