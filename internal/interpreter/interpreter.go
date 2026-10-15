package interpreter

import (
	"dap/internal/common"
	"dap/internal/lexer"
	"dap/tools"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func safePos(p *tools.Position) tools.Position {
	if p == nil {
		return tools.Position{}
	}
	return *p
}

// Interpreter buat ngeksekusi node-node AST dari program DAP.
type Interpreter struct{}

func valueTypeName(v common.Value) string {
	switch v.(type) {
	case common.Number:
		return "number"
	case common.String:
		return "string"
	case common.Boolean:
		return "boolean"
	case common.Array:
		return "array"
	case common.List:
		return "list"
	case common.Struct:
		return "struct"
	default:
		return "unknown"
	}
}

// Visit buat ngeksekusi AST node pake static type switch (jauh lebih ngebut daripada reflection).
func (i *Interpreter) Visit(node common.Expr, context *common.Context) common.Value {
	if node == nil {
		return common.Null{}
	}

	switch n := node.(type) {
	case common.NullNode:
		return i.VisitNullNode(n, context)
	case common.NumberNode:
		return i.VisitNumberNode(n, context)
	case common.StringNode:
		return i.VisitStringNode(n, context)
	case common.DictionaryNode:
		return i.VisitDictionaryNode(n, context)
	case common.ListNode:
		return i.VisitListNode(n, context)
	case common.VarAccessNode:
		return i.VisitVarAccessNode(n, context)
	case common.VarAssignNode:
		return i.VisitVarAssignNode(n, context)
	case common.BinOpNode:
		return i.VisitBinOpNode(n, context)
	case common.UnaryOpNode:
		return i.VisitUnaryOpNode(n, context)
	case common.IfNode:
		return i.VisitIfNode(n, context)
	case common.ForNode:
		return i.VisitForNode(n, context)
	case common.WhileNode:
		return i.VisitWhileNode(n, context)
	case common.RepeatNode:
		return i.VisitRepeatNode(n, context)
	case common.FuncNode:
		return i.VisitFuncNode(n, context)
	case common.CallNode:
		return i.VisitCallNode(n, context)
	case common.ReturnNode:
		return i.VisitReturnNode(n, context)
	case common.ContinueNode:
		return i.VisitContinueNode(n, context)
	case common.BreakNode:
		return i.VisitBreakNode(n, context)
	case common.ArrayTypeNode:
		return i.VisitArrayTypeNode(n, context)
	case common.ArrayIndexNode:
		return i.VisitArrayIndexNode(n, context)
	case common.ArrayAssignNode:
		return i.VisitArrayAssignNode(n, context)
	case common.MemberAccessNode:
		return i.VisitMemberAccessNode(n, context)
	case common.MemberAssignNode:
		return i.VisitMemberAssignNode(n, context)
	case common.TypeAliasNode:
		return i.VisitTypeAliasNode(n, context)
	case common.StructTypeNode:
		return i.VisitStructTypeNode(n, context)
	case common.DependNode:
		return i.VisitDependNode(n, context)
	case *common.ParseResult:
		return i.Visit(n.Node, context)
	default:
		return i.noVisitMethod(node, context)
	}
}

func (i *Interpreter) noVisitMethod(_ common.Expr, _ *common.Context) common.Value {
	return common.Null{}
}

func (i *Interpreter) VisitNullNode(_ common.Expr, _ *common.Context) common.Value {
	res := &common.RTResult{}
	return res.Success(common.Null{})
}

func (i *Interpreter) VisitNumberNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	numNode, ok := node.(common.NumberNode)
	if !ok {
		return res.Failure(common.RTError(safePos(node.GetPosStart()), safePos(node.GetPosEnd()), "Expected number node", context))
	}

	nodeToken := numNode.Token
	parseFloat, err := strconv.ParseFloat(nodeToken.Value, 64)
	if err != nil {
		return res.Failure(common.RTError(safePos(nodeToken.PosStart), safePos(nodeToken.PosEnd), fmt.Sprintf("Internal error: cannot parse '%s' as a number", nodeToken.Value), context))
	}

	numberValue := common.Number{
		Value:   parseFloat,
		IsReal:  strings.Contains(nodeToken.Value, "."),
		Context: context,
	}

	return res.Success(numberValue.SetPos(nodeToken.PosStart, nodeToken.PosEnd))
}

func (i *Interpreter) VisitStringNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	strNode, ok := node.(common.StringNode)
	if !ok {
		return res.Failure(common.RTError(safePos(node.GetPosStart()), safePos(node.GetPosEnd()), "Expected string node", context))
	}

	nodeToken := strNode.Token
	val := nodeToken.Value
	if len(val) >= 2 {
		val = val[1 : len(val)-1]
	}
	return res.Success(common.String{Value: val, Context: context}.SetPos(nodeToken.PosStart, nodeToken.PosEnd))
}

func (i *Interpreter) VisitDictionaryNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	nodeDictionary, ok := node.(common.DictionaryNode)
	if !ok {
		return res.Failure(common.RTError(safePos(node.GetPosStart()), safePos(node.GetPosEnd()), "Expected dictionary node", context))
	}

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
	nodeList, ok := node.(common.ListNode)
	if !ok {
		return res.Failure(common.RTError(safePos(node.GetPosStart()), safePos(node.GetPosEnd()), "Expected list node", context))
	}

	elements := make([]common.Value, 0, len(nodeList.ElementNode))
	for _, v := range nodeList.ElementNode {
		elements = append(elements, res.Register(i.Visit(v, context)))
		if res.ShouldReturn() {
			return res
		}
	}

	listValue := common.List{
		Elements: elements,
	}
	listValue.SetPos(nodeList.Pos_Start, nodeList.Pos_End)
	listValue.SetContext(context)
	return res.Success(listValue)
}

