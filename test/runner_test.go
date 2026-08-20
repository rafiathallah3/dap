package test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestDapPrograms(t *testing.T) {
	// 1. Build the compiler binary first
	binaryName := "dap_test_bin"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	binaryPath := filepath.Join("..", binaryName)
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "../cmd/dap")
	buildCmd.Stderr = os.Stderr
	buildCmd.Stdout = os.Stdout
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build dap compiler binary: %v", err)
	}
	defer os.Remove(binaryPath) // clean up binary after tests finish

	// 2. Find all .dap files in the current directory
	files, err := filepath.Glob("*.dap")
	if err != nil {
		t.Fatalf("Failed to glob .dap files: %v", err)
	}

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			// Run the compiled binary on the .dap file
			absBinPath, err := filepath.Abs(binaryPath)
			if err != nil {
				t.Fatalf("Failed to get absolute path of binary: %v", err)
			}

			cmd := exec.Command(absBinPath, file)
			var outBuf bytes.Buffer
			cmd.Stdout = &outBuf
			cmd.Stderr = &outBuf

			// Run the process
			_ = cmd.Run() // Ignore exit code since compiler errors exits with 1 but it's expected

			// Read expected output
			expectedFile := strings.TrimSuffix(file, ".dap") + ".expected"
			expectedBytes, err := os.ReadFile(expectedFile)
			if err != nil {
				t.Fatalf("Failed to read expected file %s: %v", expectedFile, err)
			}

			// Clean up both outputs to normalize newlines and carriage returns
			actual := normalizeText(outBuf.String())
			expected := normalizeText(string(expectedBytes))

			if actual != expected {
				t.Errorf("Mismatch in %s.\nExpected:\n%s\nActual:\n%s", file, expected, actual)
			}
		})
	}
}

func normalizeText(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.TrimSpace(s)
	return s
}
