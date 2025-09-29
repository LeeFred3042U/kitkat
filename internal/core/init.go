package core

import (
	"fmt"
	"os"
)

// InitRepo sets up the .kitkat directory structure.
func InitRepo() error {
	// Create all necessary subdirectories using the public constants.
	dirs := []string{
		RepoDir,
		ObjectsDir,
		HeadsDir,
		TagsDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil && !os.IsExist(err) {
			return err
		}
	}

	// Create empty files.
	files := []string{IndexPath, CommitsPath}
	for _, file := range files {
		f, err := os.Create(file)
		if err != nil {
			return err
		}
		f.Close()
	}

	// Create the HEAD file to point to the default branch (main).
	headContent := []byte("ref: refs/heads/main")
	if err := os.WriteFile(HeadPath, headContent, 0644); err != nil {
		return err
	}

	fmt.Printf("Initialized empty KitKat repository in ./%s/\n", RepoDir)
	return nil
}