package core

import (
	"fmt"
	"strings"

	"github.com/LeeFred3042U/kitkat/internal/storage"
)

// Prints a simple line-by-line diff of two strings
func showDiff(oldContent, newContent string) {
	oldLines := strings.Split(oldContent, "\n")
	newLines := strings.Split(newContent, "\n")

	// This is a very basic diff implementation 
	// A real implementation would use
	// a more sophisticated algorithm like Myers diff
	for _, line := range oldLines {
		if !contains(newLines, line) {
			fmt.Printf("- %s\n", line)
		}
	}
	for _, line := range newLines {
		if !contains(oldLines, line) {
			fmt.Printf("+ %s\n", line)
		}
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Shows the differences between the index and the last commit
func Diff() error {
	lastCommit, err := storage.GetLastCommit()
	if err != nil {
		return err
	}

	tree, err := storage.ParseTree(lastCommit.TreeHash)
	if err != nil {
		return err
	}

	index, err := storage.LoadIndex()
	if err != nil {
		return err
	}

	for path, indexHash := range index {
		treeHash, ok := tree[path]
		if !ok {
			fmt.Printf("Added file: %s\n", path)
			continue
		}
		if indexHash != treeHash {
			fmt.Printf("Modified file: %s\n", path)
			oldContent, err := storage.ReadObject(treeHash)
			if err != nil {
				return err
			}
			newContent, err := storage.ReadObject(indexHash)
			if err != nil {
				return err
			}
			showDiff(string(oldContent), string(newContent))
		}
	}

	for path := range tree {
		if _, ok := index[path]; !ok {
			fmt.Printf("Deleted file: %s\n", path)
		}
	}

	return nil
}