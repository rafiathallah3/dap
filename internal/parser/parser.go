package parser

import (
	"dap/internal/common"
	"dap/internal/lexer"
	"dap/tools"
	"fmt"
	"slices"
)

type parser struct {
	tokens          []lexer.Token
	hasEndProgram   bool
	tok_index       int
	apakahSatuBaris bool
}

func CreateParser(tokens []lexer.Token, ApakahSatuBaris bool) *parser {
	p := &parser{
		tokens:          tokens,
		tok_index:       -1,
		apakahSatuBaris: ApakahSatuBaris,
	}

	p.advance()
	return p
}

func (p *parser) Parse(programName *string) common.Expr {
	res := &common.ParseResult{}

	if p.currentToken().Kind == lexer.NEWLINE {
		res.Register_Advancement()
		p.advance()
	}

	res.Register(p.DapatinProgram())
	if res.Error != nil && !p.apakahSatuBaris {
		return res
	}

	hasil := p.statements().(*common.ParseResult)
	if hasil.Error == nil {
		if !p.hasEndProgram && p.currentToken().Kind != lexer.ENDPROGRAM && !p.apakahSatuBaris {
			errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, fmt.Sprintf("Expected 'endprogram' got %s", p.currentToken().Value))
			return hasil.Failure(&errorNya)
		}

		if p.currentToken().Kind == lexer.ENDPROGRAM {
			hasil.Register_Advancement()
			p.advance()
		}

		for p.currentToken().Kind == lexer.NEWLINE {
			hasil.Register_Advancement()
			p.advance()
		}

		if p.currentToken().Kind != lexer.EOF {
			errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, fmt.Sprintf("Expected end of file, got %s", p.currentToken().Value))
			return hasil.Failure(&errorNya)
		}
	}

	return hasil
}

func (p *parser) DapatinProgram() common.Expr {
	res := &common.ParseResult{}

	for p.currentToken().Kind != lexer.PROGRAM {
		errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, "Program name required")
		return res.Failure(&errorNya)
	}

	res.Register_Advancement()
	p.advance()

	for p.currentToken().Kind != lexer.IDENTIFIER {
		errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, "Program name must be an identifier")
		return res.Failure(&errorNya)
	}

	res.ProgramName = p.currentToken().Value

	res.Register_Advancement()
	p.advance()

	return res
}

func (p *parser) currentToken() lexer.Token {
	return p.tokens[p.tok_index]
}

func (p *parser) advance() lexer.Token {
	p.tok_index++
	return p.currentToken()
}

func (p *parser) reverse(amount int) lexer.Token {
	p.tok_index -= amount
	return p.currentToken()
}

func (p *parser) parse_type() common.Expr {
	res := &common.ParseResult{}

	if p.currentToken().Kind == lexer.ARRAY {
		posStart := p.currentToken().Pos_Start.Copy()
		res.Register_Advancement()
		p.advance()

		if p.currentToken().Kind != lexer.OPEN_BRACKET {
			errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, "Expected '['")
			return res.Failure(&errorNya)
		}

		res.Register_Advancement()
		p.advance()

		startExpr := res.Register(p.expr())
		if res.Error != nil {
			return res
		}

		if p.currentToken().Kind != lexer.DOT_DOT {
			errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, "Expected '..'")
			return res.Failure(&errorNya)
		}

		res.Register_Advancement()
		p.advance()

		endExpr := res.Register(p.expr())
		if res.Error != nil {
			return res
		}

		if p.currentToken().Kind != lexer.CLOSE_BRACKET {
			errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, "Expected ']'")
			return res.Failure(&errorNya)
		}

		res.Register_Advancement()
		p.advance()

		if p.currentToken().Kind != lexer.OF {
			errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, "Expected 'of'")
			return res.Failure(&errorNya)
		}

		res.Register_Advancement()
		p.advance()

		ofType := res.Register(p.parse_type())
		if res.Error != nil {
			return res
		}

		return res.Success(common.ArrayTypeNode{
			StartNode: startExpr,
			EndNode:   endExpr,
			OfType:    ofType,
			Pos_Start: posStart,
			Pos_End:   ofType.GetPosEnd(),
		})
	}

	if p.currentToken().Kind == lexer.INTEGER || p.currentToken().Kind == lexer.REAL || p.currentToken().Kind == lexer.STRINGTYPE || p.currentToken().Kind == lexer.IDENTIFIER {
		tok := p.currentToken()
		res.Register_Advancement()
		p.advance()
		return res.Success(common.VarAccessNode{VarNameTok: tok, Pos_Start: tok.Pos_Start, Pos_end: tok.Pos_End})
	}

	errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, "Expected type (integer, real, string, or array)")
	return res.Failure(&errorNya)
}

