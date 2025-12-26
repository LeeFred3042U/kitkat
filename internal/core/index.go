package core

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// IndexEntry represents a file in the staging area
type IndexEntry struct {
	Path string
	Hash string
}

// LoadIndex reads the .kitkat/index file
func LoadIndex() ([]IndexEntry, error) {
	file, err := os.Open(".kitkat/index")
	if os.IsNotExist(err) {
		return []IndexEntry{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []IndexEntry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// We expect: HASH PATH
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 {
			entries = append(entries, IndexEntry{Hash: parts[0], Path: parts[1]})
		}
	}
	return entries, scanner.Err()
}

// SaveIndex writes the index back to disk
func SaveIndex(entries []IndexEntry) error {
	file, err := os.Create(".kitkat/index")
	if err != nil {
		return err
	}
	defer file.Close()

	for _, entry := range entries {
		_, err := fmt.Fprintf(file, "%s %s\n", entry.Hash, entry.Path)
		if err != nil {
			return err
		}
	}
	return nil
}
