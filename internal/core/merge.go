package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/LeeFred3042U/kitkat/internal/storage"
)

// Merge attempts to merge the given branch into the current branch
// Currently only supports strict fast-forward merges
func Merge(branchToMerge string) error {

	// Guard: ensure we're inside a kitkat repo
	if _, err := os.Stat(RepoDir); os.IsNotExist(err) {
		return errors.New("not a kitkat repository (run `kitkat init`)")
	}

	//Safety Check: Verify working directory is clean
	dirty, err := IsWorkDirDirty()
	if err != nil {
		return fmt.Errorf("failed to check working directory status: %w", err)
	}
	if dirty {
		return fmt.Errorf("error: your local changes would be overwritten by merge. Please commit or stash them")
	}

	// Getting the commit hash of the branch to merge
	branchPath := filepath.Join(HeadsDir, branchToMerge)
	featureHeadHashBytes, err := os.ReadFile(branchPath)
	if err != nil {
		return fmt.Errorf("branch '%s' not found", branchToMerge)
	}
	featureHeadHash := strings.TrimSpace(string(featureHeadHashBytes))

	// Getting the commit hash of the current branch (HEAD)
	currentHeadHash, err := readHead()
	if err != nil {
		return fmt.Errorf("could not read current HEAD: %w", err)
	}

	//  Ancestry Check: Calculate merge base
	mergeBase, err := storage.FindMergeBase(currentHeadHash, featureHeadHash)
	if err != nil {
		return fmt.Errorf("failed to calculate merge base: %w", err)
	}

	// Merge Type Determination
	switch mergeBase {
	case currentHeadHash:
		// fast-forward
		fmt.Printf("Updating %s..%s\n", currentHeadHash[:7], featureHeadHash[:7])
		fmt.Println("Fast-forward")

	case featureHeadHash:
		// already up to date
		fmt.Println("Already up to date.")
		return nil

	default:
		// diverged
		return fmt.Errorf(
			"fatal: Not possible to fast-forward, aborting.\n"+
				"Merge commits are not supported. Please rebase '%s' onto the current branch",
			branchToMerge,
		)
	}

	// Fast-Forward Execution
	if err := UpdateBranchPointer(featureHeadHash); err != nil {
		return fmt.Errorf("failed to update branch pointer: %w", err)
	}

	// Update the working directory and index to match the new HEAD state
	err = UpdateWorkspaceAndIndex(featureHeadHash)
	if err != nil {
		// Attempt to roll back the branch pointer on failure
		fmt.Printf("UpdateWorkspaceAndIndex failed: %v. Rolling back branch pointer...\n", err)
		if rollbackErr := UpdateBranchPointer(currentHeadHash); rollbackErr != nil {
			return fmt.Errorf("failed to update workspace: %w; additionally failed to rollback branch pointer: %v", err, rollbackErr)
		}
		return fmt.Errorf("failed to update workspace: %w; branch pointer rolled back to %s", err, currentHeadHash)
	}

	return nil
}
