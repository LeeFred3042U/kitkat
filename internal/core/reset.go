package core

import (
	"fmt"
	"maps"
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
		if err := resetIndex(commitHash); err != nil {
			if err = os.WriteFile(".kitcat/HEAD", []byte(oldHead), 0o644); err != nil {
				return fmt.Errorf("failed to update HEAD: %w", err)
			}
			return fmt.Errorf("failed to reset index: %w", err)
		}
		fmt.Printf("HEAD is now at %s %s\n", commitHash[:7], commit.Message)

	case ResetHard:
		if err := resetIndex(commitHash); err != nil {
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

// resetIndex populates index from target commit's tree for mixed/hard reset
func resetIndex(commitHash string) error {
	// Step 1: Get commit to find tree hash
	commit, err := storage.FindCommit(commitHash)
	if err != nil {
		return fmt.Errorf("commit not found: %w", err)
	}

	// Step 2: Parse tree into map
	tree, err := storage.ParseTree(commit.TreeHash)
	if err != nil {
		return fmt.Errorf("failed to parse tree %s: %w", commit.TreeHash, err)
	}

	// Step 3: Copy tree directly to index using maps.Copy
	index := make(map[string]string, len(tree))
	maps.Copy(index, tree)

	// Step 4: Write index file
	return storage.WriteIndex(index)
}

// resetWorkspace restores working directory from target commit using UpdateWorkspaceAndIndex
func resetWorkspace(commitHash string) error {
	// Use the same logic that UpdateWorkspaceAndIndex uses to restore files from commit
	return UpdateWorkspaceAndIndex(commitHash)
}
