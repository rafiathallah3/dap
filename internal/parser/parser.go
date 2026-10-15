package parser

import (
	"dap/internal/common"
	"dap/internal/lexer"
	"dap/tools"
	"fmt"
	"slices"
)

type parser struct {
	tokens        []lexer.Token
	hasEndProgram bool
	tokIndex      int
	isSingleLine  bool
}

func safePos(p *tools.Position) tools.Position {
	if p == nil {
		return tools.Position{}
	}
	return *p
}

// NewParser buat bikin instance parser baru dari stream token yang dikasih.
func NewParser(tokens []lexer.Token, isSingleLine bool) *parser {
	p := &parser{
		tokens:       tokens,
		tokIndex:     -1,
		isSingleLine: isSingleLine,
	}

	p.advance()
	return p
}

// CreateParser alias buat backward compatibility.
func CreateParser(tokens []lexer.Token, isSingleLine bool) *parser {
	return NewParser(tokens, isSingleLine)
}

// Parse buat nge-parse kumpulan token jadi struktur AST.
func (p *parser) Parse(programName *string) common.Expr {
	res := &common.ParseResult{}

	if p.currentToken().Kind == lexer.NEWLINE {
		res.RegisterAdvancement()
		p.advance()
	}

	res.Register(p.parseProgramHeader())
	if res.Error != nil && !p.isSingleLine {
		return res
	}

	stmtsRes := p.statements()
	hasil, ok := stmtsRes.(*common.ParseResult)
	if !ok {
		hasil = &common.ParseResult{Node: stmtsRes}
	}

	if hasil.Error == nil {
		if !p.hasEndProgram && p.currentToken().Kind != lexer.ENDPROGRAM && !p.isSingleLine {
			errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.currentToken().PosEnd), fmt.Sprintf("Expected 'endprogram' got %s", p.currentToken().Value))
			return hasil.Failure(&errorNya)
		}

		if p.currentToken().Kind == lexer.ENDPROGRAM {
			hasil.RegisterAdvancement()
			p.advance()
		}

		for p.currentToken().Kind == lexer.NEWLINE {
			hasil.RegisterAdvancement()
			p.advance()
		}

		if p.currentToken().Kind != lexer.EOF {
			errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.currentToken().PosEnd), fmt.Sprintf("Expected end of file, got %s", p.currentToken().Value))
			return hasil.Failure(&errorNya)
		}
	}

	return hasil
}

func (p *parser) parseProgramHeader() common.Expr {
	res := &common.ParseResult{}

	for p.currentToken().Kind != lexer.PROGRAM {
		errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.currentToken().PosEnd), "Program name required")
		return res.Failure(&errorNya)
	}

	res.RegisterAdvancement()
	p.advance()

	for p.currentToken().Kind != lexer.IDENTIFIER {
		errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.currentToken().PosEnd), "Program name must be an identifier")
		return res.Failure(&errorNya)
	}

	res.ProgramName = p.currentToken().Value

	res.RegisterAdvancement()
	p.advance()

	return res
}

func (p *parser) currentToken() lexer.Token {
	if p.tokIndex < 0 || p.tokIndex >= len(p.tokens) {
		return lexer.Token{Kind: lexer.EOF, Value: "EOF"}
	}
	return p.tokens[p.tokIndex]
}

func (p *parser) advance() lexer.Token {
	p.tokIndex++
	return p.currentToken()
}

func (p *parser) reverse(amount int) lexer.Token {
	p.tokIndex -= amount
	return p.currentToken()
}

func (p *parser) parseType() common.Expr {
	res := &common.ParseResult{}

	if p.currentToken().Kind == lexer.ARRAY {
		posStart := p.currentToken().PosStart.Copy()
		res.RegisterAdvancement()
		p.advance()

		if p.currentToken().Kind != lexer.OPEN_BRACKET {
			errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.currentToken().PosEnd), "Expected '['")
			return res.Failure(&errorNya)
		}

		res.RegisterAdvancement()
		p.advance()

		startExpr := res.Register(p.expr())
		if res.Error != nil {
			return res
		}

		if p.currentToken().Kind != lexer.DOT_DOT {
			errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.currentToken().PosEnd), "Expected '..'")
			return res.Failure(&errorNya)
		}

		res.RegisterAdvancement()
		p.advance()

		endExpr := res.Register(p.expr())
		if res.Error != nil {
			return res
		}

		if p.currentToken().Kind != lexer.CLOSE_BRACKET {
			errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.currentToken().PosEnd), "Expected ']'")
			return res.Failure(&errorNya)
		}

		res.RegisterAdvancement()
		p.advance()

		if p.currentToken().Kind != lexer.OF {
			errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.currentToken().PosEnd), "Expected 'of'")
			return res.Failure(&errorNya)
		}

		res.RegisterAdvancement()
		p.advance()

		ofType := res.Register(p.parseType())
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
		res.RegisterAdvancement()
		p.advance()
		return res.Success(common.VarAccessNode{VarNameTok: tok, Pos_Start: tok.PosStart, Pos_end: tok.PosEnd})
	}

	errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.currentToken().PosEnd), "Expected type (integer, real, string, or array)")
	return res.Failure(&errorNya)
}

