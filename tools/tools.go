package tools

import (
	"fmt"
	"strings"
)

// Position nyimpen posisi kursor/karakter di file source code (index, baris, kolom, nama file).
type Position struct {
	Idx  int
	Ln   int
	Col  int
	Fn   string
	Ftxt string
}

// AdvanceN buat majuin posisi kursor sebanyak n karakter.
func (p *Position) AdvanceN(n int) {
	p.Idx += n
	p.Col += n
}

// Advance buat majuin posisi kursor 1 karakter, otomatis nge-handle baris baru kalo nemu newline.
func (p *Position) Advance(currentChar string) *Position {
	p.Idx++
	p.Col++

	if currentChar == "\n" {
		p.Ln++
		p.Col = 0
	}

	return p
}

// Copy buat bikin kloningan Position yang terpisah.
func (p *Position) Copy() *Position {
	if p == nil {
		return nil
	}
	return &Position{
		Idx:  p.Idx,
		Ln:   p.Ln,
		Col:  p.Col,
		Fn:   p.Fn,
		Ftxt: p.Ftxt,
	}
}

// BuiltinFunctions nyimpen mapping keyword function bawaan DAP ke nama function eksekutornya.
var BuiltinFunctions = map[string]string{
	"print":            "Print",
	"PRINT":            "Print",
	"write":            "Print",
	"WRITE":            "Print",
	"output":           "Print",
	"OUTPUT":           "Print",
	"read":             "Input",
	"READ":             "Input",
	"input":            "Input",
	"INPUT":            "Input",
	"panjang":          "Panjang",
	"PANJANG":          "Panjang",
	"Panjang":          "Panjang",
	"length":           "Panjang",
	"LENGTH":           "Panjang",
	"len":              "Panjang",
	"bulat":            "Bulat",
	"BULAT":            "Bulat",
	"Bulat":            "Bulat",
	"round":            "Bulat",
	"ROUND":            "Bulat",
	"isNumber":         "IsNumber",
	"is_number":        "IsNumber",
	"ISNUMBER":         "IsNumber",
	"isString":         "IsString",
	"is_string":        "IsString",
	"ISSTRING":         "IsString",
	"append":           "Append",
	"APPEND":           "Append",
	"tulis":            "Print",
	"TULIS":            "Print",
	"masukan":          "Input",
	"MASUKAN":          "Input",
	"outputRet":        "PrintRet",
	"output_ret":       "PrintRet",
	"printRet":         "PrintRet",
	"print_ret":        "PrintRet",
	"inputBaris":       "InputBaris",
	"input_baris":      "InputBaris",
	"INPUTBARIS":       "InputBaris",
	"inputbaris":       "InputBaris",
	"konversiKeString": "KonversiKeString",
	"konversi_ke_string": "KonversiKeString",
	"KONVERSIKESTRING": "KonversiKeString",
	"konversikestring": "KonversiKeString",
}

// SemuaBuiltInFunction alias buat backward compatibility.
var SemuaBuiltInFunction = BuiltinFunctions

// IsBuiltinFunction ngecek apakah nama function ini terdaftar sebagai function bawaan DAP.
func IsBuiltinFunction(s string) bool {
	return BuiltinFunctions[s] != ""
}

// ApakahBuiltinFunction alias buat backward compatibility.
func ApakahBuiltinFunction(s string) bool {
	return IsBuiltinFunction(s)
}

// FormatSyntaxError buat nge-format syntax error jadi pesan traceback yang rapi plus penunjuk panah (^).
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

	var msg strings.Builder
	msg.WriteString("Traceback (most recent call last):\n")
	msg.WriteString(fmt.Sprintf("  File \"%s\", line %d\n", pos.Fn, pos.Ln+1))
	msg.WriteString(fmt.Sprintf("    %s\n", lineText))
	msg.WriteString(fmt.Sprintf("    %s^\n", caretSpacing.String()))
	if details != "" {
		msg.WriteString(fmt.Sprintf("%s: %s", errName, details))
	} else {
		msg.WriteString(errName)
	}
	return msg.String()
}
