package core

import (
	"fmt"

	"github.com/LeeFred3042U/kitkat/internal/storage"
)

// ShowLog prints all commits reachable from HEAD by walking the parent chain.
// This ensures that after a reset, only commits in HEAD's history are shown.
func ShowLog(oneline bool) error {
	// Get the commit that HEAD currently points to
	currentCommit, err := GetHeadCommit()
	if err != nil {
		if err == storage.ErrNoCommits {
			fmt.Println("No commits yet.")
			return nil
		}
		return err
	}

	// Walk back through the commit history from HEAD
	commitHash := currentCommit.ID
	for commitHash != "" {
		commit, err := storage.FindCommit(commitHash)
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

		// Move to parent commit
		commitHash = commit.Parent
	}

	return nil
}
