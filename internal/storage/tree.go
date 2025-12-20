package storage

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// CreateTree creates a tree object from the current index and stores it
// It ensures the process is deterministic by sorting the file paths
func CreateTree() (string, error) {
	index, err := LoadIndex()
	if err != nil {
		return "", err
	}

	var treeContent bytes.Buffer

	// Sort keys to ensure the tree content is always in the same order
	keys := make([]string, 0, len(index))
	for p := range index {
		keys = append(keys, p)
	}
	sort.Strings(keys)

	// Iterate over the sorted keys to build the tree content
	for _, path := range keys {
		hash := index[path]
		treeContent.WriteString(fmt.Sprintf("%s %s\n", hash, path))
	}

	// Hash the deterministic tree content to get the tree's hash
	h := sha1.New()
	h.Write(treeContent.Bytes())
	treeHash := fmt.Sprintf("%x", h.Sum(nil))

	// Store the tree object in the objects directory
	objectPath := filepath.Join(objectsDir, treeHash)
	if err := os.WriteFile(objectPath, treeContent.Bytes(), 0644); err != nil {
		return "", err
	}

	return treeHash, nil
}

// ParseTree reads a tree object from storage and returns it as a map of path -> hash
// This is the function that was missing
func ParseTree(hash string) (map[string]string, error) {
	tree := make(map[string]string)
	objectPath := filepath.Join(objectsDir, hash)
	data, err := os.ReadFile(objectPath)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		// The format is "hash path", so we split on the first space
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 {
			// The hash is parts[0], the path is parts[1]
			tree[parts[1]] = parts[0]
		}
	}

	return tree, scanner.Err()
}
