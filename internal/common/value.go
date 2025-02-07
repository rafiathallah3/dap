package common

import (
	"dap/tools"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

type SymbolTable struct {
	Symbols map[string]Value
	Parent  *SymbolTable
}

func (symbolTable *SymbolTable) Get(name string) Value {
	if val, ok := symbolTable.Symbols[name]; ok {
		return val
	}

	if symbolTable.Parent != nil {
		return symbolTable.Parent.Get(name)
	}

	return nil
}

func (symbolTable *SymbolTable) Set(name string, value Value) {
	symbolTable.Symbols[name] = value
}

func (symbolTable *SymbolTable) Remove(name string) {
	delete(symbolTable.Symbols, name)
}

type Context struct {
	DisplayName    string
	Parent         *Context
	ParentEntryPos *tools.Position
	Symbol_Table   *SymbolTable
}

func PrintValueInterpreter(n Value) string {
	switch n := n.(type) {
	case Number:
		return strconv.FormatFloat(n.Value, 'f', -1, 64)
	case String:
		return n.Value
	case List:
		if len(n.Elements) == 1 {
			return PrintValueInterpreter(n.Elements[0])
		}

		hasil := "["
		for _, v := range n.Elements {
			hasil += PrintValueInterpreter(v) + ", "
		}
		hasil = hasil[:len(hasil)-2]
		hasil += "]"
		return hasil
	case Function:
		return fmt.Sprintf("<function %s>", n.Name)
	case Null:
		return "NULL"
	case *RTResult:
		return PrintValueInterpreter(n.Value)
	}

	return ""
}

type Value interface {
	Set_pos(Pos_Start *tools.Position, Pos_End *tools.Position) Value
	Set_context(context *Context) Value
	Get_context() *Context
	Is_true() bool
	Copy() Value
	Print() string
}

type Null struct {
	Context *Context
}

func (n Null) Print() string {
	return "NULL"
}
func (n Null) Set_pos(Pos_Start *tools.Position, Pos_End *tools.Position) Value {
	return n
}
func (n Null) Set_context(context *Context) Value {
	return n
}
func (n Null) Get_context() *Context {
	return n.Context
}
func (n Null) Copy() Value {
	copy := Null{}
	return copy
}
func (n Null) Is_true() bool {
	return false
}

type List struct {
	Elements  []Value
	Context   *Context
	Pos_Start *tools.Position
	Pos_End   *tools.Position
}

func (n List) Print() string {
	return PrintValueInterpreter(n)
}

func (n List) Set_pos(Pos_Start *tools.Position, Pos_End *tools.Position) Value {
	n.Pos_Start = Pos_Start
	n.Pos_End = Pos_End
	return n
}

func (n List) Set_context(context *Context) Value {
	n.Context = context
	return n
}

func (n List) Get_context() *Context {
	return n.Context
}

func (n List) Copy() Value {
	copy := List{}
	copy.Elements = n.Elements
	copy.Pos_Start = n.Pos_Start
	copy.Pos_End = n.Pos_End
	copy.Context = n.Context
	return copy
}

func (n List) Is_true() bool {
	return len(n.Elements) != 0
}

func (n List) Added_to(other Value) (Value, *Error) {
	return List{Elements: append(n.Elements, other)}.Set_context(n.Context), nil
}

func (n List) Subbed_to(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		angka := int(other.Value)

		if angka > len(n.Elements) {
			errorNya := RTError(*other.Pos_Start, *other.Pos_End, "Elements at this index could not be removed from list because it's out of bounds", n.Context)
			return nil, &errorNya
		}

		return List{Elements: append(n.Elements[:angka], n.Elements[angka+1:]...)}.Set_context(n.Context), nil
	}

	panic("TIDAK BISA")
}

func (n List) Multed_by(other Value) (Value, *Error) {
	switch other := other.(type) {
	case List:
		return List{Elements: append(n.Elements, other.Elements...)}.Set_context(n.Context), nil
	}

	panic("TIDAK BISA")
}

func (n List) Dived_by(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		angka := int(other.Value)

		if angka > len(n.Elements) {
			errorNya := RTError(*other.Pos_Start, *other.Pos_End, "Elements at this index could not be retrieved from list because it's out of bounds", n.Context)
			return nil, &errorNya
		}

		return n.Elements[angka], nil
	}

	panic("TIDAK BISA")
}

