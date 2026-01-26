package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/LeeFred3042U/kitcat/internal/storage"
)

func RemoveFile(filename string, recursive bool) error {
	filename = filepath.Clean(filename)
	if !IsSafePath(filename) {
		return fmt.Errorf("unsafe path detected: %s", filename)
	}

	return storage.UpdateIndex(func(index map[string]string) error {
		var filesToRemove []string

		if recursive {
			// Recursive: find ALL tracked files under this directory
			for trackedFile := range index {
				if trackedFile == filename ||
					strings.HasPrefix(trackedFile, filename+string(filepath.Separator)) {
					filesToRemove = append(filesToRemove, trackedFile)
				}
			}
		} else {
			// Single file mode
			if _, ok := index[filename]; !ok {
				return fmt.Errorf("pathspec '%s' did not match any files", filename)
			}
			filesToRemove = []string{filename}
		}

		if len(filesToRemove) == 0 {
			return fmt.Errorf("pathspec '%s' did not match any files", filename)
		}

		// Step 1: Remove files from disk (non-fatal if missing)
		for _, filePath := range filesToRemove {
			if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
				fmt.Printf("Warning: could not remove %s: %v\n", filePath, err)
			}
		}

		// Step 2: Remove from index
		for _, filePath := range filesToRemove {
			delete(index, filePath)
		}

		// Success message
		if recursive && len(filesToRemove) > 1 {
			fmt.Printf("Removed %d tracked files under '%s'\n", len(filesToRemove), filename)
		} else {
			fmt.Printf("rm '%s'\n", filename)
		}

		return nil
	})
}
