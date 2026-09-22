package common

import (
	"bufio"
	"dap/tools"
	"fmt"
	"io"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

func safePos(p *tools.Position) tools.Position {
	if p == nil {
		return tools.Position{}
	}
	return *p
}

// SymbolTable nyimpen mapping variabel dan lexical scoping-nya (termasuk parent scope).
type SymbolTable struct {
	Symbols map[string]Value
	Parent  *SymbolTable
}

// NewSymbolTable helper buat inisialisasi SymbolTable baru plus optional parent scope.
func NewSymbolTable(parent *SymbolTable) *SymbolTable {
	return &SymbolTable{
		Symbols: make(map[string]Value),
		Parent:  parent,
	}
}

// Get buat ngambil nilai variabel dari scope sekarang atau scope parent kalo ada.
func (st *SymbolTable) Get(name string) Value {
	if st == nil {
		return nil
	}
	if st.Symbols != nil {
		if val, ok := st.Symbols[name]; ok {
			return val
		}
	}
	if st.Parent != nil {
		return st.Parent.Get(name)
	}
	return nil
}

// Set buat nyimpen nilai variabel ke scope sekarang (udah ada lazy map init biar ga panic).
func (st *SymbolTable) Set(name string, value Value) {
	if st.Symbols == nil {
		st.Symbols = make(map[string]Value)
	}
	st.Symbols[name] = value
}

// Remove buat ngehapus variabel dari scope sekarang.
func (st *SymbolTable) Remove(name string) {
	if st.Symbols != nil {
		delete(st.Symbols, name)
	}
}

// Context nyimpen info scope runtime dan metadata eksekusi program.
type Context struct {
	DisplayName    string
	Parent         *Context
	ParentEntryPos *tools.Position
	Symbol_Table   *SymbolTable
	SymbolTable    *SymbolTable
}

// GetSymbolTable helper buat dapetin SymbolTable aktif (kompatibel dua-duanya).
func (c *Context) GetSymbolTable() *SymbolTable {
	if c == nil {
		return nil
	}
	if c.SymbolTable != nil {
		return c.SymbolTable
	}
	return c.Symbol_Table
}

// PrintValueInterpreter buat nge-format runtime Value jadi string yang rapi pas di-print.
func PrintValueInterpreter(n Value) string {
	switch n := n.(type) {
	case Boolean:
		if n.Value {
			return "true"
		}
		return "false"
	case Number:
		s := strconv.FormatFloat(n.Value, 'f', -1, 64)
		if n.IsReal && !strings.Contains(s, ".") && s != "NaN" && s != "+Inf" && s != "-Inf" && s != "Inf" {
			return s + ".0"
		}
		return s
	case String:
		return n.Value
	case List:
		if len(n.Elements) == 1 {
			return PrintValueInterpreter(n.Elements[0])
		}
		var sb strings.Builder
		sb.WriteString("[")
		for i, v := range n.Elements {
			sb.WriteString(PrintValueInterpreter(v))
			if i < len(n.Elements)-1 {
				sb.WriteString(", ")
			}
		}
		sb.WriteString("]")
		return sb.String()
	case Array:
		var sb strings.Builder
		sb.WriteString("[")
		for i, v := range n.Elements {
			sb.WriteString(PrintValueInterpreter(v))
			if i < len(n.Elements)-1 {
				sb.WriteString(", ")
			}
		}
		sb.WriteString("]")
		return sb.String()
	case Function:
		return fmt.Sprintf("<function %s>", n.Name)
	case Struct:
		var sb strings.Builder
		sb.WriteString("<")
		keys := make([]string, 0, len(n.Fields))
		for k := range n.Fields {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		for i, k := range keys {
			sb.WriteString(fmt.Sprintf("%s: %v", k, PrintValueInterpreter(n.Fields[k])))
			if i < len(keys)-1 {
				sb.WriteString(", ")
			}
		}
		sb.WriteString(">")
		return sb.String()
	case Type:
		return "<type>"
	case Null:
		return "NULL"
	case *RTResult:
		return PrintValueInterpreter(n.Value)
	}

	return ""
}

// Value represents a runtime value in DAP.
type Value interface {
	SetPos(posStart *tools.Position, posEnd *tools.Position) Value
	SetContext(ctx *Context) Value
	GetContext() *Context
	IsTrue() bool
	Copy() Value
	Print() string

	// Backward compatibility aliases
	Set_pos(posStart *tools.Position, posEnd *tools.Position) Value
	Set_context(ctx *Context) Value
	Get_context() *Context
	Is_true() bool
}

// Null represents the null value.
type Null struct {
	Context *Context
}

func (n Null) Print() string                                             { return "NULL" }
func (n Null) SetPos(posStart *tools.Position, posEnd *tools.Position) Value { return n }
func (n Null) SetContext(ctx *Context) Value                             { n.Context = ctx; return n }
func (n Null) GetContext() *Context                                      { return n.Context }
func (n Null) Copy() Value                                               { return Null{Context: n.Context} }
func (n Null) IsTrue() bool                                              { return false }
func (n Null) Set_pos(posStart *tools.Position, posEnd *tools.Position) Value {
	return n.SetPos(posStart, posEnd)
}
func (n Null) Set_context(ctx *Context) Value { return n.SetContext(ctx) }
func (n Null) Get_context() *Context          { return n.GetContext() }
func (n Null) Is_true() bool                  { return n.IsTrue() }

// Boolean represents a boolean value.
type Boolean struct {
	Value     bool
	Context   *Context
	Pos_Start *tools.Position
	Pos_End   *tools.Position
}

func (b Boolean) Print() string {
	if b.Value {
		return "true"
	}
	return "false"
}
func (b Boolean) SetPos(posStart *tools.Position, posEnd *tools.Position) Value {
	b.Pos_Start = posStart
	b.Pos_End = posEnd
	return b
}
func (b Boolean) SetContext(ctx *Context) Value { b.Context = ctx; return b }
func (b Boolean) GetContext() *Context          { return b.Context }
func (b Boolean) Copy() Value {
	return Boolean{
		Value:     b.Value,
		Pos_Start: b.Pos_Start,
		Pos_End:   b.Pos_End,
		Context:   b.Context,
	}
}
func (b Boolean) IsTrue() bool                                              { return b.Value }
func (b Boolean) Set_pos(posStart *tools.Position, posEnd *tools.Position) Value {
	return b.SetPos(posStart, posEnd)
}
func (b Boolean) Set_context(ctx *Context) Value { return b.SetContext(ctx) }
func (b Boolean) Get_context() *Context          { return b.GetContext() }
func (b Boolean) Is_true() bool                  { return b.IsTrue() }

func (b Boolean) Equal(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Boolean:
		return Boolean{Value: b.Value == other.Value, Context: b.Context}.SetPos(b.Pos_Start, b.Pos_End), nil
	case Number:
		return Boolean{Value: b.Value == (other.Value != 0), Context: b.Context}.SetPos(b.Pos_Start, b.Pos_End), nil
	}
	err := RTError(safePos(b.Pos_Start), safePos(b.Pos_End), "Illegal operation: cannot compare equality with the given type", b.Context)
	return nil, &err
}

func (b Boolean) NotEqual(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Boolean:
		return Boolean{Value: b.Value != other.Value, Context: b.Context}.SetPos(b.Pos_Start, b.Pos_End), nil
	case Number:
		return Boolean{Value: b.Value != (other.Value != 0), Context: b.Context}.SetPos(b.Pos_Start, b.Pos_End), nil
	}
	err := RTError(safePos(b.Pos_Start), safePos(b.Pos_End), "Illegal operation: cannot compare inequality with the given type", b.Context)
	return nil, &err
}

