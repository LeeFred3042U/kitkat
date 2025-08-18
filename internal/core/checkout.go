package core

import (
	"os"
	"fmt"
	"errors"
	"path/filepath"

	"github.com/LeeFred3042U/kitkat/internal/storage"
)

// Restore a file in the working directory to its state in the last commit
func CheckoutFile(filePath string) error {
	lastCommit, err := storage.GetLastCommit()
	if err != nil {
		return err
	}

	tree, err := storage.ParseTree(lastCommit.TreeHash)
	if err != nil {
		return err
	}

	blobHash, ok := tree[filePath]
	if !ok {
		return errors.New("file not found in the last commit")
	}

	content, err := storage.ReadObject(blobHash)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, content, 0644)
}


// Switch the current HEAD to the named branch and updates the working directory.
func CheckoutBranch(name string) error {
	branchPath := filepath.Join(headsDir, name)
	commitHashBytes, err := os.ReadFile(branchPath)
	if err != nil {
		return fmt.Errorf("branch '%s' not found", name)
	}
	commitHash := string(commitHashBytes)

	// Get the tree of the target commit
	// We need to find the commit object to get its tree hash
	commit, err := storage.FindCommit(commitHash)
	if err != nil {
		return err
	}
	targetTree, err := storage.ParseTree(commit.TreeHash)
	if err != nil {
		return err
	}
	
	// Before making changes, you should check if the user has unstaged work
	// that would be overwritten
	// So the real Git would abort here
	// For now, this is what i have done

	// Update the working directory to match the target tree
	// First, delete files that are not in the target tree
	currentIndex, _ := storage.LoadIndex()
	for path := range currentIndex {
		if _, existsInTarget := targetTree[path]; !existsInTarget {
			os.Remove(path)
		}
	}

	// Now, write/update files from the target tree
	for path, hash := range targetTree {
		content, err := storage.ReadObject(hash)
		if err != nil {
			return err
		}
		// Ensure directory exists before writing file
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(path, content, 0644); err != nil {
			return err
		}
	}

	// Update the index to match the new tree
	if err := storage.WriteIndex(targetTree); err != nil {
		return err
	}
	
	// Update HEAD to point to the new branch
	newHEADContent := fmt.Sprintf("ref: refs/heads/%s", name)
	return os.WriteFile(".kitkat/HEAD", []byte(newHEADContent), 0644)
}