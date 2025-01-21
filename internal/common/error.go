package common

import (
	"dap/tools"
	"fmt"
)

type Error struct {
	PosStart  tools.Position
	PosEnd    tools.Position
	ErrorName string
	Details   string
	Context   *Context
}

type ErrorInterface interface {
	As_string() string
}

func (error Error) As_string() string {
	if error.ErrorName == "Runtime Error" {
		hasil := error.generate_traceback()
		hasil += fmt.Sprintf("%s: %s", error.ErrorName, error.Details)
		return hasil
	}

	hasil := fmt.Sprintf("%s: %s\n", error.ErrorName, error.Details)
	hasil += fmt.Sprintf("File %s, line %d", error.PosStart.Fn, error.PosStart.Ln+1)
	return hasil
}

func (error Error) generate_traceback() string {
	hasil := ""
	pos := error.PosStart
	ctx := error.Context

	for ctx != nil {
		hasil = fmt.Sprintf("File: %s, line %d, in %s\n%s", pos.Fn, pos.Ln+1, ctx.DisplayName, hasil)

		if ctx.ParentEntryPos == nil || ctx.Parent == nil {
			break
		}

		pos = *ctx.ParentEntryPos
		ctx = ctx.Parent
	}

	return hasil
}

func IllegalCharError(PosStart tools.Position, PosEnd tools.Position) Error {
	return Error{
		PosStart:  PosStart,
		PosEnd:    PosEnd,
		ErrorName: "Illegal Character Error",
		Details:   "",
	}
}

func InvalidSyntax(PosStart tools.Position, PosEnd tools.Position, details string) Error {
	return Error{
		PosStart:  PosStart,
		PosEnd:    PosEnd,
		ErrorName: "Invalid Syntax",
		Details:   details,
	}
}

func RTError(PosStart tools.Position, PosEnd tools.Position, details string, context *Context) Error {
	return Error{
		PosStart:  PosStart,
		PosEnd:    PosEnd,
		ErrorName: "Runtime Error",
		Details:   details,
		Context:   context,
	}
}
