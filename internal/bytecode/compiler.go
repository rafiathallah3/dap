package bytecode

import (
	"dap/internal/common"
	"dap/internal/lexer"
	"dap/tools"
	"fmt"
	"strconv"
	"strings"
)

// Compiler compiles AST expressions into bytecode Chunks.
type Compiler struct {
	chunk *Chunk
}

// NewCompiler creates a new Compiler targeting a new Chunk.
func NewCompiler(name string) *Compiler {
	return &Compiler{
		chunk: NewChunk(name),
	}
}

// Compile compiles the given AST root expression and returns the compiled Chunk.
func Compile(node common.Expr) (*Chunk, error) {
	c := NewCompiler("<main>")
	if err := c.compileNode(node); err != nil {
		return nil, err
	}
	c.chunk.EmitOp(OP_HALT, 0)
	return c.chunk, nil
}

func (c *Compiler) getLine(node common.Expr) int {
	if node != nil && node.GetPosStart() != nil {
		return node.GetPosStart().Ln + 1
	}
	return 1
}

func (c *Compiler) compileNode(node common.Expr) error {
	if node == nil {
		return nil
	}

	line := c.getLine(node)

	switch n := node.(type) {
	case *common.ParseResult:
		return c.compileNode(n.Node)

	case common.ListNode:
		for _, stmt := range n.ElementNode {
			if err := c.compileNode(stmt); err != nil {
				return err
			}
		}
	case common.DictionaryNode:
		for _, decl := range n.VariableDiBuat {
			if assign, ok := decl.(common.VarAssignNode); ok {
				typeTarget := assign.ValueNode
				if assign.TypeNode != nil {
					typeTarget = assign.TypeNode
				}
				if typeAccess, ok := typeTarget.(common.VarAccessNode); ok {
					c.chunk.VarTypes[assign.VarName.Value] = typeAccess.VarNameTok.Value
				}
			}
		}
		return nil

	case common.NumberNode:
		val, err := strconv.ParseFloat(n.Token.Value, 64)
		if err != nil {
			return fmt.Errorf("invalid number literal: %s", n.Token.Value)
		}
		isReal := strings.Contains(n.Token.Value, ".")
		idx := c.chunk.AddConstant(common.Number{Value: val, IsReal: isReal})
		c.chunk.Emit(OP_CONST, idx, line)
		return nil

	case common.StringNode:
		idx := c.chunk.AddConstant(common.String{Value: n.Token.Value})
		c.chunk.Emit(OP_CONST, idx, line)
		return nil

	case common.NullNode:
		c.chunk.EmitOp(OP_NIL, line)
		return nil

	case common.VarAccessNode:
		varName := n.VarNameTok.Value
		if strings.EqualFold(varName, "true") {
			c.chunk.EmitOp(OP_TRUE, line)
			return nil
		} else if strings.EqualFold(varName, "false") {
			c.chunk.EmitOp(OP_FALSE, line)
			return nil
		}
		nameIdx := c.chunk.AddName(varName)
		c.chunk.Emit(OP_LOAD, nameIdx, line)
		return nil

	case common.VarAssignNode:
		if n.ValueNode != nil {
			if err := c.compileNode(n.ValueNode); err != nil {
				return err
			}
		} else {
			c.chunk.EmitOp(OP_NIL, line)
		}
		nameIdx := c.chunk.AddName(n.VarName.Value)
		c.chunk.Emit(OP_STORE, nameIdx, line)
		return nil

	case common.BinOpNode:
		if err := c.compileNode(n.Left); err != nil {
			return err
		}
		if err := c.compileNode(n.Right); err != nil {
			return err
		}

		switch n.Operator.Kind {
		case lexer.PLUS:
			c.chunk.EmitOp(OP_ADD, line)
		case lexer.DASH:
			c.chunk.EmitOp(OP_SUB, line)
		case lexer.STAR:
			c.chunk.EmitOp(OP_MUL, line)
		case lexer.SLASH:
			c.chunk.EmitOp(OP_DIV, line)
		case lexer.DIV:
			c.chunk.EmitOp(OP_IDIV, line)
		case lexer.PERCENT:
			c.chunk.EmitOp(OP_MOD, line)
		case lexer.POWER, lexer.DOUBLE_STAR:
			c.chunk.EmitOp(OP_POW, line)
		case lexer.BIT_AND:
			c.chunk.EmitOp(OP_BIT_AND, line)
		case lexer.BIT_OR:
			c.chunk.EmitOp(OP_BIT_OR, line)
		case lexer.BIT_XOR:
			c.chunk.EmitOp(OP_BIT_XOR, line)
		case lexer.SHIFT_LEFT:
			c.chunk.EmitOp(OP_SHL, line)
		case lexer.SHIFT_RIGHT:
			c.chunk.EmitOp(OP_SHR, line)
		case lexer.EQUALS:
			c.chunk.EmitOp(OP_EQ, line)
		case lexer.NOT_EQUALS:
			c.chunk.EmitOp(OP_NEQ, line)
		case lexer.LESS:
			c.chunk.EmitOp(OP_LT, line)
		case lexer.LESS_EQUALS:
			c.chunk.EmitOp(OP_LTE, line)
		case lexer.GREATER:
			c.chunk.EmitOp(OP_GT, line)
		case lexer.GREATER_EQUALS:
			c.chunk.EmitOp(OP_GTE, line)
		case lexer.AND:
			// AND operator (both popped and checked)
			c.chunk.EmitOp(OP_EQ, line)
		default:
			return fmt.Errorf("unsupported binary operator: %s", n.Operator.Value)
		}
		return nil

	case common.UnaryOpNode:
		if err := c.compileNode(n.Node); err != nil {
			return err
		}
		switch n.Operator.Kind {
		case lexer.DASH:
			c.chunk.EmitOp(OP_NEG, line)
		case lexer.NOT:
			c.chunk.EmitOp(OP_NOT, line)
		case lexer.BIT_NOT:
			c.chunk.EmitOp(OP_BIT_NOT, line)
		case lexer.PLUS:
			// No-op for unary plus
		default:
			return fmt.Errorf("unsupported unary operator: %s", n.Operator.Value)
		}
		return nil

	case common.IfNode:
		var endJumps []int

		for _, ifCase := range n.Cases {
			if err := c.compileNode(ifCase.Kondisi); err != nil {
				return err
			}
			falseJump := c.chunk.Emit(OP_JUMP_IF_FALSE, 0, line)

			if err := c.compileNode(ifCase.Isi); err != nil {
				return err
			}
			exitJump := c.chunk.Emit(OP_JUMP, 0, line)
			endJumps = append(endJumps, exitJump)

			c.chunk.PatchJump(falseJump, c.chunk.CurrentOffset())
		}

		if n.Else_case != nil && n.Else_case.Isi != nil {
			if err := c.compileNode(n.Else_case.Isi); err != nil {
				return err
			}
		}

		for _, jmp := range endJumps {
			c.chunk.PatchJump(jmp, c.chunk.CurrentOffset())
		}
		return nil

	case common.WhileNode:
		loopStart := c.chunk.CurrentOffset()
		if err := c.compileNode(n.KondisiNode); err != nil {
			return err
		}
		exitJump := c.chunk.Emit(OP_JUMP_IF_FALSE, 0, line)

		if err := c.compileNode(n.BodyNode); err != nil {
			return err
		}
		c.chunk.Emit(OP_LOOP, loopStart, line)
		c.chunk.PatchJump(exitJump, c.chunk.CurrentOffset())
		return nil

	case common.RepeatNode:
		loopStart := c.chunk.CurrentOffset()
		if err := c.compileNode(n.BodyNode); err != nil {
			return err
		}
		if err := c.compileNode(n.KondisiNode); err != nil {
			return err
		}
		c.chunk.Emit(OP_JUMP_IF_FALSE, loopStart, line)
		return nil

	case common.ForNode:
		// 1. Initial assignment: var = start
		if err := c.compileNode(n.StartValueNode); err != nil {
			return err
		}
		varIdx := c.chunk.AddName(n.VarNameTok.Value)
		c.chunk.Emit(OP_STORE, varIdx, line)

		// 2. Loop start condition
		loopStart := c.chunk.CurrentOffset()
		c.chunk.Emit(OP_LOAD, varIdx, line)
		if err := c.compileNode(n.EndValueNode); err != nil {
			return err
		}
		if n.IsDown {
			c.chunk.EmitOp(OP_GTE, line)
		} else {
			c.chunk.EmitOp(OP_LTE, line)
		}
		exitJump := c.chunk.Emit(OP_JUMP_IF_FALSE, 0, line)

		// 3. Body
		if err := c.compileNode(n.BodyNode); err != nil {
			return err
		}

		// 4. Step increment/decrement: var = var + step
		c.chunk.Emit(OP_LOAD, varIdx, line)
		if n.StepValueNode != nil {
			if _, isNull := n.StepValueNode.(common.NullNode); !isNull {
				if err := c.compileNode(n.StepValueNode); err != nil {
					return err
				}
			} else {
				oneIdx := c.chunk.AddConstant(common.Number{Value: 1, IsReal: false})
				c.chunk.Emit(OP_CONST, oneIdx, line)
			}
		} else {
			oneIdx := c.chunk.AddConstant(common.Number{Value: 1, IsReal: false})
			c.chunk.Emit(OP_CONST, oneIdx, line)
		}

		if n.IsDown {
			c.chunk.EmitOp(OP_SUB, line)
		} else {
			c.chunk.EmitOp(OP_ADD, line)
		}
		c.chunk.Emit(OP_STORE, varIdx, line)

		c.chunk.Emit(OP_LOOP, loopStart, line)
		c.chunk.PatchJump(exitJump, c.chunk.CurrentOffset())
		return nil

	case common.CallNode:
		callTargetName := ""
		if varAcc, ok := n.NodeToCall.(common.VarAccessNode); ok {
			callTargetName = varAcc.VarNameTok.Value
		}

		// Handle built-in output / print calls directly
		if tools.IsBuiltinFunction(callTargetName) {
			fnType := tools.BuiltinFunctions[callTargetName]
			if fnType == "Print" {
				for _, arg := range n.ArgNodes {
					if err := c.compileNode(arg); err != nil {
						return err
					}
					c.chunk.EmitOp(OP_PRINT, line)
				}
				return nil
			}
		}

		// Generic Function Call
		for _, arg := range n.ArgNodes {
			if err := c.compileNode(arg); err != nil {
				return err
			}
		}
		nameIdx := c.chunk.AddName(callTargetName)
		c.chunk.Emit(OP_LOAD, nameIdx, line)
		c.chunk.Emit(OP_CALL, len(n.ArgNodes), line)
		return nil

	case common.FuncNode:
		funcName := "<anonymous>"
		if n.VarNameTok != nil {
			funcName = n.VarNameTok.Value
		}

		fnCompiler := NewCompiler(funcName)
		if err := fnCompiler.compileNode(n.BodyNode); err != nil {
			return err
		}
		fnCompiler.chunk.EmitOp(OP_RETURN, line)

		var argNames []string
		for _, argTok := range n.ArgNameToks {
			argNames = append(argNames, argTok.Value)
		}

		fnChunk := &FunctionChunk{
			Name:     funcName,
			ArgNames: argNames,
			Chunk:    fnCompiler.chunk,
		}
		c.chunk.Functions[funcName] = fnChunk
		return nil

	case common.ArrayIndexNode:
		if err := c.compileNode(n.Left); err != nil {
			return err
		}
		if err := c.compileNode(n.Index); err != nil {
			return err
		}
		c.chunk.EmitOp(OP_ARRAY_GET, line)
		return nil

	case common.ArrayAssignNode:
		if err := c.compileNode(n.ArrayAccess.Left); err != nil {
			return err
		}
		if err := c.compileNode(n.ArrayAccess.Index); err != nil {
			return err
		}
		if err := c.compileNode(n.ValueNode); err != nil {
			return err
		}
		c.chunk.EmitOp(OP_ARRAY_SET, line)
		return nil

	case common.MemberAccessNode:
		if err := c.compileNode(n.Object); err != nil {
			return err
		}
		fieldIdx := c.chunk.AddName(n.MemberTok.Value)
		c.chunk.Emit(OP_STRUCT_GET, fieldIdx, line)
		return nil

	case common.MemberAssignNode:
		if err := c.compileNode(n.MemberAccess.Object); err != nil {
			return err
		}
		if err := c.compileNode(n.ValueNode); err != nil {
			return err
		}
		fieldIdx := c.chunk.AddName(n.MemberAccess.MemberTok.Value)
		c.chunk.Emit(OP_STRUCT_SET, fieldIdx, line)
		return nil

	case common.TypeAliasNode, common.StructTypeNode:
		// Type declarations do not emit runtime opcodes
		return nil
	}

	return nil
}