type String struct {
	Value     string
	Context   *Context
	Pos_Start *tools.Position
	Pos_End   *tools.Position
}

func (n String) Print() string {
	return n.Value
}

func (n String) Set_pos(Pos_Start *tools.Position, Pos_End *tools.Position) Value {
	n.Pos_Start = Pos_Start
	n.Pos_End = Pos_End
	return n
}

func (n String) Set_context(context *Context) Value {
	n.Context = context
	return n
}

func (n String) Get_context() *Context {
	return n.Context
}

func (n String) Copy() Value {
	copy := String{}
	copy.Value = n.Value
	copy.Set_pos(n.Pos_Start, n.Pos_End)
	copy.Set_context(n.Context)
	return copy
}

func (n String) Is_true() bool {
	return n.Value != ""
}

func (s String) Added_to(other Value) (Value, *Error) {
	switch other := other.(type) {
	case String:
		return String{Value: s.Value + other.Value}.Set_context(s.Context), nil
	}

	panic("TIDAK BISA")
}

func (s String) Multed_by(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		return String{Value: strings.Repeat(s.Value, int(other.Value))}.Set_context(s.Context), nil
	}

	panic("TIDAK BISA")
}

type Number struct {
	Value     float64
	Context   *Context
	Pos_Start *tools.Position
	Pos_End   *tools.Position
}

func (n Number) Print() string {
	return fmt.Sprintf("%f", n.Value)
}

func (n Number) Set_pos(Pos_Start *tools.Position, Pos_End *tools.Position) Value {
	n.Pos_Start = Pos_Start
	n.Pos_End = Pos_End
	return n
}

func (n Number) Set_context(context *Context) Value {
	n.Context = context
	return n
}

func (n Number) Get_context() *Context {
	return n.Context
}

func (n Number) Copy() Value {
	copy := Number{}
	copy.Value = n.Value
	copy.Set_pos(n.Pos_Start, n.Pos_End)
	copy.Set_context(n.Context)
	return copy
}

func (n Number) Added_to(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		return Number{Value: n.Value + other.Value}.Set_context(n.Context), nil
	}

	panic("TIDAK BISA")
}

func (n Number) Subbed_by(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		return Number{Value: n.Value - other.Value}.Set_context(n.Context), nil
	}

	panic("TIDAK BISA")
}

func (n Number) Multed_by(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		return Number{Value: n.Value * other.Value}.Set_context(n.Context), nil
	}

	panic("TIDAK BISA")
}

func (n Number) Divided_by(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		if other.Value == 0 {
			err := RTError(*other.Pos_Start, *other.Pos_End, "Division by zero", n.Context)
			return other, &err
		}
		return Number{Value: n.Value / other.Value}.Set_context(n.Context), nil
	}

	panic("TIDAK BISA")
}

func (n Number) Powered_by(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		return Number{Value: math.Pow(n.Value, other.Value)}.Set_context(n.Context), nil
	}

	panic("TIDAK BISA")
}

func (n Number) Get_comparison_eq(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		return Number{Value: float64(tools.GetComparison(n.Value == other.Value))}.Set_context(n.Context), nil
	}

	panic("TIDAK BISA")
}

func (n Number) Get_comparison_nq(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		return Number{Value: float64(tools.GetComparison(n.Value != other.Value))}.Set_context(n.Context), nil
	}

	panic("TIDAK BISA")
}

func (n Number) Get_comparison_lt(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		return Number{Value: float64(tools.GetComparison(n.Value < other.Value))}.Set_context(n.Context), nil
	}

	panic("TIDAK BISA")
}

func (n Number) Get_comparison_lte(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		return Number{Value: float64(tools.GetComparison(n.Value <= other.Value))}.Set_context(n.Context), nil
	}

	panic("TIDAK BISA")
}

func (n Number) Get_comparison_gt(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		return Number{Value: float64(tools.GetComparison(n.Value > other.Value))}.Set_context(n.Context), nil
	}

	panic("TIDAK BISA")
}

func (n Number) Get_comparison_gte(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		return Number{Value: float64(tools.GetComparison(n.Value >= other.Value))}.Set_context(n.Context), nil
	}

	panic("TIDAK BISA")
}

func (n Number) Anded_by(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		return Number{Value: float64(tools.GetComparison(n.Value != 0 && other.Value != 0))}.Set_context(n.Context), nil
	}

	panic("TIDAK BISA")
}

