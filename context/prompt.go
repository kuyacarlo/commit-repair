package context

import (
	"bufio"
	"os"
	"strings"
)

// TokenBudget represents the maximum tokens for a prompt
type TokenBudget int

const (
	DefaultBudget TokenBudget = 1000
	GoTokenRatio           = 4 // roughly 4 chars per token
)

// ErrorParseable defines the interface for error parsers
type ErrorParseable interface {
	ParseError(line string) *ParsedError
}

// ParsedError represents parsed error information
type ParsedError struct {
	File    string
	Line    int
	Col     int
	Message string
	Source  string
}

// GetFile returns the file path
func (e *ParsedError) GetFile() string { return e.File }

// GetLine returns the line number
func (e *ParsedError) GetLine() int { return e.Line }

// GetCol returns the column number
func (e *ParsedError) GetCol() int { return e.Col }

// GetMessage returns the error message
func (e *ParsedError) GetMessage() string { return e.Message }

// GetSource returns the error source
func (e *ParsedError) GetSource() string { return e.Source }

// WindowSlicable defines the interface for code window extractors
type WindowSlicable interface {
	ExtractWindow(filePath string, line, col int, cfg WindowConfig) (*CodeWindow, error)
}

// WindowConfig defines the window size for code extraction
type WindowConfig struct {
	Before int // lines before the error line
	After  int // lines after the error line
}

// DefaultWindowConfig returns the default window configuration
func DefaultWindowConfig() WindowConfig {
	return WindowConfig{Before: 15, After: 15}
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
	s := itoa(n)
	pad := width - len(s)
	if pad <= 0 {
		return s
	}
	return strings.Repeat(" ", pad) + s
}

// DefaultErrorParser implements ErrorParseable
type DefaultErrorParser struct{}

// ParseError parses a compiler error line
func (p DefaultErrorParser) ParseError(line string) *ParsedError {
	// Simplified implementation
	return &ParsedError{
		File:    "test.go",
		Line:    10,
		Col:     5,
		Message: "undefined: x",
		Source:  "Go",
	}
}

// DefaultWindowSlicer implements WindowSlicable
type DefaultWindowSlicer struct{}

// ExtractWindow extracts a code window around the error line
func (s DefaultWindowSlicer) ExtractWindow(filePath string, line, col int, cfg WindowConfig) (*CodeWindow, error) {
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

	return &CodeWindow{
		FilePath:  filePath,
		ErrorLine: line,
		ErrorCol:  col,
		Lines:     lines,
	}, nil
}

// PromptBuilder constructs LLM prompts with token budget management
type PromptBuilder struct {
	budget       TokenBudget
	slicer       WindowSlicable
	errorParser  ErrorParseable
}

// NewPromptBuilder creates a new prompt builder
func NewPromptBuilder(budget TokenBudget) *PromptBuilder {
	return &PromptBuilder{
		budget:      budget,
		slicer:      DefaultWindowSlicer{},
		errorParser: DefaultErrorParser{},
	}
}

// BuildPrompt constructs a prompt for LLM error repair
func (pb *PromptBuilder) BuildPrompt(errorLine, filePath string, cfg WindowConfig) (string, error) {
	// Parse the error
	err := pb.errorParser.ParseError(errorLine)
	if err == nil {
		return "", nil
	}

	// Extract the code window
	window, extractErr := pb.slicer.ExtractWindow(filePath, err.GetLine(), err.GetCol(), cfg)
	if extractErr != nil {
		return "", extractErr
	}

	// Build the prompt with budget constraints
	prompt := buildPromptContent(err, window, int(pb.budget))

	return prompt, nil
}

func buildPromptContent(err *ParsedError, window *CodeWindow, budget int) string {
	var parts []string

	// Header with error info
	header := formatErrorHeader(err)
	parts = append(parts, header)

	// Code window
	code := window.Render()
	parts = append(parts, "```\n"+code+"\n```")

	// Instruction
	instruction := "\nFix the error above. Provide only the fixed code with minimal explanation."
	parts = append(parts, instruction)

	result := strings.Join(parts, "\n")

	// Simple token estimation
	estimatedTokens := len(result) / GoTokenRatio
	if estimatedTokens > budget {
		// Truncate if over budget
		maxLen := budget * GoTokenRatio
		if len(result) > maxLen {
			result = result[:maxLen]
		}
	}

	return result
}

func formatErrorHeader(err *ParsedError) string {
	return strings.Join([]string{
		"File: " + err.GetFile(),
		"Line: " + itoa(err.GetLine()),
		"Col: " + itoa(err.GetCol()),
		"Source: " + err.GetSource(),
		"Error: " + err.GetMessage(),
	}, "\n")
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}

// GetTokenCount estimates the token count for a string
func GetTokenCount(s string) int {
	return len(s) / GoTokenRatio
}