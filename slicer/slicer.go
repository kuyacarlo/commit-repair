package slicer

import (
	"bufio"
	"os"
	"strings"
)

// WindowConfig defines the window size for code extraction
type WindowConfig struct {
	Before int // lines before the error line
	After  int // lines after the error line
}

// DefaultWindowConfig returns the default window configuration
func DefaultWindowConfig() WindowConfig {
	return WindowConfig{Before: 15, After: 15}
}

// ExtractWindow extracts a code window around the error line
func ExtractWindow(filePath string, line, col int, cfg WindowConfig) (*CodeWindow, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	currentLine := 0

	for scanner.Scan() {
		currentLine++
		if currentLine >= line-cfg.Before && currentLine <= line+cfg.After {
			lines = append(lines, scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Calculate relative position within the window
	relativeLine := line - (line - cfg.Before)
	if relativeLine < 0 {
		relativeLine = 0
	}

	return &CodeWindow{
		FilePath:  filePath,
		ErrorLine: line,
		ErrorCol:  col,
		Lines:     lines,
	}, nil
}

// CodeWindow represents a code window around an error
type CodeWindow struct {
	FilePath  string
	ErrorLine int
	ErrorCol  int
	Lines     []string
}

// Render renders the code window with line numbers
func (w *CodeWindow) Render() string {
	var builder strings.Builder
	
	for i, line := range w.Lines {
		lineNum := w.ErrorLine - len(w.Lines)/2 + i
		builder.WriteString(npad(lineNum, 4) + " | " + line + "\n")
	}
	
	return builder.String()
}

func npad(n, width int) string {
	s := string(rune(n + '0'))
	pad := width - len(s)
	if pad <= 0 {
		return s
	}
	return strings.Repeat(" ", pad) + s
}