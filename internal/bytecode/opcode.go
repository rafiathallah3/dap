package bytecode

import "fmt"

// OpCode represents a bytecode instruction operation.
type OpCode byte

const (
	OP_NOP OpCode = iota

	// Stack & Literal Operations
	OP_CONST // Push constant from pool [arg: const_index]
	OP_NIL   // Push null/nil
	OP_TRUE  // Push boolean true
	OP_FALSE // Push boolean false
	OP_POP   // Pop top of stack
	OP_DUP   // Duplicate top of stack

	// Variable Scopes
	OP_LOAD  // Load variable by name index [arg: name_index]
	OP_STORE // Store into variable by name index [arg: name_index]

	// Arithmetic Operations
	OP_ADD  // +
	OP_SUB  // -
	OP_MUL  // *
	OP_DIV  // /
	OP_IDIV // div
	OP_MOD  // %
	OP_POW  // **
	OP_NEG  // Unary minus

	// Bitwise Operations
	OP_BIT_AND // &
	OP_BIT_OR  // |
	OP_BIT_XOR // ^
	OP_BIT_NOT // ~
	OP_SHL     // <<
	OP_SHR     // >>

	// Comparison & Logical Operations
	OP_EQ  // ==
	OP_NEQ // !=
	OP_LT  // <
	OP_LTE // <=
	OP_GT  // >
	OP_GTE // >=
	OP_NOT // not / !

	// Control Flow & Branching
	OP_JUMP          // Unconditional relative jump [arg: target_offset]
	OP_JUMP_IF_FALSE // Jump if top of stack is falsy [arg: target_offset]
	OP_JUMP_IF_TRUE  // Jump if top of stack is truthy [arg: target_offset]
	OP_LOOP          // Unconditional backward jump [arg: backward_offset]

	// Composite Structures
	OP_ARRAY_MAKE  // Create array with bounds [arg: count]
	OP_ARRAY_GET   // Index array access
	OP_ARRAY_SET   // Store array index
	OP_STRUCT_GET  // Get struct field [arg: field_name_index]
	OP_STRUCT_SET  // Set struct field [arg: field_name_index]

	// Subroutines & I/O
	OP_CALL   // Call function [arg: arg_count]
	OP_RETURN // Return from function
	OP_PRINT  // Print / output top of stack
	OP_READ   // Read / input into variable
	OP_HALT   // Terminate execution
)

// OpCodeString returns the mnemonic string representation of an OpCode.
func (op OpCode) String() string {
	switch op {
	case OP_NOP:
		return "OP_NOP"
	case OP_CONST:
		return "OP_CONST"
	case OP_NIL:
		return "OP_NIL"
	case OP_TRUE:
		return "OP_TRUE"
	case OP_FALSE:
		return "OP_FALSE"
	case OP_POP:
		return "OP_POP"
	case OP_DUP:
		return "OP_DUP"
	case OP_LOAD:
		return "OP_LOAD"
	case OP_STORE:
		return "OP_STORE"
	case OP_ADD:
		return "OP_ADD"
	case OP_SUB:
		return "OP_SUB"
	case OP_MUL:
		return "OP_MUL"
	case OP_DIV:
		return "OP_DIV"
	case OP_IDIV:
		return "OP_IDIV"
	case OP_MOD:
		return "OP_MOD"
	case OP_POW:
		return "OP_POW"
	case OP_NEG:
		return "OP_NEG"
	case OP_BIT_AND:
		return "OP_BIT_AND"
	case OP_BIT_OR:
		return "OP_BIT_OR"
	case OP_BIT_XOR:
		return "OP_BIT_XOR"
	case OP_BIT_NOT:
		return "OP_BIT_NOT"
	case OP_SHL:
		return "OP_SHL"
	case OP_SHR:
		return "OP_SHR"
	case OP_EQ:
		return "OP_EQ"
	case OP_NEQ:
		return "OP_NEQ"
	case OP_LT:
		return "OP_LT"
	case OP_LTE:
		return "OP_LTE"
	case OP_GT:
		return "OP_GT"
	case OP_GTE:
		return "OP_GTE"
	case OP_NOT:
		return "OP_NOT"
	case OP_JUMP:
		return "OP_JUMP"
	case OP_JUMP_IF_FALSE:
		return "OP_JUMP_IF_FALSE"
	case OP_JUMP_IF_TRUE:
		return "OP_JUMP_IF_TRUE"
	case OP_LOOP:
		return "OP_LOOP"
	case OP_ARRAY_MAKE:
		return "OP_ARRAY_MAKE"
	case OP_ARRAY_GET:
		return "OP_ARRAY_GET"
	case OP_ARRAY_SET:
		return "OP_ARRAY_SET"
	case OP_STRUCT_GET:
		return "OP_STRUCT_GET"
	case OP_STRUCT_SET:
		return "OP_STRUCT_SET"
	case OP_CALL:
		return "OP_CALL"
	case OP_RETURN:
		return "OP_RETURN"
	case OP_PRINT:
		return "OP_PRINT"
	case OP_READ:
		return "OP_READ"
	case OP_HALT:
		return "OP_HALT"
	default:
		return fmt.Sprintf("OP_UNKNOWN(%d)", op)
	}
}
