package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/LeeFred3042U/kitcat/internal/storage"
)

// Removes untracked files from the working directory
// If includeIgnored is false, ignored files are preserved
// If includeIgnored is true, ignored files are also removed
func Clean(dryRun bool, includeIgnored bool) error {
	// Guard: ensure we're inside a kitkat repo
	if _, err := os.Stat(RepoDir); os.IsNotExist(err) {
		return errors.New("not a kitkat repository (run `kitkat init`)")
	}

	index, err := storage.LoadIndex()
	if err != nil {
		return err
	}

	// Load ignore patterns
	ignorePatterns, err := LoadIgnorePatterns()
	if err != nil {
		return err
	}

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		clean := filepath.Clean(path)

		// skip the repo dir and everything under it
		if clean == RepoDir || strings.HasPrefix(clean, RepoDir+string(os.PathSeparator)) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// skip directories and the root marker "."
		if info.IsDir() || clean == "." {
			return nil
		}

		// if not tracked, remove (or print if dry run)
		if _, tracked := index[clean]; !tracked {
			// Check if file is ignored
			isIgnored := ShouldIgnore(clean, ignorePatterns, index)

			// Skip ignored files unless -x flag is set
			if isIgnored && !includeIgnored {
				return nil
			}

			if dryRun {
				if isIgnored {
					fmt.Printf("Would remove (ignored) %s\n", clean)
				} else {
					fmt.Printf("Would remove %s\n", clean)
				}
				return nil
			}
			fmt.Printf("Removing %s\n", clean)
			return os.Remove(clean)
		}
		return nil
	})
	return err
}
