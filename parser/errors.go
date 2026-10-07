package parser

import (
	"regexp"
	"strings"
)

// ErrorMatch represents a parsed compiler/linter error
type ErrorMatch struct {
	File   string
	Line   int
	Col    int
	Message string
	Source string // Go, Python, Rust, GCC
}

// ParseError parses a compiler/linter error line and extracts coordinates
func ParseError(line string) *ErrorMatch {
	trimmed := strings.TrimSpace(line)
	
	// Go format: file.go:line:col: message or file.go:line:col message
	goRe := regexp.MustCompile(`^(.+?\.(go|mod|sum)):(\d+):(\d+):?\s*(.*)$`)
	if matches := goRe.FindStringSubmatch(trimmed); matches != nil {
		return &ErrorMatch{
			File:    matches[1],
			Line:    parseInt(matches[3]),
			Col:     parseInt(matches[4]),
			Message: matches[5],
			Source:  "Go",
		}
	}
	
	// Python format: file.py:line:col: message or file.py:line: message
	pythonRe := regexp.MustCompile(`^(.+?\.py):(\d+)(?::(\d+))?:?\s*(.*)$`)
	if matches := pythonRe.FindStringSubmatch(trimmed); matches != nil {
		return &ErrorMatch{
			File:    matches[1],
			Line:    parseInt(matches[2]),
			Col:     parseIntOrDefault(matches[3], 1),
			Message: matches[4],
			Source:  "Python",
		}
	}
	
	// Rust format: file.rs:line:col: message
	rustRe := regexp.MustCompile(`^(.+?\.(rs|toml)):(\d+):(\d+):\s*(.*)$`)
	if matches := rustRe.FindStringSubmatch(trimmed); matches != nil {
		return &ErrorMatch{
			File:    matches[1],
			Line:    parseInt(matches[3]),
			Col:     parseInt(matches[4]),
			Message: matches[5],
			Source:  "Rust",
		}
	}
	
	// GCC/Clang format: file.c:line:col: message or file.c:line:col-col: message
	gccRe := regexp.MustCompile(`^(.+?\.(c|cpp|cc|cxx|h|hpp|c\+\+)):(\d+):(\d+)(?:-(\d+))?:?\s*(.*)$`)
	if matches := gccRe.FindStringSubmatch(trimmed); matches != nil {
		return &ErrorMatch{
			File:    matches[1],
			Line:    parseInt(matches[3]),
			Col:     parseInt(matches[4]),
			Message: matches[6],
			Source:  "GCC",
		}
	}
	
	return nil
}

func parseInt(s string) int {
	if s == "" {
		return 0
	}
	n := 0
	for _, c := range s {
		n = n*10 + int(c-'0')
	}
	return n
}

func parseIntOrDefault(s string, def int) int {
	if s == "" {
		return def
	}
	return parseInt(s)
}
