package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/LeeFred3042U/kitkat/internal/storage"
)

// Merge attempts to merge the given branch into the current branch
// Currently only supports fast-forward merges
func Merge(branchToMerge string) error {
	// Getting the commit hash of the branch to merge
	branchPath := filepath.Join(headsDir, branchToMerge)
	featureHeadHashBytes, err := os.ReadFile(branchPath)
	if err != nil {
		return fmt.Errorf("branch '%s' not found", branchToMerge)
	}
	featureHeadHash := strings.TrimSpace(string(featureHeadHashBytes))

	// Getting the commit hash of the current branch (HEAD)
	currentHeadHash, err := readHEAD()
	if err != nil {
		return fmt.Errorf("could not read current HEAD: %w", err)
	}

	// Checking if the current HEAD is an ancestor of the branch to merge
	isAncestor, err := storage.IsAncestor(currentHeadHash, featureHeadHash)
	if err != nil {
		return err
	}

	if !isAncestor {
		// If the feature branch commit is an ancestor of HEAD, we are up to date.
		isAlreadyMerged, _ := storage.IsAncestor(featureHeadHash, currentHeadHash)
		if isAlreadyMerged {
			fmt.Println("Already up to date.")
			return nil
		}
		return fmt.Errorf("not a fast-forward merge. Please rebase your branch")
	}

	// This is a fast-forward merge
	// Get the path to the current branch file (e.g., .kitkat/refs/heads/main)
	headData, _ := os.ReadFile(".kitkat/HEAD")
	refPath := strings.TrimSpace(strings.TrimPrefix(string(headData), "ref: "))
	currentBranchFile := filepath.Join(".kitkat", refPath)

	// Update the current branch pointer to the new commit
	if err := os.WriteFile(currentBranchFile, featureHeadHashBytes, 0644); err != nil {
		return fmt.Errorf("failed to update branch pointer: %w", err)
	}

	// Update the working directory and index to match the new HEAD state
	fmt.Printf("Updating files to match %s...\n", featureHeadHash[:7])
	err = UpdateWorkspaceAndIndex(featureHeadHash)
	if err != nil {
		// Attempt to roll back the branch pointer on failure
		os.WriteFile(currentBranchFile, []byte(currentHeadHash), 0644)
		return fmt.Errorf("failed to update workspace: %w", err)
	}

	fmt.Printf("Merge successful. Fast-forwarded to %s\n", featureHeadHash)
	return nil
}

// UpdateWorkspaceAndIndex resets the working directory and index to a specific commit
// Shared logic between checkout, merge, and reset
func UpdateWorkspaceAndIndex(commitHash string) error {
	commit, err := storage.FindCommit(commitHash)
	if err != nil {
		return err
	}
	targetTree, err := storage.ParseTree(commit.TreeHash)
	if err != nil {
		return err
	}

	// Delete files from the current index that are not in the target tree
	currentIndex, _ := storage.LoadIndex()
	for path := range currentIndex {
		if _, existsInTarget := targetTree[path]; !existsInTarget {
			os.Remove(path)
		}
	}

	// Write/update files from the target tree
	for path, hash := range targetTree {
		content, err := storage.ReadObject(hash)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(path, content, 0644); err != nil {
			return err
		}
	}

	// Update the index to match the new tree
	return storage.WriteIndex(targetTree)
}