func (p *parser) parseDictionary() common.Expr {
	res := &common.ParseResult{}

	if p.currentToken().Kind == lexer.NEWLINE {
		res.RegisterAdvancement()
		p.advance()
	}

	posStart := p.currentToken().PosStart.Copy()
	isiNode := make([]common.Expr, 0)

	for p.currentToken().Kind != lexer.ALGORITHM {
		listVarNameToks := make([]lexer.Token, 0)
		if p.currentToken().Kind == lexer.IDENTIFIER {
			listVarNameToks = append(listVarNameToks, p.currentToken())
			res.RegisterAdvancement()
			p.advance()

			for p.currentToken().Kind == lexer.COMMA {
				res.RegisterAdvancement()
				p.advance()

				if p.currentToken().Kind != lexer.IDENTIFIER {
					errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), "Expected identifier")
					return res.Failure(&errorNya)
				}

				listVarNameToks = append(listVarNameToks, p.currentToken())
				res.RegisterAdvancement()
				p.advance()
			}

			if p.currentToken().Kind != lexer.COLON {
				errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), fmt.Sprintf("Expected ':', got %s", p.currentToken().Value))
				return res.Failure(&errorNya)
			}

			res.RegisterAdvancement()
			p.advance()

			tipeDataNode := res.Register(p.parseType())
			if res.Error != nil {
				return res
			}

			for _, val := range listVarNameToks {
				isiNode = append(isiNode, common.VarAssignNode{
					VarName:         val,
					ValueNode:       tipeDataNode,
					ApakahDeklarasi: true,
					Pos_Start:       val.PosStart,
					Pos_end:         tipeDataNode.GetPosEnd(),
				})
			}
		} else if p.currentToken().Kind == lexer.CONST {
			res.RegisterAdvancement()
			p.advance()

			if p.currentToken().Kind != lexer.IDENTIFIER {
				errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), fmt.Sprintf("Expected identifier, got %s", p.currentToken().Value))
				return res.Failure(&errorNya)
			}

			varName := p.currentToken()
			res.RegisterAdvancement()
			p.advance()

			var constType common.Expr
			if p.currentToken().Kind == lexer.COLON {
				res.RegisterAdvancement()
				p.advance()

				constType = res.Register(p.parseType())
				if res.Error != nil {
					return res
				}
			}

			if p.currentToken().Kind != lexer.ASSIGNMENT && p.currentToken().Kind != lexer.LEFT_ARROW {
				errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), fmt.Sprintf("Expected '=' or '<-', got %s", p.currentToken().Value))
				return res.Failure(&errorNya)
			}

			res.RegisterAdvancement()
			p.advance()

			expr := res.Register(p.expr())
			if res.Error != nil {
				return res
			}

			isiNode = append(isiNode, common.VarAssignNode{
				VarName:         varName,
				ValueNode:       expr,
				TypeNode:        constType,
				ApakahConst:     true,
				ApakahDeklarasi: constType != nil,
				Pos_Start:       varName.PosStart,
				Pos_end:         varName.PosEnd,
			})
		} else if p.currentToken().Kind == lexer.TYPE {
			res.RegisterAdvancement()
			p.advance()

			if p.currentToken().Kind != lexer.IDENTIFIER {
				errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), "Expected identifier after 'type'")
				return res.Failure(&errorNya)
			}
			typeName := p.currentToken()
			res.RegisterAdvancement()
			p.advance()

			if p.currentToken().Kind == lexer.COLON {
				res.RegisterAdvancement()
				p.advance()
				targetType := res.Register(p.parseType())
				if res.Error != nil {
					return res
				}
				isiNode = append(isiNode, common.TypeAliasNode{
					AliasName:  typeName,
					TargetType: targetType,
					Pos_Start:  typeName.PosStart,
					Pos_End:    targetType.GetPosEnd(),
				})
			} else if p.currentToken().Kind == lexer.LESS {
				res.RegisterAdvancement()
				p.advance()

				var fields []common.VarAssignNode
				for p.currentToken().Kind != lexer.GREATER {
					if p.currentToken().Kind == lexer.NEWLINE {
						res.RegisterAdvancement()
						p.advance()
						continue
					}

					if p.currentToken().Kind != lexer.IDENTIFIER {
						errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), "Expected field name")
						return res.Failure(&errorNya)
					}
					fieldName := p.currentToken()
					res.RegisterAdvancement()
					p.advance()

					if p.currentToken().Kind != lexer.COLON {
						errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), "Expected ':' after field name")
						return res.Failure(&errorNya)
					}
					res.RegisterAdvancement()
					p.advance()

					fieldType := res.Register(p.parseType())
					if res.Error != nil {
						return res
					}

					fields = append(fields, common.VarAssignNode{
						VarName:   fieldName,
						ValueNode: fieldType,
						Pos_Start: fieldName.PosStart,
						Pos_end:   fieldType.GetPosEnd(),
					})
				}
				posEnd := p.currentToken().PosEnd.Copy()
				res.RegisterAdvancement()
				p.advance()

				isiNode = append(isiNode, common.StructTypeNode{
					StructName: typeName,
					Fields:     fields,
					Pos_Start:  typeName.PosStart,
					Pos_End:    posEnd,
				})
			} else {
				errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), "Expected ':' or '<'")
				return res.Failure(&errorNya)
			}
		} else if p.currentToken().Kind == lexer.FUNCTION {
			funcNode := res.Register(p.parseFunctionDef())
			if res.Error != nil {
				return res
			}
			isiNode = append(isiNode, funcNode)
		} else if p.currentToken().Kind == lexer.NEWLINE {
			res.RegisterAdvancement()
			p.advance()
		} else {
			errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), fmt.Sprintf("Expected identifier, got %s", p.currentToken().Value))
			return res.Failure(&errorNya)
		}

		if p.currentToken().Kind == lexer.NEWLINE {
			res.RegisterAdvancement()
			p.advance()
		}
	}

	res.RegisterAdvancement()
	p.advance()

	return res.Success(common.DictionaryNode{
		VariableDiBuat: isiNode,
		Pos_Start:      posStart,
		Pos_end:        p.currentToken().PosEnd.Copy(),
	})
}

