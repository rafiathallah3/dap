package bytecode

import "dap/internal/common"

// Instruction represents a single compiled bytecode instruction.
type Instruction struct {
	Op   OpCode
	Arg  int
	Line int
}

// FunctionChunk represents a compiled user-defined function.
type FunctionChunk struct {
	Name     string
	ArgNames []string
	Chunk    *Chunk
}

// Chunk stores a sequence of bytecode instructions and metadata.
type Chunk struct {
	Name         string
	Instructions []Instruction
	Constants    []common.Value
	Names        []string
	Functions    map[string]*FunctionChunk
}

// NewChunk creates a new empty bytecode Chunk.
func NewChunk(name string) *Chunk {
	return &Chunk{
		Name:         name,
		Instructions: make([]Instruction, 0),
		Constants:    make([]common.Value, 0),
		Names:        make([]string, 0),
		Functions:    make(map[string]*FunctionChunk),
	}
}

// Emit appends an instruction with an argument to the chunk.
func (c *Chunk) Emit(op OpCode, arg int, line int) int {
	inst := Instruction{
		Op:   op,
		Arg:  arg,
		Line: line,
	}
	c.Instructions = append(c.Instructions, inst)
	return len(c.Instructions) - 1
}

// EmitOp appends an instruction without an argument (arg=0).
func (c *Chunk) EmitOp(op OpCode, line int) int {
	return c.Emit(op, 0, line)
}

// AddConstant adds a constant value to the pool and returns its index.
func (c *Chunk) AddConstant(val common.Value) int {
	// Look for existing identical scalar constant to deduplicate
	for i, existing := range c.Constants {
		if numA, okA := existing.(common.Number); okA {
			if numB, okB := val.(common.Number); okB && numA.Value == numB.Value && numA.IsReal == numB.IsReal {
				return i
			}
		} else if strA, okA := existing.(common.String); okA {
			if strB, okB := val.(common.String); okB && strA.Value == strB.Value {
				return i
			}
		}
	}
	c.Constants = append(c.Constants, val)
	return len(c.Constants) - 1
}

// AddName adds an identifier name to the names pool and returns its index.
func (c *Chunk) AddName(name string) int {
	for i, existing := range c.Names {
		if existing == name {
			return i
		}
	}
	c.Names = append(c.Names, name)
	return len(c.Names) - 1
}

// PatchJump updates the argument of an existing jump instruction at offset.
func (c *Chunk) PatchJump(offset int, target int) {
	if offset >= 0 && offset < len(c.Instructions) {
		c.Instructions[offset].Arg = target
	}
}

// CurrentOffset returns the index of the next instruction to be emitted.
func (c *Chunk) CurrentOffset() int {
	return len(c.Instructions)
}
