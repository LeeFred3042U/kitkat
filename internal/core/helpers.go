package core

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/LeeFred3042U/kitkat/internal/storage"
)

 

func IsSafePath(path string) bool {
	// Clean the path to resolve ".." patterns
	cleanedPath := filepath.Clean(path)

	// A safe path must not be absolute and must not try to go "up" the directory tree
	return !filepath.IsAbs(cleanedPath) && !strings.HasPrefix(cleanedPath, "..")
}

// IsWorkDirDirty checks for any tracked files that have been modified or deleted
// in the working directory but not yet staged
// This is crucial for preventing
// data loss during operations like checkout or merge
func IsWorkDirDirty() (bool, error) {
	index, err := storage.LoadIndex()
	if err != nil {
		return false, err
	}

	for path, indexHash := range index {
		// Check if a tracked file has been deleted from the working directory
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return true, nil // Dirty: file in index is missing from disk
		}

		// Check if a tracked file has been modified.
		currentHash, err := storage.HashFile(path)
		if err != nil {
			// Can't hash the file, might be a permissions issue
			// Treat as an error rather than a dirty state
			return false, err
		}

		if currentHash != indexHash {
			return true, nil // Dirty: hashes don't match
		}
	}

	return false, nil // Not dirty
}
