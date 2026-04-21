// This file is called by `git push` after the remote status has been checked with two parameters: 
//   1. The name of the remote to which the push is being done 
//   2. The URL to which the push is being done
//
// The hook receives ref information via stdin in the format:
//   <local-ref> <local-sha> <remote-ref> <remote-sha>
// All pushes to the main branch will be rejected because we want merge requests to go through a formal pull 
// request before being added into main.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	protectedBranch    = "refs/heads/main"
	expectedRefFields  = 4
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Parse the stdin format: <local-ref> <local-sha> <remote-ref> <remote-sha>
		parts := strings.Fields(line)
		if len(parts) < expectedRefFields {
			continue
		}

		remoteRef := parts[2]

		// Check if pushing to protected branch
		if remoteRef == protectedBranch {
			fmt.Printf("Failed to push to '%s'. Please create a new branch and open up a pull request.\n", protectedBranch)
			os.Exit(1)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading stdin: %v\n", err)
		os.Exit(1)
	}
}
