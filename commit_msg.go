// This file is called by `git commit` with one argument which the name of the file holding the commit message.
// Regex pattern matching is used on the commit message to enforce conventional commit formatting.
// If the pattern matches the message, the program will exit successfully. Otherwise, it will exit with 1 and
// the commit will be aborted.
package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Match all commit messages that follow common convention types
const pattern = `^(build|chore|ci|docs|feat|fix|perf|refactor|revert|style|test)(\(.+\))?!?: .+$`

// Let the user know the commit didn't pass the hook and what the accepted inputs are
const helpString = `
Failed to commit. Commit message does not follow the conventional commit format.

Commit Types:
  build:    Changes that affect the build system or external dependencies.
  chore:    Routine maintenance tasks that do not modify source or test files.
  ci:       Changes to CI configuration files and scripts.
  docs:     Documentation only changes.
  feat:     A new feature or enhancement to existing functionality.
  fix:      A bug fix or correction to existing functionality.
  perf:     Code changes that improve application performance.
  refactor: Restructuring code without changing its behavior or fixing bugs.
  revert:   Reverts a previous commit to undo changes.
  style:    Formatting changes that do not code functionality.
  test:     Adding missing tests or improving existing tests.

Optional Attributes:
  (): Scope of the change, including active ticket IDs.
  !:  Indicate a breaking change.

Examples:
  fix: Prevent infinite looping condition
  feat(api): Added profile metrics api
  chore: Moves metadata into new dir
  build!: Build system changed from meson to cmake
`

// Extract the first non-empty, non-comment line from the commit message
func parseCommitMessage(raw string) string {
	for line := range strings.SplitSeq(raw, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			return trimmed
		}
	}
	return ""
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: No file name provided")
		os.Exit(1)
	}

	// The first argument is the name of the file holding the commit message
	filename := os.Args[1]

	// Open up the file and read the contents into memory so we can validate it follows our pattern
	msg, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Failed to read %v, %v", filename, err)
		os.Exit(1)
	}

	// Trim trailing whitespace/newlines added by git
	trimmed := parseCommitMessage(string(msg))

	// Match the commit message against the regex pattern
	isMatch, err := regexp.MatchString(pattern, string(trimmed))
	if err != nil {
		fmt.Printf("Failed to match commit message, %v\n", err)
		os.Exit(1)
	} else if !isMatch {
		fmt.Printf("%v\n", helpString)
		os.Exit(1)
	}
}