func (b Boolean) And(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Boolean:
		return Boolean{Value: b.Value && other.Value, Context: b.Context}.SetPos(b.Pos_Start, b.Pos_End), nil
	case Number:
		return Boolean{Value: b.Value && (other.Value != 0), Context: b.Context}.SetPos(b.Pos_Start, b.Pos_End), nil
	}
	err := RTError(safePos(b.Pos_Start), safePos(b.Pos_End), "Illegal operation: cannot apply 'and' with the given type", b.Context)
	return nil, &err
}

func (b Boolean) Or(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Boolean:
		return Boolean{Value: b.Value || other.Value, Context: b.Context}.SetPos(b.Pos_Start, b.Pos_End), nil
	case Number:
		return Boolean{Value: b.Value || (other.Value != 0), Context: b.Context}.SetPos(b.Pos_Start, b.Pos_End), nil
	}
	err := RTError(safePos(b.Pos_Start), safePos(b.Pos_End), "Illegal operation: cannot apply 'or' with the given type", b.Context)
	return nil, &err
}

// Forwarding aliases for Boolean
func (b Boolean) Get_comparison_eq(other Value) (Value, *Error) { return b.Equal(other) }
func (b Boolean) Get_comparison_nq(other Value) (Value, *Error) { return b.NotEqual(other) }
func (b Boolean) Anded_by(other Value) (Value, *Error)          { return b.And(other) }
func (b Boolean) Ored_by(other Value) (Value, *Error)           { return b.Or(other) }

// Array represents a fixed-size indexed array.
type Array struct {
	Elements  []Value
	Start     int
	End       int
	Context   *Context
	Pos_Start *tools.Position
	Pos_End   *tools.Position
}

func (n Array) Print() string { return PrintValueInterpreter(n) }
func (n Array) SetPos(posStart *tools.Position, posEnd *tools.Position) Value {
	n.Pos_Start = posStart
	n.Pos_End = posEnd
	return n
}
func (n Array) SetContext(ctx *Context) Value { n.Context = ctx; return n }
func (n Array) GetContext() *Context          { return n.Context }
func (n Array) Copy() Value {
	return Array{
		Elements:  slices.Clone(n.Elements), // Deep slice copy to prevent aliasing
		Start:     n.Start,
		End:       n.End,
		Context:   n.Context,
		Pos_Start: n.Pos_Start,
		Pos_End:   n.Pos_End,
	}
}
func (n Array) IsTrue() bool                                              { return len(n.Elements) != 0 }
func (n Array) Set_pos(posStart *tools.Position, posEnd *tools.Position) Value {
	return n.SetPos(posStart, posEnd)
}
func (n Array) Set_context(ctx *Context) Value { return n.SetContext(ctx) }
func (n Array) Get_context() *Context          { return n.GetContext() }
func (n Array) Is_true() bool                  { return n.IsTrue() }

// List represents a dynamic list of values.
type List struct {
	Elements  []Value
	Context   *Context
	Pos_Start *tools.Position
	Pos_End   *tools.Position
}

func (n List) Print() string { return PrintValueInterpreter(n) }
func (n List) SetPos(posStart *tools.Position, posEnd *tools.Position) Value {
	n.Pos_Start = posStart
	n.Pos_End = posEnd
	return n
}
func (n List) SetContext(ctx *Context) Value { n.Context = ctx; return n }
func (n List) GetContext() *Context          { return n.Context }
func (n List) Copy() Value {
	return List{
		Elements:  slices.Clone(n.Elements), // Deep slice copy to prevent aliasing
		Pos_Start: n.Pos_Start,
		Pos_End:   n.Pos_End,
		Context:   n.Context,
	}
}
func (n List) IsTrue() bool                                              { return len(n.Elements) != 0 }
func (n List) Set_pos(posStart *tools.Position, posEnd *tools.Position) Value {
	return n.SetPos(posStart, posEnd)
}
func (n List) Set_context(ctx *Context) Value { return n.SetContext(ctx) }
func (n List) Get_context() *Context          { return n.GetContext() }
func (n List) Is_true() bool                  { return n.IsTrue() }

func (n List) Add(other Value) (Value, *Error) {
	return List{Elements: append(slices.Clone(n.Elements), other), Context: n.Context}, nil
}

func (n List) Subtract(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		idx := int(other.Value)
		if idx < 0 || idx >= len(n.Elements) {
			errNya := RTError(safePos(other.Pos_Start), safePos(other.Pos_End), "Elements at this index could not be removed from list because it's out of bounds", n.Context)
			return nil, &errNya
		}
		newElements := make([]Value, 0, len(n.Elements)-1)
		newElements = append(newElements, n.Elements[:idx]...)
		newElements = append(newElements, n.Elements[idx+1:]...)
		return List{Elements: newElements, Context: n.Context}, nil
	}
	err := RTError(safePos(n.Pos_Start), safePos(n.Pos_End), "Illegal operation: cannot perform subtraction on List with the given type", n.Context)
	return nil, &err
}

func (n List) Multiply(other Value) (Value, *Error) {
	switch other := other.(type) {
	case List:
		return List{Elements: append(slices.Clone(n.Elements), other.Elements...), Context: n.Context}, nil
	}
	err := RTError(safePos(n.Pos_Start), safePos(n.Pos_End), "Illegal operation: cannot perform multiplication on List with the given type", n.Context)
	return nil, &err
}

func (n List) Divide(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		idx := int(other.Value)
		if idx < 0 || idx >= len(n.Elements) {
			errNya := RTError(safePos(other.Pos_Start), safePos(other.Pos_End), "Elements at this index could not be retrieved from list because it's out of bounds", n.Context)
			return nil, &errNya
		}
		return n.Elements[idx], nil
	}
	err := RTError(safePos(n.Pos_Start), safePos(n.Pos_End), "Illegal operation: cannot perform division on List with the given type", n.Context)
	return nil, &err
}