func (p *parser) parseList() common.Expr {
	res := &common.ParseResult{}
	elementNodes := make([]common.Expr, 0)

	if p.currentToken().Kind != lexer.OPEN_BRACKET {
		errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.currentToken().PosEnd), "Expected '['")
		return res.Failure(&errorNya)
	}
	posStart := p.currentToken().PosStart.Copy()

	res.RegisterAdvancement()
	p.advance()

	if p.currentToken().Kind == lexer.CLOSE_BRACKET {
		res.RegisterAdvancement()
		p.advance()
	} else {
		elementNodes = append(elementNodes, res.Register(p.expr()))
		if res.Error != nil {
			errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.currentToken().PosEnd), "Expected ']', 'if', 'for', 'while', 'function', a number, or an identifier")
			return res.Failure(&errorNya)
		}

		for p.currentToken().Kind == lexer.COMMA {
			res.RegisterAdvancement()
			p.advance()

			elementNodes = append(elementNodes, res.Register(p.expr()))
			if res.Error != nil {
				return res
			}
		}

		if p.currentToken().Kind != lexer.CLOSE_BRACKET {
			errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.currentToken().PosEnd), "Expected ']' or ','")
			return res.Failure(&errorNya)
		}

		res.RegisterAdvancement()
		p.advance()
	}

	return res.Success(common.ListNode{
		ElementNode: elementNodes,
		Pos_Start:   posStart,
		Pos_End:     p.currentToken().PosEnd.Copy(),
	})
}

func (p *parser) parseIf() common.Expr {
	res := &common.ParseResult{}
	allCases := res.Register(p.parseIfCases(lexer.IF))
	if res.Error != nil {
		return res
	}

	ifNode, ok := allCases.(common.IfNode)
	if !ok {
		return res
	}

	return res.Success(common.IfNode{
		Cases:     ifNode.Cases,
		Else_case: ifNode.Else_case,
	})
}

func (p *parser) parseIfB() common.Expr {
	return p.parseIfCases(lexer.ELIF)
}

func (p *parser) parseIfC() common.Expr {
	res := &common.ParseResult{}
	var elseCase common.ElseCase

	if p.currentToken().Kind == lexer.ELSE {
		res.RegisterAdvancement()
		p.advance()

		if p.currentToken().Kind == lexer.NEWLINE {
			res.RegisterAdvancement()
			p.advance()

			statements := res.Register(p.statements())
			if res.Error != nil {
				return res
			}
			elseCase = common.ElseCase{
				Isi:              statements,
				ShouldReturnNull: true,
			}

			if p.currentToken().Kind == lexer.END || p.currentToken().Kind == lexer.ENDIF {
				res.RegisterAdvancement()
				p.advance()
			} else {
				errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), "expected 'end' or 'endif'")
				return res.Failure(&errorNya)
			}
		} else {
			expr := res.Register(p.statement())
			if res.Error != nil {
				return res
			}

			elseCase = common.ElseCase{
				Isi:              expr,
				ShouldReturnNull: false,
			}
		}
	}

	return res.Success(elseCase)
}

func (p *parser) parseIfBOrC() common.Expr {
	res := &common.ParseResult{}
	cases := make([]common.IfCase, 0)
	var elseCase common.ElseCase

	if p.currentToken().Kind == lexer.ELIF {
		allCases := res.Register(p.parseIfB())
		if res.Error != nil {
			return res
		}

		ifNode := allCases.(common.IfNode)
		cases = ifNode.Cases
		if ifNode.Else_case != nil {
			elseCase = *ifNode.Else_case
		}
	} else {
		resElseCase := res.Register(p.parseIfC())
		if res.Error != nil {
			return res
		}

		elseCase = resElseCase.(common.ElseCase)
	}

	return res.Success(common.IfNode{
		Cases:     cases,
		Else_case: &elseCase,
	})
}

func (p *parser) parseIfCases(caseKeyword lexer.TokenKind) common.Expr {
	res := &common.ParseResult{}
	var cases []common.IfCase
	var elseCase common.ElseCase

	if p.currentToken().Kind != caseKeyword {
		errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), fmt.Sprintf("Expected '%s'", lexer.TokenKindString(caseKeyword)))
		return res.Failure(&errorNya)
	}

	res.RegisterAdvancement()
	p.advance()

	condition := res.Register(p.expr())
	if res.Error != nil {
		return res
	}

	if p.currentToken().Kind != lexer.THEN {
		errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), "Expected 'then'")
		return res.Failure(&errorNya)
	}

	res.RegisterAdvancement()
	p.advance()

	if p.currentToken().Kind == lexer.NEWLINE {
		res.RegisterAdvancement()
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
			res.RegisterAdvancement()
			p.advance()
		} else {
			allCases := res.Register(p.parseIfBOrC())
			if res.Error != nil {
				return res
			}

			if ifNode, ok := allCases.(common.IfNode); ok {
				if ifNode.Else_case != nil {
					elseCase = *ifNode.Else_case
				}
				cases = append(cases, ifNode.Cases...)
			}
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

		allCases := res.Register(p.parseIfBOrC())
		if ifNode, ok := allCases.(common.IfNode); ok {
			if ifNode.Else_case != nil {
				elseCase = *ifNode.Else_case
			}
			cases = append(cases, ifNode.Cases...)
		}
	}

	return res.Success(common.IfNode{
		Cases:     cases,
		Else_case: &elseCase,
	})
}

