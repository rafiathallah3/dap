package parser

import (
	"dap/internal/common"
	"dap/internal/lexer"
	"fmt"
	"slices"
)

type parser struct {
	tokens    []lexer.Token
	tok_index int
}

func CreateParser(tokens []lexer.Token) *parser {
	p := &parser{
		tokens:    tokens,
		tok_index: -1,
	}

	p.advance()
	return p
}

func (p *parser) Parse() common.Expr {
	res := p.statements()

	switch res := res.(type) {
	case *common.ParseResult:
		if res.Error == nil && p.currentToken().Kind != lexer.EOF {
			errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, "Expected +, -, *, / or ^")
			return res.Failure(&errorNya)
		}
	}

	return res
}

func (p *parser) currentToken() lexer.Token {
	return p.tokens[p.tok_index]
}

func (p *parser) advance() lexer.Token {
	var tk lexer.Token
	p.tok_index++

	if p.tok_index < len(p.tokens) {
		tk = p.currentToken()
	} else {
		p.tok_index--
		tk = p.currentToken()
	}

	return tk
}

func (p *parser) reverse(amount int) lexer.Token {
	p.tok_index -= amount
	return p.currentToken()
}

func (p *parser) list_expr() common.Expr {
	res := &common.ParseResult{}
	elementNodes := make([]common.Expr, 0)

	if p.currentToken().Kind != lexer.OPEN_BRACKET {
		errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, "Expected '['")
		return res.Failure(&errorNya)
	}
	posStart := p.currentToken().Pos_Start.Copy()

	res.Register_Advancement()
	p.advance()

	if p.currentToken().Kind == lexer.CLOSE_BRACKET {
		res.Register_Advancement()
		p.advance()
	} else {
		elementNodes = append(elementNodes, res.Register(p.expr()))
		if res.Error != nil {
			errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, "Expected, ']', 'var', 'if', 'for', 'while', 'fun', number, identifier")
			return res.Failure(&errorNya)
		}

		for p.currentToken().Kind == lexer.COMMA {
			res.Register_Advancement()
			p.advance()

			elementNodes = append(elementNodes, res.Register(p.expr()))
			if res.Error != nil {
				return res
			}
		}

		if p.currentToken().Kind != lexer.CLOSE_BRACKET {
			errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, "Expected, ']' or ','")
			return res.Failure(&errorNya)
		}

		res.Register_Advancement()
		p.advance()
	}

	return res.Success(common.ListNode{
		ElementNode: elementNodes,
		Pos_Start:   posStart,
		Pos_End:     p.currentToken().Pos_End.Copy(),
	})
}

func (p *parser) if_expr() common.Expr {
	res := &common.ParseResult{}
	all_cases := res.Register(p.if_expr_cases(lexer.IF))
	if res.Error != nil {
		return res
	}

	return res.Success(common.IfNode{
		Cases:     all_cases.(common.IfNode).Cases,
		Else_case: all_cases.(common.IfNode).Else_case,
	})
}

func (p *parser) if_expr_b() common.Expr {
	return p.if_expr_cases(lexer.ELIF)
}

func (p *parser) if_expr_c() common.Expr {
	res := &common.ParseResult{}
	var else_case common.ElseCase

	if p.currentToken().Kind == lexer.ELSE {
		res.Register_Advancement()
		p.advance()

		if p.currentToken().Kind == lexer.NEWLINE {
			res.Register_Advancement()
			p.advance()

			statements := res.Register(p.statements())
			if res.Error != nil {
				return res
			}
			else_case = common.ElseCase{
				Isi:              statements,
				ShouldReturnNull: true,
			}

			if p.currentToken().Kind == lexer.END || p.currentToken().Kind == lexer.ENDIF {
				res.Register_Advancement()
				p.advance()
			} else {
				errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, "expected 'end' or 'endif'")
				return res.Failure(&errorNya)
			}
		} else {
			expr := res.Register(p.statement())
			if res.Error != nil {
				return res
			}

			else_case = common.ElseCase{
				Isi:              expr,
				ShouldReturnNull: false,
			}
		}
	}

	return res.Success(else_case)
}

