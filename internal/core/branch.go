package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

// ListBranches lists all local branches and highlights the current one
func ListBranches() error {
	currentBranch, err := GetHeadState()
	if err != nil {
		// It's possible to be in a detached HEAD state.
		if strings.Contains(err.Error(), "invalid HEAD format") {
			// In a real git, it would show the hash, while we just note it
			currentBranch = "HEAD (detached)"
		} else {
			return err
		}
	}

	// Read all files in the refs/heads directory
	// Each file is a branch
	branches, err := os.ReadDir(headsDir)
	if err != nil {
		return err
	}

	for _, b := range branches {
		if b.Name() == currentBranch {
			// Print the current branch with a '*' and in color.
			fmt.Printf("* %s%s%s\n", colorGreen, b.Name(), colorReset)
		} else {
			fmt.Printf("  %s\n", b.Name())
		}
	}

	return nil
}
