package core

import (
	"fmt"

	"github.com/LeeFred3042U/kitkat/internal/storage"
)

// Prints the commit log.
func ShowLog() error {
	commits, err := storage.ReadCommits()
	if err != nil {
		return err
	}

	// Print in reverse chronological order
	for i := len(commits) - 1; i >= 0; i-- {
		commit := commits[i]
		fmt.Printf("commit %s\n", commit.ID)
		fmt.Printf("Parent: %s\n", commit.Parent)
		fmt.Printf("Date:   %s\n", commit.Timestamp.Format("Mon Jan 02 15:04:05 2006 -0700"))
		fmt.Printf("\n    %s\n\n", commit.Message)
	}

	return nil
}