func (i *Interpreter) VisitVarAccessNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	nodeVarAccessNode, ok := node.(common.VarAccessNode)
	if !ok {
		return res.Failure(common.RTError(safePos(node.GetPosStart()), safePos(node.GetPosEnd()), "Expected variable access node", context))
	}

	varName := nodeVarAccessNode.VarNameTok.Value
	value := context.GetSymbolTable().Get(varName)
	if value == nil {
		value = common.Null{}
	}

	if _, isNull := value.(common.Null); isNull {
		return res.Failure(common.RTError(safePos(node.GetPosStart()), safePos(node.GetPosEnd()), fmt.Sprintf("'%s' is not defined", varName), context))
	}

	value = value.Copy().SetPos(nodeVarAccessNode.Pos_Start, nodeVarAccessNode.Pos_end).SetContext(context)
	return res.Success(value)
}

func (i *Interpreter) VisitVarAssignNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	nodeVarAssignNode, ok := node.(common.VarAssignNode)
	if !ok {
		return res.Failure(common.RTError(safePos(node.GetPosStart()), safePos(node.GetPosEnd()), "Expected variable assign node", context))
	}

	value := res.Register(i.Visit(nodeVarAssignNode.ValueNode, context))
	if res.ShouldReturn() {
		return res
	}

	if typeDef, isType := value.(common.Type); isType {
		value = res.Register(i.InitializeType(typeDef.Definition, context))
		if res.ShouldReturn() {
			return res
		}
	}

	if nodeVarAssignNode.ApakahDeklarasi {
		typeTarget := nodeVarAssignNode.ValueNode
		if nodeVarAssignNode.TypeNode != nil {
			typeTarget = nodeVarAssignNode.TypeNode
		}
		baseTypeName, err := i.GetBaseTypeName(typeTarget, context)
		if err != nil {
			return res.Failure(*err)
		}
		if baseTypeName != "" {
			context.GetSymbolTable().Set("VariableType "+nodeVarAssignNode.VarName.Value, common.String{Value: baseTypeName, Context: context})
		}
	}

	res.Register(i.AssignVariable(nodeVarAssignNode.VarName.Value, value, context, nodeVarAssignNode.ApakahConst, nodeVarAssignNode.Pos_Start.Copy(), nodeVarAssignNode.Pos_end.Copy()))
	if res.Error != nil {
		return res
	}

	return res.Success(value)
}

func (i *Interpreter) VisitBinOpNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	nodeBinary, ok := node.(common.BinOpNode)
	if !ok {
		return res.Failure(common.RTError(safePos(node.GetPosStart()), safePos(node.GetPosEnd()), "Expected binary operation node", context))
	}

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

	switch l := left.(type) {
	case common.Number:
		switch nodeBinary.Operator.Kind {
		case lexer.PLUS:
			hasil, err = l.Add(right)
		case lexer.DASH:
			hasil, err = l.Subtract(right)
		case lexer.STAR:
			hasil, err = l.Multiply(right)
		case lexer.SLASH:
			hasil, err = l.Divide(right)
		case lexer.DIV:
			hasil, err = l.Div(right)
		case lexer.PERCENT:
			hasil, err = l.Mod(right)
		case lexer.POWER, lexer.DOUBLE_STAR:
			hasil, err = l.Power(right)
		case lexer.BIT_AND:
			hasil, err = l.BitwiseAnd(right)
		case lexer.BIT_OR:
			hasil, err = l.BitwiseOr(right)
		case lexer.BIT_XOR:
			hasil, err = l.BitwiseXor(right)
		case lexer.SHIFT_LEFT:
			hasil, err = l.ShiftLeft(right)
		case lexer.SHIFT_RIGHT:
			hasil, err = l.ShiftRight(right)
		case lexer.EQUALS:
			hasil, err = l.Equal(right)
		case lexer.NOT_EQUALS:
			hasil, err = l.NotEqual(right)
		case lexer.LESS:
			hasil, err = l.Less(right)
		case lexer.LESS_EQUALS:
			hasil, err = l.LessEqual(right)
		case lexer.GREATER:
			hasil, err = l.Greater(right)
		case lexer.GREATER_EQUALS:
			hasil, err = l.GreaterEqual(right)
		case lexer.AND:
			hasil, err = l.And(right)
		case lexer.OR:
			hasil, err = l.Or(right)
		}
	case common.String:
		switch nodeBinary.Operator.Kind {
		case lexer.PLUS:
			hasil, err = l.Add(right)
		case lexer.STAR:
			hasil, err = l.Multiply(right)
		case lexer.EQUALS:
			hasil, err = l.Equal(right)
		case lexer.NOT_EQUALS:
			hasil, err = l.NotEqual(right)
		}
	case common.List:
		switch nodeBinary.Operator.Kind {
		case lexer.PLUS:
			hasil, err = l.Add(right)
		case lexer.DASH:
			hasil, err = l.Subtract(right)
		case lexer.STAR:
			hasil, err = l.Multiply(right)
		case lexer.SLASH:
			hasil, err = l.Divide(right)
		}
	case common.Boolean:
		switch nodeBinary.Operator.Kind {
		case lexer.EQUALS:
			hasil, err = l.Equal(right)
		case lexer.NOT_EQUALS:
			hasil, err = l.NotEqual(right)
		case lexer.AND:
			hasil, err = l.And(right)
		case lexer.OR:
			hasil, err = l.Or(right)
		}
	}

	if err != nil {
		return res.Failure(*err)
	}

	if hasil == nil {
		hasil = common.Null{}
	}

	return res.Success(hasil.SetPos(nodeBinary.Pos_Start, nodeBinary.Pos_End))
}

