package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/LeeFred3042U/kitkat/internal/storage"
)

// --- GLOBAL VARIABLE ---
var headsDir = filepath.Join(".kitkat", "refs", "heads")

// CreateBranch creates a new branch reference pointing to the current commit
func CreateBranch(name string) error {
	lastCommit, err := storage.GetLastCommit()
	if err != nil {
		return errors.New("cannot create branch: no commits yet")
	}
	commitHash := lastCommit.ID

	if err := os.MkdirAll(headsDir, 0755); err != nil {
		return err
	}

	branchPath := filepath.Join(headsDir, name)
	return os.WriteFile(branchPath, []byte(strings.TrimSpace(commitHash)), 0644)
}

// IsBranch checks if a branch with the given name exists.
func IsBranch(name string) bool {
	branchPath := filepath.Join(headsDir, name)
	if _, err := os.Stat(branchPath); err == nil {
		return true
	}
	return false
}

// ListBranches lists all local branches and highlights the current one
func ListBranches() error {
	files, err := os.ReadDir(headsDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No branches found.")
			return nil
		}
		return fmt.Errorf("failed to read branches: %v", err)
	}

	headPath := filepath.Join(".kitkat", "HEAD")
	headContent, _ := os.ReadFile(headPath)
	currentRef := strings.TrimSpace(string(headContent))

	for _, file := range files {
		prefix := "  "
		branchRef := "ref: refs/heads/" + file.Name()
		if currentRef == branchRef {
			prefix = "* "
		}
		fmt.Println(prefix + file.Name())
	}
	return nil
}

// GetHeadState returns the current branch name or hash
// GetHeadState returns the current branch name or hash
func GetHeadState() (string, error) {
	headPath := filepath.Join(".kitkat", "HEAD")
	headData, err := os.ReadFile(headPath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(headData)), nil
}

// --- YOUR FEATURE: DeleteBranch ---

// DeleteBranch removes a branch reference safely
func DeleteBranch(branchName string) error {
	branchPath := filepath.Join(headsDir, branchName)
	headPath := filepath.Join(".kitkat", "HEAD")

	if _, err := os.Stat(branchPath); os.IsNotExist(err) {
		return fmt.Errorf("branch '%s' not found", branchName)
	}

	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return fmt.Errorf("could not read HEAD: %v", err)
	}

	currentHeadRef := strings.TrimSpace(string(headContent))
	targetRef := "ref: refs/heads/" + branchName

	if currentHeadRef == targetRef {
		return fmt.Errorf("cannot delete branch '%s' checked out at '%s'", branchName, currentHeadRef)
	}

	if err := os.Remove(branchPath); err != nil {
		return fmt.Errorf("failed to delete branch: %v", err)
	}

	fmt.Printf("Deleted branch %s\n", branchName)
	return nil
}

// readHEAD is a helper function to get the current HEAD reference or hash
func readHEAD() (string, error) {
	return GetHeadState()
}
