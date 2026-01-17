package core

import (
	"fmt"
	"os"
	"strings"

	"github.com/LeeFred3042U/kitcat/internal/storage"
)

const (
	ResetHard  = "hard"
	ResetSoft  = "soft"
	ResetMixed = "mixed"
)

// Reset performs reset operation with specified mode
// Modes: "soft", "mixed", "hard"
func Reset(commitHash string, mode string) error {
	if !IsRepoInitialized() {
		return fmt.Errorf("not a kitcat repository (or any of the parent directories): .kitcat")
	}

	// Step 1: Validate commit exists
	commit, err := storage.FindCommit(commitHash)
	if err != nil {
		return fmt.Errorf("fatal: invalid commit: %s", commitHash)
	}

	// Step 2: Backup current HEAD
	headData, err := os.ReadFile(".kitcat/HEAD")
	if err != nil {
		return fmt.Errorf("fatal: unable to read HEAD: %w", err)
	}
	oldHead := strings.TrimSpace(string(headData))

	// Step 3: Move HEAD (ALL modes)
	if err := os.WriteFile(".kitcat/HEAD", []byte(commitHash), 0o644); err != nil {
		return fmt.Errorf("failed to update HEAD: %w", err)
	}

	// Step 4: Mode-specific operations
	switch mode {
	case ResetSoft:
		fmt.Printf("HEAD is now at %s %s\n", commitHash[:7], commit.Message)

	case ResetMixed:
		if err := resetIndex(); err != nil {
			if err = os.WriteFile(".kitcat/HEAD", []byte(oldHead), 0o644); err != nil {
				return fmt.Errorf("failed to update HEAD: %w", err)
			}
			return fmt.Errorf("failed to reset index: %w", err)
		}
		fmt.Printf("HEAD is now at %s %s\n", commitHash[:7], commit.Message)

	case ResetHard:
		if err := resetIndex(); err != nil {
			if err = os.WriteFile(".kitcat/HEAD", []byte(oldHead), 0o644); err != nil {
				return fmt.Errorf("failed to update HEAD: %w", err)
			}
			return fmt.Errorf("failed to reset index: %w", err)
		}
		if err := resetWorkspace(commitHash); err != nil {
			if err = os.WriteFile(".kitcat/HEAD", []byte(oldHead), 0o644); err != nil {
				return fmt.Errorf("failed to update HEAD: %w", err)
			}
			return fmt.Errorf("failed to reset workspace: %w", err)
		}
		fmt.Printf("HEAD is now at %s %s\n", commitHash[:7], commit.Message)

	default:
		if err = os.WriteFile(".kitcat/HEAD", []byte(oldHead), 0o644); err != nil {
			return fmt.Errorf("failed to update HEAD: %w", err)
		}
		return fmt.Errorf("unknown reset mode: %s. Use --soft, --mixed, or --hard", mode)
	}

	return nil
}

// resetIndex clears index for mixed/hard reset
func resetIndex() error {
	// Clear index completely to match target tree state
	return os.WriteFile(".kitcat/index", []byte{}, 0o644)
}

// resetWorkspace restores working directory from target commit using UpdateWorkspaceAndIndex
func resetWorkspace(commitHash string) error {
	// Use the same logic that UpdateWorkspaceAndIndex uses to restore files from commit
	return UpdateWorkspaceAndIndex(commitHash)
}