// Forwarding aliases for List
func (n List) Added_to(other Value) (Value, *Error)  { return n.Add(other) }
func (n List) Subbed_to(other Value) (Value, *Error) { return n.Subtract(other) }
func (n List) Multed_by(other Value) (Value, *Error) { return n.Multiply(other) }
func (n List) Dived_by(other Value) (Value, *Error)  { return n.Divide(other) }

// String represents a string value.
type String struct {
	Value     string
	Context   *Context
	Pos_Start *tools.Position
	Pos_End   *tools.Position
}

func (s String) Print() string { return s.Value }
func (s String) SetPos(posStart *tools.Position, posEnd *tools.Position) Value {
	s.Pos_Start = posStart
	s.Pos_End = posEnd
	return s
}
func (s String) SetContext(ctx *Context) Value { s.Context = ctx; return s }
func (s String) GetContext() *Context          { return s.Context }
func (s String) Copy() Value {
	return String{
		Value:     s.Value,
		Pos_Start: s.Pos_Start,
		Pos_End:   s.Pos_End,
		Context:   s.Context,
	}
}
func (s String) IsTrue() bool                                              { return s.Value != "" }
func (s String) Set_pos(posStart *tools.Position, posEnd *tools.Position) Value {
	return s.SetPos(posStart, posEnd)
}
func (s String) Set_context(ctx *Context) Value { return s.SetContext(ctx) }
func (s String) Get_context() *Context          { return s.GetContext() }
func (s String) Is_true() bool                  { return s.IsTrue() }

func (s String) Add(other Value) (Value, *Error) {
	switch other := other.(type) {
	case String:
		return String{Value: s.Value + other.Value, Context: s.Context}, nil
	}
	err := RTError(safePos(s.Pos_Start), safePos(s.Pos_End), "Illegal operation: cannot concatenate String with the given type", s.Context)
	return nil, &err
}

func (s String) Multiply(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		return String{Value: strings.Repeat(s.Value, int(other.Value)), Context: s.Context}, nil
	}
	err := RTError(safePos(s.Pos_Start), safePos(s.Pos_End), "Illegal operation: cannot multiply String with the given type", s.Context)
	return nil, &err
}

func (s String) Equal(other Value) (Value, *Error) {
	switch other := other.(type) {
	case String:
		return Boolean{Value: s.Value == other.Value, Context: s.Context}.SetPos(s.Pos_Start, s.Pos_End), nil
	}
	err := RTError(safePos(s.Pos_Start), safePos(s.Pos_End), "Illegal operation: cannot compare equality with the given type", s.Context)
	return nil, &err
}

func (s String) NotEqual(other Value) (Value, *Error) {
	switch other := other.(type) {
	case String:
		return Boolean{Value: s.Value != other.Value, Context: s.Context}.SetPos(s.Pos_Start, s.Pos_End), nil
	}
	err := RTError(safePos(s.Pos_Start), safePos(s.Pos_End), "Illegal operation: cannot compare inequality with the given type", s.Context)
	return nil, &err
}

// Forwarding aliases for String
func (s String) Added_to(other Value) (Value, *Error)          { return s.Add(other) }
func (s String) Multed_by(other Value) (Value, *Error)         { return s.Multiply(other) }
func (s String) Get_comparison_eq(other Value) (Value, *Error) { return s.Equal(other) }
func (s String) Get_comparison_nq(other Value) (Value, *Error) { return s.NotEqual(other) }

// Number represents integer and real numerical values.
type Number struct {
	Value     float64
	IsReal    bool
	Context   *Context
	Pos_Start *tools.Position
	Pos_End   *tools.Position
}

func (n Number) Print() string { return PrintValueInterpreter(n) }
func (n Number) SetPos(posStart *tools.Position, posEnd *tools.Position) Value {
	n.Pos_Start = posStart
	n.Pos_End = posEnd
	return n
}
func (n Number) SetContext(ctx *Context) Value { n.Context = ctx; return n }
func (n Number) GetContext() *Context          { return n.Context }
func (n Number) Copy() Value {
	return Number{
		Value:     n.Value,
		IsReal:    n.IsReal,
		Pos_Start: n.Pos_Start,
		Pos_End:   n.Pos_End,
		Context:   n.Context,
	}
}
func (n Number) IsTrue() bool                                              { return n.Value != 0 }
func (n Number) Set_pos(posStart *tools.Position, posEnd *tools.Position) Value {
	return n.SetPos(posStart, posEnd)
}
func (n Number) Set_context(ctx *Context) Value { return n.SetContext(ctx) }
func (n Number) Get_context() *Context          { return n.GetContext() }
func (n Number) Is_true() bool                  { return n.IsTrue() }

func (n Number) Add(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		return Number{Value: n.Value + other.Value, IsReal: n.IsReal || other.IsReal, Context: n.Context}, nil
	}
	err := RTError(safePos(n.Pos_Start), safePos(n.Pos_End), "Illegal operation: cannot perform addition with the given type", n.Context)
	return nil, &err
}

func (n Number) Subtract(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		return Number{Value: n.Value - other.Value, IsReal: n.IsReal || other.IsReal, Context: n.Context}, nil
	}
	err := RTError(safePos(n.Pos_Start), safePos(n.Pos_End), "Illegal operation: cannot perform subtraction with the given type", n.Context)
	return nil, &err
}

func (n Number) Multiply(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		return Number{Value: n.Value * other.Value, IsReal: n.IsReal || other.IsReal, Context: n.Context}, nil
	}
	err := RTError(safePos(n.Pos_Start), safePos(n.Pos_End), "Illegal operation: cannot perform multiplication with the given type", n.Context)
	return nil, &err
}

func (n Number) Divide(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		if other.Value == 0 {
			err := RTError(safePos(other.Pos_Start), safePos(other.Pos_End), "Division by zero", n.Context)
			return other, &err
		}
		return Number{Value: n.Value / other.Value, IsReal: true, Context: n.Context}, nil
	}
	err := RTError(safePos(n.Pos_Start), safePos(n.Pos_End), "Illegal operation: cannot perform division with the given type", n.Context)
	return nil, &err
}

func (n Number) Div(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		if n.IsReal || other.IsReal {
			err := RTError(safePos(n.Pos_Start), safePos(other.Pos_End), "Illegal operation: 'div' can only be used with integers, use '/' for real division", n.Context)
			return nil, &err
		}
		if other.Value == 0 {
			err := RTError(safePos(other.Pos_Start), safePos(other.Pos_End), "Division by zero", n.Context)
			return other, &err
		}
		return Number{Value: math.Floor(n.Value / other.Value), IsReal: false, Context: n.Context}, nil
	}
	err := RTError(safePos(n.Pos_Start), safePos(n.Pos_End), "Illegal operation: cannot perform division with the given type", n.Context)
	return nil, &err
}

