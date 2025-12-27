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

// ListTags prints all existing tags and the commit they point to
func ListTags() error {
	if _, err := os.Stat(tagsDir); os.IsNotExist(err) {
		fmt.Println("No tags found.")
		return nil
	}

	entries, err := os.ReadDir(tagsDir)
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		fmt.Println("No tags found.")
		return nil
	}

	fmt.Println("Tags:")
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		tagName := entry.Name()
		tagPath := filepath.Join(tagsDir, tagName)
		commitBytes, err := os.ReadFile(tagPath)
		if err != nil {
			fmt.Printf("  %s (error reading commit ID)\n", tagName)
			continue
		}

		fmt.Printf("  %s -> %s\n", tagName, string(commitBytes))
	}

	return nil
}