func (p *parser) if_expr_b_or_c() common.Expr {
	res := &common.ParseResult{}
	cases := make([]common.IfCase, 0)
	var else_case common.ElseCase

	if p.currentToken().Kind == lexer.ELIF {
		all_cases := res.Register(p.if_expr_b())
		if res.Error != nil {
			return res
		}

		cases = all_cases.(common.IfNode).Cases
		ifNodeElse := all_cases.(common.IfNode).Else_case
		else_case = *ifNodeElse
	} else {
		resElseCase := res.Register(p.if_expr_c())
		if res.Error != nil {
			return res
		}

		else_case = resElseCase.(common.ElseCase)
	}

	return res.Success(common.IfNode{
		Cases:     cases,
		Else_case: &else_case,
	})
}

func (p *parser) if_expr_cases(caseKeyword lexer.TokenKind) common.Expr {
	res := &common.ParseResult{}
	var cases []common.IfCase
	var else_case common.ElseCase

	if p.currentToken().Kind != caseKeyword {
		errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, fmt.Sprintf("Expected '%v'", caseKeyword))
		return res.Failure(&errorNya)
	}

	res.Register_Advancement()
	p.advance()

	condition := res.Register(p.expr())
	if res.Error != nil {
		return res
	}

	if p.currentToken().Kind != lexer.THEN {
		errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, "Expected 'then'")
		return res.Failure(&errorNya)
	}

	res.Register_Advancement()
	p.advance()

	if p.currentToken().Kind == lexer.NEWLINE {
		res.Register_Advancement()
		p.advance()

		statements := res.Register(p.statements())
		if res.Error != nil {
			return res
		}
		cases = append(cases, common.IfCase{
			Kondisi: condition,
			ElseCase: common.ElseCase{
				Isi:              statements,
				ShouldReturnNull: true,
			},
		})

		if p.currentToken().Kind == lexer.END || p.currentToken().Kind == lexer.ENDIF {
			res.Register_Advancement()
			p.advance()
		} else {
			all_cases := res.Register(p.if_expr_b_or_c())
			if res.Error != nil {
				return res
			}

			else_case = *all_cases.(common.IfNode).Else_case
			cases = append(cases, all_cases.(common.IfNode).Cases...)
		}
	} else {
		expr := res.Register(p.statement())

		if res.Error != nil {
			return res
		}

		cases = append(cases, common.IfCase{
			Kondisi: condition,
			ElseCase: common.ElseCase{
				Isi:              expr,
				ShouldReturnNull: false,
			},
		})

		all_cases := res.Register(p.if_expr_b_or_c())

		else_case = *all_cases.(common.IfNode).Else_case
		cases = append(cases, all_cases.(common.IfNode).Cases...)
	}

	return res.Success(common.IfNode{
		Cases:     cases,
		Else_case: &else_case,
	})

	// for p.currentToken().Kind == lexer.ELIF {
	// 	res.Register_Advancement()
	// 	p.advance()

	// 	condition = res.Register(p.expr())
	// 	if res.Error != nil {
	// 		return res
	// 	}

	// 	if p.currentToken().Kind != lexer.THEN {
	// 		errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, "Expected 'then'")
	// 		return res.Failure(&errorNya)
	// 	}

	// 	res.Register_Advancement()
	// 	p.advance()

	// 	expr = res.Register(p.expr())
	// 	if res.Error != nil {
	// 		return res
	// 	}

	// 	cases = append(cases, common.IfCase{
	// 		Kondisi: condition,
	// 		Isi:     expr,
	// 	})
	// }

	// if p.currentToken().Kind == lexer.ELSE {
	// 	res.Register_Advancement()
	// 	p.advance()

	// 	expr = res.Register(p.expr())
	// 	if res.Error != nil {
	// 		return res
	// 	}

	// 	else_case = &expr
	// }

	// pos_end := cases[len(cases)-1].Kondisi.GetPosEnd()
	// if else_case != nil {
	// 	pos_end = (*else_case).GetPosEnd()
	// }

	// return res.Success(common.IfNode{Cases: cases, Else_case: else_case, Pos_Start: cases[0].Kondisi.GetPosStart(), Pos_end: pos_end})
}

