package core

import (
	"os"
	"fmt"
	"strings"
	"path/filepath"

	"github.com/LeeFred3042U/kitkat/internal/storage"
)

// Compare the state of the working directory, index, and last commit
func Status() error {
	// Get HEAD commit tree
	headTree := make(map[string]string)
	lastCommit, err := storage.GetLastCommit()
	if err == nil { 
		// VCheck If there is a commit
		tree, parseErr := storage.ParseTree(lastCommit.TreeHash)
		if parseErr != nil {
			return parseErr
		}
		headTree = tree
	} else if err != storage.ErrNoCommits {
		return err
	}

	// Get index
	index, err := storage.LoadIndex()
	if err != nil {
		return err
	}

	// Find staged, unstaged, and untracked files
	stagedChanges := []string{}
	unstagedChanges := []string{}
	untrackedFiles := []string{}

	// Create a set of all known paths for easier iteration
	allPaths := make(map[string]bool)
	for path := range headTree {
		allPaths[path] = true
	}
	for path := range index {
		allPaths[path] = true
	}

	// Compare index vs HEAD for staged changes
	for path := range allPaths {
		headHash := headTree[path]
		indexHash := index[path]
		if headHash != indexHash {
			stagedChanges = append(stagedChanges, fmt.Sprintf("modified: %s", path))
		}
	}

	// Compare working dir vs index for unstaged changes and find untracked files
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || strings.HasPrefix(path, ".kitkat") {
			return nil // Skip directories and the repo itself
		}

		indexHash, isTracked := index[path]
		if !isTracked {
			untrackedFiles = append(untrackedFiles, path)
			return nil
		}

		// It's tracked, so check if it's been modified since staging
		currentHash, hashErr := storage.HashFile(path)
		if hashErr != nil {
			return hashErr
		}
		if currentHash != indexHash {
			unstagedChanges = append(unstagedChanges, fmt.Sprintf("modified: %s", path))
		}
		return nil
	})
	if err != nil {
		return err
	}

	fmt.Println("Changes to be committed:")
	for _, change := range stagedChanges {
		fmt.Printf("\t%s\n", change)
	}
	fmt.Println("\nChanges not staged for commit:")
	for _, change := range unstagedChanges {
		fmt.Printf("\t%s\n", change)
	}
	fmt.Println("\nUntracked files:")
	for _, file := range untrackedFiles {
		fmt.Printf("\t%s\n", file)
	}

	return nil
}