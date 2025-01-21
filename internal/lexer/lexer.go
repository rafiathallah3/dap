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

type lexer struct {
	patterns []regexPattern
	Tokens   []Token
	Pos      *tools.Position
	Source   string
}

func (lex *lexer) advanceN(n int) {
	lex.Pos.AdvanceN(n)
	// lex.Pos += n
}

func (lex *lexer) push(token Token) {
	lex.Tokens = append(lex.Tokens, token)
}

// func (lex *lexer) at() byte {
// 	return lex.Source[lex.Pos]
// }

func (lex *lexer) remainder() string {
	return lex.Source[lex.Pos.Idx:]
}

func (lex *lexer) at_eof() bool {
	return lex.Pos.Idx >= len(lex.Source)
}

func Tokenize(source string) []Token {
	lex := createLexer(source, "<stdin>")

	for !lex.at_eof() {
		matched := false
		apakahAdaNewLine := false

		for _, pattern := range lex.patterns {
			loc := pattern.regex.FindStringIndex(lex.remainder())
			if loc != nil {
				if loc[0] == 1 && loc[1] == 2 && lex.remainder()[loc[0]:loc[1]] == "\n" && !apakahAdaNewLine {
					apakahAdaNewLine = true
					lex.push(NewToken(NEWLINE, "\n", lex.Pos, nil))
				}

				if loc[0] == 0 {
					pattern.handler(lex, pattern.regex)
					matched = true
					break
				}
			}
		}

		// if apakahAdaNewLine {
		// 	lex.push(NewToken(NEWLINE, "\n", lex.Pos, nil))
		// }

		if !matched {
			panic(fmt.Sprintf("Lexer::Error -> Unrecognized token near %s\n", lex.remainder()))
		}
	}

	lex.push(NewToken(EOF, "EOF", lex.Pos, nil))
	return lex.Tokens
}

func defaultHandler(kind TokenKind, value string) regexHandler {
	return func(lex *lexer, regex *regexp.Regexp) {
		lex.advanceN(len(value))
		lex.push(NewToken(kind, value, lex.Pos, nil))
	}
}

func createLexer(source string, fn string) *lexer {
	lex := &lexer{
		Source: source,
		Pos: &tools.Position{
			Idx:  0,
			Ln:   0,
			Col:  0,
			Fn:   fn,
			Ftxt: source,
		},
		Tokens: make([]Token, 0),
		patterns: []regexPattern{
			{regexp.MustCompile(`\n|;`), defaultHandler(NEWLINE, "\n")},
			{regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`), symbolHandler},
			{regexp.MustCompile(`[0-9]+(\.[0-9]+)?`), numberHandler},
			{regexp.MustCompile(`"([^"\\]*(?:\\.[^"\\]*)*)"`), stringHandler},
			{regexp.MustCompile(`\/\/.*`), skipHandler},
			{regexp.MustCompile(`\s+`), skipHandler},
			{regexp.MustCompile(`\[`), defaultHandler(OPEN_BRACKET, "[")},
			{regexp.MustCompile(`\]`), defaultHandler(CLOSE_BRACKET, "]")},
			{regexp.MustCompile(`\{`), defaultHandler(OPEN_CURLY, "{")},
			{regexp.MustCompile(`\}`), defaultHandler(CLOSE_CURLY, "}")},
			{regexp.MustCompile(`\(`), defaultHandler(OPEN_PAREN, "(")},
			{regexp.MustCompile(`\)`), defaultHandler(CLOSE_PAREN, ")")},
			{regexp.MustCompile(`==`), defaultHandler(EQUALS, "==")},
			{regexp.MustCompile(`!=`), defaultHandler(NOT_EQUALS, "!=")},
			{regexp.MustCompile(`=`), defaultHandler(ASSIGNMENT, "=")},
			{regexp.MustCompile(`!`), defaultHandler(NOT, "!")},
			{regexp.MustCompile(`<=`), defaultHandler(LESS_EQUALS, "<=")},
			{regexp.MustCompile(`<`), defaultHandler(LESS, "<")},
			{regexp.MustCompile(`>=`), defaultHandler(GREATER_EQUALS, ">=")},
			{regexp.MustCompile(`>`), defaultHandler(GREATER, ">")},
			{regexp.MustCompile(`->`), defaultHandler(RIGHT_ARROW, "->")},
			{regexp.MustCompile(`\|\|`), defaultHandler(OR, "||")},
			{regexp.MustCompile(`&&`), defaultHandler(AND, "&&")},
			{regexp.MustCompile(`\.\.`), defaultHandler(DOT_DOT, "..")},
			{regexp.MustCompile(`\.`), defaultHandler(DOT, ".")},
			{regexp.MustCompile(`:`), defaultHandler(COLON, ":")},
			{regexp.MustCompile(`\?`), defaultHandler(QUESTION, "?")},
			{regexp.MustCompile(`,`), defaultHandler(COMMA, ",")},
			// {regexp.MustCompile(`\+\+`), defaultHandler(PLUS_PLUS, "++")},
			// {regexp.MustCompile(`--`), defaultHandler(MINUS_MINUS, "--")},
			{regexp.MustCompile(`\+=`), defaultHandler(PLUS_EQUALS, "+=")},
			{regexp.MustCompile(`-=`), defaultHandler(MINUS_EQUALS, "-=")},
			{regexp.MustCompile(`\+`), defaultHandler(PLUS, "+")},
			{regexp.MustCompile(`-`), defaultHandler(DASH, "-")},
			{regexp.MustCompile(`/`), defaultHandler(SLASH, "/")},
			{regexp.MustCompile(`\*`), defaultHandler(STAR, "*")},
			{regexp.MustCompile(`%`), defaultHandler(PERCENT, "%")},
			{regexp.MustCompile(`^`), defaultHandler(POWER, "^")},
		},
	}

	return lex
}

func numberHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.push(NewToken(NUMBER, match, lex.Pos, nil))
	lex.advanceN(len(match))
}

func skipHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	lex.advanceN(match[1])
}

func stringHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	stringLiteral := lex.remainder()[match[0]:match[1]]

	lex.push(NewToken(STRING, stringLiteral, lex.Pos, nil))
	lex.advanceN(len(stringLiteral))
}

func symbolHandler(lex *lexer, regex *regexp.Regexp) {
	value := regex.FindString(lex.remainder())

	if kind, exists := reserved_lu[value]; exists {
		lex.push(NewToken(kind, value, lex.Pos, nil))
	} else {
		lex.push(NewToken(IDENTIFIER, value, lex.Pos, nil))
	}

	lex.advanceN(len(value))
}
