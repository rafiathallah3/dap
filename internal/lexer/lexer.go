package lexer

import (
	"dap/tools"
	"fmt"
	"regexp"
)

type regexHandler func(lex *lexer, regex *regexp.Regexp)

type regexPattern struct {
	regex   *regexp.Regexp
	handler regexHandler
}

// Precompile semua regex pattern di awal biar ga recompilation terus tiap manggil Tokenize.
var defaultPatterns = []regexPattern{
	{regexp.MustCompile(`\n|;`), defaultHandler(NEWLINE, "\n")},
	{regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`), symbolHandler},
	{regexp.MustCompile(`[0-9]+(\.[0-9]+)?`), numberHandler},
	{regexp.MustCompile(`"([^"\\]*(?:\\.[^"\\]*)*)"`), stringHandler},
	{regexp.MustCompile(`'([^'\\]*(?:\\.[^'\\]*)*)'`), stringHandler},
	{regexp.MustCompile(`"([^"\n\\]*(?:\\.[^"\n\\]*)*)`), unterminatedStringHandler},
	{regexp.MustCompile(`\/\/.*`), skipHandler},
	{regexp.MustCompile(`(?s)\{\{.*?\}\}|\{[^}]*\}`), skipHandler},
	{regexp.MustCompile(`[ \t\r]+`), skipHandler},
	{regexp.MustCompile(`\[`), defaultHandler(OPEN_BRACKET, "[")},
	{regexp.MustCompile(`\]`), defaultHandler(CLOSE_BRACKET, "]")},
	{regexp.MustCompile(`\{`), defaultHandler(OPEN_CURLY, "{")},
	{regexp.MustCompile(`\}`), defaultHandler(CLOSE_CURLY, "}")},
	{regexp.MustCompile(`\(`), defaultHandler(OPEN_PAREN, "(")},
	{regexp.MustCompile(`\)`), defaultHandler(CLOSE_PAREN, ")")},
	{regexp.MustCompile(`==`), defaultHandler(EQUALS, "==")},
	{regexp.MustCompile(`!=`), defaultHandler(NOT_EQUALS, "!=")},
	{regexp.MustCompile(`->`), defaultHandler(RIGHT_ARROW, "->")},
	{regexp.MustCompile(`<-`), defaultHandler(LEFT_ARROW, "<-")},
	{regexp.MustCompile(`=`), defaultHandler(ASSIGNMENT, "=")},
	{regexp.MustCompile(`!`), defaultHandler(NOT, "!")},
	{regexp.MustCompile(`<=`), defaultHandler(LESS_EQUALS, "<=")},
	{regexp.MustCompile(`<<`), defaultHandler(SHIFT_LEFT, "<<")},
	{regexp.MustCompile(`<`), defaultHandler(LESS, "<")},
	{regexp.MustCompile(`>=`), defaultHandler(GREATER_EQUALS, ">=")},
	{regexp.MustCompile(`>>`), defaultHandler(SHIFT_RIGHT, ">>")},
	{regexp.MustCompile(`>`), defaultHandler(GREATER, ">")},
	{regexp.MustCompile(`\|\|`), defaultHandler(OR, "||")},
	{regexp.MustCompile(`\|`), defaultHandler(BIT_OR, "|")},
	{regexp.MustCompile(`&&`), defaultHandler(AND, "&&")},
	{regexp.MustCompile(`&`), defaultHandler(BIT_AND, "&")},
	{regexp.MustCompile(`~`), defaultHandler(BIT_NOT, "~")},
	{regexp.MustCompile(`\.\.`), defaultHandler(DOT_DOT, "..")},
	{regexp.MustCompile(`\.`), defaultHandler(DOT, ".")},
	{regexp.MustCompile(`:`), defaultHandler(COLON, ":")},
	{regexp.MustCompile(`\?`), defaultHandler(QUESTION, "?")},
	{regexp.MustCompile(`,`), defaultHandler(COMMA, ",")},
	{regexp.MustCompile(`\+=`), defaultHandler(PLUS_EQUALS, "+=")},
	{regexp.MustCompile(`-=`), defaultHandler(MINUS_EQUALS, "-=")},
	{regexp.MustCompile(`\+`), defaultHandler(PLUS, "+")},
	{regexp.MustCompile(`-`), defaultHandler(DASH, "-")},
	{regexp.MustCompile(`/`), defaultHandler(SLASH, "/")},
	{regexp.MustCompile(`\*\*`), defaultHandler(DOUBLE_STAR, "**")},
	{regexp.MustCompile(`\*`), defaultHandler(STAR, "*")},
	{regexp.MustCompile(`%`), defaultHandler(PERCENT, "%")},
	{regexp.MustCompile(`\^`), defaultHandler(BIT_XOR, "^")},
}

type lexer struct {
	patterns []regexPattern
	Tokens   []Token
	Pos      *tools.Position
	Source   string
	Err      error
}

func (lex *lexer) advanceText(text string) {
	for _, char := range text {
		lex.Pos.Idx++
		if char == '\n' {
			lex.Pos.Ln++
			lex.Pos.Col = 0
		} else {
			lex.Pos.Col++
		}
	}
}

func (lex *lexer) push(token Token) {
	lex.Tokens = append(lex.Tokens, token)
}

func (lex *lexer) remainder() string {
	return lex.Source[lex.Pos.Idx:]
}

func (lex *lexer) atEOF() bool {
	return lex.Pos.Idx >= len(lex.Source)
}

// Tokenize buat ngubah source code DAP jadi kumpulan slice Token.
func Tokenize(source string, fileName string) ([]Token, error) {
	if fileName == "" {
		fileName = "<stdin>"
	}

	lex := createLexer(source, fileName)

	for !lex.atEOF() {
		matched := false

		for _, pattern := range lex.patterns {
			loc := pattern.regex.FindStringIndex(lex.remainder())
			if loc != nil && loc[0] == 0 {
				pattern.handler(lex, pattern.regex)
				if lex.Err != nil {
					return nil, lex.Err
				}
				matched = true
				break
			}
		}

		if !matched {
			badChar := string(lex.remainder()[0])
			errMsg := tools.FormatSyntaxError(*lex.Pos, "Illegal Character Error", fmt.Sprintf("Unexpected character '%s'", badChar))
			return nil, fmt.Errorf("%s", errMsg)
		}
	}

	lex.push(NewToken(EOF, "EOF", lex.Pos, nil))
	return lex.Tokens, nil
}

func defaultHandler(kind TokenKind, value string) regexHandler {
	return func(lex *lexer, regex *regexp.Regexp) {
		lex.push(NewToken(kind, value, lex.Pos, nil))
		lex.advanceText(value)
	}
}

func createLexer(source string, fn string) *lexer {
	return &lexer{
		Source: source,
		Pos: &tools.Position{
			Idx:  0,
			Ln:   0,
			Col:  0,
			Fn:   fn,
			Ftxt: source,
		},
		Tokens:   make([]Token, 0),
		patterns: defaultPatterns,
	}
}

func numberHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.push(NewToken(NUMBER, match, lex.Pos, nil))
	lex.advanceText(match)
}

func skipHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.advanceText(match)
}

func stringHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.push(NewToken(STRING, match, lex.Pos, nil))
	lex.advanceText(match)
}

func unterminatedStringHandler(lex *lexer, regex *regexp.Regexp) {
	errMsg := tools.FormatSyntaxError(*lex.Pos, "SyntaxError", "unterminated string literal")
	lex.Err = fmt.Errorf("%s", errMsg)
}

func symbolHandler(lex *lexer, regex *regexp.Regexp) {
	value := regex.FindString(lex.remainder())

	if kind, exists := ReservedKeywords[value]; exists {
		lex.push(NewToken(kind, value, lex.Pos, nil))
	} else {
		lex.push(NewToken(IDENTIFIER, value, lex.Pos, nil))
	}

	lex.advanceText(value)
}
