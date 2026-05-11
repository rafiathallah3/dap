package interpreter

import (
	"dap/internal/common"
	"dap/internal/lexer"
	"dap/tools"
	"fmt"
	"reflect"
	"strconv"
)

type Interpreter struct{}

func (i *Interpreter) Visit(node common.Expr, context *common.Context) common.Value {
	methodName := "Visit" + reflect.TypeOf(node).Name()
	method := reflect.ValueOf(i).MethodByName(methodName)

	if method.IsValid() {
		results := method.Call([]reflect.Value{reflect.ValueOf(node), reflect.ValueOf(context)})

		if len(results) > 0 {
			return results[0].Interface().(common.Value)
		}
	}

	return i.noVisitMethod(node, context)
}

func (i *Interpreter) noVisitMethod(_ common.Expr, _ *common.Context) common.Value {
	// fmt.Println("No Visit" + reflect.TypeOf(node).Name() + " Method defined")
	return common.Null{}
}

func (i *Interpreter) VisitNullNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	return res.Success(common.Null{})
}

func (i *Interpreter) VisitNumberNode(node common.Expr, context *common.Context) common.Value {
	nodeToken := node.(common.NumberNode).Token
	parseFloat, err := strconv.ParseFloat(nodeToken.Value, 64)

	if err != nil {
		panic(fmt.Sprintf("Internal error: cannot parse '%s' as a number", nodeToken.Value))
	}

	numberValue := common.Number{
		Value:   parseFloat,
		Context: context,
	}

	res := &common.RTResult{}
	return res.Success(numberValue.Set_pos(nodeToken.Pos_Start, nodeToken.Pos_End))
}

func (i *Interpreter) VisitStringNode(node common.Expr, context *common.Context) common.Value {
	nodeToken := node.(common.StringNode).Token
	res := &common.RTResult{}
	return res.Success(common.String{Value: nodeToken.Value[1 : len(nodeToken.Value)-1], Context: context}.Set_pos(nodeToken.Pos_Start, nodeToken.Pos_End))
}

func (i *Interpreter) VisitDictionaryNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	nodeDictionary := node.(common.DictionaryNode)

	for _, v := range nodeDictionary.VariableDiBuat {
		res.Register(i.Visit(v, context))
		if res.ShouldReturn() {
			return res
		}
	}

	return res.Success(common.Null{})
}

func (i *Interpreter) VisitListNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	nodeList := node.(common.ListNode)
	elements := make([]common.Value, 0)

	for _, v := range nodeList.ElementNode {
		elements = append(elements, res.Register(i.Visit(v, context)))
		if res.ShouldReturn() {
			return res
		}
	}

	Listvalue := common.List{
		Elements: elements,
	}
	Listvalue.Set_pos(nodeList.Pos_Start, nodeList.Pos_End)
	Listvalue.Set_context(context)
	return res.Success(Listvalue)
}

func (i *Interpreter) VisitVarAccessNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	nodeVarAccessNode := node.(common.VarAccessNode)

	var_name := nodeVarAccessNode.VarNameTok.Value
	value := context.Symbol_Table.Get(var_name)
	if value == nil {
		value = common.Null{}
	}

	switch value.(type) {
	case common.Null:
		return res.Failure(common.RTError(*node.GetPosStart(), *node.GetPosEnd(), fmt.Sprintf("'%s' is not defined", var_name), context))
	}

	value = value.Copy().Set_pos(nodeVarAccessNode.Pos_Start, nodeVarAccessNode.Pos_end).Set_context(context)

	return res.Success(value)
}

