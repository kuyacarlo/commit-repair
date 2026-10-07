package slicer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultWindowConfig(t *testing.T) {
	cfg := DefaultWindowConfig()
	if cfg.Before != 15 {
		t.Errorf("Before = %d, want 15", cfg.Before)
	}
	if cfg.After != 15 {
		t.Errorf("After = %d, want 15", cfg.After)
	}
}

func TestExtractWindow_Success(t *testing.T) {
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

func helper() {
	_ = "test"
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cfg := WindowConfig{Before: 2, After: 2}
	window, err := ExtractWindow(testFile, 5, 10, cfg)
	if err != nil {
		t.Fatalf("ExtractWindow failed: %v", err)
	}

	if window.FilePath != testFile {
		t.Errorf("FilePath = %q, want %q", window.FilePath, testFile)
	}
	if window.ErrorLine != 5 {
		t.Errorf("ErrorLine = %d, want 5", window.ErrorLine)
	}
	if window.ErrorCol != 10 {
		t.Errorf("ErrorCol = %d, want 10", window.ErrorCol)
	}
	// Should have 5 lines: 3-7 (2 before, 2 after, and line 5)
	if len(window.Lines) != 5 {
		t.Errorf("len(Lines) = %d, want 5", len(window.Lines))
	}
}

func TestExtractWindow_FileNotFound(t *testing.T) {
	cfg := DefaultWindowConfig()
	_, err := ExtractWindow("/nonexistent/file.go", 10, 5, cfg)
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestExtractWindow_StartOfFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")
	content := `package main
func main() {}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Error at line 1 with large before window
	cfg := WindowConfig{Before: 15, After: 2}
	window, err := ExtractWindow(testFile, 1, 1, cfg)
	if err != nil {
		t.Fatalf("ExtractWindow failed: %v", err)
	}

	// Should handle start of file gracefully
	if window.ErrorLine != 1 {
		t.Errorf("ErrorLine = %d, want 1", window.ErrorLine)
	}
}

func TestCodeWindow_Render(t *testing.T) {
	window := &CodeWindow{
		FilePath:  "test.go",
		ErrorLine: 5,
		ErrorCol:  10,
		Lines: []string{
			"func main() {",
			"\tx := 1",
			"\ty := 2",
			"\tz := x + y",
			"\t_ = z",
			"}",
		},
	}

	rendered := window.Render()
	// Check that line numbers are included
	if !contains(rendered, "   5") && !contains(rendered, "5 |") {
		t.Error("Rendered output should contain line numbers")
	}
}

func TestCodeWindow_EmptyLines(t *testing.T) {
	window := &CodeWindow{
		FilePath:  "empty.go",
		ErrorLine: 1,
		ErrorCol:  1,
		Lines:     []string{},
	}

	rendered := window.Render()
	if rendered != "" {
		t.Errorf("Empty window should render to empty string, got %q", rendered)
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && 
		(s == substr || len(s) >= len(substr) && 
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
		containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestWindowConfig_DefaultValues(t *testing.T) {
	cfg := WindowConfig{}
	if cfg.Before != 0 {
		t.Errorf("Default Before = %d, want 0", cfg.Before)
	}
	if cfg.After != 0 {
		t.Errorf("Default After = %d, want 0", cfg.After)
	}
}