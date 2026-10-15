package interpreter

import (
	"dap/internal/common"
	"dap/internal/lexer"
	"dap/internal/parser"
	"testing"
)

func runEval(t *testing.T, source string) common.Value {
	t.Helper()
	tokens, err := lexer.Tokenize(source, "test.dap")
	if err != nil {
		t.Fatalf("tokenization failed: %v", err)
	}

	p := parser.NewParser(tokens, true)
	progName := "<program>"
	parsed := p.Parse(&progName)
	pr, ok := parsed.(*common.ParseResult)
	if !ok {
		t.Fatalf("expected *common.ParseResult, got %T", parsed)
	}
	if pr.Error != nil {
		t.Fatalf("parser error: %v", pr.Error)
	}

	inter := &Interpreter{}
	ctx := &common.Context{
		DisplayName: progName,
		SymbolTable: common.NewSymbolTable(nil),
	}
	ctx.Symbol_Table = ctx.SymbolTable

	res := inter.Visit(pr.Node, ctx)
	if rtRes, ok := res.(*common.RTResult); ok {
		if rtRes.Error != nil {
			t.Fatalf("runtime error: %v", rtRes.Error)
		}
		res = rtRes.Value
	}

	// Statements list unwrap
	if list, isList := res.(common.List); isList && len(list.Elements) > 0 {
		res = list.Elements[len(list.Elements)-1]
	}

	return res
}

func TestInterpreterArithmetic(t *testing.T) {
	tests := []struct {
		name           string
		source         string
		expected       float64
		expectedIsReal bool
	}{
		{
			name:           "simple addition",
			source:         "10 + 5",
			expected:       15,
			expectedIsReal: false,
		},
		{
			name:           "multiplication precedence",
			source:         "2 + 3 * 4",
			expected:       14,
			expectedIsReal: false,
		},
		{
			name:           "integer division",
			source:         "14 div 3",
			expected:       4,
			expectedIsReal: false,
		},
		{
			name:           "real division with integers",
			source:         "7 / 2",
			expected:       3.5,
			expectedIsReal: true,
		},
		{
			name:           "real division with evenly divisible integers",
			source:         "4 / 2",
			expected:       2,
			expectedIsReal: true,
		},
		{
			name:           "real division mixed int and real",
			source:         "5 / 2.5",
			expected:       2,
			expectedIsReal: true,
		},
		{
			name:           "real division mixed real and int",
			source:         "7.5 / 3",
			expected:       2.5,
			expectedIsReal: true,
		},
		{
			name:           "modulo",
			source:         "14 % 3",
			expected:       2,
			expectedIsReal: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := runEval(t, tt.source)
			num, ok := val.(common.Number)
			if !ok {
				t.Fatalf("expected Number value, got %T (%v)", val, val)
			}
			if num.Value != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, num.Value)
			}
			if num.IsReal != tt.expectedIsReal {
				t.Errorf("expected IsReal %v, got %v", tt.expectedIsReal, num.IsReal)
			}
		})
	}
}

func TestInterpreterBooleanLogic(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected bool
	}{
		{
			name:     "comparison equality",
			source:   "5 == 5",
			expected: true,
		},
		{
			name:     "comparison inequality",
			source:   "5 != 3",
			expected: true,
		},
		{
			name:     "logical and",
			source:   "(5 > 2) and (3 < 4)",
			expected: true,
		},
		{
			name:     "logical or",
			source:   "(5 < 2) or (3 < 4)",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := runEval(t, tt.source)
			b, ok := val.(common.Boolean)
			if !ok {
				t.Fatalf("expected Boolean value, got %T (%v)", val, val)
			}
			if b.Value != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, b.Value)
			}
		})
	}
}