func (i *Interpreter) VisitVarAssignNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	nodeVarAssignNode := node.(common.VarAssignNode)

	// var_name := nodeVarAssignNode.VarName.Value
	value := res.Register(i.Visit(nodeVarAssignNode.ValueNode, context))

	if res.ShouldReturn() {
		return res
	}

	// If the value is a Type definition, initialize it to get the default value
	if typeDef, ok := value.(common.Type); ok {
		value = res.Register(i.InitializeType(typeDef.Definition, context))
		if res.ShouldReturn() {
			return res
		}
	}

	res.Register(i.GantiVariable(nodeVarAssignNode.VarName.Value, value, context, nodeVarAssignNode.ApakahConst, nodeVarAssignNode.Pos_Start.Copy(), nodeVarAssignNode.Pos_end.Copy()))
	if res.Error != nil {
		return res
	}
	// context.Symbol_Table.Set(var_name, value)

	return res.Success(value)
}

func (i *Interpreter) VisitBinOpNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	nodeBinary := node.(common.BinOpNode)
	left := res.Register(i.Visit(nodeBinary.Left, context))

	if res.ShouldReturn() {
		return res
	}

	right := res.Register(i.Visit(nodeBinary.Right, context))

	if res.ShouldReturn() {
		return res
	}

	var hasil common.Value = common.Null{}
	var err *common.Error
	switch left := left.(type) {
	case common.Number:
		switch nodeBinary.Operator.Kind {
		case lexer.PLUS:
			hasil, err = left.Added_to(right)
		case lexer.DASH:
			hasil, err = left.Subbed_by(right)
		case lexer.STAR:
			hasil, err = left.Multed_by(right)
		case lexer.SLASH:
			hasil, err = left.Divided_by(right)
		case lexer.POWER:
			hasil, err = left.Powered_by(right)
		case lexer.EQUALS:
			hasil, err = left.Get_comparison_eq(right)
		case lexer.NOT_EQUALS:
			hasil, err = left.Get_comparison_nq(right)
		case lexer.LESS:
			hasil, err = left.Get_comparison_lt(right)
		case lexer.LESS_EQUALS:
			hasil, err = left.Get_comparison_lte(right)
		case lexer.GREATER:
			hasil, err = left.Get_comparison_gt(right)
		case lexer.GREATER_EQUALS:
			hasil, err = left.Get_comparison_gte(right)
		case lexer.AND:
			hasil, err = left.Anded_by(right)
		case lexer.OR:
			hasil, err = left.Ored_by(right)
		}
	case common.String:
		switch nodeBinary.Operator.Kind {
		case lexer.PLUS:
			hasil, err = left.Added_to(right)
		case lexer.STAR:
			hasil, err = left.Multed_by(right)
		}
	case common.List:
		switch nodeBinary.Operator.Kind {
		case lexer.PLUS:
			hasil, err = left.Added_to(right)
		case lexer.DASH:
			hasil, err = left.Subbed_to(right)
		case lexer.STAR:
			hasil, err = left.Multed_by(right)
		case lexer.SLASH:
			hasil, err = left.Dived_by(right)
		}
	}

	if err != nil {
		return res.Failure(*err)
	}

	return res.Success(hasil.Set_pos(nodeBinary.Pos_Start, nodeBinary.Pos_End))
}

func (i *Interpreter) VisitUnaryOpNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	nodeUnary := node.(common.UnaryOpNode)

	number := res.Register(i.Visit(nodeUnary.Node, context))
	if res.ShouldReturn() {
		return res
	}

	var error *common.Error
	switch nodeUnary.Operator.Kind {
	case lexer.DASH:
		angka := common.Number{
			Value: -1,
		}
		number, error = number.(common.Number).Multed_by(angka)
	case lexer.NOT:
		number, error = number.(common.Number).Notted()
	}

	if error != nil {
		return res.Failure(*error)
	}

	return res.Success(number.Set_pos(nodeUnary.Pos_Start, nodeUnary.Pos_End))
}