func (i *Interpreter) VisitUnaryOpNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	nodeUnary, ok := node.(common.UnaryOpNode)
	if !ok {
		return res.Failure(common.RTError(safePos(node.GetPosStart()), safePos(node.GetPosEnd()), "Expected unary op node", context))
	}

	val := res.Register(i.Visit(nodeUnary.Node, context))
	if res.ShouldReturn() {
		return res
	}

	var err *common.Error
	var resultVal common.Value

	switch nodeUnary.Operator.Kind {
	case lexer.BIT_NOT:
		if num, isNum := val.(common.Number); isNum {
			resultVal, err = num.BitwiseNot()
		} else {
			errNya := common.RTError(safePos(nodeUnary.Operator.PosStart), safePos(nodeUnary.Operator.PosEnd), "Illegal operation: bitwise '~' expects an integer", context)
			return res.Failure(errNya)
		}
	case lexer.DASH:
		if num, isNum := val.(common.Number); isNum {
			angka := common.Number{
				Value:  -1,
				IsReal: num.IsReal,
			}
			resultVal, err = num.Multiply(angka)
		} else {
			errNya := common.RTError(safePos(nodeUnary.Operator.PosStart), safePos(nodeUnary.Operator.PosEnd), "Illegal operation: unary minus expects a number", context)
			return res.Failure(errNya)
		}
	case lexer.NOT:
		switch v := val.(type) {
		case common.Number:
			resultVal, err = v.Not()
		case common.Boolean:
			resultVal = common.Boolean{Value: !v.Value, Context: context, Pos_Start: v.Pos_Start, Pos_End: v.Pos_End}
		default:
			errNya := common.RTError(safePos(nodeUnary.Operator.PosStart), safePos(nodeUnary.Operator.PosEnd), "Illegal operation: logic not expects a boolean or number", context)
			return res.Failure(errNya)
		}
	}

	if err != nil {
		return res.Failure(*err)
	}

	return res.Success(resultVal.SetPos(nodeUnary.Pos_Start, nodeUnary.Pos_End))
}

func (i *Interpreter) VisitIfNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	nodeIf, ok := node.(common.IfNode)
	if !ok {
		return res.Failure(common.RTError(safePos(node.GetPosStart()), safePos(node.GetPosEnd()), "Expected if node", context))
	}

	for _, ifCase := range nodeIf.Cases {
		conditionVal := res.Register(i.Visit(ifCase.Kondisi, context))
		if res.ShouldReturn() {
			return res
		}

		if conditionVal.IsTrue() {
			exprVal := res.Register(i.Visit(ifCase.Isi, context))
			if res.ShouldReturn() {
				return res
			}

			if ifCase.ShouldReturnNull {
				return res.Success(common.Null{})
			}

			return res.Success(exprVal)
		}
	}

	if nodeIf.Else_case != nil && nodeIf.Else_case.Isi != nil {
		elseVal := res.Register(i.Visit(nodeIf.Else_case.Isi, context))
		if res.ShouldReturn() {
			return res
		}

		if nodeIf.Else_case.ShouldReturnNull {
			return res.Success(common.Null{})
		}
		return res.Success(elseVal)
	}

	return res.Success(common.Null{})
}

func (i *Interpreter) VisitForNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	nodeFor, ok := node.(common.ForNode)
	if !ok {
		return res.Failure(common.RTError(safePos(node.GetPosStart()), safePos(node.GetPosEnd()), "Expected for node", context))
	}

	startRes := res.Register(i.Visit(nodeFor.StartValueNode, context))
	if res.ShouldReturn() {
		return res
	}
	startValue, ok := startRes.(common.Number)
	if !ok {
		return res.Failure(common.RTError(safePos(nodeFor.StartValueNode.GetPosStart()), safePos(nodeFor.StartValueNode.GetPosEnd()), "Start value in 'for' loop must be a number", context))
	}

	endRes := res.Register(i.Visit(nodeFor.EndValueNode, context))
	if res.ShouldReturn() {
		return res
	}
	endValue, ok := endRes.(common.Number)
	if !ok {
		return res.Failure(common.RTError(safePos(nodeFor.EndValueNode.GetPosStart()), safePos(nodeFor.EndValueNode.GetPosEnd()), "End value in 'for' loop must be a number", context))
	}

	var stepValue common.Number
	switch nodeFor.StepValueNode.(type) {
	case common.NullNode:
		if nodeFor.IsDown {
			stepValue = common.Number{Value: -1, IsReal: false}
		} else {
			stepValue = common.Number{Value: 1, IsReal: false}
		}
	default:
		stepRes := res.Register(i.Visit(nodeFor.StepValueNode, context))
		if res.ShouldReturn() {
			return res
		}
		stepNum, isNum := stepRes.(common.Number)
		if !isNum {
			return res.Failure(common.RTError(safePos(nodeFor.StepValueNode.GetPosStart()), safePos(nodeFor.StepValueNode.GetPosEnd()), "Step value in 'for' loop must be a number", context))
		}
		stepValue = stepNum
		if nodeFor.IsDown && stepValue.Value > 0 {
			stepValue.Value = -stepValue.Value
		}
	}

	elements := make([]common.Value, 0)
	iteration := int(startValue.Value)

	var kondisi func(iteration int, endValue int) bool
	if stepValue.Value >= 0 {
		kondisi = func(iter, end int) bool {
			return iter <= end
		}
	} else {
		kondisi = func(iter, end int) bool {
			return iter >= end
		}
	}

	for kondisi(iteration, int(endValue.Value)) {
		res.Register(i.AssignVariable(nodeFor.VarNameTok.Value, common.Number{Value: float64(iteration), IsReal: false}, context, false, nodeFor.VarNameTok.PosStart.Copy(), nodeFor.VarNameTok.PosEnd.Copy()))
		if res.Error != nil {
			return res
		}

		iteration += int(stepValue.Value)

		val := res.Register(i.Visit(nodeFor.BodyNode, context))
		if res.ShouldReturn() && !res.LoopShouldContinue && !res.LoopShouldBreak {
			return res
		}

		if res.LoopShouldContinue {
			continue
		}

		if res.LoopShouldBreak {
			break
		}

		elements = append(elements, val)
	}

	if nodeFor.ShouldReturnNull {
		return res.Success(common.Null{})
	}

	listVal := common.List{
		Elements: elements,
	}
	listVal.SetPos(nodeFor.Pos_Start, nodeFor.Pos_end)
	listVal.SetContext(context)
	return res.Success(listVal)
}

