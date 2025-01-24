package main

import (
	"dap/internal/common"
	"dap/internal/interpreter"
	"dap/internal/lexer"
	"dap/internal/parser"
	"dap/tools"
	"fmt"
	"os"
)

/*
program GiveMeArray

dictionary
    i, n, total : integer
algorithm
    total <- 1
    input i, n

    while ((n != -99999) and (total < i)) do
        total <- total + 1
        input n
    endwhile

    if ((total <= i) and (n != -99999)) then
        output n
    else
        output "EMPTY"
    endif
endprogram
*/

func main() {
	fileName := "program.dap"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	bytes, err := os.ReadFile(fmt.Sprintf("./%s", fileName))

	if err != nil {
		panic(fmt.Sprintf("File %s not found!", fileName))
	}

	source := string(bytes)

	globalSymbolTable := &common.SymbolTable{
		Symbols: make(map[string]common.Value),
	}
	globalSymbolTable.Set("null", common.Null{})
	globalSymbolTable.Set("true", common.Number{Value: 1})
	globalSymbolTable.Set("false", common.Number{Value: 0})

	for Keyword, NamaFunction := range tools.SemuaBuiltInFunction {
		globalSymbolTable.Set(Keyword, common.BuiltInFunction{
			BaseFunction: common.BaseFunction{
				Name: NamaFunction,
			},
		})
	}

	tokens := lexer.Tokenize(source, fileName)

	// for _, token := range tokens {
	// 	token.Debug()
	// }

	ProgramName := "<program>"
	Parser := parser.CreateParser(tokens)
	Ast := Parser.Parse(&ProgramName).(*common.ParseResult)
	// fmt.Println("########   AST   #########")

	if Ast.Error != nil {
		fmt.Println(Ast.Error.As_string())
	} else {
		// Ast.Print()
		// litter.Dump(Ast.Node)

		inter := interpreter.Interpreter{}

		context := &common.Context{
			DisplayName: ProgramName,
		}

		context.Symbol_Table = globalSymbolTable

		// fmt.Println("########   RESULT   #########")
		hasil := inter.Visit(Ast.Node, context).(*common.RTResult)
		if hasil.Error != nil {
			fmt.Println(hasil.Error.As_string())
		} else {
			// hasil.Print()
		}
	}
	// }
}
