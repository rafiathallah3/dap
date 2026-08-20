package common

import (
	"dap/internal/lexer"
	"encoding/json"
	"fmt"
)

// PrintTreeAST buat nampilin visualisasi tree AST dalam bentuk ASCII tree.
func PrintTreeAST(node Expr, indent string, last bool) {
	if node == nil {
		return
	}

	// Kalo node-nya ParseResult, langsung print node dalemnya aja
	if pr, ok := node.(*ParseResult); ok {
		PrintTreeAST(pr.Node, indent, last)
		return
	}

	marker := "├── "
	if last {
		marker = "└── "
	}

	fmt.Print(indent)
	fmt.Print(marker)

	info := ""
	var children []Expr

	switch n := node.(type) {
	case NumberNode:
		info = fmt.Sprintf(": %v", n.Token.Value)
	case StringNode:
		info = fmt.Sprintf(": \"%s\"", n.Token.Value)
	case VarAccessNode:
		info = fmt.Sprintf(": %s", n.VarNameTok.Value)
	case BinOpNode:
		info = fmt.Sprintf(": %s", n.Operator.Value)
		children = []Expr{n.Left, n.Right}
	case UnaryOpNode:
		info = fmt.Sprintf(": %s", n.Operator.Value)
		children = []Expr{n.Node}
	case VarAssignNode:
		info = fmt.Sprintf(": %s", n.VarName.Value)
		children = []Expr{n.ValueNode}
	case ListNode:
		children = append(children, n.ElementNode...)
	case DictionaryNode:
		children = append(children, n.VariableDiBuat...)
	case IfNode:
		for _, c := range n.Cases {
			children = append(children, c.Kondisi, c.Isi)
		}
		if n.Else_case != nil && n.Else_case.Isi != nil {
			children = append(children, n.Else_case.Isi)
		}
	case WhileNode:
		children = []Expr{n.KondisiNode, n.BodyNode}
	case ForNode:
		info = fmt.Sprintf(": %s", n.VarNameTok.Value)
		children = []Expr{n.StartValueNode, n.EndValueNode, n.StepValueNode, n.BodyNode}
	case RepeatNode:
		children = []Expr{n.BodyNode, n.KondisiNode}
	case FuncNode:
		name := "<anonymous>"
		if n.VarNameTok != nil {
			name = n.VarNameTok.Value
		}
		args := ""
		for i, arg := range n.ArgNameToks {
			args += arg.Value
			if i < len(n.ArgNameToks)-1 {
				args += ", "
			}
		}
		info = fmt.Sprintf(": %s(%s)", name, args)
		children = []Expr{n.BodyNode}
	case CallNode:
		children = append(children, n.NodeToCall)
		children = append(children, n.ArgNodes...)
	case ArrayIndexNode:
		children = []Expr{n.Left, n.Index}
	case ArrayAssignNode:
		children = []Expr{n.ArrayAccess, n.ValueNode}
	case ArrayTypeNode:
		children = []Expr{n.StartNode, n.EndNode, n.OfType}
	case MemberAccessNode:
		info = fmt.Sprintf(": .%s", n.MemberTok.Value)
		children = []Expr{n.Object}
	case MemberAssignNode:
		children = []Expr{n.MemberAccess, n.ValueNode}
	case TypeAliasNode:
		info = fmt.Sprintf(": %s", n.AliasName.Value)
		children = []Expr{n.TargetType}
	case StructTypeNode:
		info = fmt.Sprintf(": %s", n.StructName.Value)
		for _, field := range n.Fields {
			children = append(children, field)
		}
	case DependNode:
		if n.Value != nil {
			children = append(children, n.Value)
		}
		for _, c := range n.Cases {
			children = append(children, c.Value, c.Body)
		}
		if n.DefaultCase != nil {
			children = append(children, n.DefaultCase)
		}
	case ReturnNode:
		if n.NodeToReturn != nil {
			children = []Expr{n.NodeToReturn}
		}
	}

	fmt.Printf("%s%s\n", node.Name(), info)

	newIndent := indent
	if last {
		newIndent += "    "
	} else {
		newIndent += "│   "
	}

	for i, child := range children {
		isLast := i == len(children)-1
		PrintTreeAST(child, newIndent, isLast)
	}
}

// PrintJSONAST buat nge-print AST dalam format JSON yang rapi.
func PrintJSONAST(node Expr) {
	m := ASTToMap(node)
	if m == nil {
		fmt.Println("null")
		return
	}
	bytes, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling AST to JSON:", err)
		return
	}
	fmt.Println(string(bytes))
}

