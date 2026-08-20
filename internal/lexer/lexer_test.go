package lexer

import (
	"testing"
)

func TestTokenizeBasic(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected []TokenKind
	}{
		{
			name:   "integers and identifiers",
			source: "x <- 42",
			expected: []TokenKind{
				IDENTIFIER,
				LEFT_ARROW,
				NUMBER,
				EOF,
			},
		},
		{
			name:   "keywords recognition",
			source: "program Test dictionary x : integer algorithm x <- 10 endprogram",
			expected: []TokenKind{
				PROGRAM,
				IDENTIFIER,
				DICTIONARY,
				IDENTIFIER,
				COLON,
				INTEGER,
				ALGORITHM,
				IDENTIFIER,
				LEFT_ARROW,
				NUMBER,
				ENDPROGRAM,
				EOF,
			},
		},
		{
			name:   "arithmetic and comparison operators",
			source: "a + b * c == d div e",
			expected: []TokenKind{
				IDENTIFIER,
				PLUS,
				IDENTIFIER,
				STAR,
				IDENTIFIER,
				EQUALS,
				IDENTIFIER,
				DIV,
				IDENTIFIER,
				EOF,
			},
		},
		{
			name:   "strings and comments",
			source: "\"hello\"\n// comment\n'world'",
			expected: []TokenKind{
				STRING,
				NEWLINE,
				NEWLINE,
				STRING,
				EOF,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := Tokenize(tt.source, "test.dap")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(tokens) != len(tt.expected) {
				t.Fatalf("expected %d tokens, got %d (tokens: %v)", len(tt.expected), len(tokens), tokens)
			}

			for i, exp := range tt.expected {
				if tokens[i].Kind != exp {
					t.Errorf("token %d: expected %v, got %v", i, exp, tokens[i].Kind)
				}
			}
		})
	}
}

func TestTokenizeUnterminatedString(t *testing.T) {
	_, err := Tokenize(`"unterminated`, "test.dap")
	if err == nil {
		t.Fatal("expected error on unterminated string, got nil")
	}
}
