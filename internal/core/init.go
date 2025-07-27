package core

import (
	"os"
)

const (
	repoDir    = ".kitkat"
	indexPath  = ".kitkat/index"
	objectsDir = ".kitkat/objects"
)

// InitRepo sets up the .kitkat structure: required dirs and files.
// Safe to re-run; won't overwrite unless files are missing.
func InitRepo() error {
	err := os.Mkdir(".kitkat", 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	// objectsDir stores hashed file contents (content-addressed)
	dirs := []string{".kitkat/objects"}
	for _, dir := range dirs {
		if err := os.Mkdir(dir, 0755); err != nil && !os.IsExist(err) {
			return err
		}
	}

	// index: file path -> hash map
	// commits.log: future use for log entries or snapshots
	files := []string{".kitkat/index", ".kitkat/commits.log"}
	for _, file := range files {
		f, err := os.Create(file)
		if err != nil {
			return err
		}
		f.Close()
	}

	return nil
}