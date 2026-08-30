package bytecode

import (
	"fmt"
	"strings"
)

// Disassemble converts a Chunk into a human-readable disassembly text.
func Disassemble(c *Chunk) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("== DISASSEMBLY %s ==\n", c.Name))
	sb.WriteString(fmt.Sprintf("Constants: %d, Names: %d, Instructions: %d\n", len(c.Constants), len(c.Names), len(c.Instructions)))
	sb.WriteString("--------------------------------------------------\n")

	for i := 0; i < len(c.Instructions); i++ {
		inst := c.Instructions[i]
		lineStr := fmt.Sprintf("%4d", inst.Line)
		if i > 0 && c.Instructions[i-1].Line == inst.Line {
			lineStr = "   |"
		}

		opName := inst.Op.String()
		detail := ""

		switch inst.Op {
		case OP_CONST:
			if inst.Arg >= 0 && inst.Arg < len(c.Constants) {
				detail = fmt.Sprintf("%d (%v)", inst.Arg, c.Constants[inst.Arg].Print())
			}
		case OP_LOAD, OP_STORE, OP_STRUCT_GET, OP_STRUCT_SET:
			if inst.Arg >= 0 && inst.Arg < len(c.Names) {
				detail = fmt.Sprintf("%d ('%s')", inst.Arg, c.Names[inst.Arg])
			}
		case OP_JUMP, OP_JUMP_IF_FALSE, OP_JUMP_IF_TRUE, OP_LOOP:
			detail = fmt.Sprintf("-> %04d", inst.Arg)
		case OP_CALL:
			detail = fmt.Sprintf("args: %d", inst.Arg)
		}

		sb.WriteString(fmt.Sprintf("%04d  %s  %-18s %s\n", i, lineStr, opName, detail))
	}

	for _, fn := range c.Functions {
		sb.WriteString("\n")
		sb.WriteString(Disassemble(fn.Chunk))
	}

	return sb.String()
}
