package api

import (
	"testing"
)

func TestCleanProfanity(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no profanity",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "single profane word",
			input:    "this is a kerfuffle",
			expected: "this is a ****",
		},
		{
			name:     "multiple profane words",
			input:    "sharbert and fornax",
			expected: "**** and ****",
		},
		{
			name:     "case insensitive",
			input:    "KERFUFFLE",
			expected: "****",
		},
		{
			name:     "profane word with punctuation",
			input:    "kerfuffle!",
			expected: "****!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := cleanProfanity(tt.input)
			if actual != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}
