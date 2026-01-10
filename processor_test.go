package main

import "testing"

func TestStandardizeInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "simple text",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "text with leading/trailing whitespace",
			input:    "  hello world  ",
			expected: "hello world",
		},
		{
			name:     "text with newlines",
			input:    "hello\nworld",
			expected: "hello world",
		},
		{
			name:     "text with multiple newlines",
			input:    "hello\n\n\nworld",
			expected: "hello world",
		},
		{
			name:     "text with carriage returns and newlines",
			input:    "hello\r\nworld",
			expected: "hello world",
		},
		{
			name:     "text with multiple spaces",
			input:    "hello    world",
			expected: "hello world",
		},
		{
			name:     "text with tabs",
			input:    "hello\tworld",
			expected: "hello world",
		},
		{
			name:     "text with mixed whitespace",
			input:    "  hello\n\t  world  \n",
			expected: "hello world",
		},
		{
			name:     "multiline text",
			input:    "line one\nline two\nline three",
			expected: "line one line two line three",
		},
		{
			name:     "text with only whitespace",
			input:    "   \n\t  \n  ",
			expected: "",
		},
		{
			name:     "text with Unicode line separator (U+2028)",
			input:    "hello\u2028world",
			expected: "hello world",
		},
		{
			name:     "text with Unicode paragraph separator (U+2029)",
			input:    "hello\u2029world",
			expected: "hello world",
		},
		{
			name:     "text with both Unicode separators",
			input:    "hello\u2028\u2029world",
			expected: "hello world",
		},
		{
			name:     "text with mixed newlines and Unicode separators",
			input:    "hello\n\u2028world\u2029test",
			expected: "hello world test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StandardizeInput(tt.input)
			if result != tt.expected {
				t.Errorf("StandardizeInput(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestProcessInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "lowercase text",
			input:    "hello world",
			expected: "HELLO WORLD",
		},
		{
			name:     "uppercase text",
			input:    "HELLO WORLD",
			expected: "HELLO WORLD",
		},
		{
			name:     "mixed case text",
			input:    "Hello World",
			expected: "HELLO WORLD",
		},
		{
			name:     "text with numbers",
			input:    "hello 123 world",
			expected: "HELLO 123 WORLD",
		},
		{
			name:     "text with special characters",
			input:    "hello, world!",
			expected: "HELLO, WORLD!",
		},
		{
			name:     "text with unicode",
			input:    "héllo wörld",
			expected: "HÉLLO WÖRLD",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessInput(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessInput(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