func (p *parser) dictionary_expr() common.Expr {
	res := &common.ParseResult{}

	if p.currentToken().Kind == lexer.NEWLINE {
		res.Register_Advancement()
		p.advance()
	}

	posStart := p.currentToken().Pos_Start.Copy()
	IsiNode := make([]common.Expr, 0)

	for p.currentToken().Kind != lexer.ALGORITHM {
		ListVarNameToks := make([]lexer.Token, 0)
		if p.currentToken().Kind == lexer.IDENTIFIER {
			ListVarNameToks = append(ListVarNameToks, p.currentToken())
			res.Register_Advancement()
			p.advance()

			for p.currentToken().Kind == lexer.COMMA {
				res.Register_Advancement()
				p.advance()

				if p.currentToken().Kind != lexer.IDENTIFIER {
					errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, "Expected identifier")
					return res.Failure(&errorNya)
				}

				ListVarNameToks = append(ListVarNameToks, p.currentToken())
				res.Register_Advancement()
				p.advance()
			}

			if p.currentToken().Kind != lexer.COLON {
				errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, fmt.Sprintf("Expected ':', got %s", p.currentToken().Value))
				return res.Failure(&errorNya)
			}

			res.Register_Advancement()
			p.advance()

			tipeDataNode := res.Register(p.parse_type())
			if res.Error != nil {
				return res
			}

			for _, val := range ListVarNameToks {
				IsiNode = append(IsiNode, common.VarAssignNode{
					VarName:   val,
					ValueNode: tipeDataNode,
					Pos_Start: val.Pos_Start,
					Pos_end:   tipeDataNode.GetPosEnd(),
				})
			}
		} else if p.currentToken().Kind == lexer.CONST {
			res.Register_Advancement()
			p.advance()

			if p.currentToken().Kind != lexer.IDENTIFIER {
				errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, fmt.Sprintf("Expected identifier, got %s", p.currentToken().Value))
				return res.Failure(&errorNya)
			}

			varName := p.currentToken()
			res.Register_Advancement()
			p.advance()
			// fmt.Fprintf(os.Stderr, "EXPR CHECK: %s (%s)\n", lexer.TokenKindString(p.currentToken().Kind), p.currentToken().Value)
			if p.currentToken().Kind != lexer.ASSIGNMENT && p.currentToken().Kind != lexer.LEFT_ARROW {
				errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, fmt.Sprintf("Expected '=' or '<-', got %s", p.currentToken().Value))
				return res.Failure(&errorNya)
			}

			res.Register_Advancement()
			p.advance()

			expr := res.Register(p.expr())
			if res.Error != nil {
				return res
			}

			IsiNode = append(IsiNode, common.VarAssignNode{
				VarName:     varName,
				ValueNode:   expr,
				ApakahConst: true,
				Pos_Start:   varName.Pos_Start,
				Pos_end:     varName.Pos_End,
			})
		} else if p.currentToken().Kind == lexer.TYPE {
			res.Register_Advancement()
			p.advance()

			if p.currentToken().Kind != lexer.IDENTIFIER {
				errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, "Expected identifier after 'type'")
				return res.Failure(&errorNya)
			}
			typeName := p.currentToken()
			res.Register_Advancement()
			p.advance()

			if p.currentToken().Kind == lexer.COLON {
				// Type Alias: type Alias: TargetType
				res.Register_Advancement()
				p.advance()
				targetType := res.Register(p.parse_type())
				if res.Error != nil {
					return res
				}
				IsiNode = append(IsiNode, common.TypeAliasNode{
					AliasName:  typeName,
					TargetType: targetType,
					Pos_Start:  typeName.Pos_Start,
					Pos_End:    targetType.GetPosEnd(),
				})
			} else if p.currentToken().Kind == lexer.LESS {
				// Struct: type Name < field: type ... >
				res.Register_Advancement()
				p.advance()

				var fields []common.VarAssignNode
				for p.currentToken().Kind != lexer.GREATER {
					if p.currentToken().Kind == lexer.NEWLINE {
						res.Register_Advancement()
						p.advance()
						continue
					}

					// Parse field: type
					if p.currentToken().Kind != lexer.IDENTIFIER {
						errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, "Expected field name")
						return res.Failure(&errorNya)
					}
					fieldName := p.currentToken()
					res.Register_Advancement()
					p.advance()

					if p.currentToken().Kind != lexer.COLON {
						errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, "Expected ':' after field name")
						return res.Failure(&errorNya)
					}
					res.Register_Advancement()
					p.advance()

					fieldType := res.Register(p.parse_type())
					if res.Error != nil {
						return res
					}

					fields = append(fields, common.VarAssignNode{
						VarName:   fieldName,
						ValueNode: fieldType,
						Pos_Start: fieldName.Pos_Start,
						Pos_end:   fieldType.GetPosEnd(),
					})
				}
				posEnd := p.currentToken().Pos_End.Copy()
				res.Register_Advancement()
				p.advance()

				IsiNode = append(IsiNode, common.StructTypeNode{
					StructName: typeName,
					Fields:     fields,
					Pos_Start:  typeName.Pos_Start,
					Pos_End:    posEnd,
				})
			} else {
				errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, "Expected ':' or '<'")
				return res.Failure(&errorNya)
			}
		} else if p.currentToken().Kind == lexer.NEWLINE {
			res.Register_Advancement()
			p.advance()
		} else {
			errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, fmt.Sprintf("Expected identifier, got %s", p.currentToken().Value))
			return res.Failure(&errorNya)
		}

		if p.currentToken().Kind == lexer.NEWLINE {
			res.Register_Advancement()
			p.advance()
		}
	}

	res.Register_Advancement()
	p.advance()

	return res.Success(common.DictionaryNode{
		VariableDiBuat: IsiNode,
		Pos_Start:      posStart,
		Pos_end:        p.currentToken().Pos_End.Copy(),
	})
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

	if p.currentToken().Kind != lexer.ASSIGNMENT && p.currentToken().Kind != lexer.LEFT_ARROW {
		errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, "Expected '=' or '<-'")
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

