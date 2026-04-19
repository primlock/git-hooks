package main

import (
	"strings"
	"testing"
)

// Validates that scanContent correctly identifies secrets and does not flag clean or placeholder content as 
// findings. Each test case targets a specific pattern by name so regressions are easy to locate.
func TestScanContent(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		expectMatch   bool
		expectPattern string
	}{
		// Generic Token
		{
			name:          "generic token assignment with double quotes",
			content:       `token = "ghp_abcdefghijklmnopqrstuvwxyz"`,
			expectMatch:   true,
			expectPattern: "Generic Token",
		},
		{
			name:          "generic token assignment with single quotes",
			content:       `token = 'ghp_abcdefghijklmnopqrstuvwxyz'`,
			expectMatch:   true,
			expectPattern: "Generic Token",
		},

		// AWS Secret Key
		{
			name:          "AWS secret access key",
			content:       `aws_secret_access_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"`,
			expectMatch:   true,
			expectPattern: "AWS Secret Key",
		},
		{
			name:          "AWS secret access key without quotes",
			content:       `aws_secret_access_key = wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY`,
			expectMatch:   true,
			expectPattern: "AWS Secret Key",
		},

		// GCP API Key
		{
			name:          "GCP API key",
			content:       `GCP_KEY="AIzaSyD-9tSrke72I6gKXr9AsI9RUheZDyKZNd4"`,
			expectMatch:   true,
			expectPattern: "GCP API Key",
		},

		// Anthropic API Key
		{
			name:          "Anthropic API key",
			content:       `ANTHROPIC_API_KEY="sk-ant-abcdefghijklmnopqrstuvwxyz012345"`,
			expectMatch:   true,
			expectPattern: "Anthropic API Key",
		},

		// OpenAI API Key
		{
			name:          "OpenAI API key",
			content:       `OPENAI_API_KEY="sk-abcdefghijklmnopqrstuvwxyz012345"`,
			expectMatch:   true,
			expectPattern: "OpenAI API Key",
		},

		// Should not match
		{
			name:        "clean Go source file",
			content:     `func main() { fmt.Println("hello world") }`,
			expectMatch: false,
		},
		{
			name:        "placeholder token value",
			content:     `token = "your-token-here"`,
			expectMatch: true,
		},
		{
			name:        "commented out key",
			content:     `// token = "sk-abcdefghijklmnopqrstuvwxyz012345"`,
			expectMatch: true,
			expectPattern: "OpenAI API Key",
		},
		{
			name:        "empty file",
			content:     ``,
			expectMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := scanContent("test_file.go", []byte(tt.content))
			hasMatch := len(findings) > 0

			if hasMatch != tt.expectMatch {
				t.Errorf("scanContent() match = %v, want %v\nfindings: %v", hasMatch, tt.expectMatch, findings)
			}

			// When we expect a specific pattern, confirm it appears in the findings
			if tt.expectMatch && tt.expectPattern != "" {
				found := false
				for _, f := range findings {
					if strings.Contains(f, tt.expectPattern) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected pattern %q in findings, got: %v", tt.expectPattern, findings)
				}
			}
		})
	}
}
