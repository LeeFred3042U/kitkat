package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestMoveFile(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Change working directory into temp repo
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Initialize kitkat repository
	if err := InitRepo(); err != nil {
		t.Fatal(err)
	}

	oldPath := "old_test.txt"
	newPath := "new_test.txt"

	// Create old file
	if err := os.WriteFile(oldPath, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	// Stage old file
	if err := AddFile(oldPath); err != nil {
		t.Fatal(err)
	}

	// Move file
	if err := MoveFile(oldPath, newPath); err != nil {
		t.Fatal(err)
	}

	// Old file should be gone
	if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
		t.Fatalf("expected old file to be removed")
	}

	// New file should exist
	if _, err := os.Stat(newPath); err != nil {
		t.Fatalf("expected new file to exist")
	}

	// Index should contain new file
	idx, err := loadIndexForTest(filepath.Join(tmpDir, ".kitkat", "index"))
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := idx[newPath]; !ok {
		t.Fatalf("expected new file to be staged")
	}
}

func loadIndexForTest(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var idx map[string]string
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, err
	}

	return idx, nil
}
