package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const indexPath = ".kitkat/index"

// LoadIndex reads the .kitkat/index file (in JSON format) and returns it as a map
// It returns an empty map if the file doesn't exist, which is normal for a new repository
func LoadIndex() (map[string]string, error) {
	index := make(map[string]string)

	content, err := os.ReadFile(indexPath)
	if os.IsNotExist(err) {
		// File doesn't exist, return empty index. This is not an error ^-^
		return index, nil
	}
	if err != nil {
		return nil, fmt.Errorf("could not read index file: %w", err)
	}
	
	// If the file is empty, avoid a JSON error.
	if len(content) == 0 {
		return index, nil
	}

	if err := json.Unmarshal(content, &index); err != nil {
		return nil, fmt.Errorf("could not parse index file: %w", err)
	}

	return index, nil
}

// WriteIndex writes the index map to the .kitkat/index file atomically using a JSON format
// It uses a temporary file and an atomic rename to prevent corruption 'o'
func WriteIndex(index map[string]string) error {
	// Ensure the parent directory (.kitkat) exists.
	if err := os.MkdirAll(filepath.Dir(indexPath), 0755); err != nil {
		return err
	}

	// Lock the file to prevent concurrent writes
	l, err := lock(indexPath)
	if err != nil {
		return err
	}
	defer unlock(l)

	// Use a temporary file for the initial write
	tmpPath := indexPath + ".tmp"
	file, err := os.Create(tmpPath)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	
	err = encoder.Encode(index)
	// Must close the file before renaming it
	if closeErr := file.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}

	// Atomically rename the temporary file to the final index file
	return os.Rename(tmpPath, indexPath)
}