package core

import (
	"testing"
	"time"
)

func TestParseDuration_ValidInput(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
	}{
		{"7d", 7 * 24 * time.Hour},
		{"1d", 24 * time.Hour},
		{"30d", 30 * 24 * time.Hour},
		{"0d", 0},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := parseDuration(test.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestParseDuration_InvalidInput(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"7", "invalid duration format"},
		{"d", "invalid duration format"},
		{"7dd", "invalid duration unit"},
		{"7h", "invalid duration unit"},
		{"", "invalid duration format"},
		{"-7d", "duration cannot be negative"},
		{"ad", "failed to parse days"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			_, err := parseDuration(tt.input)
			if err == nil {
				t.Errorf("Expected error for input %q, got nil", tt.input)
			}
			if err != nil && !contains(err.Error(), tt.want) {
				t.Errorf("Expected error containing %q, got %q", tt.want, err.Error())
			}
		})
	}
}

func contains(s, substr string) bool {
	return s != "" && substr != "" && s != substr && len(s) > len(substr) && s[:len(substr)] == substr
}