func (n Number) Mod(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		if other.Value == 0 {
			err := RTError(safePos(other.Pos_Start), safePos(other.Pos_End), "Division by zero", n.Context)
			return other, &err
		}
		return Number{Value: math.Mod(n.Value, other.Value), IsReal: n.IsReal || other.IsReal, Context: n.Context}, nil
	}
	err := RTError(safePos(n.Pos_Start), safePos(n.Pos_End), "Illegal operation: cannot perform modulo with the given type", n.Context)
	return nil, &err
}

func (n Number) Power(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		return Number{Value: math.Pow(n.Value, other.Value), IsReal: n.IsReal || other.IsReal, Context: n.Context}, nil
	}
	err := RTError(safePos(n.Pos_Start), safePos(n.Pos_End), "Illegal operation: cannot perform exponentiation with the given type", n.Context)
	return nil, &err
}

func (n Number) BitwiseAnd(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		if n.IsReal || other.IsReal {
			err := RTError(safePos(n.Pos_Start), safePos(other.Pos_End), "Illegal operation: bitwise '&' requires integer operands", n.Context)
			return nil, &err
		}
		return Number{Value: float64(int64(n.Value) & int64(other.Value)), IsReal: false, Context: n.Context}, nil
	}
	err := RTError(safePos(n.Pos_Start), safePos(n.Pos_End), "Illegal operation: cannot perform bitwise AND with the given type", n.Context)
	return nil, &err
}

func (n Number) BitwiseOr(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		if n.IsReal || other.IsReal {
			err := RTError(safePos(n.Pos_Start), safePos(other.Pos_End), "Illegal operation: bitwise '|' requires integer operands", n.Context)
			return nil, &err
		}
		return Number{Value: float64(int64(n.Value) | int64(other.Value)), IsReal: false, Context: n.Context}, nil
	}
	err := RTError(safePos(n.Pos_Start), safePos(n.Pos_End), "Illegal operation: cannot perform bitwise OR with the given type", n.Context)
	return nil, &err
}

func (n Number) BitwiseXor(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		if n.IsReal || other.IsReal {
			return Number{Value: math.Pow(n.Value, other.Value), IsReal: true, Context: n.Context}, nil
		}
		return Number{Value: float64(int64(n.Value) ^ int64(other.Value)), IsReal: false, Context: n.Context}, nil
	}
	err := RTError(safePos(n.Pos_Start), safePos(n.Pos_End), "Illegal operation: cannot perform bitwise XOR with the given type", n.Context)
	return nil, &err
}

func (n Number) BitwiseNot() (Value, *Error) {
	if n.IsReal {
		err := RTError(safePos(n.Pos_Start), safePos(n.Pos_End), "Illegal operation: bitwise '~' requires integer operand", n.Context)
		return nil, &err
	}
	return Number{Value: float64(^int64(n.Value)), IsReal: false, Context: n.Context}, nil
}

func (n Number) ShiftLeft(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		if n.IsReal || other.IsReal {
			err := RTError(safePos(n.Pos_Start), safePos(other.Pos_End), "Illegal operation: bitwise '<<' requires integer operands", n.Context)
			return nil, &err
		}
		if other.Value < 0 {
			err := RTError(safePos(other.Pos_Start), safePos(other.Pos_End), "Illegal operation: negative shift count", n.Context)
			return nil, &err
		}
		return Number{Value: float64(int64(n.Value) << uint64(other.Value)), IsReal: false, Context: n.Context}, nil
	}
	err := RTError(safePos(n.Pos_Start), safePos(n.Pos_End), "Illegal operation: cannot perform bitwise shift left with the given type", n.Context)
	return nil, &err
}

func (n Number) ShiftRight(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		if n.IsReal || other.IsReal {
			err := RTError(safePos(n.Pos_Start), safePos(other.Pos_End), "Illegal operation: bitwise '>>' requires integer operands", n.Context)
			return nil, &err
		}
		if other.Value < 0 {
			err := RTError(safePos(other.Pos_Start), safePos(other.Pos_End), "Illegal operation: negative shift count", n.Context)
			return nil, &err
		}
		return Number{Value: float64(int64(n.Value) >> uint64(other.Value)), IsReal: false, Context: n.Context}, nil
	}
	err := RTError(safePos(n.Pos_Start), safePos(n.Pos_End), "Illegal operation: cannot perform bitwise shift right with the given type", n.Context)
	return nil, &err
}

func (n Number) Equal(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		return Boolean{Value: n.Value == other.Value, Context: n.Context}.SetPos(n.Pos_Start, n.Pos_End), nil
	case Boolean:
		return Boolean{Value: (n.Value != 0) == other.Value, Context: n.Context}.SetPos(n.Pos_Start, n.Pos_End), nil
	}
	err := RTError(safePos(n.Pos_Start), safePos(n.Pos_End), "Illegal operation: cannot compare equality with the given type", n.Context)
	return nil, &err
}

func (n Number) NotEqual(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		return Boolean{Value: n.Value != other.Value, Context: n.Context}.SetPos(n.Pos_Start, n.Pos_End), nil
	case Boolean:
		return Boolean{Value: (n.Value != 0) != other.Value, Context: n.Context}.SetPos(n.Pos_Start, n.Pos_End), nil
	}
	err := RTError(safePos(n.Pos_Start), safePos(n.Pos_End), "Illegal operation: cannot compare inequality with the given type", n.Context)
	return nil, &err
}

func (n Number) Less(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		return Boolean{Value: n.Value < other.Value, Context: n.Context}.SetPos(n.Pos_Start, n.Pos_End), nil
	}
	err := RTError(safePos(n.Pos_Start), safePos(n.Pos_End), "Illegal operation: cannot compare with the given type", n.Context)
	return nil, &err
}

func (n Number) LessEqual(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		return Boolean{Value: n.Value <= other.Value, Context: n.Context}.SetPos(n.Pos_Start, n.Pos_End), nil
	}
	err := RTError(safePos(n.Pos_Start), safePos(n.Pos_End), "Illegal operation: cannot compare with the given type", n.Context)
	return nil, &err
}

func (n Number) Greater(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		return Boolean{Value: n.Value > other.Value, Context: n.Context}.SetPos(n.Pos_Start, n.Pos_End), nil
	}
	err := RTError(safePos(n.Pos_Start), safePos(n.Pos_End), "Illegal operation: cannot compare with the given type", n.Context)
	return nil, &err
}

