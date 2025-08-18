package storage

import (
	"io"
	"os"
	"fmt"
	"crypto/sha1"
	"path/filepath"
)

const (
	objectsDir = ".kitkat/objects"
)

func HashAndStoreFile(path string) (string, error) {
    content, err := os.ReadFile(path) // Reads the file once
    if err != nil {
        return "", err
    }

    h := sha1.New()
    h.Write(content) // Hashes the content from memory
    hash := fmt.Sprintf("%x", h.Sum(nil))

    objectPath := filepath.Join(objectsDir, hash)
    if _, err := os.Stat(objectPath); os.IsNotExist(err) {
        if err := os.WriteFile(objectPath, content, 0644); err != nil { 
			// Writes the content from memory
            return "", err
        }
    }
    return hash, nil
}

// Reads an object from the objects directory
func ReadObject(hash string) ([]byte, error) {
	objectPath := filepath.Join(objectsDir, hash)
	return os.ReadFile(objectPath)
}

// Computes the SHA-1 hash of a file's content
// does not store the file in the object database
func HashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha1.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}