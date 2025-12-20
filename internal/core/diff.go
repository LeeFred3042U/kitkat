package core

import (
	"fmt"
	"strings"

	"github.com/LeeFred3042U/kitkat/internal/diff"
	"github.com/LeeFred3042U/kitkat/internal/storage"
)

// ANSI color codes for formatting terminal output
const (
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorBlue  = "\033[1;34m"
)

// displayDiff formats and prints the structured diff output from the Myers algorithm.
// It iterates through each change (insertion, deletion, or equal) and applies the appropriate color
func displayDiff(diffs []diff.Diff[string]) {
	// Loop over each diff "chunk" provided by the algorithm
	for _, d := range diffs {
		lines := d.Text
		switch d.Operation {
		// If the operation was an INSERT, print each line in green with a '+' prefix
		case diff.INSERT:
			for _, line := range lines {
				fmt.Printf("%s+ %s%s\n", colorGreen, line, colorReset)
			}
		// If the operation was a DELETE, print each line in red with a '-' prefix
		case diff.DELETE:
			for _, line := range lines {
				fmt.Printf("%s- %s%s\n", colorRed, line, colorReset)
			}
		// If the lines are EQUAL, print them with the default color and two spaces for context
		case diff.EQUAL:
			for _, line := range lines {
				fmt.Printf("  %s\n", line)
			}
		}
	}
}

// Diff calculates and displays the differences between the last commit and the current staging area (index)
// It identifies which files have been added, deleted, or modified.
func Diff() error {
	// Retrieve the metadata for the most recent commit.
	lastCommit, err := storage.GetLastCommit()
	if err != nil {
		// If there are no commits yet, there's nothing to compare against.
		if err == storage.ErrNoCommits {
			fmt.Println("No commits yet. Nothing to diff against.")
			return nil
		}
		return err
	}

	// From the commit, get the tree object which represents the state of the repository at that time
	// This is a map of `filePath -> contentHash`
	tree, err := storage.ParseTree(lastCommit.TreeHash)
	if err != nil {
		return err
	}

	// Load the current staging area into a map. This represents what will be in the *next* commit
	index, err := storage.LoadIndex()
	if err != nil {
		return err
	}

	// First Loop: Iterate through files in the index to find additions and modifications
	for path, indexHash := range index {
		treeHash, ok := tree[path]
		// If a file is in the index but not in the old tree, it's a new file.
		if !ok {
			fmt.Printf("%sAdded file: %s%s\n", colorBlue, path, colorReset)
			// We could optionally show the full content of the new file here.
			continue
		}

		// If the file exists in both, but the content hash is different, it has been modified
		if indexHash != treeHash {
			fmt.Printf("%sModified file: %s%s\n", colorBlue, path, colorReset)

			// Read the old and new content from the object store.
			oldContent, err := storage.ReadObject(treeHash)
			if err != nil {
				return err
			}
			newContent, err := storage.ReadObject(indexHash)
			if err != nil {
				return err
			}

			// Split file content into lines to prepare for the diff algorithm
			oldLines := strings.Split(string(oldContent), "\n")
			newLines := strings.Split(string(newContent), "\n")

			// Using the Myers algorithm to compute the differences
			myers := diff.NewMyersDiff(oldLines, newLines)
			diffs := myers.Diffs()

			// Display the computed differences with color
			displayDiff(diffs)
		}
	}

	// Next Loop: Iterate through files in the old tree to find deletions.
	for path := range tree {
		// If a file was in the old tree but is no longer in the index, it has been deleted.
		if _, ok := index[path]; !ok {
			fmt.Printf("%sDeleted file: %s%s\n", colorBlue, path, colorReset)
		}
	}

	return nil
}
