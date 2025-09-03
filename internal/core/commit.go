package core

import (
	"os"
	"fmt"
	"time"
	"strings"
	"crypto/sha1"
	"encoding/hex"
	"path/filepath"

	"github.com/LeeFred3042U/kitkat/internal/models"
	"github.com/LeeFred3042U/kitkat/internal/storage"
)

// Computes a SHA-1 over the commit contents, including the parent
func hashCommit(c models.Commit) string {
	h := sha1.New()
	h.Write([]byte(c.TreeHash))
	h.Write([]byte(c.Parent))
	h.Write([]byte(c.Message))
	h.Write([]byte(c.Timestamp.UTC().Format(time.RFC3339Nano)))
	return hex.EncodeToString(h.Sum(nil))
}

// Commit creates a new commit object with the current index contents
func Commit(message string) (string, error) {
	// Create a tree from the current index
	treeHash, err := storage.CreateTree()
	if err != nil {
		return "", err
	}

	// Get the last commit to set as the parent of this new commit
	var parentID string
	// We read the current branch's head to find the parent
	currentBranchRefPath, err := getCurrentBranchRefPath()
	if err != nil {
		// If there's no ref path, it could be the first commit
		// Let's check if any commits exist at all
		_, err := storage.GetLastCommit()
		if err != storage.ErrNoCommits {
			return "", fmt.Errorf("could not resolve HEAD: %w", err)
		}
		// It is the first commit, parentID remains empty
	} else {
		branchPath := filepath.Join(".kitkat", currentBranchRefPath)
			if parentBytes, err := os.ReadFile(branchPath); err == nil {
			    parentID = strings.TrimSpace(string(parentBytes))
			}
		// If the file doesn't exist, it's the first commit on this branch. parentID remains empty.
	}


	// Create the commit object
	commit := models.Commit{
		Parent:    parentID,
		Message:   message,
		Timestamp: time.Now().UTC(),
		TreeHash:  treeHash,
	}

	// Now compute a content addressed ID for the new commit
	commit.ID = hashCommit(commit)

	// Append the commit to the commit log
	if err := storage.AppendCommit(commit); err != nil {
		return "", err
	}

	// After saving the commit, update the current branch to point to it
	refPath, err := getCurrentBranchRefPath()
	if err != nil {
		return "", fmt.Errorf("could not get current branch reference for update: %w", err)
	}

	branchFilePath := filepath.Join(".kitkat", refPath)

	// Write the new commit ID to the branch file (e.g., .kitkat/refs/heads/main)
	if err := os.WriteFile(branchFilePath, []byte(commit.ID), 0644); err != nil {
		return "", fmt.Errorf("failed to update branch pointer: %w", err)
	}

	return commit.ID, nil
}


// getCurrentBranchRefPath reads the HEAD file to find the path to the current branch file
func getCurrentBranchRefPath() (string, error) {
	headData, err := os.ReadFile(".kitkat/HEAD")
	if err != nil {
		return "", err
	}
	ref := strings.TrimSpace(string(headData))
	if !strings.HasPrefix(ref, "ref: ") {
		return "", fmt.Errorf("invalid HEAD format: %s", ref)
	}
	return strings.TrimPrefix(ref, "ref: "), nil
}

// CommitAll stages all changes and then creates a new commit
// This is the logic for the "commit -am" shortcut
func CommitAll(message string) (string, error) {
	// Stage all changes in the repository, just like 'add -A'
	if err := AddAll(); err != nil { // We need to call AddAll from the core package
		return "", fmt.Errorf("failed to stage changes before committing: %w", err)
	}

	// Call the existing Commit function to create the snapshot
	// By now, the index is up-to-date with all changes
	return Commit(message)
}