func (p *parser) repeat_expr() common.Expr {
	res := &common.ParseResult{}

	if p.currentToken().Kind != lexer.REPEAT {
		errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.advance().Pos_End, "Expected 'repeat'")
		return res.Failure(&errorNya)
	}
	pos_start := p.currentToken().Pos_Start.Copy()

	res.Register_Advancement()
	p.advance()

	if p.currentToken().Kind == lexer.NEWLINE {
		res.Register_Advancement()
		p.advance()

		body := res.Register(p.statements())
		if res.Error != nil {
			return res
		}

		if p.currentToken().Kind != lexer.UNTIL {
			errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, "Expected 'until'")
			return res.Failure(&errorNya)
		}

		res.Register_Advancement()
		p.advance()

		kondisi := res.Register(p.expr())
		if res.Error != nil {
			return res
		}

		return res.Success(common.RepeatNode{
			KondisiNode:      kondisi,
			BodyNode:         body,
			ShouldReturnNull: true,
			Pos_Start:        pos_start,
			Pos_end:          kondisi.GetPosEnd(),
		})
	}

	IsiNode := res.Register(p.statement())
	if res.Error != nil {
		return res
	}

	if p.currentToken().Kind != lexer.END && p.currentToken().Kind != lexer.UNTIL {
		errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, "Expected 'until'")
		return res.Failure(&errorNya)
	}

	res.Register_Advancement()
	p.advance()

	kondisi := res.Register(p.expr())
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
	case lexer.REPEAT:
		repeat_expr := res.Register(p.repeat_expr())
		if res.Error != nil {
			return res
		}

		return res.Success(repeat_expr)
	case lexer.FUNCTION:
		function_def := res.Register(p.function_def())
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

	for p.currentToken().Kind == lexer.OPEN_PAREN || p.currentToken().Kind == lexer.OPEN_BRACKET || p.currentToken().Kind == lexer.DOT || tools.ApakahBuiltinFunction(atom.Print()) {
		if p.currentToken().Kind == lexer.OPEN_PAREN || tools.ApakahBuiltinFunction(atom.Print()) {
			var argNodes []common.Expr
			if p.currentToken().Kind == lexer.OPEN_PAREN {
				res.Register_Advancement()
				p.advance()

				if p.currentToken().Kind == lexer.CLOSE_PAREN {
					res.Register_Advancement()
					p.advance()
				} else {
					argNodes = append(argNodes, res.Register(p.expr()))
					if res.Error != nil {
						return res
					}

					for p.currentToken().Kind == lexer.COMMA {
						res.Register_Advancement()
						p.advance()

						argNodes = append(argNodes, res.Register(p.expr()))
						if res.Error != nil {
							return res
						}
					}

					if p.currentToken().Kind != lexer.CLOSE_PAREN {
						errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, "Expected ',' or ')'")
						return res.Failure(&errorNya)
					}

					res.Register_Advancement()
					p.advance()
				}
			} else if tools.ApakahBuiltinFunction(atom.Print()) {
				// Builtin function without parens (e.g., write M[0][0])
				// Try to parse expressions until end of line or keyword
				for {
					if p.currentToken().Kind == lexer.NEWLINE || p.currentToken().Kind == lexer.EOF || p.currentToken().Kind == lexer.ENDPROGRAM || p.currentToken().Kind == lexer.ENDWHILE || p.currentToken().Kind == lexer.ENDFOR || p.currentToken().Kind == lexer.ENDIF || p.currentToken().Kind == lexer.ELSE || p.currentToken().Kind == lexer.ELIF {
						break
					}
					argNodes = append(argNodes, res.Register(p.expr()))
					if res.Error != nil {
						return res
					}
					if p.currentToken().Kind == lexer.COMMA {
						res.Register_Advancement()
						p.advance()
					} else {
						break
					}
				}
			}

			atom = common.CallNode{
				NodeToCall: atom,
				ArgNodes:   argNodes,
				Pos_Start:  atom.GetPosStart(),
				Pos_end:    p.currentToken().Pos_End.Copy(),
			}
		} else if p.currentToken().Kind == lexer.OPEN_BRACKET {
			res.Register_Advancement()
			p.advance()

			indexExpr := res.Register(p.expr())
			if res.Error != nil {
				return res
			}

			if p.currentToken().Kind != lexer.CLOSE_BRACKET {
				errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, "Expected ']'")
				return res.Failure(&errorNya)
			}

			res.Register_Advancement()
			p.advance()

			atom = common.ArrayIndexNode{
				Left:      atom,
				Index:     indexExpr,
				Pos_Start: atom.GetPosStart(),
				Pos_End:   p.tokens[p.tok_index-1].Pos_End.Copy(),
			}
		} else if p.currentToken().Kind == lexer.DOT {
			res.Register_Advancement()
			p.advance()

			if p.currentToken().Kind != lexer.IDENTIFIER {
				errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, "Expected field name after '.'")
				return res.Failure(&errorNya)
			}
			memberTok := p.currentToken()
			res.Register_Advancement()
			p.advance()

			atom = common.MemberAccessNode{
				Object:    atom,
				MemberTok: memberTok,
				Pos_Start: atom.GetPosStart(),
				Pos_End:   memberTok.Pos_End.Copy(),
			}
		}
	}
	return res.Success(atom)
}

