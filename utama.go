package main

import (
	"bufio"
	"dap/internal/common"
	"dap/internal/interpreter"
	"dap/internal/lexer"
	"dap/internal/parser"
	"dap/tools"
	"fmt"
	"os"

	"github.com/sanity-io/litter"
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

var TunjuinToken = false
var TunjuinAST = false
var globalSymbolTable = &common.SymbolTable{
	Symbols: make(map[string]common.Value),
}

func JalaninProgram(source string, fileName string, ApakahSatuBaris bool) {
	tokens := lexer.Tokenize(source, fileName)

	if TunjuinToken {
		fmt.Println("########   TOKEN   #########")
		for _, token := range tokens {
			token.Debug()
		}
	}

	ProgramName := "<program>"
	Parser := parser.CreateParser(tokens, ApakahSatuBaris)
	Ast := Parser.Parse(&ProgramName).(*common.ParseResult)

	if TunjuinAST {
		fmt.Println("########   AST   #########")
	}

	if Ast.Error != nil {
		fmt.Println(Ast.Error.As_string())
	} else {
		if TunjuinAST {
			litter.Dump(Ast.Node)
		}

		inter := interpreter.Interpreter{}
		context := &common.Context{
			DisplayName: ProgramName,
		}

		context.Symbol_Table = globalSymbolTable

		if TunjuinAST || TunjuinToken {
			fmt.Println("########   RESULT   #########")
		}

		hasil := inter.Visit(Ast.Node, context).(*common.RTResult)
		if hasil.Error != nil {
			fmt.Println(hasil.Error.As_string())
		}
	}
}

func main() {
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

	fileName := ""
	for i, command := range os.Args {
		if i == 1 && command[:2] != "--" {
			fileName = command
		}

		if command == "--show-token" {
			TunjuinToken = true
		}

		if command == "--show-ast" {
			TunjuinAST = true
		}
	}

	if fileName == "" {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Welcome to DAP. Friendly Psuedocode.")
		fmt.Println("Type `help` for information.")

		for {
			fmt.Print(">>> ")
			scanner.Scan()
			text := scanner.Text()

			if text == "" {
				continue
			}

			if text == "exit" || text == "exit()" {
				os.Exit(1)
			}

			if text == "help" {
				fmt.Println("This is DAP! A friendly Psuedocode for you to learn basic logic made by Rafi")
				fmt.Println("To run the command `dap [NameOfAFile.dap]` Without the namefile would access to console")
				fmt.Println("Extra command")
				fmt.Println(" --show-token | Showing the Token")
				fmt.Println(" --show-AST | Showing the AST")
				continue
			}

			JalaninProgram(text, "<program>", true)
		}
	}

	bytes, err := os.ReadFile(fmt.Sprintf("./%s", fileName))

	if err != nil {
		panic(fmt.Sprintf("File %s not found!", fileName))
	}

	source := string(bytes)
	JalaninProgram(source, fileName, false)
}
