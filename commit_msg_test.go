package main

import (
	"regexp"
	"testing"
)

// parseCommitMessage is tested indirectly via the full input tests, but also directly here to cover amend/comment stripping behaviour.
func TestParseCommitMessage(t *testing.T) {
	tests := []struct {
		name     string
		raw      string
		expected string
	}{
		{
			name:     "plain message",
			raw:      "feat: add login\n",
			expected: "feat: add login",
		},
		{
			name:     "strips leading comment lines",
			raw:      "# This is a comment\nfix: correct typo\n",
			expected: "fix: correct typo",
		},
		{
			name:     "amend with comment block",
			raw:      "feat(api): new endpoint\n# Please enter the commit message for your changes.\n# Lines starting with '#' will be ignored.\n",
			expected: "feat(api): new endpoint",
		},
		{
			name:     "empty message",
			raw:      "",
			expected: "",
		},
		{
			name:     "only comments",
			raw:      "# comment one\n# comment two\n",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCommitMessage(tt.raw)
			if got != tt.expected {
				t.Errorf("parseCommitMessage(%q) = %q, want %q", tt.raw, got, tt.expected)
			}
		})
	}
}

func TestPatternMatching(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		matches bool
	}{
		// Valid commit message without scope
		{name: "build type", msg: "build: update dependencies", matches: true},
		{name: "chore type", msg: "chore: clean up temp files", matches: true},
		{name: "ci type", msg: "ci: add github actions workflow", matches: true},
		{name: "docs type", msg: "docs: update README", matches: true},
		{name: "feat type", msg: "feat: add dark mode", matches: true},
		{name: "fix type", msg: "fix: prevent infinite loop", matches: true},
		{name: "perf type", msg: "perf: cache db results", matches: true},
		{name: "refactor type", msg: "refactor: extract helper function", matches: true},
		{name: "revert type", msg: "revert: undo broken deploy", matches: true},
		{name: "style type", msg: "style: fix indentation", matches: true},
		{name: "test type", msg: "test: add unit tests for auth", matches: true},

		// Valid commit message with scope
		{name: "feat with scope", msg: "feat(api): new profile endpoint", matches: true},
		{name: "fix with scope", msg: "fix(auth): handle expired token", matches: true},
		{name: "chore with ticket ID", msg: "chore(PROJ-123): update lockfile", matches: true},
		{name: "feat with issue number", msg: "feat(#42): implement search", matches: true},
		{name: "fix with jira ticket", msg: "fix(ISS-7): null pointer on startup", matches: true},

		// Valid commit message with breaking change
		{name: "breaking change no scope", msg: "feat!: remove deprecated api", matches: true},
		{name: "breaking change with scope", msg: "feat(api)!: remove v1 endpoints", matches: true},
		{name: "fix breaking change", msg: "fix!: change default config format", matches: true},
		{name: "build breaking change", msg: "build!: switch build system to cmake", matches: true},

		// Invalid commit message with unknown type
		{name: "unknown type", msg: "update: something", matches: false},
		{name: "unknown type 'wip'", msg: "wip: work in progress", matches: false},
		{name: "empty type", msg: ": missing type", matches: false},

		// Invalid commit message with incorrect formatting
		{name: "missing colon", msg: "feat add new feature", matches: false},
		{name: "missing space after colon", msg: "feat:no space", matches: false},
		{name: "missing description", msg: "feat: ", matches: false},
		{name: "uppercase type", msg: "Fix: something", matches: false},
		{name: "extra space before colon", msg: "feat : something", matches: false},
		{name: "double colon", msg: "feat:: something", matches: false},

		// Invalid empty or whitespace commit message
		{name: "empty string", msg: "", matches: false},
		{name: "whitespace only", msg: "   ", matches: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := regexp.MatchString(pattern, tt.msg)
			if err != nil {
				t.Fatalf("regexp.MatchString error: %v", err)
			}
			if got != tt.matches {
				t.Errorf("pattern match for %q = %v, want %v", tt.msg, got, tt.matches)
			}
		})
	}
}
