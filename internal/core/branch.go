package core

import (
	"os"
	"fmt"
	"errors"
	"strings"
	"path/filepath"

	"github.com/LeeFred3042U/kitkat/internal/storage"
)

const headsDir = ".kitkat/refs/heads"

// Resolves the current commit hash by following the HEAD reference
func readHEAD() (string, error) {
	headData, err := os.ReadFile(".kitkat/HEAD")
	if err != nil {
		return "", err
	}
	ref := strings.TrimSpace(string(headData))
	if !strings.HasPrefix(ref, "ref: ") {
		return "", fmt.Errorf("invalid HEAD format")
	}
	refPath := strings.TrimPrefix(ref, "ref: ")
	
	commitHash, err := os.ReadFile(filepath.Join(".kitkat", refPath))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(commitHash)), nil
}

// Create a new branch pointing to the current HEAD commit
func CreateBranch(name string) error {
	commitHash, err := readHEAD()
	if err != nil {
		// If HEAD can't be read, maybe there are no commits yet
		lastCommit, err := storage.GetLastCommit()
		if err != nil {
			return errors.New("cannot create branch: no commits yet")
		}
		commitHash = lastCommit.ID
	}
	
	if err := os.MkdirAll(headsDir, 0755); err != nil {
		return err
	}
	
	branchPath := filepath.Join(headsDir, name)
	return os.WriteFile(branchPath, []byte(strings.TrimSpace(commitHash)), 0644)
}

// Checks if a branch with the given name exists.
func IsBranch(name string) bool {
	branchPath := filepath.Join(headsDir, name)
	if _, err := os.Stat(branchPath); err == nil {
		return true
	}
	return false
}