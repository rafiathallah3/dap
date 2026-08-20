package parser

import (
	"dap/internal/common"
	"dap/internal/lexer"
	"testing"
)

func TestParserExpression(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{
			name:   "binary addition",
			source: "1 + 2",
		},
		{
			name:   "precedence multiplication before addition",
			source: "1 + 2 * 3",
		},
		{
			name:   "assignment expression",
			source: "x <- 10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := lexer.Tokenize(tt.source, "test.dap")
			if err != nil {
				t.Fatalf("tokenization failed: %v", err)
			}

			p := NewParser(tokens, true)
			programName := "<program>"
			res := p.Parse(&programName)
			pr, ok := res.(*common.ParseResult)
			if !ok {
				t.Fatalf("expected *common.ParseResult, got %T", res)
			}
			if pr.Error != nil {
				t.Fatalf("unexpected parser error: %v", pr.Error)
			}
		})
	}
}

func TestParserFullProgram(t *testing.T) {
	source := `program Test
dictionary
    x : integer
algorithm
    x <- 100
endprogram`

	tokens, err := lexer.Tokenize(source, "test.dap")
	if err != nil {
		t.Fatalf("tokenization failed: %v", err)
	}

	p := NewParser(tokens, false)
	programName := "<program>"
	res := p.Parse(&programName)
	pr, ok := res.(*common.ParseResult)
	if !ok {
		t.Fatalf("expected *common.ParseResult, got %T", res)
	}
	if pr.Error != nil {
		t.Fatalf("unexpected parser error: %v", pr.Error)
	}
}