func (i *Interpreter) VisitWhileNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	nodeWhile, ok := node.(common.WhileNode)
	if !ok {
		return res.Failure(common.RTError(safePos(node.GetPosStart()), safePos(node.GetPosEnd()), "Expected while node", context))
	}
	elements := make([]common.Value, 0)

	for {
		kondisi := res.Register(i.Visit(nodeWhile.KondisiNode, context))
		if res.ShouldReturn() {
			return res
		}

		if !kondisi.IsTrue() {
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

	listValue := common.List{
		Elements: elements,
	}
	listValue.SetPos(nodeWhile.Pos_Start, nodeWhile.Pos_end)
	listValue.SetContext(context)
	return res.Success(listValue)
}

func (i *Interpreter) VisitRepeatNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	nodeRepeat, ok := node.(common.RepeatNode)
	if !ok {
		return res.Failure(common.RTError(safePos(node.GetPosStart()), safePos(node.GetPosEnd()), "Expected repeat node", context))
	}
	elements := make([]common.Value, 0)

	for {
		value := res.Register(i.Visit(nodeRepeat.BodyNode, context))
		if res.ShouldReturn() && !res.LoopShouldContinue && !res.LoopShouldBreak {
			return res
		}

		if res.LoopShouldContinue {
			continue
		}

		if res.LoopShouldBreak {
			break
		}

		kondisi := res.Register(i.Visit(nodeRepeat.KondisiNode, context))
		if res.ShouldReturn() {
			return res
		}

		if kondisi.IsTrue() {
			break
		}

		elements = append(elements, value)
	}

	if nodeRepeat.ShouldReturnNull {
		return res.Success(common.Null{})
	}

	listValue := common.List{
		Elements: elements,
	}
	listValue.SetPos(nodeRepeat.Pos_Start, nodeRepeat.Pos_end)
	listValue.SetContext(context)
	return res.Success(listValue)
}

func (i *Interpreter) VisitFuncNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	nodeFunc, ok := node.(common.FuncNode)
	if !ok {
		return res.Failure(common.RTError(safePos(node.GetPosStart()), safePos(node.GetPosEnd()), "Expected func node", context))
	}
	var posStart *tools.Position
	var posEnd *tools.Position

	funcName := "anonymous"
	if nodeFunc.VarNameTok != nil {
		funcName = nodeFunc.VarNameTok.Value
		posStart = nodeFunc.VarNameTok.PosStart.Copy()
		posEnd = nodeFunc.VarNameTok.PosEnd.Copy()
	}

	argNames := make([]string, 0, len(nodeFunc.ArgNameToks))
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
		res.Register(i.AssignVariable(funcName, funcValue, context, false, posStart, posEnd))
		if res.Error != nil {
			return res
		}
	}

	if nodeFunc.ShouldAutoReturn {
		return res.Success(common.Null{})
	}

	return res.Success(funcValue)
}

func (i *Interpreter) VisitCallNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	nodeCall, ok := node.(common.CallNode)
	if !ok {
		return res.Failure(common.RTError(safePos(node.GetPosStart()), safePos(node.GetPosEnd()), "Expected call node", context))
	}

	valueToCall := res.Register(i.Visit(nodeCall.NodeToCall, context))
	if res.ShouldReturn() {
		return res
	}
	valueToCall = valueToCall.Copy().SetPos(nodeCall.Pos_Start, nodeCall.Pos_end)

	rawArgs := make([]common.Expr, 0, len(nodeCall.ArgNodes))
	args := make([]common.Value, 0, len(nodeCall.ArgNodes))
	for _, argNode := range nodeCall.ArgNodes {
		rawArgs = append(rawArgs, argNode)
		args = append(args, res.Register(i.Visit(argNode, context)))
		if res.ShouldReturn() {
			return res
		}
	}

	var returnValue common.Value
	switch fn := valueToCall.(type) {
	case common.BuiltInFunction:
		returnValue = res.Register(fn.Execute(args, rawArgs))
		if res.ShouldReturn() {
			return res
		}
		if returnValue != nil {
			if retCtx := returnValue.GetContext(); retCtx != nil {
				for _, v := range rawArgs {
					if va, ok := v.(common.VarAccessNode); ok {
						newVal := retCtx.GetSymbolTable().Get(va.VarNameTok.Value)
						if newVal != nil {
							res.Register(i.AssignVariable(va.VarNameTok.Value, newVal, context, false, nil, nil))
							if res.Error != nil {
								return res
							}
						}
					}
				}
			}
		}
	default:
		returnValue = res.Register(i.Execute(fn, context, args))
	}
	if res.ShouldReturn() {
		return res
	}
	returnValue = returnValue.Copy().SetPos(nodeCall.Pos_Start, nodeCall.Pos_end).SetContext(context)
	return res.Success(returnValue)
}