func (p *parser) for_expr() common.Expr {
	res := &common.ParseResult{}

	if p.currentToken().Kind != lexer.FOR {
		errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, "Expected 'for'")
		return res.Failure(&errorNya)
	}

	res.Register_Advancement()
	p.advance()

	if p.currentToken().Kind != lexer.IDENTIFIER {
		errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, "Expected identifier")
		return res.Failure(&errorNya)
	}

	varName := p.currentToken()
	res.Register_Advancement()
	p.advance()

	if p.currentToken().Kind != lexer.ASSIGNMENT {
		errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, "Expected '='")
		return res.Failure(&errorNya)
	}

	res.Register_Advancement()
	p.advance()

	startValue := res.Register(p.expr())
	if res.Error != nil {
		return res
	}

	if p.currentToken().Kind != lexer.TO {
		errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, "Expected 'to'")
		return res.Failure(&errorNya)
	}

	res.Register_Advancement()
	p.advance()

	EndValue := res.Register(p.expr())
	if res.Error != nil {
		return res
	}

	var StepValue common.Expr = common.NullNode{}
	if p.currentToken().Kind == lexer.STEP {
		res.Register_Advancement()
		p.advance()

		StepValue = res.Register(p.expr())
		if res.Error != nil {
			return res
		}
	}

	if p.currentToken().Kind != lexer.DO {
		errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, "Expected 'do'")
		return res.Failure(&errorNya)
	}

	res.Register_Advancement()
	p.advance()

	if p.currentToken().Kind == lexer.NEWLINE {
		res.Register_Advancement()
		p.advance()

		body := res.Register(p.statements())
		if res.Error != nil {
			return res
		}

		if p.currentToken().Kind != lexer.END && p.currentToken().Kind != lexer.ENDFOR {
			errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, "Expected 'end' or 'endfor'")
			return res.Failure(&errorNya)
		}

		res.Register_Advancement()
		p.advance()

		return res.Success(common.ForNode{
			VarNameTok:       varName,
			StartValueNode:   startValue,
			EndValueNode:     EndValue,
			StepValueNode:    StepValue,
			BodyNode:         body,
			ShouldReturnNull: true,
			Pos_Start:        varName.Pos_Start,
			Pos_end:          body.GetPosEnd(),
		})
	}

	IsiValue := res.Register(p.statement())
	if res.Error != nil {
		return res
	}

	return res.Success(common.ForNode{
		VarNameTok:       varName,
		StartValueNode:   startValue,
		EndValueNode:     EndValue,
		StepValueNode:    StepValue,
		BodyNode:         IsiValue,
		ShouldReturnNull: false,
		Pos_Start:        varName.Pos_Start,
		Pos_end:          IsiValue.GetPosEnd(),
	})
}

func (p *parser) while_expr() common.Expr {
	res := &common.ParseResult{}

	if p.currentToken().Kind != lexer.WHILE {
		errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, "Expected 'while'")
		return res.Failure(&errorNya)
	}

	res.Register_Advancement()
	p.advance()

	kondisi := res.Register(p.expr())
	if res.Error != nil {
		return res
	}

	if p.currentToken().Kind != lexer.DO {
		errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, "Expected 'do'")
		return res.Failure(&errorNya)
	}

	res.Register_Advancement()
	p.advance()

	if p.currentToken().Kind == lexer.NEWLINE {
		res.Register_Advancement()
		p.advance()

		body := res.Register(p.statements())
		if res.Error != nil {
			return res
		}

		if p.currentToken().Kind != lexer.END && p.currentToken().Kind != lexer.ENDWHILE {
			errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, "Expected 'end' or 'endwhile'")
			return res.Failure(&errorNya)
		}

		res.Register_Advancement()
		p.advance()

		return res.Success(common.WhileNode{
			KondisiNode:      kondisi,
			BodyNode:         body,
			ShouldReturnNull: true,
			Pos_Start:        kondisi.GetPosStart(),
			Pos_end:          body.GetPosEnd(),
		})
	}

	IsiNode := res.Register(p.statement())
	if res.Error != nil {
		return res
	}

	return res.Success(common.WhileNode{
		KondisiNode:      kondisi,
		BodyNode:         IsiNode,
		ShouldReturnNull: false,
		Pos_Start:        kondisi.GetPosStart(),
		Pos_end:          IsiNode.GetPosEnd(),
	})
}

