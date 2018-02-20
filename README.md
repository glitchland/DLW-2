TODO: 

- Change store, allow immediate for src otherwise we cannot set up state
- Change load, allow immediate for src otherwise we cannot set up state
- Emulate both of these first


``` 
$ export GOPATH=$(pwd)
$ go run src/main.go test.asm 
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
