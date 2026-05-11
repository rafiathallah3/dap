package common

import (
	"fmt"
)

func PrintTreeAST(node Expr, indent string, last bool) {
	if node == nil {
		return
	}

	// If it's a ParseResult, just print its inner node
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

	// Print name and some info
	info := ""
	var children []Expr
	var childNames []string

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
		childNames = []string{"Left", "Right"}
	case UnaryOpNode:
		info = fmt.Sprintf(": %s", n.Operator.Value)
		children = []Expr{n.Node}
		childNames = []string{"Node"}
	case VarAssignNode:
		info = fmt.Sprintf(": %s", n.VarName.Value)
		children = []Expr{n.ValueNode}
		childNames = []string{"Value"}
	case ListNode:
		for i, v := range n.ElementNode {
			children = append(children, v)
			childNames = append(childNames, fmt.Sprintf("[%d]", i))
		}
	case DictionaryNode:
		for i, v := range n.VariableDiBuat {
			children = append(children, v)
			childNames = append(childNames, fmt.Sprintf("Var[%d]", i))
		}
	case IfNode:
		for i, c := range n.Cases {
			children = append(children, c.Kondisi, c.Isi)
			childNames = append(childNames, fmt.Sprintf("Cond[%d]", i), fmt.Sprintf("Body[%d]", i))
		}
		if n.Else_case != nil {
			children = append(children, n.Else_case.Isi)
			childNames = append(childNames, "Else")
		}
	case WhileNode:
		children = []Expr{n.KondisiNode, n.BodyNode}
		childNames = []string{"Condition", "Body"}
	case ForNode:
		info = fmt.Sprintf(": %s", n.VarNameTok.Value)
		children = []Expr{n.StartValueNode, n.EndValueNode, n.StepValueNode, n.BodyNode}
		childNames = []string{"Start", "End", "Step", "Body"}
	case RepeatNode:
		children = []Expr{n.BodyNode, n.KondisiNode}
		childNames = []string{"Body", "Condition"}
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
		childNames = []string{"Body"}
	case CallNode:
		children = append(children, n.NodeToCall)
		childNames = append(childNames, "Callee")
		for i, arg := range n.ArgNodes {
			children = append(children, arg)
			childNames = append(childNames, fmt.Sprintf("Arg[%d]", i))
		}
	case ArrayIndexNode:
		children = []Expr{n.Left, n.Index}
		childNames = []string{"Array", "Index"}
	case ArrayAssignNode:
		children = []Expr{n.ArrayAccess, n.ValueNode}
		childNames = []string{"Target", "Value"}
	case ArrayTypeNode:
		children = []Expr{n.StartNode, n.EndNode, n.OfType}
		childNames = []string{"Start", "End", "Type"}
	case MemberAccessNode:
		info = fmt.Sprintf(": .%s", n.MemberTok.Value)
		children = []Expr{n.Object}
		childNames = []string{"Object"}
	case MemberAssignNode:
		children = []Expr{n.MemberAccess, n.ValueNode}
		childNames = []string{"Target", "Value"}
	case TypeAliasNode:
		info = fmt.Sprintf(": %s", n.AliasName.Value)
		children = []Expr{n.TargetType}
		childNames = []string{"TargetType"}
	case StructTypeNode:
		info = fmt.Sprintf(": %s", n.StructName.Value)
		for i, field := range n.Fields {
			children = append(children, field)
			childNames = append(childNames, fmt.Sprintf("Field[%d]", i))
		}
	case ReturnNode:
		children = []Expr{n.NodeToReturn}
		childNames = []string{"Value"}
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
