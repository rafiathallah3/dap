package common

import (
	"dap/internal/lexer"
	"dap/tools"
	"fmt"
)

func PrintValueAST(n Expr) string {
	switch n := n.(type) {
	// case lexer.Token:
	// 	return n.Value
	case NumberNode:
		return n.Token.Value
	case StringNode:
		return n.Token.Value
	case ListNode:
		hasil := "["
		for _, v := range n.ElementNode {
			hasil += PrintValueAST(v)
		}
		return hasil + "]"
	case NullNode:
		return "NULL"
	case BinOpNode:
		return fmt.Sprintf("(%v, %s, %v)", PrintValueAST(n.Left), n.Operator.Value, PrintValueAST(n.Right))
	case UnaryOpNode:
		return fmt.Sprintf("(%s, %s)", n.Operator.Value, PrintValueAST(n.Node))
	case VarAssignNode:
		return fmt.Sprintf("%s = %s", n.VarName.Value, PrintValueAST(n.ValueNode))
	case VarAccessNode:
		return n.VarNameTok.Value
	case IfNode:
		return fmt.Sprintf("IF %s THEN %s", PrintValueAST(n.Cases[0].Kondisi), PrintValueAST(n.Cases[0].Isi))
	case ForNode:
		switch n.StepValueNode.(type) {
		case NullNode:
			return fmt.Sprintf("FOR %s TO %s DO %s", PrintValueAST(n.StartValueNode), PrintValueAST(n.EndValueNode), PrintValueAST(n.BodyNode))
		}

		return fmt.Sprintf("FOR %s TO %s STEP %s DO %s", PrintValueAST(n.StartValueNode), PrintValueAST(n.EndValueNode), PrintValueAST(n.StepValueNode), PrintValueAST(n.BodyNode))
	case WhileNode:
		return fmt.Sprintf("WHIEL %s DO %s", PrintValueAST(n.KondisiNode), PrintValueAST(n.BodyNode))
	case *ParseResult:
		return PrintValueAST(n.Node)
	}

	return ""
}

type NumberNode struct {
	Token     lexer.Token
	Pos_Start *tools.Position
	Pos_End   *tools.Position
}

func (n NumberNode) expr() {}
func (n NumberNode) Print() string {
	return PrintValueAST(n)
}
func (n NumberNode) Name() string {
	return "NumberNode"
}
func (n NumberNode) GetPosStart() *tools.Position {
	return n.Pos_Start
}
func (n NumberNode) GetPosEnd() *tools.Position {
	return n.Pos_Start
}

type StringNode struct {
	Token     lexer.Token
	Pos_Start *tools.Position
	Pos_End   *tools.Position
}

func (n StringNode) expr() {}
func (n StringNode) Print() string {
	return PrintValueAST(n)
}
func (n StringNode) Name() string {
	return "StringNode"
}
func (n StringNode) GetPosStart() *tools.Position {
	return n.Pos_Start
}
func (n StringNode) GetPosEnd() *tools.Position {
	return n.Pos_Start
}

type ListNode struct {
	ElementNode []Expr
	Pos_Start   *tools.Position
	Pos_End     *tools.Position
}

func (n ListNode) expr() {}
func (n ListNode) Print() string {
	return PrintValueAST(n)
}
func (n ListNode) Name() string {
	return "ListNode"
}
func (n ListNode) GetPosStart() *tools.Position {
	return n.Pos_Start
}
func (n ListNode) GetPosEnd() *tools.Position {
	return n.Pos_Start
}

type NullNode struct {
	Token     any
	Pos_Start *tools.Position
	Pos_End   *tools.Position
}

func (n NullNode) expr() {}
func (n NullNode) Print() string {
	return PrintValueAST(n)
}
func (n NullNode) Name() string {
	return "NullNode"
}
func (n NullNode) GetPosStart() *tools.Position {
	return n.Pos_Start
}
func (n NullNode) GetPosEnd() *tools.Position {
	return n.Pos_Start
}

type BinOpNode struct {
	Operator  lexer.Token
	Left      Expr
	Right     Expr
	Pos_Start *tools.Position
	Pos_End   *tools.Position
}

