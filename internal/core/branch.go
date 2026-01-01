package core

import (
<<<<<<< HEAD
=======
	"errors"
>>>>>>> upstream/main
	"fmt"
	"os"
	"path/filepath"
	"strings"
<<<<<<< HEAD
=======

	"github.com/LeeFred3042U/kitkat/internal/storage"
>>>>>>> upstream/main
)

// --- GLOBAL VARIABLE (Fixes undefined: headsDir) ---
// This allows checkout.go and other files to find the heads directory
var headsDir = filepath.Join(".kitkat", "refs", "heads")

// CreateBranch creates a new branch reference pointing to the current commit
func CreateBranch(branchName string) error {
	// 1. Get the current commit ID (HEAD)
	headPath := filepath.Join(".kitkat", "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return fmt.Errorf("failed to read HEAD: %v", err)
	}
<<<<<<< HEAD
=======
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
>>>>>>> upstream/main

	// Resolve HEAD to a commit hash if it's a ref
	ref := strings.TrimSpace(string(headContent))
	var commitHash string

	if strings.HasPrefix(ref, "ref: ") {
		refPath := strings.TrimPrefix(ref, "ref: ")
		fullRefPath := filepath.Join(".kitkat", refPath)
		hashBytes, err := os.ReadFile(fullRefPath)
		if err != nil {
<<<<<<< HEAD
			return fmt.Errorf("failed to resolve HEAD ref: %v", err)
=======
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
>>>>>>> upstream/main
		}
		commitHash = strings.TrimSpace(string(hashBytes))
	} else {
		commitHash = ref // Detached HEAD
	}

	// 2. Create the new branch file
	// We use the global variable headsDir here too
	branchPath := filepath.Join(headsDir, branchName)
	if err := os.WriteFile(branchPath, []byte(commitHash), 0644); err != nil {
		return fmt.Errorf("failed to create branch file: %v", err)
	}

	fmt.Printf("Created branch '%s'\n", branchName)
	return nil
}
<<<<<<< HEAD

// ListBranches displays all local branches
func ListBranches() error {
	// Use the global variable
	files, err := os.ReadDir(headsDir)
	if err != nil {
		return fmt.Errorf("failed to read branches: %v", err)
	}

	headPath := filepath.Join(".kitkat", "HEAD")
	headContent, _ := os.ReadFile(headPath)
	currentRef := strings.TrimSpace(string(headContent))

	for _, file := range files {
		prefix := "  "
		branchRef := "ref: refs/heads/" + file.Name()
		if currentRef == branchRef {
			prefix = "* " // Mark current branch
		}
		fmt.Println(prefix + file.Name())
	}
	return nil
}

// DeleteBranch removes a branch reference safely
func DeleteBranch(branchName string) error {
	// 1. Define paths using the global variable
	branchPath := filepath.Join(headsDir, branchName)
	headPath := filepath.Join(".kitkat", "HEAD")

	// 2. Check if branch exists
	if _, err := os.Stat(branchPath); os.IsNotExist(err) {
		return fmt.Errorf("branch '%s' not found", branchName)
	}

	// 3. Safety Check: Prevent deletion if it is the current HEAD
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return fmt.Errorf("could not read HEAD: %v", err)
	}

	currentHeadRef := strings.TrimSpace(string(headContent))
	targetRef := "ref: refs/heads/" + branchName

	if currentHeadRef == targetRef {
		return fmt.Errorf("cannot delete branch '%s' checked out at '%s'", branchName, currentHeadRef)
	}

	// 4. Delete the branch file
	if err := os.Remove(branchPath); err != nil {
		return fmt.Errorf("failed to delete branch: %v", err)
	}

	fmt.Printf("Deleted branch %s\n", branchName)
	return nil
}

// IsBranch checks if a branch with the given name exists
func IsBranch(name string) bool {
	path := filepath.Join(headsDir, name)
	_, err := os.Stat(path)
	return err == nil
}

// readHEAD is a helper function to get the current HEAD reference or hash
func readHEAD() (string, error) {
	headPath := filepath.Join(".kitkat", "HEAD")
	data, err := os.ReadFile(headPath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}
=======
>>>>>>> upstream/main