func (p *parser) function_def() common.Expr {
	res := &common.ParseResult{}

	if p.currentToken().Kind != lexer.FUNCTION {
		errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, "Expected 'function'")
		return res.Failure(&errorNya)
	}

	res.Register_Advancement()
	p.advance()

	var VarNameTok *lexer.Token
	var errorNya common.Error
	if p.currentToken().Kind == lexer.IDENTIFIER {
		tok := p.currentToken()
		VarNameTok = &tok

		if p.currentToken().Kind != lexer.OPEN_PAREN {
			errorNya = common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, "Expected '('")
		}
	} else {
		VarNameTok = nil

		if p.currentToken().Kind != lexer.OPEN_PAREN {
			errorNya = common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, "Expected identifier or '('")
		}
	}

	if p.currentToken().Kind != lexer.OPEN_PAREN {
		return res.Failure(&errorNya)
	}

	res.Register_Advancement()
	p.advance()

	ArgNameToks := make([]lexer.Token, 0)

	if p.currentToken().Kind == lexer.IDENTIFIER {
		ArgNameToks = append(ArgNameToks, p.currentToken())
		res.Register_Advancement()
		p.advance()

		for p.currentToken().Kind == lexer.COMMA {
			res.Register_Advancement()
			p.advance()

			if p.currentToken().Kind != lexer.IDENTIFIER {
				errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, "Expected identifier")
				return res.Failure(&errorNya)
			}

			ArgNameToks = append(ArgNameToks, p.currentToken())
			res.Register_Advancement()
			p.advance()
		}

		if p.currentToken().Kind != lexer.CLOSE_PAREN {
			errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, "Expected ')' or ','")
			return res.Failure(&errorNya)
		}
	} else {
		if p.currentToken().Kind != lexer.CLOSE_PAREN {
			errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, "Expected ')' or identifier")
			return res.Failure(&errorNya)
		}
	}

	res.Register_Advancement()
	p.advance()

	if p.currentToken().Kind == lexer.RIGHT_ARROW {
		res.Register_Advancement()
		p.advance()

		nodeToReturn := res.Register(p.expr())
		if res.Error != nil {
			return res
		}

		funcNode := common.FuncNode{
			VarNameTok:       VarNameTok,
			ArgNameToks:      ArgNameToks,
			BodyNode:         nodeToReturn,
			ShouldAutoReturn: true,
			Pos_Start:        VarNameTok.Pos_Start,
			Pos_end:          nodeToReturn.GetPosEnd(),
		}
		funcNode.Pos_Start = funcNode.GetPosStart()
		funcNode.Pos_end = funcNode.GetPosEnd()

		return res.Success(funcNode)
	}

	if p.currentToken().Kind != lexer.NEWLINE {
		errorNya := common.InvalidSyntax(errorNya.PosStart, errorNya.PosEnd, "Expected '->' or new line")
		return res.Failure(&errorNya)
	}

	res.Register_Advancement()
	p.advance()

	body := res.Register(p.statements())
	if res.Error != nil {
		return res
	}

	if p.currentToken().Kind != lexer.END && p.currentToken().Kind != lexer.ENDIF {
		errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, "Expected 'end' or 'endif'")
		return res.Failure(&errorNya)
	}

	res.Register_Advancement()
	p.advance()

	return res.Success(common.FuncNode{
		VarNameTok:       VarNameTok,
		ArgNameToks:      ArgNameToks,
		BodyNode:         body,
		ShouldAutoReturn: false,
		Pos_Start:        VarNameTok.Pos_Start,
		Pos_end:          body.GetPosEnd(),
	})
}

func (p *parser) atom() common.Expr {
	res := &common.ParseResult{}
	tok := p.currentToken()

	switch tok.Kind {
	case lexer.NUMBER:
		res.Register_Advancement()
		p.advance()

		return res.Success(common.NumberNode{
			Token:     tok,
			Pos_Start: tok.Pos_Start,
			Pos_End:   tok.Pos_End,
		})
	case lexer.STRING:
		res.Register_Advancement()
		p.advance()

		return res.Success(common.StringNode{
			Token:     tok,
			Pos_Start: tok.Pos_Start,
			Pos_End:   tok.Pos_End,
		})
	case lexer.IDENTIFIER:
		res.Register_Advancement()
		p.advance()

		return res.Success(common.VarAccessNode{VarNameTok: tok, Pos_Start: tok.Pos_Start, Pos_end: tok.Pos_End})
	case lexer.OPEN_PAREN:
		res.Register_Advancement()
		p.advance()

		expr := res.Register(p.expr())
		if p.currentToken().Kind == lexer.CLOSE_PAREN {
			res.Register_Advancement()
			p.advance()

			return res.Success(expr)
		}

		errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, "Expected ')'")
		return res.Failure(&errorNya)
	case lexer.OPEN_BRACKET:
		list_expr := res.Register(p.list_expr())
		if res.Error != nil {
			return res
		}

		return res.Success(list_expr)
	case lexer.IF:
		if_expr := res.Register(p.if_expr())
		if res.Error != nil {
			return res
		}

		return res.Success(if_expr)
	case lexer.FOR:
		for_expr := res.Register(p.for_expr())
		if res.Error != nil {
			return res
		}

		return res.Success(for_expr)
	case lexer.WHILE:
		while_expr := res.Register(p.while_expr())
		if res.Error != nil {
			return res
		}

		return res.Success(while_expr)
	case lexer.FUNCTION:
		funcDef := p.function_def()
		function_def := res.Register(funcDef)

		if res.Error != nil {
			return res
		}

		return res.Success(function_def)
	}

	errorNya := common.InvalidSyntax(*tok.Pos_Start, *tok.Pos_End, "Expected identifier, int, float, '+', '-', '[', '(', 'if', 'for', 'while', 'function'")
	return res.Failure(&errorNya)
}

