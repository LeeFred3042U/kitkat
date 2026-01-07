package core

import (
	"encoding/json"
	"fmt"
	"os"
)

// IndexEntry represents a file in the staging area
type IndexEntry struct {
	Path string
	Hash string
}

// LoadIndex reads the .kitkat/index file
func LoadIndex() ([]IndexEntry, error) {
	data, err := os.ReadFile(IndexPath)
	if os.IsNotExist(err) {
		return []IndexEntry{}, nil
	}
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return []IndexEntry{}, nil
	}

	var entryMap map[string]string
	err = json.Unmarshal(data, &entryMap)
	if err != nil {
		return nil, fmt.Errorf("index file corrupted")
	}

	var entries []IndexEntry
	for key, value := range entryMap {
		entries = append(entries, IndexEntry{Path: key, Hash: value})
	}
	return entries, nil
}

// SaveIndex writes the index back to disk
func SaveIndex(entries []IndexEntry) error {
	file, err := os.Create(IndexPath)
	if err != nil {
		return err
	}
	defer file.Close()

	entryMap := make(map[string]string)
	for _, entry := range entries {
		entryMap[entry.Path] = entry.Hash
	}

	data, err := json.Marshal(entryMap)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("unable to write to index")
	}
	return nil
}
