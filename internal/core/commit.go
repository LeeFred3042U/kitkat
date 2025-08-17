package core

import (
	"crypto/sha1"
	"encoding/hex"
	"time"

	"github.com/LeeFred3042U/kitkat/internal/models"
	"github.com/LeeFred3042U/kitkat/internal/storage"
)

// Computes a SHA-1 over the commit contents.
func hashCommit(c models.Commit) string {
	h := sha1.New()
	// Minimal fields to hash; add Author/Parent later
	h.Write([]byte(c.TreeHash))
	h.Write([]byte(c.Message))
	h.Write([]byte(c.Timestamp.UTC().String()))
	return hex.EncodeToString(h.Sum(nil))
}

func Commit(message string) (string, error) {
	// Create a tree from the current index
	treeHash, err := storage.CreateTree()
	if err != nil {
		return "", err
	}

	// Create the commit object without ID first
	commit := models.Commit{
		Message:   message,
		Timestamp: time.Now().UTC(),
		TreeHash:  treeHash,
	}

	// Now compute a content addressed ID
	commit.ID = hashCommit(commit)

	// Append the commit to storage
	if err := storage.AppendCommit(commit); err != nil {
		return "", err
	}

	return commit.ID, nil
}