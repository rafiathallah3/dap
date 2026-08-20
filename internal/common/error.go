package common

import (
	"dap/tools"
	"fmt"
)

// Error nyimpen detail error pas lexing, parsing, ataupun pas runtime DAP jalan.
type Error struct {
	PosStart  tools.Position
	PosEnd    tools.Position
	ErrorName string
	Details   string
	Context   *Context
}

// Error implement standard Go error interface biar bisa dipake kayak error Go biasa.
func (e Error) Error() string {
	return e.AsString()
}

// AsString buat nge-format pesan error biar rapi dan gampang dibaca.
func (e Error) AsString() string {
	if e.ErrorName == "Runtime Error" {
		traceback := e.generateTraceback()
		return fmt.Sprintf("%s%s: %s", traceback, e.ErrorName, e.Details)
	}

	return tools.FormatSyntaxError(e.PosStart, e.ErrorName, e.Details)
}

// As_string alias buat backward compatibility.
func (e Error) As_string() string {
	return e.AsString()
}

func (e Error) generateTraceback() string {
	result := ""
	pos := e.PosStart
	ctx := e.Context

	for ctx != nil {
		result = fmt.Sprintf("File: %s, line %d, in %s\n%s", pos.Fn, pos.Ln+1, ctx.DisplayName, result)

		if ctx.ParentEntryPos == nil || ctx.Parent == nil {
			break
		}

		pos = *ctx.ParentEntryPos
		ctx = ctx.Parent
	}

	return result
}

// IllegalCharError helper buat bikin error karakter aneh/ilegal.
func IllegalCharError(posStart tools.Position, posEnd tools.Position) Error {
	return Error{
		PosStart:  posStart,
		PosEnd:    posEnd,
		ErrorName: "Illegal Character Error",
		Details:   "",
	}
}

// InvalidSyntax helper buat bikin syntax error pas parsing.
func InvalidSyntax(posStart tools.Position, posEnd tools.Position, details string) Error {
	return Error{
		PosStart:  posStart,
		PosEnd:    posEnd,
		ErrorName: "Invalid Syntax",
		Details:   details,
	}
}

// RTError helper buat bikin runtime error pas program lagi dieksekusi.
func RTError(posStart tools.Position, posEnd tools.Position, details string, context *Context) Error {
	return Error{
		PosStart:  posStart,
		PosEnd:    posEnd,
		ErrorName: "Runtime Error",
		Details:   details,
		Context:   context,
	}
}