func (p *parser) power() common.Expr {
	return p.bin_op(p.call, []lexer.TokenKind{lexer.POWER}, p.factor)
}

func (p *parser) call() common.Expr {
	res := &common.ParseResult{}
	atom := res.Register(p.atom())
	if res.Error != nil {
		return res
	}

	if p.currentToken().Kind == lexer.OPEN_PAREN {
		res.Register_Advancement()
		p.advance()

		ArgNodes := make([]common.Expr, 0)

		if p.currentToken().Kind == lexer.CLOSE_PAREN {
			res.Register_Advancement()
			p.advance()
		} else {
			ArgNodes = append(ArgNodes, res.Register(p.expr()))
			if res.Error != nil {
				errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, "Expected, ')', 'var', 'if', 'for', 'while', 'fun', number, identifier")
				return res.Failure(&errorNya)
			}

			for p.currentToken().Kind == lexer.COMMA {
				res.Register_Advancement()
				p.advance()

				ArgNodes = append(ArgNodes, res.Register(p.expr()))
				if res.Error != nil {
					return res
				}
			}

			if p.currentToken().Kind != lexer.CLOSE_PAREN {
				errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, "Expected, ')' or ','")
				return res.Failure(&errorNya)
			}

			res.Register_Advancement()
			p.advance()
		}

		callNodes := common.CallNode{
			NodeToCall: atom,
			ArgNodes:   ArgNodes,
		}
		callNodes.Pos_Start = callNodes.GetPosStart()
		callNodes.Pos_end = callNodes.GetPosEnd()
		return res.Success(callNodes)
	}

	return res.Success(atom)
}

func (p *parser) factor() common.Expr {
	res := &common.ParseResult{Error: nil, Node: common.NullNode{}}
	tok := p.currentToken()

	if tok.Kind == lexer.PLUS || tok.Kind == lexer.DASH {
		res.Register_Advancement()
		p.advance()
		factor := res.Register(p.factor())

		if res.Error != nil {
			return res
		}

		return res.Success(common.UnaryOpNode{
			Operator:  tok,
			Node:      factor,
			Pos_Start: tok.Pos_Start,
			Pos_End:   factor.GetPosEnd(),
		})
	}

	return p.power()
}

func (p *parser) term() common.Expr {
	return p.bin_op(p.factor, []lexer.TokenKind{lexer.STAR, lexer.SLASH}, p.factor)
}

func (p *parser) arith_expr() common.Expr {
	return p.bin_op(p.term, []lexer.TokenKind{lexer.PLUS, lexer.DASH}, p.term)
}

func (p *parser) comp_expr() common.Expr {
	res := &common.ParseResult{}

	if p.currentToken().Kind == lexer.NOT {
		op_tok := p.currentToken()

		res.Register_Advancement()
		p.advance()

		node := res.Register(p.comp_expr())
		if res.Error != nil {
			return res
		}

		return res.Success(common.UnaryOpNode{
			Operator: op_tok,
			Node:     node,
		})
	}

	node := res.Register(p.bin_op(p.arith_expr, []lexer.TokenKind{lexer.EQUALS, lexer.NOT_EQUALS, lexer.LESS, lexer.LESS_EQUALS, lexer.GREATER, lexer.GREATER_EQUALS}, p.arith_expr))
	if res.Error != nil {
		errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, "Expected int, float, identifier, NOT, '+', '-', or '(', '[")
		return res.Failure(&errorNya)
	}

	return res.Success(node)
}

