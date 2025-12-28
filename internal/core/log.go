package core

import (
	"fmt"

	"github.com/LeeFred3042U/kitkat/internal/storage"
)

// ShowLog prints the commit log. It accepts a boolean to switch to a compact, one-line format.
func ShowLog(oneline bool) error {
	// Determine where to start traversing
	headHash, err := readHEAD()
	if err != nil {
		return err
	}

	currentHash := headHash
	for currentHash != "" {
		commit, err := storage.FindCommit(currentHash)
		if err != nil {
			return err
		}

		if oneline {
			fmt.Printf("%s %s\n", commit.ID[:7], commit.Message)
		} else {
			fmt.Printf("commit %s\n", commit.ID)
			fmt.Printf("Author: %s <%s>\n", commit.AuthorName, commit.AuthorEmail)
			fmt.Printf("Date:   %s\n", commit.Timestamp.Local().Format("Mon Jan 02 15:04:05 2006 -0700"))
			fmt.Printf("\n    %s\n\n", commit.Message)
		}

		// Move to parent
		currentHash = commit.Parent
	}
	return nil
}