func (p *parser) parseFor() common.Expr {
	res := &common.ParseResult{}

	if p.currentToken().Kind != lexer.FOR {
		errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), "Expected 'for'")
		return res.Failure(&errorNya)
	}

	res.RegisterAdvancement()
	p.advance()

	if p.currentToken().Kind != lexer.IDENTIFIER {
		errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), "Expected identifier")
		return res.Failure(&errorNya)
	}

	varName := p.currentToken()
	res.RegisterAdvancement()
	p.advance()

	if p.currentToken().Kind != lexer.ASSIGNMENT && p.currentToken().Kind != lexer.LEFT_ARROW {
		errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), "Expected '=' or '<-'")
		return res.Failure(&errorNya)
	}

	res.RegisterAdvancement()
	p.advance()

	startValue := res.Register(p.expr())
	if res.Error != nil {
		return res
	}

	isDown := false
	if p.currentToken().Kind != lexer.TO && p.currentToken().Kind != lexer.DOWNTO {
		errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), "Expected 'to' or 'downto'")
		return res.Failure(&errorNya)
	}
	if p.currentToken().Kind == lexer.DOWNTO {
		isDown = true
	}

	res.RegisterAdvancement()
	p.advance()

	endValue := res.Register(p.expr())
	if res.Error != nil {
		return res
	}

	var stepValue common.Expr = common.NullNode{}
	if p.currentToken().Kind == lexer.STEP {
		res.RegisterAdvancement()
		p.advance()

		stepValue = res.Register(p.expr())
		if res.Error != nil {
			return res
		}
	}

	if p.currentToken().Kind != lexer.DO {
		errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), "Expected 'do'")
		return res.Failure(&errorNya)
	}

	res.RegisterAdvancement()
	p.advance()

	if p.currentToken().Kind == lexer.NEWLINE {
		res.RegisterAdvancement()
		p.advance()

		body := res.Register(p.statements())
		if res.Error != nil {
			return res
		}

		if p.currentToken().Kind != lexer.END && p.currentToken().Kind != lexer.ENDFOR {
			errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), "Expected 'end' or 'endfor'")
			return res.Failure(&errorNya)
		}

		res.RegisterAdvancement()
		p.advance()

		return res.Success(common.ForNode{
			VarNameTok:       varName,
			StartValueNode:   startValue,
			EndValueNode:     endValue,
			StepValueNode:    stepValue,
			BodyNode:         body,
			ShouldReturnNull: true,
			IsDown:           isDown,
			Pos_Start:        varName.PosStart,
			Pos_end:          body.GetPosEnd(),
		})
	}

	isiValue := res.Register(p.statement())
	if res.Error != nil {
		return res
	}

	return res.Success(common.ForNode{
		VarNameTok:       varName,
		StartValueNode:   startValue,
		EndValueNode:     endValue,
		StepValueNode:    stepValue,
		BodyNode:         isiValue,
		ShouldReturnNull: false,
		IsDown:           isDown,
		Pos_Start:        varName.PosStart,
		Pos_end:          isiValue.GetPosEnd(),
	})
}

func (p *parser) parseWhile() common.Expr {
	res := &common.ParseResult{}

	if p.currentToken().Kind != lexer.WHILE {
		errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), "Expected 'while'")
		return res.Failure(&errorNya)
	}

	res.RegisterAdvancement()
	p.advance()

	kondisi := res.Register(p.expr())
	if res.Error != nil {
		return res
	}

	if p.currentToken().Kind != lexer.DO {
		errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), "Expected 'do'")
		return res.Failure(&errorNya)
	}

	res.RegisterAdvancement()
	p.advance()

	if p.currentToken().Kind == lexer.NEWLINE {
		res.RegisterAdvancement()
		p.advance()

		body := res.Register(p.statements())
		if res.Error != nil {
			return res
		}

		if p.currentToken().Kind != lexer.END && p.currentToken().Kind != lexer.ENDWHILE {
			errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), "Expected 'end' or 'endwhile'")
			return res.Failure(&errorNya)
		}

		res.RegisterAdvancement()
		p.advance()

		return res.Success(common.WhileNode{
			KondisiNode:      kondisi,
			BodyNode:         body,
			ShouldReturnNull: true,
			Pos_Start:        kondisi.GetPosStart(),
			Pos_end:          body.GetPosEnd(),
		})
	}

	isiNode := res.Register(p.statement())
	if res.Error != nil {
		return res
	}

	return res.Success(common.WhileNode{
		KondisiNode:      kondisi,
		BodyNode:         isiNode,
		ShouldReturnNull: false,
		Pos_Start:        kondisi.GetPosStart(),
		Pos_end:          isiNode.GetPosEnd(),
	})
}

func (p *parser) isDependCaseStart() bool {
	if p.currentToken().Kind == lexer.IDENTIFIER && p.currentToken().Value == "default" {
		if p.tokIndex+1 < len(p.tokens) && p.tokens[p.tokIndex+1].Kind == lexer.COLON {
			return true
		}
	}

	lookaheadCount := 0
	parenLevel := 0
	bracketLevel := 0
	for {
		if p.tokIndex+lookaheadCount >= len(p.tokens) {
			return false
		}
		tok := p.tokens[p.tokIndex+lookaheadCount]
		if tok.Kind == lexer.NEWLINE || tok.Kind == lexer.EOF || tok.Kind == lexer.ENDDEPENDON || tok.Kind == lexer.END {
			return false
		}
		if tok.Kind == lexer.OPEN_PAREN {
			parenLevel++
		} else if tok.Kind == lexer.CLOSE_PAREN {
			parenLevel--
		} else if tok.Kind == lexer.OPEN_BRACKET {
			bracketLevel++
		} else if tok.Kind == lexer.CLOSE_BRACKET {
			bracketLevel--
		}
		if parenLevel == 0 && bracketLevel == 0 && tok.Kind == lexer.COLON {
			return true
		}
		lookaheadCount++
	}
}

func (p *parser) parseDependCaseBody() (common.Expr, *common.ParseResult) {
	res := &common.ParseResult{}

	if p.currentToken().Kind != lexer.NEWLINE {
		stmt := res.Register(p.statement())
		if res.Error != nil {
			return nil, res
		}
		return stmt, res
	}

	for p.currentToken().Kind == lexer.NEWLINE {
		res.RegisterAdvancement()
		p.advance()
	}

	stmts := make([]common.Expr, 0)
	posStart := p.currentToken().PosStart.Copy()

	for {
		if p.currentToken().Kind == lexer.ENDDEPENDON || p.currentToken().Kind == lexer.END || p.currentToken().Kind == lexer.EOF {
			break
		}
		if p.isDependCaseStart() {
			break
		}

		stmt := res.Register(p.statement())
		if res.Error != nil {
			return nil, res
		}
		stmts = append(stmts, stmt)

		for p.currentToken().Kind == lexer.NEWLINE {
			res.RegisterAdvancement()
			p.advance()
		}
	}

	return common.ListNode{
		ElementNode: stmts,
		Pos_Start:   posStart,
		Pos_End:     p.currentToken().PosEnd.Copy(),
	}, res
}

