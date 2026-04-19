// This file is called by `git commit` with no arguments. It gives the programmer the ability to interact with
// a staged file before it's committed. In in this example, the staged files are parsed for API key exposure
// in the file contents by using regex to match on common patterns like AWS keys, GCP keys, Anthropic keys and
// the like. If a pattern matches on the content of a file, the commit is aborted by returning a non-zero
// status.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// Maps a human-readable label to a compiled regex pattern. Each pattern targets a specific type of secret or 
// credential that should never be committed to source control.
var secretPatterns = map[string]*regexp.Regexp{
	// Matches common token assignment patterns, e.g. token = "abc123"
	"Generic Token": regexp.MustCompile(`(?i)token\s*=\s*["'].+["']`),

	// Matches AWS secret access keys assigned in config or source files
	"AWS Secret Key": regexp.MustCompile(`(?i)aws_secret_access_key\s*=\s*["']?[A-Za-z0-9/+=]{40}["']?`),

	// Matches GCP and Google Gemini API keys, which share the AIza prefix
	"GCP API Key": regexp.MustCompile(`AIza[0-9A-Za-z\-_]{35}`),

	// Matches Anthropic API keys, which are prefixed with sk-ant-
	"Anthropic API Key": regexp.MustCompile(`sk-ant-[a-zA-Z0-9\-_]{32,}`),

	// Matches OpenAI API keys, which are prefixed with sk-
	"OpenAI API Key": regexp.MustCompile(`sk-[a-zA-Z0-9]{20,}`),
}

// Return a list of files staged for commit by invoking `git diff` with arguments or an error indicating why
// the operation failed.
// Example: git diff --cached --name-only
//          pre-commit/pre_commit.go
func getStagedFiles() ([]string, error) {
	out, err := exec.Command("git", "diff", "--cached", "--name-only").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get staged files: %w", err)
	}

	var files []string
	for line := range strings.SplitSeq(string(out), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			files = append(files, trimmed)
		}
	}
	return files, nil
}

// Check the parameter `content` against all the secret patterns. If a secret is found, return a string with
// the type of secret identified and file it was discovered in.
func scanContent(path string, content []byte) []string {
	var findings []string
	for name, pattern := range secretPatterns {
		if pattern.Match(content) {
			findings = append(findings, fmt.Sprintf("  %s found in %s", name, path))
		}
	}
	return findings
}

// Read the contents of each `path` from disk and invoke scanContent to handle pattern matching.
// If a file cannot be read, notify the user and continue.
func scanFile(path string) []string {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Warning: could not read staged file %s: %v\n", path, err)
		return nil
	}
	return scanContent(path, content)
}

func main() {
	// Collect all the files staged for commit
	files, err := getStagedFiles()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Scan the collection of staged files for secrets and add any file names matching defined patterns to the
	// list of findings
	var findings []string
	for _, file := range files {
		findings = append(findings, scanFile(file)...)
	}

	// Abort the commit if any secrets are found
	if len(findings) > 0 {
		fmt.Println("\nFailed to commit. Secrets discovered in staged files:")
		for _, f := range findings {
			fmt.Println(f)
		}
		fmt.Println("\nCommit aborted. Please remove secrets before committing.")
		os.Exit(1)
	}
}
