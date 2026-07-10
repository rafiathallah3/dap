package tools

import (
	"fmt"
	"strings"
)

type Position struct {
	Idx  int
	Ln   int
	Col  int
	Fn   string
	Ftxt string
}

func (position *Position) AdvanceN(berapa int) {
	position.Idx += berapa
	position.Col += berapa
}

func (position *Position) Advance(CurrentChar string) *Position {
	position.Idx++
	position.Col++

	if CurrentChar == "\n" {
		position.Ln++
		position.Col = 0
	}

	return position
}

func (position *Position) Copy() *Position {
	return &Position{
		Idx:  position.Idx,
		Ln:   position.Ln,
		Col:  position.Col,
		Fn:   position.Fn,
		Ftxt: position.Ftxt,
	}
}

func GetComparison(hasil bool) int {
	if hasil {
		return 1
	}

	return 0
}

var SemuaBuiltInFunction map[string]string = map[string]string{
	"print":  "Print",
	"PRINT":  "Print",
	"write":  "Print",
	"WRITE":  "Print",
	"output": "Print",
	"OUTPUT": "Print",
	"read":   "Input",
	"READ":   "Input",
	"input":  "Input",
	"INPUT":  "Input",
	"panjang": "Panjang",
	"PANJANG": "Panjang",
	"Panjang": "Panjang",
	"bulat":   "Bulat",
	"BULAT":   "Bulat",
	"Bulat":   "Bulat",
	"inputBaris":       "InputBaris",
	"INPUTBARIS":       "InputBaris",
	"inputbaris":       "InputBaris",
	"konversiKeString": "KonversiKeString",
	"KONVERSIKESTRING": "KonversiKeString",
	"konversikestring": "KonversiKeString",
}

func ApakahBuiltinFunction(s string) bool {
	return SemuaBuiltInFunction[s] != ""
}

func FormatSyntaxError(pos Position, errorName string, details string) string {
	lines := strings.Split(strings.ReplaceAll(pos.Ftxt, "\r\n", "\n"), "\n")
	lineText := ""
	if pos.Ln >= 0 && pos.Ln < len(lines) {
		lineText = lines[pos.Ln]
	}

	var caretSpacing strings.Builder
	for i, r := range lineText {
		if i >= pos.Col {
			break
		}
		if r == '\t' {
			caretSpacing.WriteRune('\t')
		} else {
			caretSpacing.WriteRune(' ')
		}
	}
	if pos.Col > len(lineText) {
		padding := pos.Col - len(lineText)
		caretSpacing.WriteString(strings.Repeat(" ", padding))
	}

	// 3. Map DAP error names to Python-like equivalents and simplify verbose syntax error details
	errName := errorName
	switch errName {
	case "Invalid Syntax":
		errName = "SyntaxError"
		if details == "" {
			details = "invalid syntax"
		}
	case "Illegal Character Error":
		errName = "SyntaxError"
		if details == "" {
			details = "illegal character"
		}
	}

	if strings.Contains(details, "Expected identifier, int, float") ||
		strings.Contains(details, "Expected 'for', 'while', 'function', a number") ||
		strings.Contains(details, "Expected a number, an identifier, 'not'") ||
		(strings.Contains(details, "Got '") && strings.Contains(details, "expected 'continue', 'break'")) {
		details = "invalid syntax"
	}

	msg := "Traceback (most recent call last):\n"
	msg += fmt.Sprintf("  File \"%s\", line %d\n", pos.Fn, pos.Ln+1)
	msg += fmt.Sprintf("    %s\n", lineText)
	msg += fmt.Sprintf("    %s^\n", caretSpacing.String())
	if details != "" {
		msg += fmt.Sprintf("%s: %s", errName, details)
	} else {
		msg += fmt.Sprintf("%s", errName)
	}
	return msg
}
