package core

import (
	"errors"
	"os"

	"github.com/LeeFred3042U/kitkat/internal/storage"
)

func AddFile(path string) error {
	// Guard: ensure we're inside a kitkat repo
	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		return errors.New("not a kitkat repository (run `kitkat init`)")
	}

	hash, err := storage.HashAndStoreFile(path)
	if err != nil {
		return err
	}

	index, err := storage.LoadIndex()
	if err != nil {
		return err
	}

	// Skip if already tracked with same hash
	if existing, ok := index[path]; ok && existing == hash {
		return nil
	}

	index[path] = hash
	return storage.WriteIndex(index)
}