func (i *Interpreter) VisitReturnNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	nodeReturn, ok := node.(common.ReturnNode)
	if !ok {
		return res.Failure(common.RTError(safePos(node.GetPosStart()), safePos(node.GetPosEnd()), "Expected return node", context))
	}

	var value common.Value = common.Null{}
	if nodeReturn.NodeToReturn != nil {
		value = res.Register(i.Visit(nodeReturn.NodeToReturn, context))
		if res.ShouldReturn() {
			return res
		}
	}

	return res.SuccessReturn(value)
}

func (i *Interpreter) VisitContinueNode(_ common.Expr, _ *common.Context) common.Value {
	res := &common.RTResult{}
	return res.SuccessContinue()
}

func (i *Interpreter) VisitBreakNode(_ common.Expr, _ *common.Context) common.Value {
	res := &common.RTResult{}
	return res.SuccessBreak()
}

func (i *Interpreter) Execute(node common.Value, context *common.Context, args []common.Value) common.Value {
	res := &common.RTResult{}
	nodeFunc, ok := node.(common.BaseFunctionInterface)
	if !ok {
		return res.Failure(common.RTError(tools.Position{}, tools.Position{}, "Object is not callable", context))
	}

	execCtx := nodeFunc.GenerateNewContext()
	res.Register(nodeFunc.CheckAndPopulateArgs(nodeFunc.GetArgsName(), args, &execCtx))
	if res.ShouldReturn() {
		return res
	}

	value := res.Register(i.Visit(nodeFunc.GetBodyNode(), &execCtx))
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

// AssignVariable buat nge-assign nilai variabel plus proteksi konstanta & validasi tipe data.
func (i *Interpreter) AssignVariable(varName string, value common.Value, context *common.Context, isConst bool, posStart *tools.Position, posEnd *tools.Position) common.Value {
	res := &common.RTResult{}
	st := context.GetSymbolTable()
	isConstVal := st.Get("ApakahKonstant " + varName)

	if isConstVal == nil {
		flag := 0.0
		if isConst {
			flag = 1.0
			if posStart != nil {
				st.Set("ConstDeclLine "+varName, common.Number{Value: float64(posStart.Ln + 1), IsReal: false, Context: context})
			}
		}
		st.Set("ApakahKonstant "+varName, common.Number{Value: flag, IsReal: false, Context: context})
	}

	if isConstVal != nil {
		if num, ok := isConstVal.(common.Number); ok && num.Value == 1 {
			var start, end tools.Position
			if posStart != nil {
				start = *posStart
			}
			if posEnd != nil {
				end = *posEnd
			}
			declLineMsg := ""
			declLineVal := st.Get("ConstDeclLine " + varName)
			if declLineVal != nil {
				if dNum, ok := declLineVal.(common.Number); ok {
					declLineMsg = fmt.Sprintf(" (declared at line %d)", int(dNum.Value))
				}
			}
			return res.Failure(common.RTError(start, end, fmt.Sprintf("Constant variable '%s' can not be assigned!%s", varName, declLineMsg), context))
		}
	}

	// Cek deklarasi tipe variabel biar ga salah assign tipe
	declaredTypeVal := st.Get("VariableType " + varName)
	if declaredTypeVal != nil {
		if declaredTypeStr, ok := declaredTypeVal.(common.String); ok {
			switch declaredTypeStr.Value {
			case "integer":
				if num, ok := value.(common.Number); ok {
					if num.IsReal || num.Value != math.Floor(num.Value) {
						var start, end tools.Position
						if posStart != nil {
							start = *posStart
						}
						if posEnd != nil {
							end = *posEnd
						}
						return res.Failure(common.RTError(start, end, fmt.Sprintf("Type mismatch: cannot assign real number to integer variable '%s'", varName), context))
					}
					num.IsReal = false
					value = num
				} else {
					var start, end tools.Position
					if posStart != nil {
						start = *posStart
					}
					if posEnd != nil {
						end = *posEnd
					}
					return res.Failure(common.RTError(start, end, fmt.Sprintf("Type mismatch: cannot assign %s to integer variable '%s'", valueTypeName(value), varName), context))
				}
			case "real":
				if num, ok := value.(common.Number); ok {
					num.IsReal = true
					value = num
				} else {
					var start, end tools.Position
					if posStart != nil {
						start = *posStart
					}
					if posEnd != nil {
						end = *posEnd
					}
					return res.Failure(common.RTError(start, end, fmt.Sprintf("Type mismatch: cannot assign %s to real variable '%s'", valueTypeName(value), varName), context))
				}
			}
		}
	}

	st.Set(varName, value)
	return res.Success(common.Null{})
}

// GantiVariable alias buat backward compatibility.
func (i *Interpreter) GantiVariable(varName string, value common.Value, context *common.Context, isConst bool, posStart *tools.Position, posEnd *tools.Position) common.Value {
	return i.AssignVariable(varName, value, context, isConst, posStart, posEnd)
}

func (i *Interpreter) VisitArrayTypeNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	arrayNode, ok := node.(common.ArrayTypeNode)
	if !ok {
		return res.Failure(common.RTError(safePos(node.GetPosStart()), safePos(node.GetPosEnd()), "Expected array type node", context))
	}

	startVal := res.Register(i.Visit(arrayNode.StartNode, context))
	if res.ShouldReturn() {
		return res
	}

	endVal := res.Register(i.Visit(arrayNode.EndNode, context))
	if res.ShouldReturn() {
		return res
	}

	startNum, ok1 := startVal.(common.Number)
	if !ok1 {
		return res.Failure(common.RTError(safePos(arrayNode.StartNode.GetPosStart()), safePos(arrayNode.StartNode.GetPosEnd()), "Start index must be a number", context))
	}
	endNum, ok2 := endVal.(common.Number)
	if !ok2 {
		return res.Failure(common.RTError(safePos(arrayNode.EndNode.GetPosStart()), safePos(arrayNode.EndNode.GetPosEnd()), "End index must be a number", context))
	}

	start := int(startNum.Value)
	end := int(endNum.Value)
	size := end - start + 1

	if size < 0 {
		return res.Failure(common.RTError(safePos(node.GetPosStart()), safePos(node.GetPosEnd()), "Array start index must be less than or equal to end index", context))
	}

	elements := make([]common.Value, size)
	for j := 0; j < size; j++ {
		elemVal := res.Register(i.InitializeType(arrayNode.OfType, context))
		if res.ShouldReturn() {
			return res
		}
		elements[j] = elemVal
	}

	elemTypeName, _ := i.GetBaseTypeName(arrayNode.OfType, context)
	return res.Success(common.Array{
		Elements:    elements,
		ElementType: elemTypeName,
		Start:       start,
		End:         end,
		Context:     context,
		Pos_Start:   node.GetPosStart(),
		Pos_End:     node.GetPosEnd(),
	})
}

