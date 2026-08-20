package common

import (
	"dap/internal/lexer"
	"dap/tools"
	"fmt"
)

// Expr itu interface utama yang wajib diimplement sama semua AST node.
type Expr interface {
	expr()
	GetPosStart() *tools.Position
	GetPosEnd() *tools.Position
	Print() string
	Name() string
}

// Token ini wrapper AST buat token dari lexer.
type Token struct {
	Kind      lexer.TokenKind
	Value     string
	Pos_Start *tools.Position // Dijaga biar tetep kompatibel sama struktur AST
	Pos_End   *tools.Position
}

func (t Token) expr() {}
func (t Token) Print() string {
	return t.Value
}
func (t Token) Name() string {
	return "Token"
}
func (t Token) GetPosStart() *tools.Position {
	return t.Pos_Start
}
func (t Token) GetPosEnd() *tools.Position {
	return t.Pos_End
}

// ParseResult buat nyimpen hasil parsing di tiap step (node, error, plus tracking advance token).
type ParseResult struct {
	Error                      *Error
	Node                       Expr
	ProgramName                string
	ToReverseCount             int
	AdvanceCount               int
	LastRegisteredAdvanceCount int
	Pos_Start                  *tools.Position
	Pos_End                    *tools.Position
}

// RegisterAdvancement buat nyatet kalo parser udah advance 1 token.
func (pr *ParseResult) RegisterAdvancement() {
	pr.LastRegisteredAdvanceCount = 1
	pr.AdvanceCount++
}

// Register_Advancement alias buat backward compatibility.
func (pr *ParseResult) Register_Advancement() {
	pr.RegisterAdvancement()
}

// Register buat nge-record progress dan error dari ParseResult lain atau Token ke dalam ParseResult ini.
func (pr *ParseResult) Register(res any) Expr {
	switch r := res.(type) {
	case *ParseResult:
		pr.LastRegisteredAdvanceCount = r.AdvanceCount
		pr.AdvanceCount += r.AdvanceCount
		if r.Error != nil {
			pr.Error = r.Error
		}
		return r.Node
	case lexer.Token:
		tk := Token{
			Kind:      r.Kind,
			Value:     r.Value,
			Pos_Start: r.PosStart,
			Pos_End:   r.PosEnd,
		}
		return tk
	case Expr:
		return r
	}
	return nil
}

// TryRegister buat nyoba register result. Kalo error, cuma nyatet backtrack count tanpa langsung commit error-nya.
func (pr *ParseResult) TryRegister(res any) Expr {
	if r, ok := res.(*ParseResult); ok {
		if r.Error != nil {
			pr.ToReverseCount = r.AdvanceCount
			return nil
		}
	}
	return pr.Register(res)
}

// Try_register alias buat backward compatibility.
func (pr *ParseResult) Try_register(res any) Expr {
	return pr.TryRegister(res)
}

// Success buat nge-set node hasil parsing yang valid.
func (pr *ParseResult) Success(node Expr) Expr {
	pr.Node = node
	return pr
}

// Failure buat nempelin error pas proses parsing gagal.
func (pr *ParseResult) Failure(err *Error) Expr {
	if pr.Error == nil || pr.AdvanceCount == 0 {
		pr.Error = err
	}
	return pr
}

func (pr *ParseResult) expr() {}
func (pr *ParseResult) Print() string {
	fmt.Println(PrintValueAST(pr))
	return ""
}
func (pr *ParseResult) Name() string {
	return "ParseResult"
}
func (pr *ParseResult) GetPosStart() *tools.Position {
	return pr.Pos_Start
}
func (pr *ParseResult) GetPosEnd() *tools.Position {
	return pr.Pos_End
}
