package core

import (
	"fmt"
	"os"
	"path/filepath"
)

const tagsDir = ".kitkat/refs/tags"

// Creates a new lightweight tag pointing to a specific commit
func CreateTag(tagName, commitID string) error {
	if err := os.MkdirAll(tagsDir, 0755); err != nil {
		return err
	}

	tagPath := filepath.Join(tagsDir, tagName)
	// Checks if tag already exists.
	if _, err := os.Stat(tagPath); err == nil {
		return fmt.Errorf("Error: tag %s already exists", tagName)
	} else if !os.IsNotExist(err) {
		return err
	}

	// Creates a new tag.
	if err := os.WriteFile(tagPath, []byte(commitID), 0644); err != nil {
		return err
	}

	fmt.Printf("Tag '%s' created for commit %s\n", tagName, commitID)
	return nil
}