func (p *parser) statements() common.Expr {
	res := &common.ParseResult{}
	statements := make([]common.Expr, 0)
	pos_start := p.currentToken().Pos_Start.Copy()

	for p.currentToken().Kind == lexer.NEWLINE {
		res.Register_Advancement()
		p.advance()
	}

	statement := res.Register(p.statement())
	if res.Error != nil {
		return res
	}
	statements = append(statements, statement)

	moreStatement := true

	for {
		NewLineCount := 0
		for p.currentToken().Kind == lexer.NEWLINE {
			res.Register_Advancement()
			p.advance()
			NewLineCount++
		}

		if NewLineCount == 0 {
			moreStatement = false
		}

		if !moreStatement {
			break
		}

		statement := res.Try_register(p.statement())
		if statement == nil {
			p.reverse(res.ToReverseCount)
			moreStatement = false
			continue
		}

		statements = append(statements, statement)
	}

	return res.Success(common.ListNode{
		ElementNode: statements,
		Pos_Start:   pos_start,
		Pos_End:     p.currentToken().Pos_End.Copy(),
	})
}

func (p *parser) statement() common.Expr {
	res := &common.ParseResult{}
	pos_Start := p.currentToken().Pos_Start.Copy()

	if p.currentToken().Kind == lexer.RETURN {
		res.Register_Advancement()
		p.advance()

		expr := res.Try_register(p.expr())

		if expr == nil {
			p.reverse(res.ToReverseCount)
		}

		return res.Success(common.ReturnNode{NodeToReturn: expr, Pos_Start: pos_Start, Pos_end: p.currentToken().Pos_End.Copy()})
	}

	if p.currentToken().Kind == lexer.CONTINUE {
		res.Register_Advancement()
		p.advance()

		return res.Success(common.ContinueNode{Pos_Start: pos_Start, Pos_end: p.currentToken().Pos_Start.Copy()})
	}

	if p.currentToken().Kind == lexer.BREAK {
		res.Register_Advancement()
		p.advance()

		return res.Success(common.BreakNode{Pos_Start: pos_Start, Pos_end: p.currentToken().Pos_Start.Copy()})
	}

	expr := res.Register(p.expr())
	if res.Error != nil {
		errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, "Expected 'continue', 'break', 'return' 'var', 'for', 'while', 'function', int, float, identifier, '+', '-', '(', '['")
		return res.Failure(&errorNya)
	}

	return res.Success(expr)
}

func (p *parser) expr() common.Expr {
	res := &common.ParseResult{}
	if p.currentToken().Matches(lexer.VAR) {
		res.Register_Advancement()
		p.advance()

		if p.currentToken().Kind != lexer.IDENTIFIER {
			errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, "Expected Identifier")
			return res.Failure(&errorNya)
		}

		var_name := p.currentToken()
		res.Register_Advancement()
		p.advance()

		if p.currentToken().Kind != lexer.ASSIGNMENT {
			errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, "Expected '='")
			return res.Failure(&errorNya)
		}

		res.Register_Advancement()
		p.advance()

		expr := res.Register(p.expr())

		if res.Error != nil {
			return res
		}

		return res.Success(common.VarAssignNode{VarName: var_name, ValueNode: expr, Pos_Start: var_name.Pos_Start, Pos_end: expr.GetPosEnd()})
	}

	node := res.Register(p.bin_op(p.comp_expr, []lexer.TokenKind{lexer.AND, lexer.OR}, p.factor))

	if res.Error != nil {
		errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, "Expected 'var', 'for', 'while', 'function', int, float, identifier, '+', '-', '(', '[")
		return res.Failure(&errorNya)
	}

	return res.Success(node)
}

func (p *parser) bin_op(fungsi_a func() common.Expr, ops []lexer.TokenKind, fungsi_b func() common.Expr) common.Expr {
	res := &common.ParseResult{Error: nil, Node: common.NullNode{}}
	left := res.Register(fungsi_a())

	if res.Error != nil {
		return res
	}

	for slices.Contains(ops, p.currentToken().Kind) {
		op_tok := p.currentToken()
		res.Register_Advancement()
		p.advance()

		right := res.Register(fungsi_b())

		if res.Error != nil {
			return res
		}

		left = common.BinOpNode{
			Operator:  op_tok,
			Left:      left,
			Right:     right,
			Pos_Start: op_tok.Pos_Start,
			Pos_End:   right.GetPosEnd(),
		}
	}

	return res.Success(left)
}
