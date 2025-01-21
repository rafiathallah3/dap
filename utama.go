package main

import (
	"dap/internal/common"
	"dap/internal/interpreter"
	"dap/internal/lexer"
	"dap/internal/parser"
	"fmt"
	"os"
	"regexp"

	"github.com/sanity-io/litter"
)

func main() {
	bytes, _ := os.ReadFile("./program.dap")
	source := string(bytes)

	re := regexp.MustCompile(`\n`)
	newlines := re.FindAllStringIndex(source, -1)

	fmt.Printf("Newline Positions: %v\n", newlines)

	globalSymbolTable := &common.SymbolTable{
		Symbols: make(map[string]common.Value),
	}
	globalSymbolTable.Set("null", common.Number{Value: 0})
	globalSymbolTable.Set("true", common.Number{Value: 1})
	globalSymbolTable.Set("false", common.Number{Value: 0})
	globalSymbolTable.Set("print", common.BuiltInFunction{
		BaseFunction: common.BaseFunction{
			Name: "Print",
		},
	})
	globalSymbolTable.Set("input", common.BuiltInFunction{
		BaseFunction: common.BaseFunction{
			Name: "Input",
		},
	})

	// reader := bufio.NewReader(os.Stdin)
	// for {
	// fmt.Print("stdin > ")
	// source, _ := reader.ReadString('\n')

	// if source == "" {
	// 	continue
	// }

	tokens := lexer.Tokenize(source)

	for _, token := range tokens {
		token.Debug()
	}

	Parser := parser.CreateParser(tokens)
	Ast := Parser.Parse().(*common.ParseResult)
	fmt.Println("########   AST   #########")

	if Ast.Error != nil {
		fmt.Println(Ast.Error.As_string())
	} else {
		Ast.Print()
		litter.Dump(Ast.Node)

		inter := interpreter.Interpreter{}

		context := &common.Context{
			DisplayName: "<program>",
		}

		context.Symbol_Table = globalSymbolTable

		fmt.Println("########   RESULT   #########")
		hasil := inter.Visit(Ast.Node, context).(*common.RTResult)
		if hasil.Error != nil {
			fmt.Println(hasil.Error.As_string())
		} else {
			// hasil.Print()
		}
	}
	// }
}