func (p *parser) factor() common.Expr {
	res := &common.ParseResult{}
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
	for moreStatement {
		NewLineCount := 0
		for p.currentToken().Kind == lexer.NEWLINE {
			res.Register_Advancement()
			p.advance()
			NewLineCount++
		}

		if NewLineCount == 0 && len(statements) > 0 {
			break
		}

		if p.currentToken().Kind == lexer.EOF || p.currentToken().Kind == lexer.ENDPROGRAM {
			break
		}

		// Use a local result for this statement to avoid reversing everything if it fails
		statementRes := &common.ParseResult{}
		statement := statementRes.Try_register(p.statement())
		res.Register_Advancement() 
		
		if statement == nil {
			p.reverse(statementRes.ToReverseCount)
			moreStatement = false
			continue
		}
		res.AdvanceCount += statementRes.AdvanceCount
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

	switch p.currentToken().Kind {
	case lexer.RETURN:
		res.Register_Advancement()
		p.advance()

		expr := res.Try_register(p.expr())

		if expr == nil {
			p.reverse(res.ToReverseCount)
		}

		return res.Success(common.ReturnNode{NodeToReturn: expr, Pos_Start: pos_Start, Pos_end: p.currentToken().Pos_End.Copy()})
	case lexer.CONTINUE:
		res.Register_Advancement()
		p.advance()

		return res.Success(common.ContinueNode{Pos_Start: pos_Start, Pos_end: p.currentToken().Pos_Start.Copy()})
	case lexer.BREAK:
		res.Register_Advancement()
		p.advance()

		return res.Success(common.BreakNode{Pos_Start: pos_Start, Pos_end: p.currentToken().Pos_Start.Copy()})
	case lexer.DICTIONARY:
		res.Register_Advancement()
		p.advance()

		expr := res.Register(p.dictionary_expr())
		if res.Error != nil {
			return res
		}

		return res.Success(expr)
	case lexer.ENDPROGRAM:
		res.Register_Advancement()
		p.advance()

		p.hasEndProgram = true
		return res.Success(nil)
	case lexer.ALGORITHM:
		res.Register_Advancement()
		p.advance()

		return res.Success(common.NullNode{})
	}

	expr := res.Register(p.expr())
	if res.Error != nil {
		errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, fmt.Sprintf("Got, %s | Expected 'continue', 'break', 'return' 'var', 'for', 'while', 'function', int, float, identifier, '+', '-', '(', '['", p.currentToken().Value))
		return res.Failure(&errorNya)
	}

	return res.Success(expr)
}

