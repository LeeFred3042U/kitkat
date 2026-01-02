package core

import (
	"fmt"

	"github.com/LeeFred3042U/kitkat/internal/storage"
)

// ShowLog prints the commit log. It accepts a boolean for oneline format
// and an optional limit to restrict the number of commits shown (use -1 or 0 for no limit)
func ShowLog(oneline bool, limit int) error {
	// 1Start from HEAD (Architecture from reset-hard branch)
	// We must walk backwards from HEAD, otherwise 'reset' changes won't be reflected
	currentCommit, err := GetHeadCommit()
	if err != nil {
		// Handle the case where the repo is empty or HEAD is invalid
		return nil
	}

	commitHash := currentCommit.ID
	count := 0

	// Walk the graph (Architecture from reset-hard branch)
	for commitHash != "" {
		// Apply the Limit Check (Feature from main branch)
		if limit > 0 && count >= limit {
			break
		}

		commit, err := storage.FindCommit(commitHash)
		if err != nil {
			return err
		}

		// Print Logic
		if oneline {
			fmt.Printf("%s %s\n", commit.ID[:7], commit.Message)
		} else {
			fmt.Printf("commit %s\n", commit.ID)
			fmt.Printf("Author: %s <%s>\n", commit.AuthorName, commit.AuthorEmail)
			fmt.Printf("Date:   %s\n", commit.Timestamp.Local().Format("Mon Jan 02 15:04:05 2006 -0700"))
			fmt.Printf("\n    %s\n\n", commit.Message)
		}

		// Move to parent pointer
		commitHash = commit.Parent
		count++
	}

	return nil
}
