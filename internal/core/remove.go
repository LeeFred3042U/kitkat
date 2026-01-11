package core

import (
	"fmt"
	"os"
	"path/filepath"
)

func RemoveFile(filename string) error {
	filename = filepath.Clean(filename)
	index, err := LoadIndex()
	if err != nil {
		return fmt.Errorf("failed to load index: %w", err)
	}

	found := false
	newIndex := []IndexEntry{}
	for _, entry := range index {
		if filepath.Clean(entry.Path) == filename {
			found = true
			continue
		}
		newIndex = append(newIndex, entry)
	}

	if !found {
		return fmt.Errorf("pathspec '%s' did not match any files", filename)
	}

	if err := SaveIndex(newIndex); err != nil {
		return fmt.Errorf("failed to save index: %w", err)
	}

	if err := os.Remove(filename); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	return nil
}