func (i *Interpreter) VisitArrayIndexNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	indexNode, ok := node.(common.ArrayIndexNode)
	if !ok {
		return res.Failure(common.RTError(safePos(node.GetPosStart()), safePos(node.GetPosEnd()), "Expected array index node", context))
	}

	left := res.Register(i.Visit(indexNode.Left, context))
	if res.ShouldReturn() {
		return res
	}

	indexVal := res.Register(i.Visit(indexNode.Index, context))
	if res.ShouldReturn() {
		return res
	}

	switch l := left.(type) {
	case common.Array:
		num, ok := indexVal.(common.Number)
		if !ok {
			return res.Failure(common.RTError(safePos(indexNode.Index.GetPosStart()), safePos(indexNode.Index.GetPosEnd()), "Array index must be a number", context))
		}

		index := int(num.Value)
		if index < l.Start || index > l.End {
			return res.Failure(common.RTError(safePos(indexNode.Index.GetPosStart()), safePos(indexNode.Index.GetPosEnd()), fmt.Sprintf("Index %d out of bounds [%d..%d]", index, l.Start, l.End), context))
		}

		return res.Success(l.Elements[index-l.Start])
	case common.String:
		num, ok := indexVal.(common.Number)
		if !ok {
			return res.Failure(common.RTError(safePos(indexNode.Index.GetPosStart()), safePos(indexNode.Index.GetPosEnd()), "String index must be a number", context))
		}

		index := int(num.Value)
		if index < 1 || index > len(l.Value) {
			return res.Failure(common.RTError(safePos(indexNode.Index.GetPosStart()), safePos(indexNode.Index.GetPosEnd()), fmt.Sprintf("String index %d out of bounds [1..%d]", index, len(l.Value)), context))
		}

		return res.Success(common.String{Value: string(l.Value[index-1]), Context: context})
	default:
		return res.Failure(common.RTError(safePos(indexNode.Left.GetPosStart()), safePos(indexNode.Left.GetPosEnd()), "Left hand side is not indexable", context))
	}
}

func (i *Interpreter) VisitArrayAssignNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	assignNode, ok := node.(common.ArrayAssignNode)
	if !ok {
		return res.Failure(common.RTError(safePos(node.GetPosStart()), safePos(node.GetPosEnd()), "Expected array assign node", context))
	}

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
		return res.Failure(common.RTError(safePos(assignNode.ArrayAccess.Left.GetPosStart()), safePos(assignNode.ArrayAccess.Left.GetPosEnd()), "Left hand side is not an array", context))
	}

	num, ok := indexVal.(common.Number)
	if !ok {
		return res.Failure(common.RTError(safePos(assignNode.ArrayAccess.Index.GetPosStart()), safePos(assignNode.ArrayAccess.Index.GetPosEnd()), "Array index must be a number", context))
	}

	index := int(num.Value)
	if index < array.Start || index > array.End {
		return res.Failure(common.RTError(safePos(assignNode.ArrayAccess.Index.GetPosStart()), safePos(assignNode.ArrayAccess.Index.GetPosEnd()), fmt.Sprintf("Index %d out of bounds [%d..%d]", index, array.Start, array.End), context))
	}

	elemType := array.ElementType
	if elemType == "" && len(array.Elements) > 0 {
		if firstNum, ok := array.Elements[0].(common.Number); ok {
			if firstNum.IsReal {
				elemType = "real"
			} else {
				elemType = "integer"
			}
		}
	}

	switch elemType {
	case "integer":
		if numVal, ok := value.(common.Number); ok {
			if numVal.IsReal || numVal.Value != math.Floor(numVal.Value) {
				return res.Failure(common.RTError(safePos(assignNode.ValueNode.GetPosStart()), safePos(assignNode.ValueNode.GetPosEnd()), "Type mismatch: cannot assign real number to integer array element", context))
			}
			numVal.IsReal = false
			value = numVal
		} else {
			return res.Failure(common.RTError(safePos(assignNode.ValueNode.GetPosStart()), safePos(assignNode.ValueNode.GetPosEnd()), fmt.Sprintf("Type mismatch: cannot assign %s to integer array element", valueTypeName(value)), context))
		}
	case "real":
		if numVal, ok := value.(common.Number); ok {
			numVal.IsReal = true
			value = numVal
		} else {
			return res.Failure(common.RTError(safePos(assignNode.ValueNode.GetPosStart()), safePos(assignNode.ValueNode.GetPosEnd()), fmt.Sprintf("Type mismatch: cannot assign %s to real array element", valueTypeName(value)), context))
		}
	}

	array.Elements[index-array.Start] = value
	i.writeBack(assignNode.ArrayAccess.Left, array, context)
	return res.Success(value)
}

