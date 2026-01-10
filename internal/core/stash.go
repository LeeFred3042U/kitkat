package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/LeeFred3042U/kitcat/internal/models"
	"github.com/LeeFred3042U/kitcat/internal/storage"
)

// StashSave saves the current state (index and working directory) to the stash stack
// and cleans the working directory.
func StashSave(message string) error {
	if !IsRepoInitialized() {
		return fmt.Errorf("not a kitkat repository")
	}

	// 1. Check if there are changes to stash
	isDirty, err := IsWorkDirDirty()
	if err != nil {
		return fmt.Errorf("failed to check status: %w", err)
	}
	if !isDirty {
		fmt.Println("No local changes to save")
		return nil
	}

	// 2. Get current HEAD info
	headCommit, err := GetHeadCommit()
	if err != nil {
		// If no commits yet, we can't really stash against anything?
		// Git allows stashing from emptiness? No, need a parent.
		return fmt.Errorf("cannot stash without any commits")
	}
	headState, _ := GetHeadState()
	branchName := strings.TrimPrefix(headState, "refs/heads/")

	// Default message if empty
	if message == "" {
		message = fmt.Sprintf("WIP on %s: %s %s", branchName, headCommit.ID[:7], headCommit.Message)
	}

	authorName, _, _ := GetConfig("user.name")
	if authorName == "" {
		authorName = "Unknown"
	}
	authorEmail, _, _ := GetConfig("user.email")
	if authorEmail == "" {
		authorEmail = "unknown@example.com"
	}

	// 3. Create Index Commit (I)
	// Capture the current index state into a tree
	indexTreeHash, err := storage.CreateTree()
	if err != nil {
		return fmt.Errorf("failed to create index tree: %w", err)
	}

	commitI := models.Commit{
		Parent:      headCommit.ID,
		Message:     fmt.Sprintf("index on %s: %s %s", branchName, headCommit.ID[:7], headCommit.Message),
		Timestamp:   time.Now().UTC(),
		TreeHash:    indexTreeHash,
		AuthorName:  authorName,
		AuthorEmail: authorEmail,
	}
	commitI.ID = hashCommit(commitI)
	if err := storage.AppendCommit(commitI); err != nil {
		return fmt.Errorf("failed to save index commit: %w", err)
	}

	// 4. Create Workdir Commit (W)
	// We need to capture untracked/modified files that are NOT in the index yet.
	// So we effectively do 'git add -A' behind the scenes to capture the workdir state.
	// NOTE: This modifies the Index file on disk! But we will reset it later.
	if err := AddAll(); err != nil {
		return fmt.Errorf("failed to stage workdir changes: %w", err)
	}
	workTreeHash, err := storage.CreateTree()
	if err != nil {
		return fmt.Errorf("failed to create workdir tree: %w", err)
	}

	commitW := models.Commit{
		Parent:      headCommit.ID, // Parent 1: HEAD
		Message:     message,
		Timestamp:   time.Now().UTC(),
		TreeHash:    workTreeHash,
		AuthorName:  authorName,
		AuthorEmail: authorEmail,
		// We store the Index Commit ID in the message or a separate field if we had one?
		// Git stores it as a second parent (merge commit).
		// Our Commit model assumes single Parent string.
		// WORKAROUND: We will store the Index Commit hash in the message or handle it specially.
		// BETTER: For now, let's just chain them: HEAD <- I <- W.
		// Wait, if we chain them (W parent is I), then `git log` looks linear.
		// Git structure: W has parents (HEAD, I).
		// Since Kitkat Commit only has one Parent field, we can't represent a merge commit fully.
		// Adapting plan: Make W's parent = I.
		// Then W contains everything in I + Workdir.
		// So when we pop, we just restore W.
		// Logic: HEAD -> ...
		//        I (parent: HEAD)
		//        W (parent: I)
		// This makes W effectively a descendant of HEAD.
	}
	// For stash, we want W to represent the final state.
	// If W.Parent = HEAD, it's a sibling of I.
	// If Kitkat doesn't support merge commits (2 parents), we have to improvise.
	// Let's set W.Parent = I.ID.
	commitW.Parent = commitI.ID

	commitW.ID = hashCommit(commitW)
	if err := storage.AppendCommit(commitW); err != nil {
		return fmt.Errorf("failed to save workdir commit: %w", err)
	}

	// 5. Update Stash Ref and Log
	// Update .kitkat/refs/stash
	if err := os.MkdirAll(filepath.Dir(StashPath), 0755); err != nil {
		return err
	}
	if err := SafeWrite(StashPath, []byte(commitW.ID), 0644); err != nil {
		return err
	}

	// Update Stash Log (Stack)
	// Format: old_hash new_hash committer timestamp message
	// We can cheat and just store the Commit IDs line by line for now since we don't have a rigid reflog format yet.
	// Let's just append the Commit W hash to the log file -> Stack
	if err := os.MkdirAll(filepath.Dir(StashLogPath), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(StashLogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	// Storing just the hash for simplicity in MVP.
	// But `git stash list` needs a message. We can look it up from the commit.
	if _, err := fmt.Fprintln(f, commitW.ID); err != nil {
		return err
	}

	// 6. Reset Workspace to HEAD
	if err := ResetHard(headCommit.ID); err != nil {
		return fmt.Errorf("failed to reset workspace after stash: %w", err)
	}

	fmt.Printf("Saved working directory and index state WIP on %s: %s %s\n", branchName, headCommit.ID[:7], headCommit.Message)
	return nil
}

// StashList lists all stashed changes
func StashList() error {
	if !IsRepoInitialized() {
		return fmt.Errorf("not a kitkat repository")
	}

	if _, err := os.Stat(StashLogPath); os.IsNotExist(err) {
		return nil
	}

	content, err := os.ReadFile(StashLogPath)
	if err != nil {
		return err
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	// Process in reverse order (stack)
	for i := len(lines) - 1; i >= 0; i-- {
		hash := strings.TrimSpace(lines[i])
		if hash == "" {
			continue
		}
		commit, err := storage.FindCommit(hash)
		if err != nil {
			fmt.Printf("stash@{%d}: %s (commit not found)\n", len(lines)-1-i, hash)
			continue
		}
		fmt.Printf("stash@{%d}: %s\n", len(lines)-1-i, commit.Message)
	}
	return nil
}

// StashPop applies the stash at the given index and then drops it
func StashPop(index int) error {
	if err := StashApply(index); err != nil {
		return err
	}
	return StashDrop(index)
}

// StashApply applies the stash at the given index to the current working directory
func StashApply(index int) error {
	commitW, err := getStashCommit(index)
	if err != nil {
		return err
	}

	// 1. Apply Index Commit (I)
	// I is the parent of W in our simplified model
	commitIParams, err := storage.FindCommit(commitW.Parent)
	if err != nil {
		return fmt.Errorf("corrupt stash: index commit not found: %w", err)
	}

	// Apply changes from I (staged changes)
	// CherryPick(I, true) -> applies diff(I.Parent, I) -> diff(HEAD, I) if we were on same HEAD.
	// If we are on different HEAD, it applies the diff to the new HEAD.
	fmt.Println("Restoring index...")
	if err := CherryPick(commitIParams.ID, true); err != nil {
		return fmt.Errorf("conflict applying index: %w", err)
	}

	// 2. Apply Workdir Commit (W)
	// Apply changes from W (unstaged changes relative to I)
	fmt.Println("Restoring working directory...")
	if err := CherryPick(commitW.ID, true); err != nil {
		return fmt.Errorf("conflict applying working directory: %w", err)
	}

	return nil
}

// StashDrop removes the stash at the given index
func StashDrop(index int) error {
	lines, err := readStashLog()
	if err != nil {
		return err
	}

	// Convert visual index (0 = latest = last line) to actual index
	actualIndex := len(lines) - 1 - index
	if actualIndex < 0 || actualIndex >= len(lines) {
		return fmt.Errorf("stash@{%d} not found", index)
	}

	// Remove element
	lines = append(lines[:actualIndex], lines[actualIndex+1:]...)

	if err := writeStashLog(lines); err != nil {
		return err
	}

	// Update refs/stashRef to point to new latest (or delete if empty)
	if len(lines) > 0 {
		lastHash := strings.TrimSpace(lines[len(lines)-1])
		if err := SafeWrite(StashPath, []byte(lastHash), 0644); err != nil {
			return err
		}
	} else {
		os.Remove(StashPath)
	}

	fmt.Printf("Dropped stash@{%d}\n", index)
	return nil
}

// Helper to get stash commit by index (0 = latest)
func getStashCommit(index int) (models.Commit, error) {
	lines, err := readStashLog()
	if err != nil {
		return models.Commit{}, err
	}

	if len(lines) == 0 {
		return models.Commit{}, fmt.Errorf("stash list is empty")
	}

	actualIndex := len(lines) - 1 - index
	if actualIndex < 0 || actualIndex >= len(lines) {
		return models.Commit{}, fmt.Errorf("stash@{%d} not found", index)
	}

	hash := strings.TrimSpace(lines[actualIndex])
	return storage.FindCommit(hash)
}

func readStashLog() ([]string, error) {
	if _, err := os.Stat(StashLogPath); os.IsNotExist(err) {
		return []string{}, nil
	}
	content, err := os.ReadFile(StashLogPath)
	if err != nil {
		return nil, err
	}
	var lines []string
	for _, line := range strings.Split(strings.TrimSpace(string(content)), "\n") {
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}
	return lines, nil
}

func writeStashLog(lines []string) error {
	// Rebuild content
	content := strings.Join(lines, "\n")
	if len(lines) > 0 {
		content += "\n" // Trailing newline
	}
	return SafeWrite(StashLogPath, []byte(content), 0644)
}