func (p *parser) parseDepend() common.Expr {
	res := &common.ParseResult{}

	if p.currentToken().Kind != lexer.DEPEND {
		errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), "Expected 'depend'")
		return res.Failure(&errorNya)
	}
	posStart := p.currentToken().PosStart.Copy()

	res.RegisterAdvancement()
	p.advance()

	if p.currentToken().Kind != lexer.ON {
		errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), "Expected 'on'")
		return res.Failure(&errorNya)
	}
	res.RegisterAdvancement()
	p.advance()

	var switchExpr common.Expr
	if p.currentToken().Kind != lexer.NEWLINE {
		switchExpr = res.Register(p.expr())
		if res.Error != nil {
			return res
		}
	}

	if p.currentToken().Kind == lexer.NEWLINE {
		res.RegisterAdvancement()
		p.advance()
	}

	var cases []common.SwitchCase
	var defaultCase common.Expr = common.NullNode{}

	for p.currentToken().Kind != lexer.ENDDEPENDON && p.currentToken().Kind != lexer.END {
		if p.currentToken().Kind == lexer.EOF {
			errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.currentToken().PosEnd), "Expected 'enddependon' or 'end'")
			return res.Failure(&errorNya)
		}
		if p.currentToken().Kind == lexer.NEWLINE {
			res.RegisterAdvancement()
			p.advance()
			continue
		}

		if p.currentToken().Kind == lexer.IDENTIFIER && p.currentToken().Value == "default" {
			res.RegisterAdvancement()
			p.advance()

			if p.currentToken().Kind != lexer.COLON {
				errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.currentToken().PosEnd), "Expected ':' after default")
				return res.Failure(&errorNya)
			}
			res.RegisterAdvancement()
			p.advance()

			bodyExpr, bodyRes := p.parseDependCaseBody()
			res.Register(bodyRes)
			if res.Error != nil {
				return res
			}
			defaultCase = bodyExpr
		} else {
			caseVal := res.Register(p.expr())
			if res.Error != nil {
				return res
			}

			if p.currentToken().Kind != lexer.COLON {
				errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.currentToken().PosEnd), "Expected ':' after case value")
				return res.Failure(&errorNya)
			}
			res.RegisterAdvancement()
			p.advance()

			bodyExpr, bodyRes := p.parseDependCaseBody()
			res.Register(bodyRes)
			if res.Error != nil {
				return res
			}

			cases = append(cases, common.SwitchCase{Value: caseVal, Body: bodyExpr})
		}
	}

	res.RegisterAdvancement()
	p.advance()

	return res.Success(common.DependNode{
		Value:       switchExpr,
		Cases:       cases,
		DefaultCase: defaultCase,
		Pos_Start:   posStart,
		Pos_End:     p.currentToken().PosEnd.Copy(),
	})
}

func (p *parser) parseRepeat() common.Expr {
	res := &common.ParseResult{}

	if p.currentToken().Kind != lexer.REPEAT {
		errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), "Expected 'repeat'")
		return res.Failure(&errorNya)
	}
	posStart := p.currentToken().PosStart.Copy()

	res.RegisterAdvancement()
	p.advance()

	if p.currentToken().Kind == lexer.NEWLINE {
		res.RegisterAdvancement()
		p.advance()

		body := res.Register(p.statements())
		if res.Error != nil {
			return res
		}

		if p.currentToken().Kind != lexer.UNTIL {
			errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.currentToken().PosEnd), "Expected 'until'")
			return res.Failure(&errorNya)
		}

		res.RegisterAdvancement()
		p.advance()

		kondisi := res.Register(p.expr())
		if res.Error != nil {
			return res
		}

		return res.Success(common.RepeatNode{
			KondisiNode:      kondisi,
			BodyNode:         body,
			ShouldReturnNull: true,
			Pos_Start:        posStart,
			Pos_end:          kondisi.GetPosEnd(),
		})
	}

	isiNode := res.Register(p.statement())
	if res.Error != nil {
		return res
	}

	if p.currentToken().Kind != lexer.END && p.currentToken().Kind != lexer.UNTIL {
		errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.currentToken().PosEnd), "Expected 'until'")
		return res.Failure(&errorNya)
	}

	res.RegisterAdvancement()
	p.advance()

	kondisi := res.Register(p.expr())
	if res.Error != nil {
		return res
	}

	return res.Success(common.WhileNode{
		KondisiNode:      kondisi,
		BodyNode:         isiNode,
		ShouldReturnNull: false,
		Pos_Start:        kondisi.GetPosStart(),
		Pos_end:          isiNode.GetPosEnd(),
	})
}