func (i *Interpreter) VisitIfNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	nodeIf := node.(common.IfNode)

	for _, ifCase := range nodeIf.Cases {
		condition_value := res.Register(i.Visit(ifCase.Kondisi, context))
		if res.ShouldReturn() {
			return res
		}

		if condition_value.Is_true() {
			expr_value := res.Register(i.Visit(ifCase.Isi, context))
			if res.ShouldReturn() {
				return res
			}

			if ifCase.ShouldReturnNull {
				return res.Success(common.Null{})
			}

			return res.Success(expr_value)
		}
	}

	if nodeIf.Else_case.Isi != nil {
		isi := nodeIf.Else_case.Isi
		else_value := res.Register(i.Visit(isi, context))
		if res.ShouldReturn() {
			return res
		}

		if nodeIf.Else_case.ShouldReturnNull {
			return res.Success(common.Null{})
		}
		return res.Success(else_value)
	}

	return res.Success(common.Null{})
}

func (i *Interpreter) VisitForNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	elements := make([]common.Value, 0)

	nodeFor := node.(common.ForNode)
	startValue := res.Register(i.Visit(nodeFor.StartValueNode, context)).(common.Number)
	if res.ShouldReturn() {
		return res
	}

	endValue := res.Register(i.Visit(nodeFor.EndValueNode, context)).(common.Number)
	if res.ShouldReturn() {
		return res
	}

	var stepValue common.Number
	switch nodeFor.StepValueNode.(type) {
	case common.NullNode:
		stepValue = common.Number{Value: 1}
	default:
		stepValue = res.Register(i.Visit(nodeFor.StepValueNode, context)).(common.Number)

		if res.ShouldReturn() {
			return res
		}
	}

	iteration := int(startValue.Value)

	var kondisi func(iteration int, endValue int) bool
	if stepValue.Value >= 0 {
		kondisi = func(iteration, endValue int) bool {
			return iteration <= endValue
		}
	} else {
		kondisi = func(iteration, endValue int) bool {
			return iteration >= endValue
		}
	}

	for kondisi(iteration, int(endValue.Value)) {
		res.Register(i.GantiVariable(nodeFor.VarNameTok.Value, common.Number{Value: float64(iteration)}, context, false, nodeFor.VarNameTok.Pos_Start.Copy(), nodeFor.VarNameTok.Pos_End.Copy()))
		if res.Error != nil {
			return res
		}

		// context.Symbol_Table.Set(nodeFor.VarNameTok.Value, common.Number{Value: float64(iteration)})

		iteration += int(stepValue.Value)

		value := res.Register(i.Visit(nodeFor.BodyNode, context))
		if res.ShouldReturn() && !res.LoopShouldContinue && !res.LoopShouldBreak {
			return res
		}

		if res.LoopShouldContinue {
			continue
		}

		if res.LoopShouldBreak {
			break
		}

		elements = append(elements, value)
	}

	if nodeFor.ShouldReturnNull {
		return res.Success(common.Null{})
	}

	ListValue := common.List{
		Elements: elements,
	}
	ListValue.Set_pos(nodeFor.Pos_Start, nodeFor.Pos_end)
	ListValue.Set_context(context)

	return res.Success(ListValue)
}

func (i *Interpreter) VisitWhileNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	nodeWhile := node.(common.WhileNode)
	elements := make([]common.Value, 0)

	for {
		kondisi := res.Register(i.Visit(nodeWhile.KondisiNode, context))
		if res.ShouldReturn() {
			return res
		}

		if !kondisi.Is_true() {
			break
		}

		value := res.Register(i.Visit(nodeWhile.BodyNode, context))
		if res.ShouldReturn() && !res.LoopShouldContinue && !res.LoopShouldBreak {
			return res
		}

		if res.LoopShouldContinue {
			continue
		}

		if res.LoopShouldBreak {
			break
		}

		elements = append(elements, value)
	}

	if nodeWhile.ShouldReturnNull {
		return res.Success(common.Null{})
	}

	ListValue := common.List{
		Elements: elements,
	}
	ListValue.Set_pos(nodeWhile.Pos_Start, nodeWhile.Pos_end)
	ListValue.Set_context(context)
	return res.Success(ListValue)
}

