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
	if err := os.WriteFile(tagPath, []byte(commitID), 0644); err != nil {
		return err
	}

	fmt.Printf("Tag '%s' created for commit %s\n", tagName, commitID)
	return nil
}