func (n Number) Ored_by(other Value) (Value, *Error) {
	switch other := other.(type) {
	case Number:
		return Number{Value: float64(tools.GetComparison(n.Value != 0 || other.Value != 0))}.Set_context(n.Context), nil
	}

	panic("TIDAK BISA")
}

func (n Number) Notted() (Value, *Error) {
	hasil := 0
	if n.Value == 0 {
		hasil = 1
	}

	return Number{Value: float64(hasil)}.Set_context(n.Context), nil
}

func (n Number) Is_true() bool {
	return n.Value != 0
}

type BaseFunctionInterface interface {
	Value
	GenerateNewContext() Context
	CheckArgs(ArgNames []string, args []Value) Value
	PopulateArgs(ArgNames []string, args []Value, exec_ctx *Context)
	CheckAndPopulateArgs(ArgNames []string, args []Value, exec_ctx *Context) Value
	GetArgsName() []string
	GetBodyNode() Expr
	GetShouldAutoReturn() bool
}

type BaseFunction struct {
	Name             string
	ArgNames         []string
	BodyNode         Expr
	ShouldAutoReturn bool
	Pos_Start        *tools.Position
	Pos_End          *tools.Position
	Context          *Context
}

func (n BaseFunction) Print() string {
	return fmt.Sprintf("<function %s>", n.Name)
}
func (n BaseFunction) Set_pos(Pos_Start *tools.Position, Pos_End *tools.Position) Value {
	n.Pos_Start = Pos_Start
	n.Pos_End = Pos_End
	return n
}
func (n BaseFunction) Set_context(context *Context) Value {
	n.Context = context
	return n
}
func (n BaseFunction) Get_context() *Context {
	return n.Context
}
func (n BaseFunction) Copy() Value {
	copy := BaseFunction{}
	copy.Name = n.Name
	copy.ArgNames = n.ArgNames
	copy.BodyNode = n.BodyNode
	copy.ShouldAutoReturn = n.ShouldAutoReturn
	copy.Pos_Start = n.Pos_Start
	copy.Pos_End = n.Pos_End
	copy.Context = n.Context
	return copy
}
func (n BaseFunction) Is_true() bool {
	return true
}

func (n BaseFunction) GetArgsName() []string {
	return n.ArgNames
}

func (n BaseFunction) GetBodyNode() Expr {
	return n.BodyNode
}

func (n BaseFunction) GetShouldAutoReturn() bool {
	return n.ShouldAutoReturn
}

func (n BaseFunction) GenerateNewContext() Context {
	newContext := Context{
		DisplayName:    n.Name,
		Parent:         n.Context,
		ParentEntryPos: n.Pos_Start,
		Symbol_Table: &SymbolTable{
			Symbols: map[string]Value{},
			Parent:  n.Context.Symbol_Table,
		},
	}

	return newContext
}

func (n BaseFunction) CheckArgs(ArgNames []string, args []Value) Value {
	res := &RTResult{}

	if len(args) > len(ArgNames) {
		adaSpread := false
		for i, v := range ArgNames {
			if i == len(ArgNames)-1 && v[:3] == "..." {
				adaSpread = true
			}
		}

		if !adaSpread {
			return res.Failure(RTError(*n.Pos_Start, *n.Pos_End, fmt.Sprintf("%d too many args passed into '%s'", len(args)-len(ArgNames), n.Name), n.Context))
		}
	}

	if len(args) < len(ArgNames) {
		return res.Failure(RTError(*n.Pos_Start, *n.Pos_End, fmt.Sprintf("%d too few args passed into '%s'", len(ArgNames)-len(args), n.Name), n.Context))
	}

	return res.Success(Null{})
}

func (n BaseFunction) PopulateArgs(ArgNames []string, args []Value, exec_ctx *Context) {
	for i := 0; i < len(ArgNames); i++ {
		if i == len(ArgNames)-1 && ArgNames[i][:3] == "..." {
			v := strings.ReplaceAll(ArgNames[i], "...", "")

			args[i].Set_context(exec_ctx)
			exec_ctx.Symbol_Table.Set(v, List{
				Elements:  args[i:],
				Context:   exec_ctx,
				Pos_Start: n.Pos_Start,
				Pos_End:   n.Pos_End,
			})

			continue
		}

		arg_name := ArgNames[i]
		arg_value := args[i]

		arg_value.Set_context(exec_ctx)
		exec_ctx.Symbol_Table.Set(arg_name, arg_value)
	}
}