func (n BinOpNode) expr() {}
func (n BinOpNode) Print() string {
	return PrintValueAST(n)
}
func (n BinOpNode) Name() string {
	return "BinOpNode"
}
func (n BinOpNode) GetPosStart() *tools.Position {
	return n.Pos_Start
}
func (n BinOpNode) GetPosEnd() *tools.Position {
	return n.Pos_Start
}

type UnaryOpNode struct {
	Operator  lexer.Token
	Node      Expr
	Pos_Start *tools.Position
	Pos_End   *tools.Position
}

func (n UnaryOpNode) expr() {}
func (n UnaryOpNode) Print() string {
	return PrintValueAST(n)
}
func (n UnaryOpNode) Name() string {
	return "NumberNode"
}
func (n UnaryOpNode) GetPosStart() *tools.Position {
	return n.Pos_Start
}
func (n UnaryOpNode) GetPosEnd() *tools.Position {
	return n.Pos_Start
}

type VarAssignNode struct {
	VarName   lexer.Token
	ValueNode Expr
	Pos_Start *tools.Position
	Pos_end   *tools.Position
}

func (n VarAssignNode) expr() {}
func (n VarAssignNode) Print() string {
	return PrintValueAST(n)
}
func (n VarAssignNode) Name() string {
	return "VarAssignNode"
}
func (n VarAssignNode) GetPosStart() *tools.Position {
	return n.Pos_Start
}
func (n VarAssignNode) GetPosEnd() *tools.Position {
	return n.Pos_end
}

type VarAccessNode struct {
	VarNameTok lexer.Token
	Pos_Start  *tools.Position
	Pos_end    *tools.Position
}

func (n VarAccessNode) expr() {}
func (n VarAccessNode) Print() string {
	return PrintValueAST(n)
}
func (n VarAccessNode) Name() string {
	return "VarAccessToken"
}
func (n VarAccessNode) GetPosStart() *tools.Position {
	return n.Pos_Start
}
func (n VarAccessNode) GetPosEnd() *tools.Position {
	return n.Pos_end
}

type IfCase struct {
	ElseCase
	Kondisi Expr
}

type ElseCase struct {
	Isi              Expr
	ShouldReturnNull bool
}

func (n ElseCase) expr() {}
func (n ElseCase) Print() string {
	return PrintValueAST(n)
}
func (n ElseCase) Name() string {
	return "ElseCase"
}
func (n ElseCase) GetPosStart() *tools.Position {
	return n.Isi.GetPosStart()
}
func (n ElseCase) GetPosEnd() *tools.Position {
	return n.Isi.GetPosEnd()
}

type IfNode struct {
	Cases     []IfCase
	Else_case *ElseCase
	Pos_Start *tools.Position
	Pos_end   *tools.Position
}

func (n IfNode) expr() {}
func (n IfNode) Print() string {
	return PrintValueAST(n)
}
func (n IfNode) Name() string {
	return "IfNode"
}
func (n IfNode) GetPosStart() *tools.Position {
	return n.Cases[0].Isi.GetPosStart()
}
func (n IfNode) GetPosEnd() *tools.Position {
	if n.Else_case != nil {
		return n.Else_case.Isi.GetPosEnd()
	}
	return n.Cases[len(n.Cases)-1].Isi.GetPosEnd()
}

type ForNode struct {
	VarNameTok       lexer.Token
	StartValueNode   Expr
	EndValueNode     Expr
	StepValueNode    Expr
	BodyNode         Expr
	ShouldReturnNull bool
	Pos_Start        *tools.Position
	Pos_end          *tools.Position
}

func (n ForNode) expr() {}
func (n ForNode) Print() string {
	return PrintValueAST(n)
}
func (n ForNode) Name() string {
	return "ForNode"
}
func (n ForNode) GetPosStart() *tools.Position {
	return n.Pos_Start
}
func (n ForNode) GetPosEnd() *tools.Position {
	return n.Pos_end
}

type WhileNode struct {
	KondisiNode      Expr
	BodyNode         Expr
	ShouldReturnNull bool
	Pos_Start        *tools.Position
	Pos_end          *tools.Position
}