func (n Number) GreaterEqual(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		return Boolean{Value: n.Value >= other.Value, Context: n.Context}.SetPos(n.Pos_Start, n.Pos_End), nil
	}
	err := RTError(safePos(n.Pos_Start), safePos(n.Pos_End), "Illegal operation: cannot compare with the given type", n.Context)
	return nil, &err
}

func (n Number) And(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		return Boolean{Value: n.Value != 0 && other.Value != 0, Context: n.Context}.SetPos(n.Pos_Start, n.Pos_End), nil
	case Boolean:
		return Boolean{Value: n.Value != 0 && other.Value, Context: n.Context}.SetPos(n.Pos_Start, n.Pos_End), nil
	}
	err := RTError(safePos(n.Pos_Start), safePos(n.Pos_End), "Illegal operation: cannot apply 'and' with the given type", n.Context)
	return nil, &err
}

func (n Number) Or(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		return Boolean{Value: n.Value != 0 || other.Value != 0, Context: n.Context}.SetPos(n.Pos_Start, n.Pos_End), nil
	case Boolean:
		return Boolean{Value: n.Value != 0 || other.Value, Context: n.Context}.SetPos(n.Pos_Start, n.Pos_End), nil
	}
	err := RTError(safePos(n.Pos_Start), safePos(n.Pos_End), "Illegal operation: cannot apply 'or' with the given type", n.Context)
	return nil, &err
}

func (n Number) Not() (Value, *Error) {
	return Boolean{Value: n.Value == 0, Context: n.Context}.SetPos(n.Pos_Start, n.Pos_End), nil
}

// Forwarding aliases for Number
func (n Number) Added_to(other Value) (Value, *Error)          { return n.Add(other) }
func (n Number) Subbed_by(other Value) (Value, *Error)         { return n.Subtract(other) }
func (n Number) Multed_by(other Value) (Value, *Error)         { return n.Multiply(other) }
func (n Number) Divided_by(other Value) (Value, *Error)        { return n.Divide(other) }
func (n Number) DivBy(other Value) (Value, *Error)             { return n.Div(other) }
func (n Number) ModBy(other Value) (Value, *Error)             { return n.Mod(other) }
func (n Number) Powered_by(other Value) (Value, *Error)        { return n.Power(other) }
func (n Number) Get_comparison_eq(other Value) (Value, *Error) { return n.Equal(other) }
func (n Number) Get_comparison_nq(other Value) (Value, *Error) { return n.NotEqual(other) }
func (n Number) Get_comparison_lt(other Value) (Value, *Error) { return n.Less(other) }
func (n Number) Get_comparison_lte(other Value) (Value, *Error) {
	return n.LessEqual(other)
}
func (n Number) Get_comparison_gt(other Value) (Value, *Error) { return n.Greater(other) }
func (n Number) Get_comparison_gte(other Value) (Value, *Error) {
	return n.GreaterEqual(other)
}
func (n Number) Anded_by(other Value) (Value, *Error) { return n.And(other) }
func (n Number) Ored_by(other Value) (Value, *Error)  { return n.Or(other) }
func (n Number) Notted() (Value, *Error)              { return n.Not() }

// BaseFunctionInterface is implemented by functions.
type BaseFunctionInterface interface {
	Value
	GenerateNewContext() Context
	CheckArgs(argNames []string, args []Value) Value
	PopulateArgs(argNames []string, args []Value, execCtx *Context)
	CheckAndPopulateArgs(argNames []string, args []Value, execCtx *Context) Value
	GetArgsName() []string
	GetBodyNode() Expr
	GetShouldAutoReturn() bool
}

// BaseFunction provides common fields and methods for user and builtin functions.
type BaseFunction struct {
	Name             string
	ArgNames         []string
	BodyNode         Expr
	ShouldAutoReturn bool
	Pos_Start        *tools.Position
	Pos_End          *tools.Position
	Context          *Context
}

func (n BaseFunction) Print() string { return fmt.Sprintf("<function %s>", n.Name) }
func (n BaseFunction) SetPos(posStart *tools.Position, posEnd *tools.Position) Value {
	n.Pos_Start = posStart
	n.Pos_End = posEnd
	return n
}
func (n BaseFunction) SetContext(ctx *Context) Value { n.Context = ctx; return n }
func (n BaseFunction) GetContext() *Context          { return n.Context }
func (n BaseFunction) Copy() Value {
	return BaseFunction{
		Name:             n.Name,
		ArgNames:         slices.Clone(n.ArgNames),
		BodyNode:         n.BodyNode,
		ShouldAutoReturn: n.ShouldAutoReturn,
		Pos_Start:        n.Pos_Start,
		Pos_End:          n.Pos_End,
		Context:          n.Context,
	}
}
func (n BaseFunction) IsTrue() bool                                              { return true }
func (n BaseFunction) Set_pos(posStart *tools.Position, posEnd *tools.Position) Value {
	return n.SetPos(posStart, posEnd)
}
func (n BaseFunction) Set_context(ctx *Context) Value { return n.SetContext(ctx) }
func (n BaseFunction) Get_context() *Context          { return n.GetContext() }
func (n BaseFunction) Is_true() bool                  { return n.IsTrue() }
func (n BaseFunction) GetArgsName() []string          { return n.ArgNames }
func (n BaseFunction) GetBodyNode() Expr              { return n.BodyNode }
func (n BaseFunction) GetShouldAutoReturn() bool      { return n.ShouldAutoReturn }

func (n BaseFunction) GenerateNewContext() Context {
	var parentSymbols *SymbolTable
	if n.Context != nil {
		parentSymbols = n.Context.GetSymbolTable()
	}
	return Context{
		DisplayName:    n.Name,
		Parent:         n.Context,
		ParentEntryPos: n.Pos_Start,
		SymbolTable:    NewSymbolTable(parentSymbols),
		Symbol_Table:   NewSymbolTable(parentSymbols),
	}
}

func (n BaseFunction) CheckArgs(argNames []string, args []Value) Value {
	res := &RTResult{}

	if len(args) > len(argNames) {
		hasSpread := false
		for i, v := range argNames {
			if i == len(argNames)-1 && strings.HasPrefix(v, "...") {
				hasSpread = true
			}
		}

		if !hasSpread {
			return res.Failure(RTError(safePos(n.Pos_Start), safePos(n.Pos_End), fmt.Sprintf("%d too many args passed into '%s'", len(args)-len(argNames), n.Name), n.Context))
		}
	}

	if len(args) < len(argNames) {
		return res.Failure(RTError(safePos(n.Pos_Start), safePos(n.Pos_End), fmt.Sprintf("%d too few args passed into '%s'", len(argNames)-len(args), n.Name), n.Context))
	}

	return res.Success(Null{})
}

