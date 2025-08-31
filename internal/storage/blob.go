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
	
	objPath := filepath.Join(objectsDir, hash)
	// ensure objects dir exists
	if err := os.MkdirAll(objectsDir, 0755); err != nil {
	    return "", err
	}
	
	if _, err := os.Stat(objPath); os.IsNotExist(err) {
	    // write via tmp file — reuse already-open file by rewinding
	    tmp := objPath + ".tmp"
	    if _, err := f.Seek(0, 0); err != nil {
	        return "", err
	    }
	    out, err := os.Create(tmp)
	    if err != nil {
	        return "", err
	    }
	    if _, err := io.Copy(out, f); err != nil {
	        out.Close()
	        os.Remove(tmp)
	        return "", err
	    }
	    out.Close()
	    if err := os.Rename(tmp, objPath); err != nil {
	        os.Remove(tmp)
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