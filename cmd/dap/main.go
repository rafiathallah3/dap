package main

import (
	"bufio"
	"dap/internal/common"
	"dap/internal/interpreter"
	"dap/internal/lexer"
	"dap/internal/parser"
	"dap/internal/updater"
	"dap/tools"
	"fmt"
	"os"
	"strings"
)

var (
	// Version nyimpen versi aplikasi DAP sekarang.
	Version = "v1.0.0"

	showToken   = false
	showAST     = false
	showASTJSON = false

	globalSymbolTable = &common.SymbolTable{
		Symbols: make(map[string]common.Value),
	}
)

// RunProgram buat ngejalanin source code program DAP dari string.
func RunProgram(source string, fileName string, isSingleLine bool) {
	tokens, err := lexer.Tokenize(source, fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	if showToken {
		fmt.Println("########   TOKEN   #########")
		for _, token := range tokens {
			token.Debug()
		}
	}

	programName := "<program>"
	p := parser.NewParser(tokens, isSingleLine)
	parsed := p.Parse(&programName)
	ast, ok := parsed.(*common.ParseResult)
	if !ok {
		ast = &common.ParseResult{Node: parsed}
	}

	if showAST {
		fmt.Println("########   AST   #########")
	}

	if ast.Error != nil {
		fmt.Println(ast.Error.AsString())
		return
	}

	if showAST {
		common.PrintTreeAST(ast.Node, "", true)
	} else if showASTJSON {
		common.PrintJSONAST(ast.Node)
		return
	}

	inter := interpreter.Interpreter{}
	ctx := &common.Context{
		DisplayName:  programName,
		SymbolTable:  globalSymbolTable,
		Symbol_Table: globalSymbolTable,
	}

	if showAST || showASTJSON || showToken {
		fmt.Println("########   RESULT   #########")
	}

	visitRes := inter.Visit(ast.Node, ctx)
	if rtRes, ok := visitRes.(*common.RTResult); ok {
		if rtRes.Error != nil {
			fmt.Println(rtRes.Error.AsString())
		}
	}
}

// JalaninProgram alias buat backward compatibility.
func JalaninProgram(source string, fileName string, apakahSatuBaris bool) {
	RunProgram(source, fileName, apakahSatuBaris)
}

func initSymbolTable() {
	globalSymbolTable.Set("null", common.Null{})
	globalSymbolTable.Set("true", common.Boolean{Value: true})
	globalSymbolTable.Set("false", common.Boolean{Value: false})
	globalSymbolTable.Set("integer", common.Number{Value: 0, IsReal: false})
	globalSymbolTable.Set("real", common.Number{Value: 0, IsReal: true})
	globalSymbolTable.Set("string", common.String{Value: ""})
	globalSymbolTable.Set("boolean", common.Boolean{Value: false})
	globalSymbolTable.Set("char", common.String{Value: " "})

	for keyword, fnName := range tools.BuiltinFunctions {
		globalSymbolTable.Set(keyword, common.BuiltInFunction{
			BaseFunction: common.BaseFunction{
				Name: fnName,
			},
		})
	}
}

func main() {
	initSymbolTable()

	fileName := ""
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		switch arg {
		case "--show-token":
			showToken = true
		case "--show-ast":
			showAST = true
		case "--show-ast-json":
			showASTJSON = true
		case "--version", "-v":
			updater.CheckVersion(Version)
			os.Exit(0)
		case "--update", "-u":
			if err := updater.PerformSelfUpdate(Version, false); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			os.Exit(0)
		case "--help", "-h":
			fmt.Println("DAP, A friendly Pseudocode for you to learn basic logic")
			fmt.Println("Usage:")
			fmt.Println("  dap [file.dap]       Run a DAP program file")
			fmt.Println("  dap                  Enter interactive console mode")
			fmt.Println("")
			fmt.Println("Options:")
			fmt.Println("  --version, -v        Show version and check for updates")
			fmt.Println("  --update, -u         Update DAP to the latest release")
			fmt.Println("  --show-token         Show tokens during execution")
			fmt.Println("  --show-ast           Show abstract syntax tree during execution")
			fmt.Println("  --show-ast-json      Show abstract syntax tree in JSON format")
			fmt.Println("  --help, -h           Show this help message")
			os.Exit(0)
		default:
			if strings.HasPrefix(arg, "-") {
				fmt.Fprintf(os.Stderr, "Error: Unknown option '%s'. Use --help for usage.\n", arg)
				os.Exit(1)
			}
			if fileName == "" {
				fileName = arg
			}
		}
	}

	if fileName == "" {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Welcome to DAP. Friendly Pseudocode.")
		fmt.Println("Type `help` for information.")

		for {
			fmt.Print(">>> ")
			if !scanner.Scan() {
				break
			}
			text := strings.TrimSpace(strings.TrimPrefix(scanner.Text(), "\ufeff"))

			if text == "" {
				continue
			}

			if text == "exit" || text == "exit()" {
				fmt.Println("Jumpe lagi!")
				os.Exit(0)
			}

			if text == "help" {
				fmt.Println("DAP, A friendly Pseudocode for you to learn basic logic")
				fmt.Println("Interactive Commands:")
				fmt.Println("  help                 Show this help message")
				fmt.Println("  version              Show DAP version and check for updates")
				fmt.Println("  update               Update DAP to the latest release")
				fmt.Println("  exit, exit()         Exit interactive console")
				fmt.Println("")
				fmt.Println("File Execution:")
				fmt.Println("  dap [file.dap]       Run a file from your terminal")
				continue
			}

			if text == "version" {
				updater.CheckVersion(Version)
				continue
			}

			if text == "update" {
				if err := updater.PerformSelfUpdate(Version, false); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				}
				continue
			}

			RunProgram(text, "<program>", true)
		}
		return
	}

	bytes, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: File '%s' not found or cannot be read.\n", fileName)
		os.Exit(1)
	}

	source := string(bytes)
	RunProgram(source, fileName, false)
}
