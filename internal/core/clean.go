package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/LeeFred3042U/kitcat/internal/storage"
)

// Removes untracked files from the working directory
// If includeIgnored is false, ignored files are preserved
// If includeIgnored is true, ignored files are also removed
func Clean(dryRun bool, includeIgnored bool) error {
	// Guard: ensure we're inside a kitcat repo
	if _, err := os.Stat(RepoDir); os.IsNotExist(err) {
		return errors.New("not a kitcat repository (run `kitcat init`)")
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

	var visitedDirs []string

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

		// Track directories (except root)
		if info.IsDir() && clean != "." {
			visitedDirs = append(visitedDirs, clean)
			return nil
		}

		// skip root marker "." and directories (already handled above)
		if clean == "." || info.IsDir() {
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
	if err != nil {
		return err
	}

	// Post-process: remove empty directories (deepest first)
	if !dryRun {
		sort.Sort(sort.Reverse(sort.StringSlice(visitedDirs)))
		for _, dir := range visitedDirs {
			if err := os.Remove(dir); err != nil {
				// Directory still contains tracked files or files we chose not to remove
				continue // ← Fix: explicit continue instead of empty block
			}
		}
	} else {
		// In dry-run mode, show which directories would be removed
		sort.Sort(sort.Reverse(sort.StringSlice(visitedDirs)))
		for _, dir := range visitedDirs {
			fmt.Printf("Would remove directory %s\n", dir)
		}
	}

	return nil
}