func (p *parser) expr() common.Expr {
	res := &common.ParseResult{}

	apakahLeftArrow := false
	if p.currentToken().Kind == lexer.IDENTIFIER {
		var lookahead lexer.Token
		lookaheadCount := 0
		bracketLevel := 0

		for {
			lookahead = p.tokens[p.tok_index+lookaheadCount]
			if lookahead.Kind == lexer.OPEN_BRACKET {
				bracketLevel++
			} else if lookahead.Kind == lexer.CLOSE_BRACKET {
				bracketLevel--
			}

			if bracketLevel < 0 {
				break
			}

			if bracketLevel == 0 && (lookahead.Kind == lexer.ASSIGNMENT || lookahead.Kind == lexer.LEFT_ARROW) {
				apakahLeftArrow = true
				break
			}
			if lookahead.Kind != lexer.IDENTIFIER && lookahead.Kind != lexer.OPEN_BRACKET && lookahead.Kind != lexer.CLOSE_BRACKET && lookahead.Kind != lexer.NUMBER && lookahead.Kind != lexer.DOT {
				break
			}
			if lookahead.Kind == lexer.NEWLINE {
				break
			}
			lookaheadCount++
			if p.tok_index+lookaheadCount >= len(p.tokens) {
				break
			}
		}

		if apakahLeftArrow {
			lhs := res.Register(p.call())
			if res.Error != nil {
				return res
			}

			if p.currentToken().Kind != lexer.ASSIGNMENT && p.currentToken().Kind != lexer.LEFT_ARROW {
				errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, "Expected '<-' or '='")
				return res.Failure(&errorNya)
			}

			res.Register_Advancement()
			p.advance()

			expr := res.Register(p.expr())
			if res.Error != nil {
				return res
			}

			switch lhs := lhs.(type) {
			case common.VarAccessNode:
				return res.Success(common.VarAssignNode{
					VarName:   lhs.VarNameTok,
					ValueNode: expr,
					Pos_Start: lhs.Pos_Start,
					Pos_end:   expr.GetPosEnd(),
				})
			case common.ArrayIndexNode:
				return res.Success(common.ArrayAssignNode{
					ArrayAccess: lhs,
					ValueNode:   expr,
					Pos_Start:   lhs.GetPosStart(),
					Pos_End:     expr.GetPosEnd(),
				})
			case common.MemberAccessNode:
				return res.Success(common.MemberAssignNode{
					MemberAccess: lhs,
					ValueNode:    expr,
					Pos_Start:    lhs.GetPosStart(),
					Pos_End:      expr.GetPosEnd(),
				})
			default:
				errorNya := common.InvalidSyntax(*lhs.GetPosStart(), *lhs.GetPosEnd(), "Illegal assignment target")
				return res.Failure(&errorNya)
			}
		}
	}

	node := res.Register(p.bin_op(p.comp_expr, []lexer.TokenKind{lexer.AND, lexer.OR}, p.comp_expr))
	if res.Error != nil {
		errorNya := common.InvalidSyntax(*p.currentToken().Pos_Start, *p.currentToken().Pos_End, "Expected 'var', 'for', 'while', 'function', int, float, identifier, '+', '-', '(', '[")
		return res.Failure(&errorNya)
	}

	return res.Success(node)
}

func (p *parser) bin_op(fungsi_a func() common.Expr, ops []lexer.TokenKind, fungsi_b func() common.Expr) common.Expr {
	res := &common.ParseResult{}
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
