package common

import (
	"dap/internal/lexer"
	"dap/tools"
	"fmt"
)

type Expr interface {
	expr()
	GetPosStart() *tools.Position
	GetPosEnd() *tools.Position
	Print() string
	Name() string
}

type Token struct {
	Kind      lexer.TokenKind
	Value     string
	Pos_Start *tools.Position
	Pos_End   *tools.Position
}

func (token Token) expr() {}
func (token Token) Print() string {
	return token.Value
}
func (n Token) Name() string {
	return "Token"
}
func (n Token) GetPosStart() *tools.Position {
	return n.Pos_Start
}
func (n Token) GetPosEnd() *tools.Position {
	return n.Pos_Start
}

type ParseResult struct {
	Error                      *Error
	Node                       Expr
	ToReverseCount             int
	AdvanceCount               int
	LastRegisteredAdvanceCount int
	Pos_Start                  *tools.Position
	Pos_End                    *tools.Position
}

func (parseResult *ParseResult) Register_Advancement() {
	parseResult.LastRegisteredAdvanceCount = 1
	parseResult.AdvanceCount++
}

func (parseResult *ParseResult) Register(res interface{}) Expr {
	switch res := res.(type) {
	case *ParseResult:
		parseResult.LastRegisteredAdvanceCount = res.AdvanceCount
		parseResult.AdvanceCount += res.AdvanceCount
		if res.Error != nil {
			parseResult.Error = res.Error
		}

		return res.Node
	case lexer.Token:
		tk := Token{ //So ugly to do, I'll try to improve it after everything is done. [18/01/2025 18:47]
			Kind:      res.Kind,
			Value:     res.Value,
			Pos_Start: res.Pos_Start,
			Pos_End:   res.Pos_End,
		}
		return tk
	}

	return res.(Expr)
}

func (parseResult *ParseResult) Try_register(res interface{}) Expr {
	switch res := res.(type) {
	case *ParseResult:
		if res.Error != nil {
			parseResult.ToReverseCount = res.AdvanceCount
			return nil
		}

	}

	return parseResult.Register(res)
}

func (parserResult *ParseResult) Success(node Expr) Expr {
	parserResult.Node = node
	return parserResult
}

func (parserResult *ParseResult) Failure(error *Error) Expr {
	if parserResult.Error == nil || parserResult.AdvanceCount == 0 {
		parserResult.Error = error
	}

	return parserResult
}

func (parseResult *ParseResult) expr() {}
func (parseResult *ParseResult) Print() string {
	fmt.Println(PrintValueAST(parseResult))
	return ""
}
func (n *ParseResult) Name() string {
	return "ParseResult"
}
func (n *ParseResult) GetPosStart() *tools.Position {
	return n.Pos_Start
}
func (n *ParseResult) GetPosEnd() *tools.Position {
	return n.Pos_Start
}