func (i *Interpreter) VisitRepeatNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	nodeWhile := node.(common.RepeatNode)
	elements := make([]common.Value, 0)

	for {
		value := res.Register(i.Visit(nodeWhile.BodyNode, context))
		if res.ShouldReturn() && !res.LoopShouldContinue && !res.LoopShouldBreak {
			return res
		}

		if res.LoopShouldContinue {
			continue
		}

		if res.LoopShouldBreak {
			break
		}

		kondisi := res.Register(i.Visit(nodeWhile.KondisiNode, context))
		if res.ShouldReturn() {
			return res
		}

		if kondisi.Is_true() {
			break
		}

		elements = append(elements, value)
	}

	if nodeWhile.ShouldReturnNull {
		return res.Success(common.Null{})
	}

	ListValue := common.List{
		Elements: elements,
	}
	ListValue.Set_pos(nodeWhile.Pos_Start, nodeWhile.Pos_end)
	ListValue.Set_context(context)
	return res.Success(ListValue)
}

func (i *Interpreter) VisitFuncNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	nodeFunc := node.(common.FuncNode)
	var PosStart *tools.Position
	var PosEnd *tools.Position

	funcName := "anonymous"
	if nodeFunc.VarNameTok != nil {
		funcName = nodeFunc.VarNameTok.Value
		PosStart = nodeFunc.VarNameTok.Pos_Start.Copy()
		PosEnd = nodeFunc.VarNameTok.Pos_End.Copy()
	}

	argNames := make([]string, 0)
	for _, argName := range nodeFunc.ArgNameToks {
		argNames = append(argNames, argName.Value)
	}
	funcValue := common.Function{
		BaseFunction: common.BaseFunction{
			Name:             funcName,
			BodyNode:         nodeFunc.BodyNode,
			ArgNames:         argNames,
			ShouldAutoReturn: nodeFunc.ShouldAutoReturn,
			Context:          context,
			Pos_Start:        nodeFunc.Pos_Start,
			Pos_End:          nodeFunc.Pos_end,
		},
	}

	if nodeFunc.VarNameTok != nil {
		res.Register(i.GantiVariable(funcName, funcValue, context, false, PosStart, PosEnd))
		if res.Error != nil {
			return res
		}
		// context.Symbol_Table.Set(funcName, funcValue)
	}

	if nodeFunc.ShouldAutoReturn {
		return res.Success(common.Null{})
	}

	return res.Success(funcValue)
}

func (i *Interpreter) VisitCallNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	args := make([]common.Value, 0)
	rawArgs := make([]common.Expr, 0)
	nodeCall := node.(common.CallNode)

	value_to_call := res.Register(i.Visit(nodeCall.NodeToCall, context))
	if res.ShouldReturn() {
		return res
	}
	value_to_call = value_to_call.Copy().Set_pos(nodeCall.Pos_Start, nodeCall.Pos_end)

	for _, argNode := range nodeCall.ArgNodes {
		rawArgs = append(rawArgs, argNode)
		args = append(args, res.Register(i.Visit(argNode, context)))
		if res.ShouldReturn() {
			return res
		}
	}

	var returnValue common.Value
	switch value_to_call := value_to_call.(type) {
	case common.BuiltInFunction:
		returnValue = res.Register(value_to_call.Execute(args, rawArgs))
		if returnValueContext := returnValue.Get_context(); returnValueContext != nil {
			for _, v := range rawArgs {
				switch v := v.(type) {
				case common.VarAccessNode:
					res.Register(i.GantiVariable(v.VarNameTok.Value, returnValueContext.Symbol_Table.Get(v.VarNameTok.Value), context, false, nil, nil))
					if res.Error != nil {
						return res
					}
					// context.Symbol_Table.Set(v.VarNameTok.Value, returnValueContext.Symbol_Table.Get(v.VarNameTok.Value))
				}
			}
		}
	default: //Normal Function
		returnValue = res.Register(i.Execute(value_to_call, context, args))
	}
	if res.ShouldReturn() {
		return res
	}
	returnValue = returnValue.Copy().Set_pos(nodeCall.Pos_Start, nodeCall.Pos_end).Set_context(context)

	return res.Success(returnValue)
}

