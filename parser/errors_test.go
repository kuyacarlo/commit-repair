package parser

import (
	"testing"
)

func TestParseError_Go(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *ErrorMatch
	}{
		{
			name:  "Go standard error format",
			input: "main.go:10:5: undefined: x",
			expected: &ErrorMatch{
				File:    "main.go",
				Line:    10,
				Col:     5,
				Message: "undefined: x",
				Source:  "Go",
			},
		},
		{
			name:  "Go error with package path",
			input: "github.com/user/repo/main.go:15:2: undefined: y",
			expected: &ErrorMatch{
				File:    "github.com/user/repo/main.go",
				Line:    15,
				Col:     2,
				Message: "undefined: y",
				Source:  "Go",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseError(tt.input)
			if result == nil {
				t.Fatalf("ParseError(%q) returned nil", tt.input)
			}
			if result.File != tt.expected.File {
				t.Errorf("File = %q, want %q", result.File, tt.expected.File)
			}
			if result.Line != tt.expected.Line {
				t.Errorf("Line = %d, want %d", result.Line, tt.expected.Line)
			}
			if result.Col != tt.expected.Col {
				t.Errorf("Col = %d, want %d", result.Col, tt.expected.Col)
			}
			if result.Message != tt.expected.Message {
				t.Errorf("Message = %q, want %q", result.Message, tt.expected.Message)
			}
			if result.Source != tt.expected.Source {
				t.Errorf("Source = %q, want %q", result.Source, tt.expected.Source)
			}
		})
	}
}

func TestParseError_Python(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *ErrorMatch
	}{
		{
			name:  "Python error format with column",
			input: "main.py:10:5: NameError: name 'x' is not defined",
			expected: &ErrorMatch{
				File:    "main.py",
				Line:    10,
				Col:     5,
				Message: "NameError: name 'x' is not defined",
				Source:  "Python",
			},
		},
		{
			name:  "Python error format without column",
			input: "main.py:15: SyntaxError: invalid syntax",
			expected: &ErrorMatch{
				File:    "main.py",
				Line:    15,
				Col:     1,
				Message: "SyntaxError: invalid syntax",
				Source:  "Python",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseError(tt.input)
			if result == nil {
				t.Fatalf("ParseError(%q) returned nil", tt.input)
			}
			if result.File != tt.expected.File {
				t.Errorf("File = %q, want %q", result.File, tt.expected.File)
			}
			if result.Line != tt.expected.Line {
				t.Errorf("Line = %d, want %d", result.Line, tt.expected.Line)
			}
			if result.Col != tt.expected.Col {
				t.Errorf("Col = %d, want %d", result.Col, tt.expected.Col)
			}
			if result.Source != tt.expected.Source {
				t.Errorf("Source = %q, want %q", result.Source, tt.expected.Source)
			}
		})
	}
}

func TestParseError_Rust(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *ErrorMatch
	}{
		{
			name:  "Rust compilation error",
			input: "src/main.rs:10:5: error[E0425]: cannot find value `x` in this scope",
			expected: &ErrorMatch{
				File:    "src/main.rs",
				Line:    10,
				Col:     5,
				Message: "error[E0425]: cannot find value `x` in this scope",
				Source:  "Rust",
			},
		},
		{
			name:  "Rust Cargo.toml format",
			input: "Cargo.toml:5:1: error: failed to parse manifest",
			expected: &ErrorMatch{
				File:    "Cargo.toml",
				Line:    5,
				Col:     1,
				Message: "error: failed to parse manifest",
				Source:  "Rust",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseError(tt.input)
			if result == nil {
				t.Fatalf("ParseError(%q) returned nil", tt.input)
			}
			if result.File != tt.expected.File {
				t.Errorf("File = %q, want %q", result.File, tt.expected.File)
			}
			if result.Line != tt.expected.Line {
				t.Errorf("Line = %d, want %d", result.Line, tt.expected.Line)
			}
			if result.Col != tt.expected.Col {
				t.Errorf("Col = %d, want %d", result.Col, tt.expected.Col)
			}
			if result.Source != tt.expected.Source {
				t.Errorf("Source = %q, want %q", result.Source, tt.expected.Source)
			}
		})
	}
}

func TestParseError_GCC(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *ErrorMatch
	}{
		{
			name:  "GCC C++ error",
			input: "main.cpp:10:5: error: 'x' was not declared in this scope",
			expected: &ErrorMatch{
				File:    "main.cpp",
				Line:    10,
				Col:     5,
				Message: "error: 'x' was not declared in this scope",
				Source:  "GCC",
			},
		},
		{
			name:  "GCC C header error",
			input: "src/main.c:15:10: warning: implicit declaration of function",
			expected: &ErrorMatch{
				File:    "src/main.c",
				Line:    15,
				Col:     10,
				Message: "warning: implicit declaration of function",
				Source:  "GCC",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseError(tt.input)
			if result == nil {
				t.Fatalf("ParseError(%q) returned nil", tt.input)
			}
			if result.File != tt.expected.File {
				t.Errorf("File = %q, want %q", result.File, tt.expected.File)
			}
			if result.Line != tt.expected.Line {
				t.Errorf("Line = %d, want %d", result.Line, tt.expected.Line)
			}
			if result.Col != tt.expected.Col {
				t.Errorf("Col = %d, want %d", result.Col, tt.expected.Col)
			}
			if result.Source != tt.expected.Source {
				t.Errorf("Source = %q, want %q", result.Source, tt.expected.Source)
			}
		})
	}
}

func TestParseError_Invalid(t *testing.T) {
	inputs := []string{
		"",
		"not a valid error format",
		"some random text without file",
	}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			result := ParseError(input)
			if result != nil {
				t.Errorf("ParseError(%q) = %+v, want nil", input, result)
			}
		})
	}
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"", 0},
		{"0", 0},
		{"1", 1},
		{"10", 10},
		{"123", 123},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseInt(tt.input)
			if result != tt.expected {
				t.Errorf("parseInt(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseIntOrDefault(t *testing.T) {
	tests := []struct {
		input    string
		def      int
		expected int
	}{
		{"", 5, 5},
		{"0", 5, 0},
		{"10", 5, 10},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseIntOrDefault(tt.input, tt.def)
			if result != tt.expected {
				t.Errorf("parseIntOrDefault(%q, %d) = %d, want %d", tt.input, tt.def, result, tt.expected)
			}
		})
	}
}