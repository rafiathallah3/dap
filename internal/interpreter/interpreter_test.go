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
