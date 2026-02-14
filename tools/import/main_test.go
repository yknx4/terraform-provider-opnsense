package main

import (
	"testing"
)

func TestSanitizeResourceName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Allow HTTPS Traffic", "allow_https_traffic"},
		{"Test-Rule-123", "test_rule_123"},
		{"rule_with_underscores", "rule_with_underscores"},
		{"Rule With Multiple   Spaces", "rule_with_multiple_spaces"},
		{"_leading_underscore", "leading_underscore"},
		{"trailing_underscore_", "trailing_underscore"},
		{"UPPERCASE", "uppercase"},
		{"special!@#chars", "special_chars"},
		{"", ""},
		{"123numeric", "123numeric"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := sanitizeResourceName(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeResourceName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
