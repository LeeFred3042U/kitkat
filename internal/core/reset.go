package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/LeeFred3042U/kitkat/internal/storage"
)

// ResetHard resets the current branch head to the specified commit
// and updates the working directory and index to match.
// This is a destructive operation that discards local changes.
func ResetHard(commitHash string) error {
	// 1. Validate the commit exists
	_, err := storage.FindCommit(commitHash)
	if err != nil {
		return fmt.Errorf("commit '%s' not found", commitHash)
	}

	// 2. Update the Branch Pointer (The "Reset")
	headData, err := os.ReadFile(".kitkat/HEAD")
	if err != nil {
		return fmt.Errorf("could not read HEAD: %w", err)
	}
	headContent := strings.TrimSpace(string(headData))

	if strings.HasPrefix(headContent, "ref: ") {
		// Case A: On a Branch
		refPath := strings.TrimPrefix(headContent, "ref: ")
		// Normalize slashes for Windows
		refPath = filepath.FromSlash(refPath)

		fullRefPath := filepath.Join(".kitkat", refPath)
		
		// Ensure directory exists (though it should if we are on it)
		if err := os.MkdirAll(filepath.Dir(fullRefPath), 0755); err != nil {
			return fmt.Errorf("failed to ensure ref dir: %w", err)
		}

		if err := os.WriteFile(fullRefPath, []byte(commitHash), 0644); err != nil {
			return fmt.Errorf("failed to update branch ref: %w", err)
		}
	} else {
		// Case B: Detached HEAD
		// Update HEAD directly
		if err := os.WriteFile(".kitkat/HEAD", []byte(commitHash), 0644); err != nil {
			return fmt.Errorf("failed to update HEAD: %w", err)
		}
	}

	// 3. Update Working Directory & Index (The "Hard" part)
	fmt.Printf("Resetting index and working directory to %s...\n", commitHash[:7])
	if err := UpdateWorkspaceAndIndex(commitHash); err != nil {
		return fmt.Errorf("failed to update workspace: %w", err)
	}

	fmt.Printf("HEAD is now at %s\n", commitHash[:7])
	return nil
}