func (p *parser) parseFunctionDef() common.Expr {
	res := &common.ParseResult{}

	if p.currentToken().Kind != lexer.FUNCTION {
		errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), "Expected 'function'")
		return res.Failure(&errorNya)
	}

	res.RegisterAdvancement()
	p.advance()

	var varNameTok *lexer.Token
	if p.currentToken().Kind == lexer.IDENTIFIER {
		tok := p.currentToken()
		varNameTok = &tok
		res.RegisterAdvancement()
		p.advance()
	}

	if p.currentToken().Kind != lexer.OPEN_PAREN {
		errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), "Expected '('")
		return res.Failure(&errorNya)
	}

	res.RegisterAdvancement()
	p.advance()

	argNameToks := make([]lexer.Token, 0)

	if p.currentToken().Kind == lexer.IDENTIFIER {
		argNameToks = append(argNameToks, p.currentToken())
		res.RegisterAdvancement()
		p.advance()

		for p.currentToken().Kind == lexer.COMMA {
			res.RegisterAdvancement()
			p.advance()

			if p.currentToken().Kind != lexer.IDENTIFIER {
				errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), "Expected identifier")
				return res.Failure(&errorNya)
			}

			argNameToks = append(argNameToks, p.currentToken())
			res.RegisterAdvancement()
			p.advance()
		}

		if p.currentToken().Kind != lexer.CLOSE_PAREN {
			errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), "Expected ')' or ','")
			return res.Failure(&errorNya)
		}
	} else {
		if p.currentToken().Kind != lexer.CLOSE_PAREN {
			errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), "Expected ')' or identifier")
			return res.Failure(&errorNya)
		}
	}

	res.RegisterAdvancement()
	p.advance()

	if p.currentToken().Kind == lexer.RIGHT_ARROW {
		res.RegisterAdvancement()
		p.advance()

		nodeToReturn := res.Register(p.expr())
		if res.Error != nil {
			return res
		}

		var posStart *tools.Position
		if varNameTok != nil {
			posStart = varNameTok.PosStart
		}
		funcNode := common.FuncNode{
			VarNameTok:       varNameTok,
			ArgNameToks:      argNameToks,
			BodyNode:         nodeToReturn,
			ShouldAutoReturn: true,
			Pos_Start:        posStart,
			Pos_end:          nodeToReturn.GetPosEnd(),
		}
		funcNode.Pos_Start = funcNode.GetPosStart()
		funcNode.Pos_end = funcNode.GetPosEnd()

		return res.Success(funcNode)
	}

	if p.currentToken().Kind != lexer.NEWLINE {
		errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), "Expected '->' or new line")
		return res.Failure(&errorNya)
	}

	res.RegisterAdvancement()
	p.advance()

	body := res.Register(p.statements())
	if res.Error != nil {
		return res
	}

	if p.currentToken().Kind != lexer.END && p.currentToken().Kind != lexer.ENDIF {
		errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.advance().PosEnd), "Expected 'end' to close the function body")
		return res.Failure(&errorNya)
	}

	res.RegisterAdvancement()
	p.advance()

	var posStart *tools.Position
	if varNameTok != nil {
		posStart = varNameTok.PosStart
	}
	return res.Success(common.FuncNode{
		VarNameTok:       varNameTok,
		ArgNameToks:      argNameToks,
		BodyNode:         body,
		ShouldAutoReturn: false,
		Pos_Start:        posStart,
		Pos_end:          body.GetPosEnd(),
	})
}

func (p *parser) atom() common.Expr {
	res := &common.ParseResult{}
	tok := p.currentToken()

	switch tok.Kind {
	case lexer.NUMBER:
		res.RegisterAdvancement()
		p.advance()

		return res.Success(common.NumberNode{
			Token:     tok,
			Pos_Start: tok.PosStart,
			Pos_End:   tok.PosEnd,
		})
	case lexer.STRING:
		res.RegisterAdvancement()
		p.advance()

		return res.Success(common.StringNode{
			Token:     tok,
			Pos_Start: tok.PosStart,
			Pos_End:   tok.PosEnd,
		})
	case lexer.IDENTIFIER:
		res.RegisterAdvancement()
		p.advance()

		return res.Success(common.VarAccessNode{VarNameTok: tok, Pos_Start: tok.PosStart, Pos_end: tok.PosEnd})
	case lexer.OPEN_PAREN:
		res.RegisterAdvancement()
		p.advance()

		expr := res.Register(p.expr())
		if p.currentToken().Kind == lexer.CLOSE_PAREN {
			res.RegisterAdvancement()
			p.advance()

			return res.Success(expr)
		}

		errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.currentToken().PosEnd), "Expected ')'")
		return res.Failure(&errorNya)
	case lexer.OPEN_BRACKET:
		listExpr := res.Register(p.parseList())
		if res.Error != nil {
			return res
		}

		return res.Success(listExpr)
	case lexer.IF:
		ifExpr := res.Register(p.parseIf())
		if res.Error != nil {
			return res
		}

		return res.Success(ifExpr)
	case lexer.FOR:
		forExpr := res.Register(p.parseFor())
		if res.Error != nil {
			return res
		}

		return res.Success(forExpr)
	case lexer.WHILE:
		whileExpr := res.Register(p.parseWhile())
		if res.Error != nil {
			return res
		}

		return res.Success(whileExpr)
	case lexer.REPEAT:
		repeatExpr := res.Register(p.parseRepeat())
		if res.Error != nil {
			return res
		}

		return res.Success(repeatExpr)
	case lexer.DEPEND:
		dependExpr := res.Register(p.parseDepend())
		if res.Error != nil {
			return res
		}

		return res.Success(dependExpr)
	case lexer.FUNCTION:
		functionDef := res.Register(p.parseFunctionDef())
		if res.Error != nil {
			return res
		}

		return res.Success(functionDef)
	}

	errorNya := common.InvalidSyntax(safePos(tok.PosStart), safePos(tok.PosEnd), "Expected identifier, int, float, '+', '-', '[', '(', 'if', 'for', 'while', 'function'")
	return res.Failure(&errorNya)
}

