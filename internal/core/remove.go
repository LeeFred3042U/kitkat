package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/LeeFred3042U/kitcat/internal/storage"
)

func RemoveFile(filename string) error {
	filename = filepath.Clean(filename)
	if !IsSafePath(filename) {
		return fmt.Errorf("unsafe path detected: %s", filename)
	}
	index, err := storage.LoadIndex()
	if err != nil {
		return fmt.Errorf("failed to load index: %w", err)
	}

	// First, verify the file exists in the index
	indexHash, ok := index[filename]
	if !ok {
		return fmt.Errorf("pathspec '%s' did not match any files", filename)
	}

	// Check for uncommitted changes before deletion
	diskHash, err := storage.HashFile(filename)
	if err != nil {
		return fmt.Errorf("failed to hash file: %w", err)
	}

	if diskHash != indexHash {
		return fmt.Errorf("local changes present")
	}

	// Step 1: Delete file from disk FIRST
	if err := os.Remove(filename); err != nil {
		// If file doesn't exist, that's OK (already deleted)
		if !os.IsNotExist(err) {
			// Permission error or other failure - return immediately
			return err
		}
	}

	// Step 2: Only update index if deletion succeeded
	delete(index, filename)

	if err := storage.WriteIndex(index); err != nil {
		return fmt.Errorf("failed to save index: %w", err)
	}

	return nil
}