func (i *Interpreter) VisitMemberAccessNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	accessNode, ok := node.(common.MemberAccessNode)
	if !ok {
		return res.Failure(common.RTError(safePos(node.GetPosStart()), safePos(node.GetPosEnd()), "Expected member access node", context))
	}
	object := res.Register(i.Visit(accessNode.Object, context))
	if res.ShouldReturn() {
		return res
	}

	structVal, ok := object.(common.Struct)
	if !ok {
		return res.Failure(common.RTError(safePos(accessNode.Object.GetPosStart()), safePos(accessNode.Object.GetPosEnd()), "Object is not a struct", context))
	}

	val, ok := structVal.Fields[accessNode.MemberTok.Value]
	if !ok {
		return res.Failure(common.RTError(safePos(accessNode.MemberTok.PosStart), safePos(accessNode.MemberTok.PosEnd), fmt.Sprintf("Field '%s' not found in struct", accessNode.MemberTok.Value), context))
	}

	return res.Success(val)
}

func (i *Interpreter) VisitMemberAssignNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	assignNode, ok := node.(common.MemberAssignNode)
	if !ok {
		return res.Failure(common.RTError(safePos(node.GetPosStart()), safePos(node.GetPosEnd()), "Expected member assign node", context))
	}

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
		return res.Failure(common.RTError(safePos(assignNode.MemberAccess.Object.GetPosStart()), safePos(assignNode.MemberAccess.Object.GetPosEnd()), "Object is not a struct", context))
	}

	fieldName := assignNode.MemberAccess.MemberTok.Value
	if _, exists := structVal.Fields[fieldName]; !exists {
		return res.Failure(common.RTError(safePos(assignNode.MemberAccess.MemberTok.PosStart), safePos(assignNode.MemberAccess.MemberTok.PosEnd), fmt.Sprintf("Field '%s' not found in struct", fieldName), context))
	}

	fieldType := structVal.FieldTypes[fieldName]
	if fieldType == "" {
		if existingVal, exists := structVal.Fields[fieldName]; exists {
			if existingNum, ok := existingVal.(common.Number); ok {
				if existingNum.IsReal {
					fieldType = "real"
				} else {
					fieldType = "integer"
				}
			}
		}
	}

	switch fieldType {
	case "integer":
		if numVal, ok := value.(common.Number); ok {
			if numVal.IsReal || numVal.Value != math.Floor(numVal.Value) {
				return res.Failure(common.RTError(safePos(assignNode.ValueNode.GetPosStart()), safePos(assignNode.ValueNode.GetPosEnd()), fmt.Sprintf("Type mismatch: cannot assign real number to integer field '%s'", fieldName), context))
			}
			numVal.IsReal = false
			value = numVal
		} else {
			return res.Failure(common.RTError(safePos(assignNode.ValueNode.GetPosStart()), safePos(assignNode.ValueNode.GetPosEnd()), fmt.Sprintf("Type mismatch: cannot assign %s to integer field '%s'", valueTypeName(value), fieldName), context))
		}
	case "real":
		if numVal, ok := value.(common.Number); ok {
			numVal.IsReal = true
			value = numVal
		} else {
			return res.Failure(common.RTError(safePos(assignNode.ValueNode.GetPosStart()), safePos(assignNode.ValueNode.GetPosEnd()), fmt.Sprintf("Type mismatch: cannot assign %s to real field '%s'", valueTypeName(value), fieldName), context))
		}
	}

	structVal.Fields[fieldName] = value
	i.writeBack(assignNode.MemberAccess.Object, structVal, context)
	return res.Success(value)
}

func (i *Interpreter) writeBack(target common.Expr, updatedValue common.Value, context *common.Context) {
	switch t := target.(type) {
	case common.VarAccessNode:
		context.GetSymbolTable().Set(t.VarNameTok.Value, updatedValue)
	case common.ArrayIndexNode:
		parentVal := i.Visit(t.Left, context)
		if parentArr, ok := parentVal.(common.Array); ok {
			idxVal := i.Visit(t.Index, context)
			if num, ok := idxVal.(common.Number); ok {
				idx := int(num.Value)
				if idx >= parentArr.Start && idx <= parentArr.End {
					parentArr.Elements[idx-parentArr.Start] = updatedValue
					i.writeBack(t.Left, parentArr, context)
				}
			}
		}
	case common.MemberAccessNode:
		parentVal := i.Visit(t.Object, context)
		if parentStruct, ok := parentVal.(common.Struct); ok {
			parentStruct.Fields[t.MemberTok.Value] = updatedValue
			i.writeBack(t.Object, parentStruct, context)
		}
	}
}

func (i *Interpreter) VisitTypeAliasNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	aliasNode, ok := node.(common.TypeAliasNode)
	if !ok {
		return res.Failure(common.RTError(safePos(node.GetPosStart()), safePos(node.GetPosEnd()), "Expected type alias node", context))
	}
	typeValue := common.Type{
		Definition: aliasNode.TargetType,
		Context:    context,
		Pos_Start:  aliasNode.Pos_Start,
		Pos_End:    aliasNode.Pos_End,
	}
	context.GetSymbolTable().Set(aliasNode.AliasName.Value, typeValue)
	return res.Success(common.Null{})
}

func (i *Interpreter) VisitStructTypeNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	structNode, ok := node.(common.StructTypeNode)
	if !ok {
		return res.Failure(common.RTError(safePos(node.GetPosStart()), safePos(node.GetPosEnd()), "Expected struct type node", context))
	}
	typeValue := common.Type{
		Definition: structNode,
		Context:    context,
		Pos_Start:  structNode.Pos_Start,
		Pos_End:    structNode.Pos_End,
	}
	context.GetSymbolTable().Set(structNode.StructName.Value, typeValue)
	return res.Success(common.Null{})
}