func (n BaseFunction) CheckAndPopulateArgs(ArgNames []string, args []Value, exec_ctx *Context) Value {
	res := &RTResult{}

	res.Register(n.CheckArgs(ArgNames, args))
	if res.ShouldReturn() {
		return res
	}

	n.PopulateArgs(ArgNames, args, exec_ctx)
	return res.Success(Null{})
}

type Function struct {
	BaseFunction
}

func (n Function) Print() string {
	return fmt.Sprintf("<function %s>", n.Name)
}
func (n Function) Set_pos(Pos_Start *tools.Position, Pos_End *tools.Position) Value {
	n.Pos_Start = Pos_Start
	n.Pos_End = Pos_End
	return n
}
func (n Function) Set_context(context *Context) Value {
	n.Context = context
	return n
}
func (n Function) Get_context() *Context {
	return n.Context
}
func (n Function) Copy() Value {
	copy := Function{}
	copy.ArgNames = n.ArgNames
	copy.BodyNode = n.BodyNode
	copy.Name = n.Name
	copy.Pos_Start = n.Pos_Start
	copy.Pos_End = n.Pos_End
	copy.Context = n.Context

	return copy
}
func (n Function) Is_true() bool {
	return true
}

type BuiltInFunction struct {
	BaseFunction
}

func (n BuiltInFunction) Print() string {
	return fmt.Sprintf("<function %s>", n.Name)
}
func (n BuiltInFunction) Set_pos(Pos_Start *tools.Position, Pos_End *tools.Position) Value {
	n.Pos_Start = Pos_Start
	n.Pos_End = Pos_End
	return n
}
func (n BuiltInFunction) Set_context(context *Context) Value {
	n.Context = context
	return n
}
func (n BuiltInFunction) Get_context() *Context {
	return n.Context
}
func (n BuiltInFunction) Copy() Value {
	copy := BuiltInFunction{}
	copy.Name = n.Name
	copy.Pos_Start = n.Pos_Start
	copy.Pos_End = n.Pos_End
	copy.Context = n.Context

	return copy
}
func (n BuiltInFunction) Is_true() bool {
	return true
}

func (n BuiltInFunction) Execute(args []Value, rawArgs []Expr) Value {
	res := &RTResult{}
	execCtx := n.GenerateNewContext()

	methodName := "Execute" + n.Name
	method := reflect.ValueOf(n).MethodByName(methodName)

	if !method.IsValid() {
		panic("BUILT IN FUNCTION EXEcUTION IS NOT VALID")
	}

	results := method.Call([]reflect.Value{})

	if len(results) == 0 {
		panic("NO RESULTS FOR THE BUILT IN FUNCTION")
	}

	BuiltinArgs := make([]string, results[0].Len())
	for i := 0; i < results[0].Len(); i++ {
		elem := results[0].Index(i)
		BuiltinArgs[i] = elem.String()
	}

	res.Register(n.CheckAndPopulateArgs(BuiltinArgs, args, &execCtx))
	if res.ShouldReturn() {
		return res
	}

	methodCall := results[1].Call([]reflect.Value{reflect.ValueOf(&execCtx), reflect.ValueOf(rawArgs)})
	hasil := res.Register(methodCall[0].Interface().(Value))
	if res.ShouldReturn() {
		return res
	}

	if hasil.Get_context() != nil {
		n.Context = hasil.Get_context()
	}
	return res.Success(hasil)
}

func (n BuiltInFunction) ExecutePrint() ([]string, func(*Context, []Expr) Value) {
	return []string{"...value"}, func(ctx *Context, rawArgs []Expr) Value {
		res := &RTResult{}

		for _, v := range ctx.Symbol_Table.Get("value").(List).Elements {
			fmt.Printf("%v\n", PrintValueInterpreter(v))
		}

		// fmt.Printf("%v\n", PrintValueInterpreter(ctx.Symbol_Table.Get("value")))
		return res.Success(Null{})
	}
}

func (n BuiltInFunction) ExecutePrintRet() ([]string, func(*Context, []Expr) Value) {
	return []string{"value"}, func(ctx *Context, rawArgs []Expr) Value {
		res := &RTResult{}
		return res.Success(String{
			Value: fmt.Sprintf("%v", PrintValueInterpreter(ctx.Symbol_Table.Get("value"))),
		})
	}
}

