package interpreter

import (
	"dap/internal/common"
	"dap/internal/lexer"
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
		panic("ERROR! Tidak bisa parse float")
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

	var_name := nodeVarAssignNode.VarName.Value
	value := res.Register(i.Visit(nodeVarAssignNode.ValueNode, context))

	if res.ShouldReturn() {
		return res
	}

	context.Symbol_Table.Set(var_name, value)

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

	if nodeIf.Else_case != nil {
		else_value := res.Register(i.Visit(*nodeIf.Else_case, context))
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
		context.Symbol_Table.Set(nodeFor.VarNameTok.Value, common.Number{Value: float64(iteration)})
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

func (i *Interpreter) VisitFuncNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	nodeFunc := node.(common.FuncNode)

	funcName := "anonymous"
	if nodeFunc.VarNameTok != nil {
		funcName = nodeFunc.VarNameTok.Value
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
		context.Symbol_Table.Set(funcName, funcValue)
	}

	if nodeFunc.ShouldAutoReturn {
		return res.Success(common.Null{})
	}

	return res.Success(funcValue)
}

func (i *Interpreter) VisitCallNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	args := make([]common.Value, 0)
	nodeCall := node.(common.CallNode)

	value_to_call := res.Register(i.Visit(nodeCall.NodeToCall, context))
	if res.ShouldReturn() {
		return res
	}
	value_to_call = value_to_call.Copy().Set_pos(nodeCall.Pos_Start, nodeCall.Pos_end)

	for _, argNode := range nodeCall.ArgNodes {
		args = append(args, res.Register(i.Visit(argNode, context)))
		if res.ShouldReturn() {
			return res
		}
	}

	var returnValue common.Value
	switch value_to_call := value_to_call.(type) {
	case common.BuiltInFunction:
		returnValue = res.Register(value_to_call.Execute(args))
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
