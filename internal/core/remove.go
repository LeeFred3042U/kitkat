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
	// Use UpdateIndex to safely update the index transactionally
	return storage.UpdateIndex(func(index map[string]string) error {
		// First, verify the file exists in the index
		if _, ok := index[filename]; !ok {
			return fmt.Errorf("pathspec '%s' did not match any files", filename)
		}

		// Step 1: Delete file from disk FIRST
		if err := os.Remove(filename); err != nil {
			// If file doesn't exist, that's OK (already deleted)
			if !os.IsNotExist(err) {
				// Permission error or other failure - return immediately
				return err
			}
		}

		// Step 2: Only update index if deletion succeeded (or file was already gone)
		delete(index, filename)
		return nil
	})
}
