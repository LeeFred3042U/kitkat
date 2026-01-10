package core

import (
	"fmt"

	"github.com/LeeFred3042U/kitcat/internal/storage"
)

// ResetHard moves the current branch (or HEAD in detached state) to the specified commit
// and forcibly updates the working directory and index to match that commit.
// WARNING: This is a destructive operation that discards all uncommitted changes.
func ResetHard(commitHash string) error {
	// Step 1: Validate that the commit exists
	commit, err := storage.FindCommit(commitHash)
	if err != nil {
		if err == storage.ErrNoCommits {
			return fmt.Errorf("fatal: invalid commit: %s", commitHash)
		}
		return fmt.Errorf("fatal: invalid commit: %s", commitHash)
	}

	// Step 2: Save current HEAD for potential rollback
	oldHeadCommit, err := readHead()
	if err != nil {
		return fmt.Errorf("fatal: unable to read HEAD file: %w", err)
	}

	// Step 3: Update the branch pointer or HEAD
	if err := UpdateBranchPointer(commitHash); err != nil {
		return err
	}

	// Step 4: Update workspace and index to match the target commit
	// If this fails, attempt to roll back the branch pointer
	if err := UpdateWorkspaceAndIndex(commitHash); err != nil {
		// Attempt rollback
		_ = UpdateBranchPointer(oldHeadCommit)
		return fmt.Errorf("failed to update workspace: %w", err)
	}

	// Step 5: Success - print confirmation message
	fmt.Printf("HEAD is now at %s %s\n", commitHash[:7], commit.Message)
	return nil
}
