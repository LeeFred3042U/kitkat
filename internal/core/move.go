package core

import (
	"os"

	"github.com/LeeFred3042U/kitkat/internal/storage"
)

func MoveFile(oldPath, newPath string) error {
	// Read old file
	data, err := os.ReadFile(oldPath)
	if err != nil {
		return err
	}

	// Write new file
	if err := os.WriteFile(newPath, data, 0644); err != nil {
		return err
	}

	// Stage new file
	if err := AddFile(newPath); err != nil {
		return err
	}

	// Load index
	idx, err := storage.LoadIndex()
	if err != nil {
		return err
	}

	// Remove old file from index
	delete(idx, oldPath)

	// Write index
	if err := storage.WriteIndex(idx); err != nil {
		return err
	}

	// Remove old file from disk
	if err := os.Remove(oldPath); err != nil {
		return err
	}

	return nil
}