func runProgram(t *testing.T, source string) (*common.Context, *common.Error) {
	t.Helper()
	tokens, err := lexer.Tokenize(source, "test.dap")
	if err != nil {
		t.Fatalf("tokenization failed: %v", err)
	}

	p := parser.NewParser(tokens, false)
	progName := "<program>"
	parsed := p.Parse(&progName)
	pr, ok := parsed.(*common.ParseResult)
	if !ok {
		t.Fatalf("expected *common.ParseResult, got %T", parsed)
	}
	if pr.Error != nil {
		t.Fatalf("parser error: %v", pr.Error)
	}

	st := common.NewSymbolTable(nil)
	st.Set("null", common.Null{})
	st.Set("true", common.Boolean{Value: true})
	st.Set("false", common.Boolean{Value: false})
	st.Set("integer", common.Number{Value: 0, IsReal: false})
	st.Set("real", common.Number{Value: 0, IsReal: true})
	st.Set("string", common.String{Value: ""})
	st.Set("boolean", common.Boolean{Value: false})
	st.Set("char", common.String{Value: " "})

	inter := &Interpreter{}
	ctx := &common.Context{
		DisplayName: progName,
		SymbolTable: st,
	}
	ctx.Symbol_Table = ctx.SymbolTable

	res := inter.Visit(pr.Node, ctx)
	if rtRes, ok := res.(*common.RTResult); ok {
		return ctx, rtRes.Error
	}
	return ctx, nil
}

func TestVariableTypeChecking(t *testing.T) {
	t.Run("assign real to integer variable returns error", func(t *testing.T) {
		src := `program Test
dictionary
    x : integer
algorithm
    x <- 2.5
endprogram`
		_, err := runProgram(t, src)
		if err == nil {
			t.Fatalf("expected error when assigning real to integer variable, got nil")
		}
	})

	t.Run("assign real variable to integer variable returns error", func(t *testing.T) {
		src := `program Test
dictionary
    x : integer
    y : real
algorithm
    y <- 5
    x <- y
endprogram`
		_, err := runProgram(t, src)
		if err == nil {
			t.Fatalf("expected error when assigning real variable to integer variable, got nil")
		}
	})

	t.Run("assign integer to real variable promotes to real", func(t *testing.T) {
		src := `program Test
dictionary
    y : real
algorithm
    y <- 5
endprogram`
		ctx, err := runProgram(t, src)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val := ctx.GetSymbolTable().Get("y")
		num, ok := val.(common.Number)
		if !ok {
			t.Fatalf("expected Number, got %T", val)
		}
		if num.Value != 5 || !num.IsReal {
			t.Errorf("expected 5 with IsReal=true, got %f with IsReal=%v", num.Value, num.IsReal)
		}
	})

	t.Run("assign real to real variable works as intended", func(t *testing.T) {
		src := `program Test
dictionary
    y : real
algorithm
    y <- 2.5
endprogram`
		ctx, err := runProgram(t, src)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val := ctx.GetSymbolTable().Get("y")
		num, ok := val.(common.Number)
		if !ok {
			t.Fatalf("expected Number, got %T", val)
		}
		if num.Value != 2.5 || !num.IsReal {
			t.Errorf("expected 2.5 with IsReal=true, got %f with IsReal=%v", num.Value, num.IsReal)
		}
	})

	t.Run("array element integer and real type checking", func(t *testing.T) {
		errSrc := `program Test
dictionary
    arr : array [1..3] of integer
algorithm
    arr[1] <- 3.5
endprogram`
		_, err := runProgram(t, errSrc)
		if err == nil {
			t.Fatalf("expected error assigning real to integer array element, got nil")
		}

		okSrc := `program Test
dictionary
    arr : array [1..3] of real
algorithm
    arr[1] <- 5
endprogram`
		ctx, err := runProgram(t, okSrc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		arrVal := ctx.GetSymbolTable().Get("arr").(common.Array)
		elem := arrVal.Elements[0].(common.Number)
		if elem.Value != 5 || !elem.IsReal {
			t.Errorf("expected array element 5 with IsReal=true, got %f with IsReal=%v", elem.Value, elem.IsReal)
		}
	})
}