func (n BuiltInFunction) ExecuteInput() ([]string, func(*Context, []Expr) Value) {
	return []string{"...value"}, func(ctx *Context, rawArgs []Expr) Value {
		res := &RTResult{}

		//Need to improve this 20/01/2024 [23:41]
		// Gonna improve it 24/01/2024 [14:07]
		// var s string
		// fmt.Scan(&s)

		// return res.Success(String{
		// 	Value: s,
		// })

		for i := 0; i < len(ctx.Symbol_Table.Get("value").(List).Elements); i++ {
			switch rawArgs := rawArgs[i].(type) {
			case VarAccessNode:
				var s string
				fmt.Scan(&s)

				if parseFloat, err := strconv.ParseFloat(s, 64); err == nil {
					ctx.Symbol_Table.Set(rawArgs.VarNameTok.Value, Number{Value: parseFloat, Context: ctx})
					continue
				}

				ctx.Symbol_Table.Set(rawArgs.VarNameTok.Value, String{Value: s})
			default:
				return res.Failure(RTError(*n.Pos_Start, *n.Pos_End, "Parameter must be a Variable access node", ctx))
			}
		}

		return res.Success(Null{Context: ctx})
	}
}

func (n BuiltInFunction) ExecuteIsNumber() ([]string, func(*Context, []Expr) Value) {
	return []string{"value"}, func(ctx *Context, rawArgs []Expr) Value {
		res := &RTResult{}
		_, apakahNumber := ctx.Symbol_Table.Get("value").(Number)
		if apakahNumber {
			return res.Success(Number{Value: 1})
		}
		return res.Success(Number{Value: 0})
	}
}

func (n BuiltInFunction) ExecuteIsString() ([]string, func(*Context, []Expr) Value) {
	return []string{"value"}, func(ctx *Context, rawArgs []Expr) Value {
		res := &RTResult{}
		_, apakahString := ctx.Symbol_Table.Get("value").(String)
		if apakahString {
			return res.Success(Number{Value: 1})
		}
		return res.Success(Number{Value: 0})
	}
}

func (n BuiltInFunction) ExecuteAppend() ([]string, func(*Context, []Expr) Value) {
	return []string{"list", "value"}, func(ctx *Context, rawArgs []Expr) Value {
		res := &RTResult{}

		list, apakahList := ctx.Symbol_Table.Get("list").(List)
		value := ctx.Symbol_Table.Get("value")

		if !apakahList {
			return res.Failure(RTError(*n.Pos_Start, *n.Pos_End, "First Argument must be list", ctx))
		}

		list.Elements = append(list.Elements, value)
		ctx.Symbol_Table.Set("list", list)
		return res.Success(Number{Value: 0})
	}
}

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
func (n *RTResult) Set_pos(Pos_Start *tools.Position, Pos_End *tools.Position) Value {
	return n
}
func (n *RTResult) Set_context(context *Context) Value {
	return n
}
func (n RTResult) Get_context() *Context {
	return nil
}
func (n *RTResult) Copy() Value {
	return n
}

func (rtResult *RTResult) Register(res Value) Value {
	switch res := res.(type) {
	case *RTResult:
		if res.ShouldReturn() {
			rtResult.Error = res.Error
		}

		rtResult.FuncReturnValue = res.FuncReturnValue
		rtResult.LoopShouldContinue = res.LoopShouldContinue
		rtResult.LoopShouldBreak = res.LoopShouldBreak

		return res.Value
	}

	return res
}

func (rtResult *RTResult) Success(value Value) Value {
	rtResult.reset()
	rtResult.Value = value
	return rtResult
}

func (rtResult *RTResult) Success_Return(value Value) Value {
	rtResult.reset()
	rtResult.FuncReturnValue = value
	return rtResult
}

func (rtResult *RTResult) Success_Continue() Value {
	rtResult.reset()
	rtResult.LoopShouldContinue = true
	return rtResult
}

func (rtResult *RTResult) Success_Break() Value {
	rtResult.reset()
	rtResult.LoopShouldBreak = true
	return rtResult
}

func (rtResult *RTResult) Failure(error Error) Value {
	rtResult.reset()
	rtResult.Error = &error
	return rtResult
}

func (rtResult *RTResult) ShouldReturn() bool {
	return rtResult.Error != nil || rtResult.FuncReturnValue != nil || rtResult.LoopShouldBreak || rtResult.LoopShouldContinue
}

func (n *RTResult) Is_true() bool {
	return false
}
