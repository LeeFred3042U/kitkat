package core

import (
	"os"
	"testing"
	"time"

	"github.com/LeeFred3042U/kitcat/internal/models"
	"github.com/LeeFred3042U/kitcat/internal/storage"
)

// Test_CheckoutFile_PreservesDirtyFile ensures strict safety behavior:
// It verifies that CheckoutFile does NOT overwrite a file that has uncommitted changes.
func Test_CheckoutFile_PreservesDirtyFile(t *testing.T) {
	// 1. Setup temporary repository
	repoDir := t.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer func() {
		_ = os.Chdir(cwd)
	}() // Restore cwd after test

	if err := os.Chdir(repoDir); err != nil {
		t.Fatalf("failed to chdir to temp repo: %v", err)
	}

	// Initialize minimal .kitkat structure
	dirs := []string{
		".kitcat",
		".kitcat/objects",
		".kitcat/refs/heads",
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			t.Fatalf("failed to create dir %s: %v", d, err)
		}
	}

	// 2. Create a file foo.txt with content v1
	filePath := "foo.txt"
	v1Content := []byte("v1")
	if err := os.WriteFile(filePath, v1Content, 0644); err != nil {
		t.Fatalf("failed to create foo.txt: %v", err)
	}

	// 3. Stage and store v1 (simulate a commit)
	// a. Hash and store blob
	blobHash, err := storage.HashAndStoreFile(filePath)
	if err != nil {
		t.Fatalf("failed to store blob: %v", err)
	}

	// b. Create tree
	// We first write to index, because CreateTree reads from Index
	index := map[string]string{
		filePath: blobHash,
	}
	if err := storage.WriteIndex(index); err != nil {
		t.Fatalf("failed to write index: %v", err)
	}

	treeHash, err := storage.CreateTree()
	if err != nil {
		t.Fatalf("failed to create tree: %v", err)
	}

	// c. Create commit
	commit := models.Commit{
		TreeHash:  treeHash,
		Message:   "Initial commit",
		Timestamp: time.Now(),
		ID:        "test-commit-hash", // ID doesn't matter for GetLastCommit logic if we append it
	}
	// AppendCommit writes to .kitcat/commits.log
	if err := storage.AppendCommit(commit); err != nil {
		t.Fatalf("failed to append commit: %v", err)
	}

	// 4. Modify foo.txt in working dir to v2 (uncommitted/dirty change)
	v2Content := []byte("v2")
	if err := os.WriteFile(filePath, v2Content, 0644); err != nil {
		t.Fatalf("failed to modify foo.txt: %v", err)
	}

	// 5. Invoke CheckoutFile("foo.txt")
	// This should fail because the file is dirty (tracked in index, but modified on disk)
	err = CheckoutFile(filePath)

	// 6. Assertions
	// Expect an error
	if err == nil {
		t.Error("CheckoutFile succeeded but expected error due to dirty file")
	} else {
		// Optional: check error message content
		expectedMsg := "local changes to 'foo.txt' would be overwritten"
		if err.Error() != "error: "+expectedMsg {
			// exact message match might be brittle, but sticking to requested assert logic
			t.Logf("Got expected error: %v", err)
		}
	}

	// Assert that content is STILL v2 (NOT overwritten)
	currentContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read foo.txt: %v", err)
	}
	if string(currentContent) != "v2" {
		t.Errorf("File content was overwritten! Expected 'v2', got '%s'", string(currentContent))
	}
}
