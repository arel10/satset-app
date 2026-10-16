package service_test

import (
	"testing"

	"satset-app/internal/service"
)

func TestNormalizePhone(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"081234567890", "+6281234567890"},
		{"6281234567890", "+6281234567890"},
		{"+6281234567890", "+6281234567890"},
		{"0812-3456-7890", "+6281234567890"},
		{"0812 3456 7890", "+6281234567890"},
		{"(0812) 34567890", "+6281234567890"},
		{"81234567890", "+6281234567890"},
		// Invalid cases
		{"123", ""},               // too short
		{"abcdef", ""},            // not digits
		{"", ""},                  // empty
		{"08123456789012345", ""}, // too long
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := service.NormalizePhone(tt.input)
			if got != tt.want {
				t.Errorf("NormalizePhone(%q) = %q; want %q", tt.input, got, tt.want)
			}
		})
	}
}
