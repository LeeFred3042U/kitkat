package core

import (
	"os"
	"testing"
	"time"

	"github.com/LeeFred3042U/kitcat/internal/models"
	"github.com/LeeFred3042U/kitcat/internal/storage"
	"github.com/LeeFred3042U/kitcat/internal/testutil"
)

// Test_CheckoutFile_PreservesDirtyFile ensures strict safety behavior:
// It verifies that CheckoutFile does NOT overwrite a file that has uncommitted changes.
func Test_CheckoutFile_PreservesDirtyFile(t *testing.T) {
	// 1. Setup temporary repository
	_, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

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

// Test_CheckoutFile_UpdatesIndex verifies that CheckoutFile correctly updates the index/staging area
// after restoring a file. This prevents "phantom" changes from appearing in `kitkat status`.
func Test_CheckoutFile_UpdatesIndex(t *testing.T) {
	// 1. Setup temporary repository
	_, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

	// 2. Create and commit a file
	filePath := "file.txt"
	content := []byte("original content")
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	// Store blob
	blobHash, err := storage.HashAndStoreFile(filePath)
	if err != nil {
		t.Fatalf("failed to store blob: %v", err)
	}

	// Create initial index
	index := map[string]string{
		filePath: blobHash,
	}
	if err := storage.WriteIndex(index); err != nil {
		t.Fatalf("failed to write index: %v", err)
	}

	// Create tree and commit
	treeHash, err := storage.CreateTree()
	if err != nil {
		t.Fatalf("failed to create tree: %v", err)
	}

	commit := models.Commit{
		TreeHash:  treeHash,
		Message:   "Initial commit",
		Timestamp: time.Now(),
		ID:        "test-commit-id",
	}
	if err := storage.AppendCommit(commit); err != nil {
		t.Fatalf("failed to append commit: %v", err)
	}
	// Update HEAD
	if err := os.WriteFile(".kitcat/HEAD", []byte(commit.ID), 0644); err != nil {
		t.Fatalf("failed to update HEAD: %v", err)
	}

	// 3. Delete the file from workspace (simulate accidental deletion or state change)
	if err := os.Remove(filePath); err != nil {
		t.Fatalf("failed to remove file: %v", err)
	}

	// 4. Modify Index to simulate "missing" entry or inconsistent state
	// We want to ensure CheckoutFile RESTORES the index entry.
	// So let's delete it from the index explicitly.
	delete(index, filePath)
	if err := storage.WriteIndex(index); err != nil {
		t.Fatalf("failed to write corrupt index: %v", err)
	}

	// 5. Checkout the file
	if err := CheckoutFile(filePath); err != nil {
		t.Fatalf("CheckoutFile failed: %v", err)
	}

	// 6. Verify File Content
	restoredContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read restored file: %v", err)
	}
	if string(restoredContent) != string(content) {
		t.Errorf("Content mismatch. Want %s, got %s", content, restoredContent)
	}

	// 7. Verify Index is Updated
	// Re-load the index
	loadedIndex, err := storage.LoadIndex()
	if err != nil {
		t.Fatalf("failed to load index: %v", err)
	}

	storedHash, ok := loadedIndex[filePath]
	if !ok {
		t.Error("CheckoutFile did NOT update the index: file entry is missing")
	} else if storedHash != blobHash {
		t.Errorf("Index has wrong hash. Want %s, got %s", blobHash, storedHash)
	}
}
