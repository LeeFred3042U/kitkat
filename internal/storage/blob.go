package storage

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	objectsDir = ".kitkat/objects"
)

// HashAndStoreFile computes SHA-1 hash of file content,
//  saves it in .kitkat/objects/<hash> if not present
// Returns a hash string
func HashAndStoreFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	hash := fmt.Sprintf("%x", h.Sum(nil))

	// Reset file pointer before second read — needed since hashing moved it to EOF
	f.Seek(0, 0)
	content, _ := io.ReadAll(f)

	objectPath := filepath.Join(objectsDir, hash)

	// Only write if object doesn't exist
	if _, err := os.Stat(objectPath); os.IsNotExist(err) {
		err = os.WriteFile(objectPath, content, 0644)
		if err != nil {
			return "", err
		}
	}

	return hash, nil
}