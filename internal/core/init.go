package core

import (
	"fmt"
	"os"
	"path/filepath"
)

const hint = "\033[33m"

// isPathExist checks if a path exist or not
func isPathExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// InitRepo sets up the .kitkat directory structure.
func InitRepo() error {
	// Checks if Repo is already initialized or not
	if isPathExist(RepoDir) {
		return fmt.Errorf("repository already initialized")
	}

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
	fmt.Printf("%sUsing 'main' as the name for the default branch.%s\n\n", hint, hint)
	fmt.Printf("%sBranches can be renamed via this command:%s\n", hint, hint)
	fmt.Printf("%s\tkitkat branch -m <branch_name>%s\n\n", hint, hint)
	fmt.Printf("%sList all the branches via this command:%s\n", hint, hint)
	fmt.Printf("%s\tkitkat branch -l%s\n", hint, hint)
	// Generating empty main branch file.
	if err := os.WriteFile(HeadsDir+"/main", []byte(""), 0o644); err != nil {
		return err
	}

	// Create default .kitignore to prevent self-tracking
	ignoreContent := []byte(".DS_Store\nkitkat\nkitkat.exe\n")
	if err := os.WriteFile(".kitignore", ignoreContent, 0644); err != nil {
		return err
	}

	if absPath, err := filepath.Abs(RepoDir); err != nil {
		return err
	} else {
		fmt.Printf("%s\nInitialized empty kitkat repository in %s\n\n%s", hint, absPath, hint)
	}
	return nil
}