func (n BaseFunction) PopulateArgs(argNames []string, args []Value, execCtx *Context) {
	st := execCtx.GetSymbolTable()
	for i := 0; i < len(argNames); i++ {
		if i == len(argNames)-1 && strings.HasPrefix(argNames[i], "...") {
			v := strings.ReplaceAll(argNames[i], "...", "")
			args[i].SetContext(execCtx)
			st.Set(v, List{
				Elements:  args[i:],
				Context:   execCtx,
				Pos_Start: n.Pos_Start,
				Pos_End:   n.Pos_End,
			})
			continue
		}

		argName := argNames[i]
		argValue := args[i]
		argValue.SetContext(execCtx)
		st.Set(argName, argValue)
	}
}

func (n BaseFunction) CheckAndPopulateArgs(argNames []string, args []Value, execCtx *Context) Value {
	res := &RTResult{}

	res.Register(n.CheckArgs(argNames, args))
	if res.ShouldReturn() {
		return res
	}

	n.PopulateArgs(argNames, args, execCtx)
	return res.Success(Null{})
}

// Function represents a user-defined DAP function.
type Function struct {
	BaseFunction
}

func (n Function) Print() string { return fmt.Sprintf("<function %s>", n.Name) }
func (n Function) SetPos(posStart *tools.Position, posEnd *tools.Position) Value {
	n.Pos_Start = posStart
	n.Pos_End = posEnd
	return n
}
func (n Function) SetContext(ctx *Context) Value { n.Context = ctx; return n }
func (n Function) GetContext() *Context          { return n.Context }
func (n Function) Copy() Value {
	return Function{
		BaseFunction: BaseFunction{
			ArgNames:         slices.Clone(n.ArgNames),
			BodyNode:         n.BodyNode,
			Name:             n.Name,
			ShouldAutoReturn: n.ShouldAutoReturn,
			Pos_Start:        n.Pos_Start,
			Pos_End:          n.Pos_End,
			Context:          n.Context,
		},
	}
}
func (n Function) IsTrue() bool                                              { return true }
func (n Function) Set_pos(posStart *tools.Position, posEnd *tools.Position) Value {
	return n.SetPos(posStart, posEnd)
}
func (n Function) Set_context(ctx *Context) Value { return n.SetContext(ctx) }
func (n Function) Get_context() *Context          { return n.GetContext() }
func (n Function) Is_true() bool                  { return n.IsTrue() }

// BuiltinHandler represents a function signature and implementation for built-ins.
type BuiltinHandler struct {
	Args    []string
	Execute func(ctx *Context, rawArgs []Expr, posStart, posEnd *tools.Position) Value
}

