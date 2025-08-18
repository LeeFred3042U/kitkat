package core

import (
	"os"
	"fmt"
	"strings"
	"path/filepath"

	"github.com/LeeFred3042U/kitkat/internal/storage"
)

// Removes untracked files from the working directory
func Clean(dryRun bool) error {
	index, err := storage.LoadIndex()
	if err != nil {
		return err
	}

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		clean := filepath.Clean(path)

		// skip the repo dir and everything under it
		if clean == repoDir || strings.HasPrefix(clean, repoDir+string(os.PathSeparator)) {
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
			if dryRun {
				fmt.Printf("Would remove %s\n", clean)
				return nil
			}
			fmt.Printf("Removing %s\n", clean)
			return os.Remove(clean)
		}
		return nil
	})
	return err
}
