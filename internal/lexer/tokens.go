package lexer

import (
	"dap/tools"
	"fmt"
)

type TokenKind int

const (
	EOF TokenKind = iota
	NUMBER
	STRING
	IDENTIFIER

	OPEN_BRACKET
	CLOSE_BRACKET
	OPEN_CURLY
	CLOSE_CURLY
	OPEN_PAREN
	CLOSE_PAREN

	ASSIGNMENT
	EQUALS
	NOT_EQUALS

	LESS
	LESS_EQUALS
	GREATER
	GREATER_EQUALS
	RIGHT_ARROW
	LEFT_ARROW

	OR
	AND
	NOT

	DOT
	DOT_DOT
	NEWLINE
	COLON
	QUESTION
	COMMA

	PLUS_EQUALS
	MINUS_EQUALS
	SLASH_EQUALS
	STAR_EQUALS
	PERCENT_EQUALS

	PLUS
	DASH
	SLASH
	STAR
	PERCENT
	POWER

	PROGRAM
	DICTIONARY
	ALGORITHM
	ENDPROGRAM
	VAR
	CONST
	NEW
	IMPORT
	FROM
	FUNCTION
	IF
	THEN
	ELIF
	ELSE
	RETURN
	CONTINUE
	BREAK
	FOREACH
	WHILE
	FOR
	TO
	STEP
	DO
	END
	ENDWHILE
	ENDFOR
	ENDIF
	INTEGER
	REAL
	STRINGTYPE
)

var reserved_lu map[string]TokenKind = map[string]TokenKind{
	"var":        VAR,
	"newline":    NEWLINE,
	"const":      CONST,
	"program":    PROGRAM,
	"endprogram": ENDPROGRAM,
	"dictionary": DICTIONARY,
	"algorithm":  ALGORITHM,
	"new":        NEW,
	"import":     IMPORT,
	"from":       FROM,
	"function":   FUNCTION,
	"if":         IF,
	"then":       THEN,
	"elif":       ELIF,
	"else":       ELSE,
	"return":     RETURN,
	"continue":   CONTINUE,
	"break":      BREAK,
	"foreach":    FOREACH,
	"while":      WHILE,
	"for":        FOR,
	"to":         TO,
	"step":       STEP,
	"do":         DO,
	"end":        END,
	"endwhile":   ENDWHILE,
	"endfor":     ENDFOR,
	"endif":      ENDIF,
	"integer":    INTEGER,
	"real":       REAL,
	"string":     STRINGTYPE,
}

type Token struct {
	Kind      TokenKind
	Value     string
	Pos_Start *tools.Position
	Pos_End   *tools.Position
}

// func (token Token) expr() {}
// func (token Token) Print() string {
// 	return token.Value
// }
// func (n Token) Name() string {
// 	return "Token"
// }
// func (n Token) GetPosStart() *tools.Position {
// 	return n.Pos_Start
// }
// func (n Token) GetPosEnd() *tools.Position {
// 	return n.Pos_Start
// }

func (token Token) IsOneOfMany(expectedTokens ...TokenKind) bool {
	for _, expectedToken := range expectedTokens {
		if token.Kind == expectedToken {
			return true
		}
	}

	return false
}

func (token Token) Debug() {
	if token.IsOneOfMany(IDENTIFIER, NUMBER, STRING) {
		fmt.Printf("%s (%s)\n", TokenKindString(token.Kind), token.Value)
	} else {
		fmt.Printf("%s ()\n", TokenKindString(token.Kind))
	}
}

func NewToken(kind TokenKind, value string, Pos_start *tools.Position, Pos_end *tools.Position) Token {
	tk := Token{
		kind, value,
		Pos_start, Pos_end,
	}

	if Pos_start != nil {
		tk.Pos_Start = Pos_start.Copy()
		tk.Pos_End = Pos_start.Copy()
		tk.Pos_End.Advance("")
	}

	if Pos_end != nil {
		tk.Pos_End = Pos_end
	}

	return tk
}

func TokenKindString(kind TokenKind) string {
	switch kind {
	case EOF:
		return "EOF"
	case NUMBER:
		return "NUMBER"
	case STRING:
		return "STRING"
	case IDENTIFIER:
		return "IDENTIFIER"
	case OPEN_BRACKET:
		return "OPEN_BRACKET"
	case CLOSE_BRACKET:
		return "CLOSE_BRACKET"
	case OPEN_CURLY:
		return "OPEN_CURLY"
	case CLOSE_CURLY:
		return "CLOSE_CURLY"
	case OPEN_PAREN:
		return "OPEN_PAREN"
	case CLOSE_PAREN:
		return "CLOSE_PAREN"
	case ASSIGNMENT:
		return "ASSIGNMENT"
	case EQUALS:
		return "EQUALS"
	case NOT:
		return "NOT"
	case NOT_EQUALS:
		return "NOT_EQUALS"
	case LESS:
		return "LESS"
	case LESS_EQUALS:
		return "LESS_EQUALS"
	case GREATER:
		return "GREATER"
	case GREATER_EQUALS:
		return "GREATER_EQUALS"
	case RIGHT_ARROW:
		return "RIGHT_ARROW"
	case LEFT_ARROW:
		return "LEFT_ARROW"
	case OR:
		return "OR"
	case AND:
		return "AND"
	case DOT:
		return "DOT"
	case DOT_DOT:
		return "DOT_DOT"
	case NEWLINE:
		return "NEW LINE"
	case COLON:
		return "COLON"
	case QUESTION:
		return "QUESTION"
	case COMMA:
		return "COMMA"
	case PLUS_EQUALS:
		return "PLUS_EQUALS"
	case MINUS_EQUALS:
		return "MINUS_EQUALS"
	case SLASH_EQUALS:
		return "SLASH_EQUALS"
	case STAR_EQUALS:
		return "STAR_EQUALS"
	case PERCENT_EQUALS:
		return "PERCENT_EQUALS"
	case PLUS:
		return "PLUS"
	case DASH:
		return "DASH"
	case SLASH:
		return "SLASH"
	case STAR:
		return "STAR"
	case PERCENT:
		return "PERCENT"
	case POWER:
		return "POWER"
	case VAR:
		return "VAR"
	case CONST:
		return "CONST"
	case PROGRAM:
		return "PROGRAM"
	case DICTIONARY:
		return "DICTIONARY"
	case ALGORITHM:
		return "ALGORITHM"
	case ENDPROGRAM:
		return "ENDPROGRAM"
	case NEW:
		return "NEW"
	case IMPORT:
		return "IMPORT"
	case FROM:
		return "FROM"
	case FUNCTION:
		return "FUNCTION"
	case IF:
		return "IF"
	case THEN:
		return "THEN"
	case ELIF:
		return "ELIF"
	case ELSE:
		return "ELSE"
	case RETURN:
		return "RETURN"
	case CONTINUE:
		return "CONTINUE"
	case BREAK:
		return "BREAK"
	case FOREACH:
		return "FOREACH"
	case WHILE:
		return "WHILE"
	case FOR:
		return "FOR"
	case TO:
		return "TO"
	case STEP:
		return "STEP"
	case DO:
		return "DO"
	case ENDIF:
		return "ENDIF"
	case ENDWHILE:
		return "ENDWHILE"
	case ENDFOR:
		return "ENDFOR"
	case END:
		return "END"
	case INTEGER:
		return "INTEGER"
	case REAL:
		return "REAL"
	case STRINGTYPE:
		return "STRING TYPE"
	default:
		return "UNKNOWN"
	}
}
