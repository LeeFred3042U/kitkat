package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/LeeFred3042U/kitkat/internal/storage"
)

// Removes untracked files from the working directory
func Clean() error {
	index, err := storage.LoadIndex()
	if err != nil {
		return err
	}

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == ".kitkat" {
			return filepath.SkipDir
		}

		if _, tracked := index[path]; !tracked && path != "." && !info.IsDir() {
			fmt.Printf("Removing %s\n", path)
			return os.Remove(path)
		}

		return nil
	})

	return err
}