func (i *Interpreter) VisitReturnNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	nodeReturn := node.(common.ReturnNode)

	var value common.Value = common.Null{}
	if nodeReturn.NodeToReturn != nil {
		value = res.Register(i.Visit(nodeReturn.NodeToReturn, context))
		if res.ShouldReturn() {
			return res
		}
	}

	return res.Success_Return(value)
}

func (i *Interpreter) VisitContinueNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	return res.Success_Continue()
}

func (i *Interpreter) VisitBreakNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	return res.Success_Break()
}

func (i *Interpreter) Execute(node common.Value, context *common.Context, args []common.Value) common.Value {
	res := &common.RTResult{}
	inter := Interpreter{}
	nodeFunc := node.(common.BaseFunctionInterface)

	exec_ctx := nodeFunc.GenerateNewContext()

	res.Register(nodeFunc.CheckAndPopulateArgs(nodeFunc.GetArgsName(), args, &exec_ctx))
	if res.ShouldReturn() {
		return res
	}

	value := res.Register(inter.Visit(nodeFunc.GetBodyNode(), &exec_ctx))
	if res.ShouldReturn() && res.FuncReturnValue == nil {
		return res
	}

	if nodeFunc.GetShouldAutoReturn() && value != nil {
		return res.Success(value)
	}

	if res.FuncReturnValue == nil {
		return res.Success(common.Null{})
	}

	return res.Success(res.FuncReturnValue)
}

func (i *Interpreter) GantiVariable(varName string, value common.Value, context *common.Context, apakahKonst bool, posStart *tools.Position, posend *tools.Position) common.Value {
	res := &common.RTResult{}
	apakahAdaKonst := context.Symbol_Table.Get("ApakahKonstant " + varName)

	if apakahAdaKonst == nil {
		hasil := 0.0
		if apakahKonst {
			hasil = 1.0
		}

		context.Symbol_Table.Set("ApakahKonstant "+varName, common.Number{Value: hasil})
	}

	if apakahAdaKonst != nil && apakahAdaKonst.(common.Number).Value == 1 {
		return res.Failure(common.RTError(*posStart, *posend, fmt.Sprintf("Constant variable '%s' can not be assigned!", varName), context))
	}

	context.Symbol_Table.Set(varName, value)

	return res.Success(common.Null{})
}

func (i *Interpreter) VisitArrayTypeNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	arrayNode := node.(common.ArrayTypeNode)

	startVal := res.Register(i.Visit(arrayNode.StartNode, context))
	if res.ShouldReturn() {
		return res
	}

	endVal := res.Register(i.Visit(arrayNode.EndNode, context))
	if res.ShouldReturn() {
		return res
	}

	if _, ok := startVal.(common.Number); !ok {
		return res.Failure(common.RTError(*arrayNode.StartNode.GetPosStart(), *arrayNode.StartNode.GetPosEnd(), "Start index must be a number", context))
	}
	if _, ok := endVal.(common.Number); !ok {
		return res.Failure(common.RTError(*arrayNode.EndNode.GetPosStart(), *arrayNode.EndNode.GetPosEnd(), "End index must be a number", context))
	}

	start := int(startVal.(common.Number).Value)
	end := int(endVal.(common.Number).Value)
	size := end - start + 1

	if size < 0 {
		return res.Failure(common.RTError(*node.GetPosStart(), *node.GetPosEnd(), "Array start index must be less than or equal to end index", context))
	}

	elements := make([]common.Value, size)
	for j := 0; j < size; j++ {
		elemVal := res.Register(i.InitializeType(arrayNode.OfType, context))
		if res.ShouldReturn() {
			return res
		}
		elements[j] = elemVal
	}

	return res.Success(common.Array{
		Elements:  elements,
		Start:     start,
		End:       end,
		Context:   context,
		Pos_Start: node.GetPosStart(),
		Pos_End:   node.GetPosEnd(),
	})
}

