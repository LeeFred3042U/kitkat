package storage

import (
	"os"
	"fmt"
	"bufio"
	"bytes"
	"strings"
	"crypto/sha1"
	"path/filepath"
)
// Creates a tree object from the current index and stores it
func CreateTree() (string, error) {
	index, err := LoadIndex()
	if err != nil {
		return "", err
	}

	var treeContent bytes.Buffer
	for path, hash := range index {
		treeContent.WriteString(fmt.Sprintf("%s %s\n", hash, path))
	}

	// Hash the tree content to get the tree hash
	h := sha1.New()
	h.Write(treeContent.Bytes())
	treeHash := fmt.Sprintf("%x", h.Sum(nil))

	// Store the tree object
	objectPath := filepath.Join(objectsDir, treeHash)
	if err := os.WriteFile(objectPath, treeContent.Bytes(), 0644); err != nil {
		return "", err
	}

	return treeHash, nil
}

// Parses a tree object and returns a map of path -> hash
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
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 {
			tree[parts[1]] = parts[0]
		}
	}

	return tree, scanner.Err()
}