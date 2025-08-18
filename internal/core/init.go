package core

import (
	"os"
	"fmt"
)

const (
	repoDir    = ".kitkat"
	indexPath  = ".kitkat/index"
	objectsDir = ".kitkat/objects"
)

// Sets up the .kitkat structure
func InitRepo() error {
	// Create the main .kitkat directory
	if err := os.Mkdir(".kitkat", 0755); err != nil && !os.IsExist(err) {
		return err
	}

	// Create all nested subdirectories using MkdirAll
	dirs := []string{
		".kitkat/objects",
		".kitkat/refs/heads", // For branches
		".kitkat/refs/tags",  // For tags
	}

	for _, dir := range dirs {
		// Use os.MkdirAll to create parent directories as needed
		if err := os.MkdirAll(dir, 0755); err != nil && !os.IsExist(err) {
			return err
		}
	}

	// Create empty files
	files := []string{".kitkat/index", ".kitkat/commits.log"}
	for _, file := range files {
		f, err := os.Create(file)
		if err != nil {
			return err
		}
		f.Close()
	}

	// Create the HEAD file to point to the default branch (main)
	// This makes the repository immediately ready for the first commit
	headContent := []byte("ref: refs/heads/main")
	if err := os.WriteFile(".kitkat/HEAD", headContent, 0644); err != nil {
		return err
	}

	fmt.Println("Initialized empty KitKat repository in ./.kitkat/")
	return nil
}