func (i *Interpreter) VisitArrayIndexNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	indexNode := node.(common.ArrayIndexNode)

	left := res.Register(i.Visit(indexNode.Left, context))
	if res.ShouldReturn() {
		return res
	}

	indexVal := res.Register(i.Visit(indexNode.Index, context))
	if res.ShouldReturn() {
		return res
	}

	array, ok := left.(common.Array)
	if !ok {
		return res.Failure(common.RTError(*indexNode.Left.GetPosStart(), *indexNode.Left.GetPosEnd(), "Left hand side is not an array", context))
	}

	if _, ok := indexVal.(common.Number); !ok {
		return res.Failure(common.RTError(*indexNode.Index.GetPosStart(), *indexNode.Index.GetPosEnd(), "Array index must be a number", context))
	}

	index := int(indexVal.(common.Number).Value)
	if index < array.Start || index > array.End {
		return res.Failure(common.RTError(*indexNode.Index.GetPosStart(), *indexNode.Index.GetPosEnd(), fmt.Sprintf("Index %d out of bounds [%d..%d]", index, array.Start, array.End), context))
	}

	return res.Success(array.Elements[index-array.Start])
}

func (i *Interpreter) VisitArrayAssignNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	assignNode := node.(common.ArrayAssignNode)

	left := res.Register(i.Visit(assignNode.ArrayAccess.Left, context))
	if res.ShouldReturn() {
		return res
	}

	indexVal := res.Register(i.Visit(assignNode.ArrayAccess.Index, context))
	if res.ShouldReturn() {
		return res
	}

	value := res.Register(i.Visit(assignNode.ValueNode, context))
	if res.ShouldReturn() {
		return res
	}

	array, ok := left.(common.Array)
	if !ok {
		return res.Failure(common.RTError(*assignNode.ArrayAccess.Left.GetPosStart(), *assignNode.ArrayAccess.Left.GetPosEnd(), "Left hand side is not an array", context))
	}

	if _, ok := indexVal.(common.Number); !ok {
		return res.Failure(common.RTError(*assignNode.ArrayAccess.Index.GetPosStart(), *assignNode.ArrayAccess.Index.GetPosEnd(), "Array index must be a number", context))
	}

	index := int(indexVal.(common.Number).Value)
	if index < array.Start || index > array.End {
		return res.Failure(common.RTError(*assignNode.ArrayAccess.Index.GetPosStart(), *assignNode.ArrayAccess.Index.GetPosEnd(), fmt.Sprintf("Index %d out of bounds [%d..%d]", index, array.Start, array.End), context))
	}

	array.Elements[index-array.Start] = value
	return res.Success(value)
}

func (i *Interpreter) VisitMemberAccessNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	accessNode := node.(common.MemberAccessNode)
	object := res.Register(i.Visit(accessNode.Object, context))
	if res.ShouldReturn() {
		return res
	}

	structVal, ok := object.(common.Struct)
	if !ok {
		return res.Failure(common.RTError(*accessNode.Object.GetPosStart(), *accessNode.Object.GetPosEnd(), "Object is not a struct", context))
	}

	val, ok := structVal.Fields[accessNode.MemberTok.Value]
	if !ok {
		return res.Failure(common.RTError(*accessNode.MemberTok.Pos_Start, *accessNode.MemberTok.Pos_End, fmt.Sprintf("Field '%s' not found in struct", accessNode.MemberTok.Value), context))
	}

	return res.Success(val)
}

