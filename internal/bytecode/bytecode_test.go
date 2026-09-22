package bytecode_test

import (
	"bytes"
	"dap/internal/bytecode"
	"dap/internal/common"
	"dap/internal/lexer"
	"dap/internal/parser"
	"io"
	"os"
	"strings"
	"testing"
)

func parseDAP(source string) common.Expr {
	tokens, err := lexer.Tokenize(source, "<test>")
	if err != nil {
		panic(err)
	}
	progName := "test"
	p := parser.NewParser(tokens, false)
	parsed := p.Parse(&progName)
	if pr, ok := parsed.(*common.ParseResult); ok {
		if pr.Error != nil {
			panic(pr.Error.AsString())
		}
		return pr.Node
	}
	return parsed
}

func captureStdout(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return strings.TrimSpace(buf.String())
}

func TestBytecodeCompilerAndVM(t *testing.T) {
	testCases := []struct {
		name     string
		source   string
		expected string
	}{
		{
			name: "Basic Arithmetic & Variables",
			source: `program test1
kamus
    a, b, c : integer
algoritma
    a <- 10
    b <- 25
    c <- a + b * 2
    output(c)
endprogram`,
			expected: "60",
		},
		{
			name: "Bitwise Operations on VM",
			source: `program testBit
kamus
    a, b, c : integer
algoritma
    a <- 12
    b <- 10
    c <- (a & b) + (a | b) + (a ^ b)
    output(c)
endprogram`,
			expected: "28", // 8 + 14 + 6 = 28
		},
		{
			name: "Shift Operations on VM",
			source: `program testShift
kamus
    a, b : integer
algoritma
    a <- 5 << 3   // 40
    b <- a >> 2   // 10
    output(b)
endprogram`,
			expected: "10",
		},
		{
			name: "Conditionals If Else on VM",
			source: `program testIf
kamus
    x, y : integer
algoritma
    x <- 15
    if x > 10 then
        y <- 100
    else
        y <- 200
    endif
    output(y)
endprogram`,
			expected: "100",
		},
		{
			name: "While Loop on VM",
			source: `program testWhile
kamus
    i, total : integer
algoritma
    i <- 1
    total <- 0
    while i <= 5 do
        total <- total + i
        i <- i + 1
    endwhile
    output(total)
endprogram`,
			expected: "15",
		},
		{
			name: "For Loop with Step on VM",
			source: `program testFor
kamus
    i, sum : integer
algoritma
    sum <- 0
    for i <- 1 to 7 step 2 do
        sum <- sum + i
    endfor
    output(sum)
endprogram`,
			expected: "16", // 1 + 3 + 5 + 7 = 16
		},
		{
			name: "Real Division with Integer Operands on VM",
			source: `program testDiv
kamus
    a, b : integer
    c : real
algoritma
    a <- 7
    b <- 2
    c <- a / b
    output(c)
endprogram`,
			expected: "3.5",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ast := parseDAP(tc.source)
			chunk, err := bytecode.Compile(ast)
			if err != nil {
				t.Fatalf("Compile error: %v", err)
			}

			// Verify disassembler doesn't crash
			dis := bytecode.Disassemble(chunk)
			if len(dis) == 0 {
				t.Errorf("Disassembler returned empty string")
			}

			vm := bytecode.NewVM()
			got := captureStdout(func() {
				if err := vm.Run(chunk); err != nil {
					t.Fatalf("VM Runtime error: %v", err)
				}
			})

			if got != tc.expected {
				t.Errorf("Output mismatch. Expected %q, got %q", tc.expected, got)
			}
		})
	}
}