func (n WhileNode) expr() {}
func (n WhileNode) Print() string {
	return PrintValueAST(n)
}
func (n WhileNode) Name() string {
	return "WhileNode"
}
func (n WhileNode) GetPosStart() *tools.Position {
	return n.Pos_Start
}
func (n WhileNode) GetPosEnd() *tools.Position {
	return n.Pos_end
}

type FuncNode struct {
	VarNameTok       *lexer.Token
	ArgNameToks      []lexer.Token
	BodyNode         Expr
	ShouldAutoReturn bool
	Pos_Start        *tools.Position
	Pos_end          *tools.Position
}

func (n FuncNode) expr() {}
func (n FuncNode) Print() string {
	return PrintValueAST(n)
}
func (n FuncNode) Name() string {
	return "FuncNode"
}
func (n FuncNode) GetPosStart() *tools.Position {
	if n.VarNameTok != nil {
		return n.VarNameTok.Pos_Start
	}

	if len(n.ArgNameToks) > 0 {
		return n.ArgNameToks[0].Pos_Start
	}

	return n.BodyNode.GetPosStart()
}
func (n FuncNode) GetPosEnd() *tools.Position {
	return n.BodyNode.GetPosEnd()
}

type CallNode struct {
	NodeToCall Expr
	ArgNodes   []Expr
	Pos_Start  *tools.Position
	Pos_end    *tools.Position
}

func (n CallNode) expr() {}
func (n CallNode) Print() string {
	return PrintValueAST(n)
}
func (n CallNode) Name() string {
	return "CallNode"
}
func (n CallNode) GetPosStart() *tools.Position {
	return n.NodeToCall.GetPosStart()
}
func (n CallNode) GetPosEnd() *tools.Position {
	if len(n.ArgNodes) > 0 {
		return n.ArgNodes[len(n.ArgNodes)-1].GetPosEnd()
	}

	return n.NodeToCall.GetPosEnd()
}

type ReturnNode struct {
	NodeToReturn Expr
	Pos_Start    *tools.Position
	Pos_end      *tools.Position
}

func (n ReturnNode) expr() {}
func (n ReturnNode) Print() string {
	return PrintValueAST(n)
}
func (n ReturnNode) Name() string {
	return "ReturnNode"
}
func (n ReturnNode) GetPosStart() *tools.Position {
	return n.NodeToReturn.GetPosStart()
}
func (n ReturnNode) GetPosEnd() *tools.Position {
	return n.NodeToReturn.GetPosEnd()
}

type ContinueNode struct {
	Pos_Start *tools.Position
	Pos_end   *tools.Position
}

func (n ContinueNode) expr() {}
func (n ContinueNode) Print() string {
	return PrintValueAST(n)
}
func (n ContinueNode) Name() string {
	return "ContinueNode"
}
func (n ContinueNode) GetPosStart() *tools.Position {
	return n.Pos_Start
}
func (n ContinueNode) GetPosEnd() *tools.Position {
	return n.Pos_end
}

type BreakNode struct {
	Pos_Start *tools.Position
	Pos_end   *tools.Position
}

func (n BreakNode) expr() {}
func (n BreakNode) Print() string {
	return PrintValueAST(n)
}
func (n BreakNode) Name() string {
	return "BreakNode"
}
func (n BreakNode) GetPosStart() *tools.Position {
	return n.Pos_Start
}
func (n BreakNode) GetPosEnd() *tools.Position {
	return n.Pos_end
}

type DictionaryNode struct {
	VariableDiBuat []Expr
	Pos_Start      *tools.Position
	Pos_end        *tools.Position
}

func (n DictionaryNode) expr() {}
func (n DictionaryNode) Print() string {
	return PrintValueAST(n)
}
func (n DictionaryNode) Name() string {
	return "DictionaryNode"
}
func (n DictionaryNode) GetPosStart() *tools.Position {
	return n.VariableDiBuat[0].GetPosStart()
}
func (n DictionaryNode) GetPosEnd() *tools.Position {
	return n.VariableDiBuat[len(n.VariableDiBuat)-1].GetPosEnd()
}