var builtinDispatchTable = map[string]BuiltinHandler{
	"Print": {
		Args: []string{"...value"},
		Execute: func(ctx *Context, rawArgs []Expr, posStart, posEnd *tools.Position) Value {
			res := &RTResult{}
			vals := ctx.GetSymbolTable().Get("value")
			var elements []Value
			if list, ok := vals.(List); ok {
				elements = list.Elements
			} else if arr, ok := vals.(Array); ok {
				elements = arr.Elements
			} else {
				elements = []Value{vals}
			}

			for _, v := range elements {
				fmt.Printf("%v\n", PrintValueInterpreter(v))
			}
			return res.Success(Null{})
		},
	},
	"PrintRet": {
		Args: []string{"value"},
		Execute: func(ctx *Context, rawArgs []Expr, posStart, posEnd *tools.Position) Value {
			res := &RTResult{}
			return res.Success(String{
				Value: fmt.Sprintf("%v", PrintValueInterpreter(ctx.GetSymbolTable().Get("value"))),
			})
		},
	},
	"Input": {
		Args: []string{"...value"},
		Execute: func(ctx *Context, rawArgs []Expr, posStart, posEnd *tools.Position) Value {
			res := &RTResult{}
			valList, ok := ctx.GetSymbolTable().Get("value").(List)
			if !ok {
				return res.Failure(RTError(safePos(posStart), safePos(posEnd), "Internal error: invalid input parameter", ctx))
			}
			for i := 0; i < len(valList.Elements); i++ {
				if i >= len(rawArgs) {
					break
				}
				switch rawArg := rawArgs[i].(type) {
				case VarAccessNode:
					var s string
					fmt.Scan(&s)

					varName := rawArg.VarNameTok.Value
					declaredTypeVal := ctx.GetSymbolTable().Get("VariableType " + varName)
					var declaredTypeStr string
					if declaredTypeVal != nil {
						if sVal, ok := declaredTypeVal.(String); ok {
							declaredTypeStr = sVal.Value
						}
					}

					if parseFloat, err := strconv.ParseFloat(s, 64); err == nil {
						isReal := strings.Contains(s, ".")
						if declaredTypeStr == "integer" {
							if isReal || parseFloat != math.Floor(parseFloat) {
								errNya := RTError(safePos(posStart), safePos(posEnd), fmt.Sprintf("Type mismatch: cannot assign real number to integer variable '%s'", varName), ctx)
								return res.Failure(errNya)
							}
						}
						ctx.GetSymbolTable().Set(varName, Number{Value: parseFloat, IsReal: isReal, Context: ctx})
						continue
					}

					if declaredTypeStr == "integer" {
						errNya := RTError(safePos(posStart), safePos(posEnd), fmt.Sprintf("Type mismatch: cannot assign string to integer variable '%s'", varName), ctx)
						return res.Failure(errNya)
					}

					ctx.GetSymbolTable().Set(varName, String{Value: s, Context: ctx})
				default:
					return res.Failure(RTError(safePos(posStart), safePos(posEnd), "input() can only read into a variable, not an expression or literal", ctx))
				}
			}
			return res.Success(Null{Context: ctx})
		},
	},
	"Panjang": {
		Args: []string{"value"},
		Execute: func(ctx *Context, rawArgs []Expr, posStart, posEnd *tools.Position) Value {
			res := &RTResult{}
			val := ctx.GetSymbolTable().Get("value")
			if s, ok := val.(String); ok {
				return res.Success(Number{Value: float64(len(s.Value)), IsReal: false, Context: ctx})
			}
			if arr, ok := val.(Array); ok {
				return res.Success(Number{Value: float64(len(arr.Elements)), IsReal: false, Context: ctx})
			}
			if list, ok := val.(List); ok {
				return res.Success(Number{Value: float64(len(list.Elements)), IsReal: false, Context: ctx})
			}

			errNya := RTError(safePos(posStart), safePos(posEnd), "Argument to panjang must be a string, array, or list", ctx)
			return res.Failure(errNya)
		},
	},
	"Bulat": {
		Args: []string{"value"},
		Execute: func(ctx *Context, rawArgs []Expr, posStart, posEnd *tools.Position) Value {
			res := &RTResult{}
			val := ctx.GetSymbolTable().Get("value")
			if num, ok := val.(Number); ok {
				rounded := math.Floor(num.Value + 0.5)
				return res.Success(Number{Value: rounded, IsReal: false, Context: ctx})
			}
			errNya := RTError(safePos(posStart), safePos(posEnd), "Argument to Bulat must be a number", ctx)
			return res.Failure(errNya)
		},
	},
	"IsNumber": {
		Args: []string{"value"},
		Execute: func(ctx *Context, rawArgs []Expr, posStart, posEnd *tools.Position) Value {
			res := &RTResult{}
			if _, ok := ctx.GetSymbolTable().Get("value").(Number); ok {
				return res.Success(Number{Value: 1, Context: ctx})
			}
			return res.Success(Number{Value: 0, Context: ctx})
		},
	},
	"IsString": {
		Args: []string{"value"},
		Execute: func(ctx *Context, rawArgs []Expr, posStart, posEnd *tools.Position) Value {
			res := &RTResult{}
			if _, ok := ctx.GetSymbolTable().Get("value").(String); ok {
				return res.Success(Number{Value: 1, Context: ctx})
			}
			return res.Success(Number{Value: 0, Context: ctx})
		},
	},
	"Append": {
		Args: []string{"list", "value"},
		Execute: func(ctx *Context, rawArgs []Expr, posStart, posEnd *tools.Position) Value {
			res := &RTResult{}
			list, ok := ctx.GetSymbolTable().Get("list").(List)
			value := ctx.GetSymbolTable().Get("value")
			if !ok {
				return res.Failure(RTError(safePos(posStart), safePos(posEnd), "First Argument must be list", ctx))
			}
			list.Elements = append(list.Elements, value)
			ctx.GetSymbolTable().Set("list", list)
			return res.Success(Number{Value: 0, Context: ctx})
		},
	},
	"InputBaris": {
		Args: []string{"...value"},
		Execute: func(ctx *Context, rawArgs []Expr, posStart, posEnd *tools.Position) Value {
			res := &RTResult{}
			reader := getStdinReader()
			valList, ok := ctx.GetSymbolTable().Get("value").(List)
			if !ok {
				return res.Failure(RTError(safePos(posStart), safePos(posEnd), "Internal error: invalid input parameter", ctx))
			}

			for i := 0; i < len(valList.Elements); i++ {
				if i >= len(rawArgs) {
					break
				}
				switch rawArg := rawArgs[i].(type) {
				case VarAccessNode:
					line, err := reader.ReadString('\n')
					if err != nil && err != io.EOF {
						errNya := RTError(safePos(posStart), safePos(posEnd), fmt.Sprintf("Error reading input: %v", err), ctx)
						return res.Failure(errNya)
					}
					line = strings.TrimRight(line, "\r\n")

					varName := rawArg.VarNameTok.Value
					declaredTypeVal := ctx.GetSymbolTable().Get("VariableType " + varName)
					var declaredTypeStr string
					if declaredTypeVal != nil {
						if sVal, ok := declaredTypeVal.(String); ok {
							declaredTypeStr = sVal.Value
						}
					}

					if declaredTypeStr == "integer" {
						if parseFloat, err := strconv.ParseFloat(line, 64); err == nil {
							isReal := strings.Contains(line, ".")
							if isReal || parseFloat != math.Floor(parseFloat) {
								errNya := RTError(safePos(posStart), safePos(posEnd), fmt.Sprintf("Type mismatch: cannot assign real number to integer variable '%s'", varName), ctx)
								return res.Failure(errNya)
							}
							ctx.GetSymbolTable().Set(varName, Number{Value: parseFloat, IsReal: isReal, Context: ctx})
							continue
						}
						errNya := RTError(safePos(posStart), safePos(posEnd), fmt.Sprintf("Type mismatch: cannot assign string to integer variable '%s'", varName), ctx)
						return res.Failure(errNya)
					}

					ctx.GetSymbolTable().Set(varName, String{Value: line, Context: ctx})
				default:
					return res.Failure(RTError(safePos(posStart), safePos(posEnd), "inputBaris() can only read into a variable, not an expression or literal", ctx))
				}
			}
			return res.Success(Null{Context: ctx})
		},
	},
	"KonversiKeString": {
		Args: []string{"value"},
		Execute: func(ctx *Context, rawArgs []Expr, posStart, posEnd *tools.Position) Value {
			res := &RTResult{}
			val := ctx.GetSymbolTable().Get("value")
			return res.Success(String{
				Value:   PrintValueInterpreter(val),
				Context: ctx,
			})
		},
	},
}

var (
	lastStdin   *os.File
	stdinReader *bufio.Reader
)

func getStdinReader() *bufio.Reader {
	if os.Stdin != lastStdin {
		lastStdin = os.Stdin
		stdinReader = bufio.NewReader(os.Stdin)
	}
	return stdinReader
}

// BuiltInFunction represents a built-in DAP function.
type BuiltInFunction struct {
	BaseFunction
}

func (n BuiltInFunction) Print() string { return fmt.Sprintf("<function %s>", n.Name) }
func (n BuiltInFunction) SetPos(posStart *tools.Position, posEnd *tools.Position) Value {
	n.Pos_Start = posStart
	n.Pos_End = posEnd
	return n
}
func (n BuiltInFunction) SetContext(ctx *Context) Value { n.Context = ctx; return n }
func (n BuiltInFunction) GetContext() *Context          { return n.Context }
func (n BuiltInFunction) Copy() Value {
	return BuiltInFunction{
		BaseFunction: BaseFunction{
			Name:      n.Name,
			Pos_Start: n.Pos_Start,
			Pos_End:   n.Pos_End,
			Context:   n.Context,
		},
	}
}
func (n BuiltInFunction) IsTrue() bool                                              { return true }
func (n BuiltInFunction) Set_pos(posStart *tools.Position, posEnd *tools.Position) Value {
	return n.SetPos(posStart, posEnd)
}
func (n BuiltInFunction) Set_context(ctx *Context) Value { return n.SetContext(ctx) }
func (n BuiltInFunction) Get_context() *Context          { return n.GetContext() }
func (n BuiltInFunction) Is_true() bool                  { return n.IsTrue() }