// ASTToMap buat convert hierarki AST node jadi nested map biar gampang di-serialize ke JSON.
func ASTToMap(node Expr) any {
	if node == nil {
		return nil
	}

	if pr, ok := node.(*ParseResult); ok {
		return ASTToMap(pr.Node)
	}

	res := make(map[string]any)
	res["type"] = node.Name()

	switch n := node.(type) {
	case Token:
		res["value"] = n.Value
		res["kind"] = lexer.TokenKindString(n.Kind)
	case NumberNode:
		res["value"] = n.Token.Value
	case StringNode:
		res["value"] = n.Token.Value
	case VarAccessNode:
		res["name"] = n.VarNameTok.Value
	case BinOpNode:
		res["operator"] = n.Operator.Value
		res["left"] = ASTToMap(n.Left)
		res["right"] = ASTToMap(n.Right)
	case UnaryOpNode:
		res["operator"] = n.Operator.Value
		res["node"] = ASTToMap(n.Node)
	case VarAssignNode:
		res["name"] = n.VarName.Value
		res["value"] = ASTToMap(n.ValueNode)
		res["is_const"] = n.ApakahConst
		res["is_declaration"] = n.ApakahDeklarasi
	case ListNode:
		var program []any
		for _, el := range n.ElementNode {
			program = append(program, ASTToMap(el))
		}
		res["program"] = program
	case DictionaryNode:
		var variables []any
		for _, v := range n.VariableDiBuat {
			variables = append(variables, ASTToMap(v))
		}
		res["variables"] = variables
	case IfNode:
		var cases []any
		for _, c := range n.Cases {
			cases = append(cases, map[string]any{
				"condition": ASTToMap(c.Kondisi),
				"body":      ASTToMap(c.Isi),
			})
		}
		res["cases"] = cases
		if n.Else_case != nil && n.Else_case.Isi != nil {
			res["else_case"] = map[string]any{
				"body": ASTToMap(n.Else_case.Isi),
			}
		} else {
			res["else_case"] = nil
		}
	case WhileNode:
		res["condition"] = ASTToMap(n.KondisiNode)
		res["body"] = ASTToMap(n.BodyNode)
	case ForNode:
		res["var_name"] = n.VarNameTok.Value
		res["start"] = ASTToMap(n.StartValueNode)
		res["end"] = ASTToMap(n.EndValueNode)
		res["step"] = ASTToMap(n.StepValueNode)
		res["body"] = ASTToMap(n.BodyNode)
		res["is_down"] = n.IsDown
	case RepeatNode:
		res["body"] = ASTToMap(n.BodyNode)
		res["condition"] = ASTToMap(n.KondisiNode)
	case FuncNode:
		name := "<anonymous>"
		if n.VarNameTok != nil {
			name = n.VarNameTok.Value
		}
		res["name"] = name
		var args []string
		for _, arg := range n.ArgNameToks {
			args = append(args, arg.Value)
		}
		res["args"] = args
		res["body"] = ASTToMap(n.BodyNode)
	case CallNode:
		res["call"] = ASTToMap(n.NodeToCall)
		var args []any
		for _, arg := range n.ArgNodes {
			args = append(args, ASTToMap(arg))
		}
		res["args"] = args
	case ArrayIndexNode:
		res["array"] = ASTToMap(n.Left)
		res["index"] = ASTToMap(n.Index)
	case ArrayAssignNode:
		res["target"] = ASTToMap(n.ArrayAccess)
		res["value"] = ASTToMap(n.ValueNode)
	case ArrayTypeNode:
		res["start"] = ASTToMap(n.StartNode)
		res["end"] = ASTToMap(n.EndNode)
		res["type"] = ASTToMap(n.OfType)
	case MemberAccessNode:
		res["object"] = ASTToMap(n.Object)
		res["member"] = n.MemberTok.Value
	case MemberAssignNode:
		res["target"] = ASTToMap(n.MemberAccess)
		res["value"] = ASTToMap(n.ValueNode)
	case TypeAliasNode:
		res["alias_name"] = n.AliasName.Value
		res["target_type"] = ASTToMap(n.TargetType)
	case StructTypeNode:
		res["struct_name"] = n.StructName.Value
		var fields []any
		for _, field := range n.Fields {
			fields = append(fields, ASTToMap(field))
		}
		res["fields"] = fields
	case DependNode:
		if n.Value != nil {
			res["value"] = ASTToMap(n.Value)
		}
		var cases []any
		for _, c := range n.Cases {
			cases = append(cases, map[string]any{
				"value": ASTToMap(c.Value),
				"body":  ASTToMap(c.Body),
			})
		}
		res["cases"] = cases
		if n.DefaultCase != nil {
			res["default_case"] = ASTToMap(n.DefaultCase)
		}
	case ReturnNode:
		if n.NodeToReturn != nil {
			res["value"] = ASTToMap(n.NodeToReturn)
		}
	}

	return res
}
