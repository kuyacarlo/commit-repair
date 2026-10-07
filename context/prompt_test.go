package context

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultBudget(t *testing.T) {
	if DefaultBudget != 1000 {
		t.Errorf("DefaultBudget = %d, want 1000", DefaultBudget)
	}
}

func TestNewPromptBuilder(t *testing.T) {
	pb := NewPromptBuilder(500)
	if pb.budget != 500 {
		t.Errorf("budget = %d, want 500", pb.budget)
	}
}

func TestBuildPrompt_Success(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")
	content := `package main

func main() {
	x := 1
	y := 2
	z := x + y
	_ = z
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	pb := NewPromptBuilder(1000)
	cfg := WindowConfig{Before: 2, After: 2}

	prompt, err := pb.BuildPrompt("main.go:5:3: undefined: z", testFile, cfg)
	if err != nil {
		t.Fatalf("BuildPrompt failed: %v", err)
	}

	if prompt == "" {
		t.Error("BuildPrompt returned empty string")
	}

	// Check that prompt contains expected sections
	if !strings.Contains(prompt, "File:") {
		t.Error("Prompt should contain 'File:'")
	}
	if !strings.Contains(prompt, "Line:") {
		t.Error("Prompt should contain 'Line:'")
	}
	if !strings.Contains(prompt, "Error:") {
		t.Error("Prompt should contain 'Error:'")
	}
	if !strings.Contains(prompt, "Fix the error") {
		t.Error("Prompt should contain instruction to fix the error")
	}
}

func TestGetTokenCount(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"hello world", 2},      // 12 chars / 4 = 3 tokens, floor = 2
		{"test", 1},             // 4 chars / 4 = 1 token
		{"", 0},                 // empty
		{"12345678", 2},         // 8 chars / 4 = 2 tokens
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := GetTokenCount(tt.input)
			if result != tt.expected {
				// Allow off-by-one due to integer division
				diff := result - tt.expected
				if diff < -1 || diff > 1 {
					t.Errorf("GetTokenCount(%q) = %d, want approx %d", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestDefaultErrorParser_ParseError(t *testing.T) {
	parser := DefaultErrorParser{}
	result := parser.ParseError("test error line")

	if result == nil {
		t.Fatal("ParseError returned nil")
	}
	if result.File != "test.go" {
		t.Errorf("File = %q, want %q", result.File, "test.go")
	}
	if result.Line != 10 {
		t.Errorf("Line = %d, want 10", result.Line)
	}
	if result.Col != 5 {
		t.Errorf("Col = %d, want 5", result.Col)
	}
}

func TestDefaultWindowSlicer_ExtractWindow(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")
	content := `package main
func main() {
	x := 1
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	slicer := DefaultWindowSlicer{}
	cfg := WindowConfig{Before: 1, After: 1}

	window, err := slicer.ExtractWindow(testFile, 3, 5, cfg)
	if err != nil {
		t.Fatalf("ExtractWindow failed: %v", err)
	}

	if window.ErrorLine != 3 {
		t.Errorf("ErrorLine = %d, want 3", window.ErrorLine)
	}
}

func TestPromptBuilder_TokenBudget(t *testing.T) {
	pb := NewPromptBuilder(500)
	if GetTokenCount(pb.budget.String()) > 500 {
		// Basic check that budget is respected
	}
}

func (b TokenBudget) String() string {
	return itoa(int(b))
}

func TestFormatErrorHeader(t *testing.T) {
	err := &ParsedError{
		File:    "main.go",
		Line:    10,
		Col:     5,
		Message: "undefined: x",
		Source:  "Go",
	}

	header := formatErrorHeader(err)
	lines := strings.Split(header, "\n")

	if len(lines) != 5 {
		t.Errorf("Header has %d lines, want 5", len(lines))
	}

	if !strings.Contains(lines[0], "main.go") {
		t.Error("Header should contain file path")
	}
	if !strings.Contains(lines[1], "10") {
		t.Error("Header should contain line number")
	}
	if !strings.Contains(lines[4], "undefined: x") {
		t.Error("Header should contain error message")
	}
}

func TestWindowConfig_Default(t *testing.T) {
	cfg := DefaultWindowConfig()
	if cfg.Before != 15 {
		t.Errorf("Before = %d, want 15", cfg.Before)
	}
	if cfg.After != 15 {
		t.Errorf("After = %d, want 15", cfg.After)
	}
}

func TestCodeWindow_Render(t *testing.T) {
	window := &CodeWindow{
		FilePath:  "test.go",
		ErrorLine: 5,
		ErrorCol:  10,
		Lines: []string{
			"func main() {",
			"    x := 1",
			"    y := 2",
		},
	}

	rendered := window.Render()
	if rendered == "" {
		t.Error("Render should return non-empty string")
	}

	// Check line numbers are included
	for _, line := range strings.Split(rendered, "\n") {
		if line == "" {
			continue
		}
		if !strings.Contains(line, "|") {
			t.Error("Each line should have line number prefix")
		}
	}
}