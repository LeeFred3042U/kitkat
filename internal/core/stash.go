package core

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/LeeFred3042U/kitcat/internal/models"
	"github.com/LeeFred3042U/kitcat/internal/storage"
)

// Stash saves the current working directory and index state to a temporary storage area.
// It creates a "WIP" commit containing the current index state and then performs a hard
// reset to HEAD, cleaning the workspace. This allows users to switch branches or pull
// updates without losing their work-in-progress.
func Stash() error {
	// Step 1: Validate repository is initialized
	if !IsRepoInitialized() {
		return fmt.Errorf("fatal: not a kitcat repository (or any of the parent directories): .kitcat")
	}

	// Step 2: Get current HEAD commit for parent reference and message
	// This must be done before checking dirty state because IsWorkDirDirty needs commits to exist
	headCommit, err := GetHeadCommit()
	if err != nil {
		// Check if the error is because there are no commits
		if err == storage.ErrNoCommits || strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("cannot stash: no commits yet")
		}
		return fmt.Errorf("failed to get HEAD commit: %w", err)
	}

	// Step 3: Check if there are any changes to stash
	isDirty, err := IsWorkDirDirty()
	if err != nil {
		return fmt.Errorf("failed to check working directory status: %w", err)
	}
	if !isDirty {
		return fmt.Errorf("nothing to stash, working tree clean")
	}

	// Step 4: Get current branch name for WIP message
	branchName, err := GetHeadState()
	if err != nil {
		branchName = "detached HEAD"
	}

	// Step 5: Update index with current working directory state for tracked files
	// This ensures unstaged changes are included in the stash tree
	index, err := storage.LoadIndex()
	if err != nil {
		return fmt.Errorf("failed to load index: %w", err)
	}

	for path := range index {
		// If file exists in working directory, hash it and update index
		if _, err := os.Stat(path); err == nil {
			hash, err := storage.HashAndStoreFile(path)
			if err != nil {
				return fmt.Errorf("failed to hash file %s: %w", path, err)
			}
			index[path] = hash
		}
	}

	// Write updated index to disk so CreateTree uses the current state
	if err := storage.WriteIndex(index); err != nil {
		return fmt.Errorf("failed to write updated index: %w", err)
	}

	// Step 6: Create tree from current index
	treeHash, err := storage.CreateTree()
	if err != nil {
		return fmt.Errorf("failed to create tree from index: %w", err)
	}

	// Step 6: Get author information
	authorName, _, _ := GetConfig("user.name")
	if authorName == "" {
		authorName = "Unknown"
	}
	authorEmail, _, _ := GetConfig("user.email")
	if authorEmail == "" {
		authorEmail = "unknown@example.com"
	}

	// Step 7: Create WIP commit message
	// Format: "WIP on <branch>: <latest_commit_message>"
	wipMessage := fmt.Sprintf("WIP on %s: %s", branchName, headCommit.Message)

	// Step 8: Create the stash commit
	stashCommit := models.Commit{
		Parent:      headCommit.ID,
		Message:     wipMessage,
		Timestamp:   time.Now().UTC(),
		TreeHash:    treeHash,
		AuthorName:  authorName,
		AuthorEmail: authorEmail,
	}
	stashCommit.ID = hashCommit(stashCommit)

	// Step 9: Save the stash commit to commits.log
	if err := storage.AppendCommit(stashCommit); err != nil {
		return fmt.Errorf("failed to save stash commit: %w", err)
	}

	// Step 10: Write the stash reference
	stashRefPath := ".kitcat/refs/stash"
	if err := SafeWrite(stashRefPath, []byte(stashCommit.ID), 0o644); err != nil {
		return fmt.Errorf("failed to write stash reference: %w", err)
	}

	// Step 11: Perform hard reset to HEAD to clean the workspace
	if err := ResetHard(headCommit.ID); err != nil {
		// Attempt to clean up the stash reference on failure
		_ = os.Remove(stashRefPath)
		return fmt.Errorf("failed to reset workspace after stashing: %w", err)
	}

	return nil
}

// StashPop applies the most recent stash to the working directory and removes it.
// It reads the stash commit, applies it to the workspace, and deletes the stash reference.
// This operation will fail if the working directory has uncommitted changes to prevent data loss.
func StashPop() error {
	// Step 1: Validate repository is initialized
	if !IsRepoInitialized() {
		return fmt.Errorf("fatal: not a kitcat repository (or any of the parent directories): .kitcat")
	}

	// Step 2: Check if stash exists
	stashRefPath := ".kitcat/refs/stash"
	stashHashBytes, err := os.ReadFile(stashRefPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no stash found")
		}
		return fmt.Errorf("failed to read stash reference: %w", err)
	}

	stashHash := strings.TrimSpace(string(stashHashBytes))
	if stashHash == "" {
		return fmt.Errorf("stash reference is empty")
	}

	// Step 3: Verify the stash commit exists
	stashCommit, err := storage.FindCommit(stashHash)
	if err != nil {
		return fmt.Errorf("stash commit not found: %w", err)
	}

	// Step 4: Check if working directory is clean to prevent data loss
	isDirty, err := IsWorkDirDirty()
	if err != nil {
		return fmt.Errorf("failed to check working directory status: %w", err)
	}
	if isDirty {
		return fmt.Errorf("error: your local changes would be overwritten by stash pop\nPlease commit your changes or stash them before you pop")
	}

	// Step 5: Apply the stash commit to the working directory
	if err := UpdateWorkspaceAndIndex(stashHash); err != nil {
		return fmt.Errorf("failed to apply stash: %w", err)
	}

	// Step 6: Remove the stash reference
	if err := os.Remove(stashRefPath); err != nil {
		// Log warning but don't fail - the stash was already applied
		fmt.Fprintf(os.Stderr, "Warning: failed to remove stash reference: %v\n", err)
	}

	// Step 7: Print success message with commit info
	fmt.Printf("On branch %s\n", getCurrentBranchName())
	fmt.Printf("Dropped refs/stash (%s)\n", stashCommit.ID[:7])

	return nil
}

// getCurrentBranchName is a helper to get the current branch name
func getCurrentBranchName() string {
	headState, err := GetHeadState()
	if err != nil {
		return "unknown"
	}
	return headState
}