func (p *parser) call() common.Expr {
	res := &common.ParseResult{}
	atom := res.Register(p.atom())
	if res.Error != nil {
		return res
	}

	for p.currentToken().Kind == lexer.OPEN_PAREN || p.currentToken().Kind == lexer.OPEN_BRACKET || p.currentToken().Kind == lexer.DOT || tools.IsBuiltinFunction(atom.Print()) {
		if p.currentToken().Kind == lexer.OPEN_PAREN || tools.IsBuiltinFunction(atom.Print()) {
			var argNodes []common.Expr
			if p.currentToken().Kind == lexer.OPEN_PAREN {
				res.RegisterAdvancement()
				p.advance()

				if p.currentToken().Kind == lexer.CLOSE_PAREN {
					res.RegisterAdvancement()
					p.advance()
				} else {
					argNodes = append(argNodes, res.Register(p.expr()))
					if res.Error != nil {
						return res
					}

					for p.currentToken().Kind == lexer.COMMA {
						res.RegisterAdvancement()
						p.advance()

						argNodes = append(argNodes, res.Register(p.expr()))
						if res.Error != nil {
							return res
						}
					}

					if p.currentToken().Kind != lexer.CLOSE_PAREN {
						errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.currentToken().PosEnd), "Expected ',' or ')'")
						return res.Failure(&errorNya)
					}

					res.RegisterAdvancement()
					p.advance()
				}
			} else if tools.IsBuiltinFunction(atom.Print()) {
				for {
					if p.currentToken().Kind == lexer.NEWLINE || p.currentToken().Kind == lexer.EOF || p.currentToken().Kind == lexer.ENDPROGRAM || p.currentToken().Kind == lexer.ENDWHILE || p.currentToken().Kind == lexer.ENDFOR || p.currentToken().Kind == lexer.ENDIF || p.currentToken().Kind == lexer.ELSE || p.currentToken().Kind == lexer.ELIF {
						break
					}
					argNodes = append(argNodes, res.Register(p.expr()))
					if res.Error != nil {
						return res
					}
					if p.currentToken().Kind == lexer.COMMA {
						res.RegisterAdvancement()
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
				Pos_end:    p.currentToken().PosEnd.Copy(),
			}
		} else if p.currentToken().Kind == lexer.OPEN_BRACKET {
			res.RegisterAdvancement()
			p.advance()

			indexExpr := res.Register(p.expr())
			if res.Error != nil {
				return res
			}

			if p.currentToken().Kind != lexer.CLOSE_BRACKET {
				errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.currentToken().PosEnd), "Expected ']'")
				return res.Failure(&errorNya)
			}

			res.RegisterAdvancement()
			p.advance()

			atom = common.ArrayIndexNode{
				Left:      atom,
				Index:     indexExpr,
				Pos_Start: atom.GetPosStart(),
				Pos_End:   p.tokens[p.tokIndex-1].PosEnd.Copy(),
			}
		} else if p.currentToken().Kind == lexer.DOT {
			res.RegisterAdvancement()
			p.advance()

			if p.currentToken().Kind != lexer.IDENTIFIER {
				errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.currentToken().PosEnd), "Expected field name after '.'")
				return res.Failure(&errorNya)
			}
			memberTok := p.currentToken()
			res.RegisterAdvancement()
			p.advance()

			atom = common.MemberAccessNode{
				Object:    atom,
				MemberTok: memberTok,
				Pos_Start: atom.GetPosStart(),
				Pos_End:   memberTok.PosEnd.Copy(),
			}
		}
	}
	return res.Success(atom)
}

func (p *parser) power() common.Expr {
	return p.binOp(p.call, []lexer.TokenKind{lexer.POWER, lexer.DOUBLE_STAR}, p.factor)
}

func (p *parser) factor() common.Expr {
	res := &common.ParseResult{}
	tok := p.currentToken()

	if tok.Kind == lexer.PLUS || tok.Kind == lexer.DASH || tok.Kind == lexer.BIT_NOT {
		res.RegisterAdvancement()
		p.advance()
		factor := res.Register(p.factor())

		if res.Error != nil {
			return res
		}

		return res.Success(common.UnaryOpNode{
			Operator:  tok,
			Node:      factor,
			Pos_Start: tok.PosStart,
			Pos_End:   factor.GetPosEnd(),
		})
	}

	return p.power()
}

func (p *parser) term() common.Expr {
	return p.binOp(p.factor, []lexer.TokenKind{lexer.STAR, lexer.SLASH, lexer.DIV, lexer.PERCENT}, p.factor)
}

func (p *parser) arithExpr() common.Expr {
	return p.binOp(p.term, []lexer.TokenKind{lexer.PLUS, lexer.DASH}, p.term)
}

func (p *parser) shiftExpr() common.Expr {
	return p.binOp(p.arithExpr, []lexer.TokenKind{lexer.SHIFT_LEFT, lexer.SHIFT_RIGHT}, p.arithExpr)
}

func (p *parser) bitAndExpr() common.Expr {
	return p.binOp(p.shiftExpr, []lexer.TokenKind{lexer.BIT_AND}, p.shiftExpr)
}

func (p *parser) bitXorExpr() common.Expr {
	return p.binOp(p.bitAndExpr, []lexer.TokenKind{lexer.BIT_XOR}, p.bitAndExpr)
}

func (p *parser) bitOrExpr() common.Expr {
	return p.binOp(p.bitXorExpr, []lexer.TokenKind{lexer.BIT_OR}, p.bitXorExpr)
}

func (p *parser) compExpr() common.Expr {
	res := &common.ParseResult{}

	if p.currentToken().Kind == lexer.NOT {
		opTok := p.currentToken()

		res.RegisterAdvancement()
		p.advance()

		node := res.Register(p.compExpr())
		if res.Error != nil {
			return res
		}

		return res.Success(common.UnaryOpNode{
			Operator:  opTok,
			Node:      node,
			Pos_Start: opTok.PosStart,
			Pos_End:   node.GetPosEnd(),
		})
	}

	node := res.Register(p.binOp(p.bitOrExpr, []lexer.TokenKind{lexer.EQUALS, lexer.NOT_EQUALS, lexer.LESS, lexer.LESS_EQUALS, lexer.GREATER, lexer.GREATER_EQUALS}, p.bitOrExpr))
	if res.Error != nil {
		errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.currentToken().PosEnd), "Expected a number, an identifier, 'not', '+', '-', '~', '(' or '['")
		return res.Failure(&errorNya)
	}

	return res.Success(node)
}