// Execute executes the built-in function without dynamic reflection.
func (n BuiltInFunction) Execute(args []Value, rawArgs []Expr) Value {
	res := &RTResult{}
	handler, exists := builtinDispatchTable[n.Name]
	if !exists {
		return res.Failure(RTError(safePos(n.Pos_Start), safePos(n.Pos_End), fmt.Sprintf("Built-in function '%s' is not defined", n.Name), n.Context))
	}

	execCtx := n.GenerateNewContext()
	res.Register(n.CheckAndPopulateArgs(handler.Args, args, &execCtx))
	if res.ShouldReturn() {
		return res
	}

	hasil := res.Register(handler.Execute(&execCtx, rawArgs, n.Pos_Start, n.Pos_End))
	if res.ShouldReturn() {
		return res
	}

	if hasil != nil && hasil.GetContext() != nil {
		n.Context = hasil.GetContext()
	}
	return res.Success(hasil)
}

// Struct represents a structured record with named fields.
type Struct struct {
	Fields    map[string]Value
	Context   *Context
	Pos_Start *tools.Position
	Pos_End   *tools.Position
}

func (n Struct) Print() string { return PrintValueInterpreter(n) }
func (n Struct) SetPos(posStart *tools.Position, posEnd *tools.Position) Value {
	n.Pos_Start = posStart
	n.Pos_End = posEnd
	return n
}
func (n Struct) SetContext(ctx *Context) Value { n.Context = ctx; return n }
func (n Struct) GetContext() *Context          { return n.Context }
func (n Struct) Copy() Value {
	fieldsCopy := make(map[string]Value, len(n.Fields))
	for k, v := range n.Fields {
		fieldsCopy[k] = v.Copy()
	}
	return Struct{
		Fields:    fieldsCopy,
		Context:   n.Context,
		Pos_Start: n.Pos_Start,
		Pos_End:   n.Pos_End,
	}
}
func (n Struct) IsTrue() bool                                              { return true }
func (n Struct) Set_pos(posStart *tools.Position, posEnd *tools.Position) Value {
	return n.SetPos(posStart, posEnd)
}
func (n Struct) Set_context(ctx *Context) Value { return n.SetContext(ctx) }
func (n Struct) Get_context() *Context          { return n.GetContext() }
func (n Struct) Is_true() bool                  { return n.IsTrue() }

// Type represents a user-defined type definition.
type Type struct {
	Definition Expr
	Context    *Context
	Pos_Start  *tools.Position
	Pos_End    *tools.Position
}

func (n Type) Print() string { return "<type>" }
func (n Type) SetPos(posStart *tools.Position, posEnd *tools.Position) Value {
	n.Pos_Start = posStart
	n.Pos_End = posEnd
	return n
}
func (n Type) SetContext(ctx *Context) Value { n.Context = ctx; return n }
func (n Type) GetContext() *Context          { return n.Context }
func (n Type) Copy() Value {
	return Type{
		Definition: n.Definition,
		Context:    n.Context,
		Pos_Start:  n.Pos_Start,
		Pos_End:    n.Pos_End,
	}
}
func (n Type) IsTrue() bool                                              { return true }
func (n Type) Set_pos(posStart *tools.Position, posEnd *tools.Position) Value {
	return n.SetPos(posStart, posEnd)
}
func (n Type) Set_context(ctx *Context) Value { return n.SetContext(ctx) }
func (n Type) Get_context() *Context          { return n.GetContext() }
func (n Type) Is_true() bool                  { return n.IsTrue() }

// RTResult wraps the result of evaluating an AST node at runtime.
type RTResult struct {
	Value              Value
	FuncReturnValue    Value
	LoopShouldContinue bool
	LoopShouldBreak    bool
	Error              *Error
}

func (n *RTResult) reset() {
	n.Value = nil
	n.FuncReturnValue = nil
	n.LoopShouldContinue = false
	n.LoopShouldBreak = false
	n.Error = nil
}

func (n *RTResult) Print() string {
	fmt.Printf("%s\n", PrintValueInterpreter(n.Value))
	return ""
}
func (n *RTResult) SetPos(posStart *tools.Position, posEnd *tools.Position) Value { return n }
func (n *RTResult) SetContext(ctx *Context) Value                                 { return n }
func (n *RTResult) GetContext() *Context                                          { return nil }
func (n *RTResult) Copy() Value                                                   { return n }
func (n *RTResult) IsTrue() bool                                                  { return false }
func (n *RTResult) Set_pos(posStart *tools.Position, posEnd *tools.Position) Value {
	return n.SetPos(posStart, posEnd)
}
func (n *RTResult) Set_context(ctx *Context) Value { return n.SetContext(ctx) }
func (n *RTResult) Get_context() *Context          { return n.GetContext() }
func (n *RTResult) Is_true() bool                  { return n.IsTrue() }

// Register registers results from a child evaluation step.
func (rtResult *RTResult) Register(res Value) Value {
	if r, ok := res.(*RTResult); ok {
		if r.ShouldReturn() {
			rtResult.Error = r.Error
		}
		rtResult.FuncReturnValue = r.FuncReturnValue
		rtResult.LoopShouldContinue = r.LoopShouldContinue
		rtResult.LoopShouldBreak = r.LoopShouldBreak
		return r.Value
	}
	return res
}

// Success returns a successful result with a value.
func (rtResult *RTResult) Success(value Value) Value {
	rtResult.reset()
	rtResult.Value = value
	return rtResult
}

// SuccessReturn returns a function return value.
func (rtResult *RTResult) SuccessReturn(value Value) Value {
	rtResult.reset()
	rtResult.FuncReturnValue = value
	return rtResult
}

// Success_Return is an alias for SuccessReturn.
func (rtResult *RTResult) Success_Return(value Value) Value {
	return rtResult.SuccessReturn(value)
}

// SuccessContinue sets the loop continue flag.
func (rtResult *RTResult) SuccessContinue() Value {
	rtResult.reset()
	rtResult.LoopShouldContinue = true
	return rtResult
}

// Success_Continue is an alias for SuccessContinue.
func (rtResult *RTResult) Success_Continue() Value {
	return rtResult.SuccessContinue()
}

// SuccessBreak sets the loop break flag.
func (rtResult *RTResult) SuccessBreak() Value {
	rtResult.reset()
	rtResult.LoopShouldBreak = true
	return rtResult
}

// Success_Break is an alias for SuccessBreak.
func (rtResult *RTResult) Success_Break() Value {
	return rtResult.SuccessBreak()
}

// Failure attaches a runtime error.
func (rtResult *RTResult) Failure(err Error) Value {
	rtResult.reset()
	errCopy := err
	rtResult.Error = &errCopy
	return rtResult
}

// ShouldReturn reports whether execution should unwind early.
func (rtResult *RTResult) ShouldReturn() bool {
	return rtResult.Error != nil || rtResult.FuncReturnValue != nil || rtResult.LoopShouldBreak || rtResult.LoopShouldContinue
}
