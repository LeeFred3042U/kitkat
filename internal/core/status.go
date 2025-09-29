package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/LeeFred3042U/kitkat/internal/storage"
)

// Status compares the state of the working directory, index, and last commit,
// then prints a summary of the changes
func Status() error {
	// Print the current branch status at the top
	headState, err := GetHeadState()
	if err != nil {
		headState = "no commits yet"
	}
	fmt.Printf("On branch %s\n", headState)

	// Load the tree from the last commit (HEAD)
	headTree := make(map[string]string)
	lastCommit, err := storage.GetLastCommit()
	if err == nil { // A commit exists.
		tree, parseErr := storage.ParseTree(lastCommit.TreeHash)
		if parseErr != nil {
			return parseErr
		}
		headTree = tree
	} else if err != storage.ErrNoCommits {
		return err
	}

	// Load the current staging area
	index, err := storage.LoadIndex()
	if err != nil {
		return err
	}

	// Prepare slices to hold the categorized changes
	stagedChanges := []string{}
	unstagedChanges := []string{}
	untrackedFiles := []string{}

	// Create a set of all file paths from both HEAD and the index for a complete comparison
	allPaths := make(map[string]bool)
	for path := range headTree {
		allPaths[path] = true
	}
	for path := range index {
		allPaths[path] = true
	}

	// Categorize Staged Changes (Index vs. HEAD)
	for path := range allPaths {
		headHash, inHead := headTree[path]
		indexHash, inIndex := index[path]

		if inIndex && !inHead {
			stagedChanges = append(stagedChanges, fmt.Sprintf("new file:  %s", path))
		} else if !inIndex && inHead {
			stagedChanges = append(stagedChanges, fmt.Sprintf("deleted:   %s", path))
		} else if inIndex && inHead && headHash != indexHash {
			stagedChanges = append(stagedChanges, fmt.Sprintf("modified:  %s", path))
		}
	}

	// Categorize Unstaged & Untracked Changes (Working Directory vs. Index)
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		cleanPath := filepath.Clean(path)
		
		// Skip the .kitkat directory and other directories
		if info.IsDir() || strings.HasPrefix(cleanPath, RepoDir+string(os.PathSeparator)) || cleanPath == RepoDir {
			return nil
		}

		indexHash, isTracked := index[cleanPath]

		// If the file is not in the index, it's untracked
		if !isTracked {
			untrackedFiles = append(untrackedFiles, cleanPath)
			return nil
		}

		// If the file is tracked, hash it and compare with the index to see if it's been modified
		currentHash, hashErr := storage.HashFile(cleanPath)
		if hashErr != nil {
			return hashErr
		}
		if currentHash != indexHash {
			unstagedChanges = append(unstagedChanges, fmt.Sprintf("modified:  %s", cleanPath))
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Print Final Summary
	fmt.Println("\nChanges to be committed:")
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