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

// ListTags prints all tags stored in .kitkat/refs/tags/
func ListTags() error {
	if _, err := os.Stat(tagsDir); os.IsNotExist(err) {
		// No tags directory means no tags created yet
		return nil
	}

	entries, err := os.ReadDir(tagsDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			fmt.Println(entry.Name())
		}
	}
	return nil
}