func (i *Interpreter) VisitMemberAssignNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	assignNode := node.(common.MemberAssignNode)

	object := res.Register(i.Visit(assignNode.MemberAccess.Object, context))
	if res.ShouldReturn() {
		return res
	}

	value := res.Register(i.Visit(assignNode.ValueNode, context))
	if res.ShouldReturn() {
		return res
	}

	structVal, ok := object.(common.Struct)
	if !ok {
		return res.Failure(common.RTError(*assignNode.MemberAccess.Object.GetPosStart(), *assignNode.MemberAccess.Object.GetPosEnd(), "Object is not a struct", context))
	}

	if _, ok := structVal.Fields[assignNode.MemberAccess.MemberTok.Value]; !ok {
		return res.Failure(common.RTError(*assignNode.MemberAccess.MemberTok.Pos_Start, *assignNode.MemberAccess.MemberTok.Pos_End, fmt.Sprintf("Field '%s' not found in struct", assignNode.MemberAccess.MemberTok.Value), context))
	}

	structVal.Fields[assignNode.MemberAccess.MemberTok.Value] = value
	return res.Success(value)
}

func (i *Interpreter) VisitTypeAliasNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	aliasNode := node.(common.TypeAliasNode)
	typeValue := common.Type{
		Definition: aliasNode.TargetType,
		Context:    context,
		Pos_Start:  aliasNode.Pos_Start,
		Pos_End:    aliasNode.Pos_End,
	}
	context.Symbol_Table.Set(aliasNode.AliasName.Value, typeValue)
	return res.Success(common.Null{})
}

func (i *Interpreter) VisitStructTypeNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	structNode := node.(common.StructTypeNode)
	typeValue := common.Type{
		Definition: structNode,
		Context:    context,
		Pos_Start:  structNode.Pos_Start,
		Pos_End:    structNode.Pos_End,
	}
	context.Symbol_Table.Set(structNode.StructName.Value, typeValue)
	return res.Success(common.Null{})
}

func (i *Interpreter) InitializeType(typeNode common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}

	switch t := typeNode.(type) {
	case common.VarAccessNode:
		typeName := t.VarNameTok.Value
		switch typeName {
		case "integer", "real":
			return res.Success(common.Number{Value: 0, Context: context})
		case "string":
			return res.Success(common.String{Value: "", Context: context})
		default:
			val := context.Symbol_Table.Get(typeName)
			if typeDef, ok := val.(common.Type); ok {
				return i.InitializeType(typeDef.Definition, context)
			}
			return res.Failure(common.RTError(*t.GetPosStart(), *t.GetPosEnd(), fmt.Sprintf("Unknown type '%s'", typeName), context))
		}
	case common.ArrayTypeNode:
		startVal := res.Register(i.Visit(t.StartNode, context))
		if res.ShouldReturn() {
			return res
		}
		endVal := res.Register(i.Visit(t.EndNode, context))
		if res.ShouldReturn() {
			return res
		}

		start := int(startVal.(common.Number).Value)
		end := int(endVal.(common.Number).Value)
		size := end - start + 1
		elements := make([]common.Value, size)

		for idx := 0; idx < size; idx++ {
			elemVal := res.Register(i.InitializeType(t.OfType, context))
			if res.ShouldReturn() {
				return res
			}
			elements[idx] = elemVal
		}

		return res.Success(common.Array{
			Elements:  elements,
			Start:     start,
			End:       end,
			Context:   context,
			Pos_Start: t.Pos_Start,
			Pos_End:   t.Pos_End,
		})
	case common.StructTypeNode:
		fields := make(map[string]common.Value)
		for _, field := range t.Fields {
			fieldVal := res.Register(i.InitializeType(field.ValueNode, context))
			if res.ShouldReturn() {
				return res
			}
			fields[field.VarName.Value] = fieldVal
		}
		return res.Success(common.Struct{
			Fields:    fields,
			Context:   context,
			Pos_Start: t.Pos_Start,
			Pos_End:   t.Pos_End,
		})
	}

	return res.Success(common.Null{})
}