func (i *Interpreter) InitializeType(typeNode common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}

	switch t := typeNode.(type) {
	case common.VarAccessNode:
		typeName := t.VarNameTok.Value
		switch typeName {
		case "integer":
			return res.Success(common.Number{Value: 0, IsReal: false, Context: context})
		case "real":
			return res.Success(common.Number{Value: 0, IsReal: true, Context: context})
		case "string":
			return res.Success(common.String{Value: "", Context: context})
		case "boolean":
			return res.Success(common.Boolean{Value: false, Context: context})
		case "char":
			return res.Success(common.String{Value: " ", Context: context})
		default:
			val := context.GetSymbolTable().Get(typeName)
			if typeDef, ok := val.(common.Type); ok {
				return i.InitializeType(typeDef.Definition, context)
			}
			return res.Failure(common.RTError(safePos(t.GetPosStart()), safePos(t.GetPosEnd()), fmt.Sprintf("Unknown type '%s'", typeName), context))
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

		startNum, ok1 := startVal.(common.Number)
		if !ok1 {
			return res.Failure(common.RTError(safePos(t.StartNode.GetPosStart()), safePos(t.StartNode.GetPosEnd()), "Start index must be a number", context))
		}
		endNum, ok2 := endVal.(common.Number)
		if !ok2 {
			return res.Failure(common.RTError(safePos(t.EndNode.GetPosStart()), safePos(t.EndNode.GetPosEnd()), "End index must be a number", context))
		}

		start := int(startNum.Value)
		end := int(endNum.Value)
		size := end - start + 1
		elements := make([]common.Value, size)

		for idx := 0; idx < size; idx++ {
			elemVal := res.Register(i.InitializeType(t.OfType, context))
			if res.ShouldReturn() {
				return res
			}
			elements[idx] = elemVal
		}

		elemTypeName, _ := i.GetBaseTypeName(t.OfType, context)
		return res.Success(common.Array{
			Elements:    elements,
			ElementType: elemTypeName,
			Start:       start,
			End:         end,
			Context:     context,
			Pos_Start:   t.Pos_Start,
			Pos_End:     t.Pos_End,
		})
	case common.StructTypeNode:
		fields := make(map[string]common.Value)
		fieldTypes := make(map[string]string)
		for _, field := range t.Fields {
			fieldVal := res.Register(i.InitializeType(field.ValueNode, context))
			if res.ShouldReturn() {
				return res
			}
			fields[field.VarName.Value] = fieldVal
			ft, _ := i.GetBaseTypeName(field.ValueNode, context)
			fieldTypes[field.VarName.Value] = ft
		}
		return res.Success(common.Struct{
			Fields:     fields,
			FieldTypes: fieldTypes,
			Context:    context,
			Pos_Start:  t.Pos_Start,
			Pos_End:    t.Pos_End,
		})
	}

	return res.Success(common.Null{})
}

func (i *Interpreter) GetBaseTypeName(typeNode common.Expr, context *common.Context) (string, *common.Error) {
	if t, ok := typeNode.(common.VarAccessNode); ok {
		typeName := t.VarNameTok.Value
		switch typeName {
		case "integer", "real", "string", "boolean", "char":
			return typeName, nil
		default:
			val := context.GetSymbolTable().Get(typeName)
			if typeDef, ok := val.(common.Type); ok {
				return i.GetBaseTypeName(typeDef.Definition, context)
			}
			return "", nil
		}
	}
	return "", nil
}

func valuesAreEqual(v1, v2 common.Value) bool {
	if v1 == nil && v2 == nil {
		return true
	}
	if v1 == nil || v2 == nil {
		return false
	}

	n1, ok1 := v1.(common.Number)
	n2, ok2 := v2.(common.Number)
	if ok1 && ok2 {
		return n1.Value == n2.Value
	}

	b1, ok1 := v1.(common.Boolean)
	b2, ok2 := v2.(common.Boolean)
	if ok1 && ok2 {
		return b1.Value == b2.Value
	}

	s1, ok1 := v1.(common.String)
	s2, ok2 := v2.(common.String)
	if ok1 && ok2 {
		return s1.Value == s2.Value
	}

	return v1.Print() == v2.Print()
}

func (i *Interpreter) VisitDependNode(node common.Expr, context *common.Context) common.Value {
	res := &common.RTResult{}
	dependNode, ok := node.(common.DependNode)
	if !ok {
		return res.Failure(common.RTError(safePos(node.GetPosStart()), safePos(node.GetPosEnd()), "Expected depend node", context))
	}

	var switchVal common.Value
	if dependNode.Value != nil {
		switchVal = res.Register(i.Visit(dependNode.Value, context))
		if res.ShouldReturn() {
			return res
		}
	}

	matched := false
	for _, c := range dependNode.Cases {
		caseVal := res.Register(i.Visit(c.Value, context))
		if res.ShouldReturn() {
			return res
		}

		if switchVal != nil {
			if valuesAreEqual(switchVal, caseVal) {
				matched = true
				val := res.Register(i.Visit(c.Body, context))
				if res.ShouldReturn() {
					return res
				}
				return res.Success(val)
			}
		} else {
			if caseVal.IsTrue() {
				matched = true
				val := res.Register(i.Visit(c.Body, context))
				if res.ShouldReturn() {
					return res
				}
				return res.Success(val)
			}
		}
	}

	if !matched && dependNode.DefaultCase != nil {
		val := res.Register(i.Visit(dependNode.DefaultCase, context))
		if res.ShouldReturn() {
			return res
		}
		return res.Success(val)
	}

	return res.Success(common.Null{})
}