func (p *parser) statements() common.Expr {
	res := &common.ParseResult{}
	statements := make([]common.Expr, 0)
	posStart := p.currentToken().PosStart.Copy()

	for p.currentToken().Kind == lexer.NEWLINE {
		res.RegisterAdvancement()
		p.advance()
	}

	isBlockEnding := func(kind lexer.TokenKind) bool {
		return kind == lexer.EOF ||
			kind == lexer.ENDPROGRAM ||
			kind == lexer.ENDIF ||
			kind == lexer.ENDWHILE ||
			kind == lexer.ENDFOR ||
			kind == lexer.END ||
			kind == lexer.UNTIL ||
			kind == lexer.ELSE ||
			kind == lexer.ELIF
	}

	if isBlockEnding(p.currentToken().Kind) {
		return res.Success(common.ListNode{
			ElementNode: statements,
			Pos_Start:   posStart,
			Pos_End:     p.currentToken().PosEnd.Copy(),
		})
	}

	statement := res.Register(p.statement())
	if res.Error != nil {
		return res
	}
	statements = append(statements, statement)

	moreStatement := true
	for moreStatement {
		newLineCount := 0
		for p.currentToken().Kind == lexer.NEWLINE {
			res.RegisterAdvancement()
			p.advance()
			newLineCount++
		}

		if newLineCount == 0 && len(statements) > 0 {
			break
		}

		if isBlockEnding(p.currentToken().Kind) {
			break
		}

		statementExpr := p.statement()
		statementResult, ok := statementExpr.(*common.ParseResult)
		if !ok {
			statements = append(statements, statementExpr)
			continue
		}

		if statementResult.Error != nil {
			if statementResult.AdvanceCount > 0 || !isBlockEnding(p.currentToken().Kind) {
				res.Register(statementResult)
				return res
			}
			p.reverse(statementResult.AdvanceCount)
			moreStatement = false
			continue
		}
		res.Register(statementResult)
		statements = append(statements, statementResult.Node)
	}

	return res.Success(common.ListNode{
		ElementNode: statements,
		Pos_Start:   posStart,
		Pos_End:     p.currentToken().PosEnd.Copy(),
	})
}

func (p *parser) statement() common.Expr {
	res := &common.ParseResult{}
	posStart := p.currentToken().PosStart.Copy()

	switch p.currentToken().Kind {
	case lexer.RETURN:
		res.RegisterAdvancement()
		p.advance()

		expr := res.TryRegister(p.expr())
		if expr == nil {
			p.reverse(res.ToReverseCount)
		}

		return res.Success(common.ReturnNode{NodeToReturn: expr, Pos_Start: posStart, Pos_end: p.currentToken().PosEnd.Copy()})
	case lexer.CONTINUE:
		res.RegisterAdvancement()
		p.advance()

		return res.Success(common.ContinueNode{Pos_Start: posStart, Pos_end: p.currentToken().PosStart.Copy()})
	case lexer.BREAK:
		res.RegisterAdvancement()
		p.advance()

		return res.Success(common.BreakNode{Pos_Start: posStart, Pos_end: p.currentToken().PosStart.Copy()})
	case lexer.DICTIONARY:
		res.RegisterAdvancement()
		p.advance()

		expr := res.Register(p.parseDictionary())
		if res.Error != nil {
			return res
		}

		return res.Success(expr)
	case lexer.ENDPROGRAM:
		res.RegisterAdvancement()
		p.advance()

		p.hasEndProgram = true
		return res.Success(nil)
	case lexer.ALGORITHM:
		res.RegisterAdvancement()
		p.advance()

		return res.Success(common.NullNode{})
	}

	expr := res.Register(p.expr())
	if res.Error != nil {
		errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.currentToken().PosEnd), fmt.Sprintf("Got '%s', expected 'continue', 'break', 'return', 'for', 'while', 'function', a number, an identifier, '+', '-', '(' or '['", p.currentToken().Value))
		return res.Failure(&errorNya)
	}

	return res.Success(expr)
}

func (p *parser) expr() common.Expr {
	res := &common.ParseResult{}

	isLeftArrow := false
	if p.currentToken().Kind == lexer.IDENTIFIER {
		var lookahead lexer.Token
		lookaheadCount := 0
		bracketLevel := 0

		for {
			lookahead = p.tokens[p.tokIndex+lookaheadCount]
			if lookahead.Kind == lexer.OPEN_BRACKET {
				bracketLevel++
			} else if lookahead.Kind == lexer.CLOSE_BRACKET {
				bracketLevel--
			}

			if bracketLevel < 0 {
				break
			}

			if bracketLevel == 0 && (lookahead.Kind == lexer.ASSIGNMENT || lookahead.Kind == lexer.LEFT_ARROW) {
				isLeftArrow = true
				break
			}
			if lookahead.Kind == lexer.NEWLINE || lookahead.Kind == lexer.EOF {
				break
			}
			if bracketLevel == 0 && lookahead.Kind != lexer.IDENTIFIER && lookahead.Kind != lexer.OPEN_BRACKET && lookahead.Kind != lexer.CLOSE_BRACKET && lookahead.Kind != lexer.NUMBER && lookahead.Kind != lexer.DOT {
				break
			}
			lookaheadCount++
			if p.tokIndex+lookaheadCount >= len(p.tokens) {
				break
			}
		}

		if isLeftArrow {
			lhs := res.Register(p.call())
			if res.Error != nil {
				return res
			}

			if p.currentToken().Kind != lexer.ASSIGNMENT && p.currentToken().Kind != lexer.LEFT_ARROW {
				errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.currentToken().PosEnd), "Expected '<-' or '='")
				return res.Failure(&errorNya)
			}

			res.RegisterAdvancement()
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
				errorNya := common.InvalidSyntax(safePos(lhs.GetPosStart()), safePos(lhs.GetPosEnd()), "Illegal assignment target")
				return res.Failure(&errorNya)
			}
		}
	}

	node := res.Register(p.binOp(p.compExpr, []lexer.TokenKind{lexer.AND, lexer.OR}, p.compExpr))
	if res.Error != nil {
		errorNya := common.InvalidSyntax(safePos(p.currentToken().PosStart), safePos(p.currentToken().PosEnd), "Expected 'for', 'while', 'function', a number, an identifier, '+', '-', '(' or '['")
		return res.Failure(&errorNya)
	}

	return res.Success(node)
}

func (p *parser) binOp(fnA func() common.Expr, ops []lexer.TokenKind, fnB func() common.Expr) common.Expr {
	res := &common.ParseResult{}
	left := res.Register(fnA())

	if res.Error != nil {
		return res
	}

	for slices.Contains(ops, p.currentToken().Kind) {
		opTok := p.currentToken()
		res.RegisterAdvancement()
		p.advance()

		right := res.Register(fnB())
		if res.Error != nil {
			return res
		}

		left = common.BinOpNode{
			Operator:  opTok,
			Left:      left,
			Right:     right,
			Pos_Start: opTok.PosStart,
			Pos_End:   right.GetPosEnd(),
		}
	}

	return res.Success(left)
}
