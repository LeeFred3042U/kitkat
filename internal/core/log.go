package core

import (
	"fmt"

	"github.com/LeeFred3042U/kitkat/internal/storage"
)

// ShowLog prints the commit log. It accepts a boolean to switch to a compact, one-line format
// and an optional limit to restrict the number of commits shown (use -1 for no limit).
func ShowLog(oneline bool, limit int) error {
	commits, err := storage.ReadCommits()
	if err != nil {
		return err
	}

	count := 0
	// Print in reverse chronological order (newest first).
	for i := len(commits) - 1; i >= 0; i-- {
		if limit > 0 && count >= limit {
			break
		}
		commit := commits[i]
		if oneline {
			fmt.Printf("%s %s\n", commit.ID[:7], commit.Message)
		} else {
			fmt.Printf("commit %s\n", commit.ID)
			fmt.Printf("Author: %s <%s>\n", commit.AuthorName, commit.AuthorEmail)
			fmt.Printf("Date:   %s\n", commit.Timestamp.Local().Format("Mon Jan 02 15:04:05 2006 -0700"))
			fmt.Printf("\n    %s\n\n", commit.Message)
		}
		count++
	}
	return nil
}
