package core

import (
	"errors"
	"fmt"
	"os"

	"github.com/LeeFred3042U/kitkat/internal/storage"
)

// RemoveFile removes a file from the index and the working directory.
func RemoveFile(path string) error {
	// 1. Safety Checks
	if !IsSafePath(path) {
		return fmt.Errorf("unsafe path detected: %s", path)
	}

	// Ensure we are inside a kitkat repo
	if _, err := os.Stat(RepoDir); os.IsNotExist(err) {
		return errors.New("not a kitkat repository (run `kitkat init`)")
	}

	// 2. Load the Index (The list of tracked files)
	index, err := storage.LoadIndex()
	if err != nil {
		return err
	}

	// 3. Check if the file is actually tracked
	if _, ok := index[path]; !ok {
		return fmt.Errorf("pathspec '%s' did not match any files", path)
	}

	// 4. Remove from Index (Stop tracking it)
	delete(index, path)

	// 5. Remove from Disk (Delete the actual file)
	// We check if the file exists on disk before trying to delete it
	if _, err := os.Stat(path); err == nil {
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("could not remove file '%s': %v", path, err)
		}
	}

	// 6. Save the updated Index
	return storage.WriteIndex(index)
}
