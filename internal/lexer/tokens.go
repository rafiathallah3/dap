package lexer

import (
	"dap/tools"
	"fmt"
)

// TokenKind represents the type of a lexical token.
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
	DIV

	PROGRAM
	DICTIONARY
	ALGORITHM
	ENDPROGRAM
	VAR
	CONST
	NEW
	REPEAT
	UNTIL
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
	ARRAY
	OF
	TYPE
	DOWNTO
	DEPEND
	ON
	ENDDEPENDON
)

// ReservedKeywords nyimpen mapping keyword bahasa DAP ke TokenKind-nya masing-masing.
var ReservedKeywords = map[string]TokenKind{
	"var":         VAR,
	"newline":     NEWLINE,
	"const":       CONST,
	"constant":    CONST,
	"program":     PROGRAM,
	"endprogram":  ENDPROGRAM,
	"dictionary":  DICTIONARY,
	"kamus":       DICTIONARY,
	"algorithm":   ALGORITHM,
	"algoritma":   ALGORITHM,
	"new":         NEW,
	"repeat":      REPEAT,
	"until":       UNTIL,
	"function":    FUNCTION,
	"if":          IF,
	"then":        THEN,
	"elif":        ELIF,
	"else":        ELSE,
	"return":      RETURN,
	"continue":    CONTINUE,
	"break":       BREAK,
	"foreach":     FOREACH,
	"while":       WHILE,
	"for":         FOR,
	"to":          TO,
	"step":        STEP,
	"do":          DO,
	"end":         END,
	"endwhile":    ENDWHILE,
	"endfor":      ENDFOR,
	"endif":       ENDIF,
	"integer":     INTEGER,
	"real":        REAL,
	"string":      STRINGTYPE,
	"array":       ARRAY,
	"of":          OF,
	"type":        TYPE,
	"div":         DIV,
	"DIV":         DIV,
	"mod":         PERCENT,
	"MOD":         PERCENT,
	"and":         AND,
	"AND":         AND,
	"or":          OR,
	"OR":          OR,
	"not":         NOT,
	"NOT":         NOT,
	"downto":      DOWNTO,
	"DOWNTO":      DOWNTO,
	"depend":      DEPEND,
	"DEPEND":      DEPEND,
	"on":          ON,
	"ON":          ON,
	"enddependon": ENDDEPENDON,
	"ENDDEPENDON": ENDDEPENDON,
}

// reserved_lu disimpen biar tetep backward compatible.
var reserved_lu = ReservedKeywords

// Token nyimpen representasi token: jenisnya (Kind), nilai teks (Value), dan posisinya di source code.
type Token struct {
	Kind      TokenKind
	Value     string
	PosStart  *tools.Position
	PosEnd    *tools.Position
	Pos_Start *tools.Position
	Pos_End   *tools.Position
}

// IsOneOfMany ngecek apakah token ini cocok sama salah satu TokenKind yang dikasih.
func (token Token) IsOneOfMany(expectedTokens ...TokenKind) bool {
	for _, expectedToken := range expectedTokens {
		if token.Kind == expectedToken {
			return true
		}
	}
	return false
}

// Debug buat nge-print token dan nilainya pas debugging.
func (token Token) Debug() {
	if token.IsOneOfMany(IDENTIFIER, NUMBER, STRING) {
		fmt.Printf("%s (%s)\n", TokenKindString(token.Kind), token.Value)
	} else {
		fmt.Printf("%s ()\n", TokenKindString(token.Kind))
	}
}

// NewToken helper buat bikin Token baru lengkap dengan posisi awal & akhirnya.
func NewToken(kind TokenKind, value string, posStart *tools.Position, posEnd *tools.Position) Token {
	tk := Token{
		Kind:      kind,
		Value:     value,
		PosStart:  posStart,
		PosEnd:    posEnd,
		Pos_Start: posStart,
		Pos_End:   posEnd,
	}

	if posStart != nil {
		pStart := posStart.Copy()
		pEnd := posStart.Copy()
		pEnd.Advance("")
		tk.PosStart = pStart
		tk.PosEnd = pEnd
		tk.Pos_Start = pStart
		tk.Pos_End = pEnd
	}

	if posEnd != nil {
		tk.PosEnd = posEnd
		tk.Pos_End = posEnd
	}

	return tk
}

// String implements fmt.Stringer for TokenKind.
func (kind TokenKind) String() string {
	return TokenKindString(kind)
}

// TokenKindString returns the human-readable string for a TokenKind.
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
	case DIV:
		return "DIV"
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
	case UNTIL:
		return "UNTIL"
	case REPEAT:
		return "REPEAT"
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
	case ARRAY:
		return "ARRAY"
	case OF:
		return "OF"
	case TYPE:
		return "TYPE"
	case DOWNTO:
		return "DOWNTO"
	case DEPEND:
		return "DEPEND"
	case ON:
		return "ON"
	case ENDDEPENDON:
		return "ENDDEPENDON"
	default:
		return "UNKNOWN"
	}
}
