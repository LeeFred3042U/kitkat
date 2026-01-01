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

// Lists all tags stored in .kitkat/refs/tags
func ListTags() error {
	if _, err := os.Stat(RepoDir); os.IsNotExist(err) {
		return fmt.Errorf("fatal: not a kitkat repository (or any of the parent directories): %s", RepoDir)
	}

	entries, err := os.ReadDir(tagsDir)
	if err != nil {
		if os.IsNotExist(err) {
			// No tags yet, just return nil (empty list)
			return nil
		}
		return fmt.Errorf("failed to read tags directory: %w", err)
	}

	for _, entry := range entries {
		fmt.Println(entry.Name())
	}
	return nil
}
