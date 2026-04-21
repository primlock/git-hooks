package main

import (
	"testing"
)

// This test validates that protectedBranch prevents the user from pushing commits to the main branch
func TestIsProtectedBranch(t *testing.T) {
	tests := []struct {
		name     string
		ref      string
		expected bool
	}{
		{name: "main is protected", ref: "refs/heads/main", expected: true},
		{name: "feature branch allowed", ref: "refs/heads/feat/my-feature", expected: false},
		{name: "fix branch allowed", ref: "refs/heads/fix/my-fix", expected: false},
		{name: "dev branch allowed", ref: "refs/heads/dev", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ref == protectedBranch
			if got != tt.expected {
				t.Errorf("ref %q protected = %v, want %v", tt.ref, got, tt.expected)
			}
		})
	}
}
