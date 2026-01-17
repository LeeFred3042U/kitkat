package core_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/LeeFred3042U/kitcat/internal/core"
	"github.com/LeeFred3042U/kitcat/internal/storage"
	"github.com/LeeFred3042U/kitcat/internal/testutil"
)

func readStashStack(t *testing.T, stashRefPath string) []string {
	t.Helper()
	data, err := os.ReadFile(stashRefPath)
	if os.IsNotExist(err) {
		return []string{}
	}
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	var out []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			out = append(out, l)
		}
	}
	return out
}

// TestStash_BasicWorkflow tests the basic stash save workflow
func TestStash_BasicWorkflow(t *testing.T) {
	tmpDir, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

	// Create and commit initial file
	testFile := "test.txt"
	if err := os.WriteFile(testFile, []byte("initial content"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(testFile); err != nil {
		t.Fatal(err)
	}
	if _, _, err := core.Commit("initial commit"); err != nil {
		t.Fatal(err)
	}

	// Modify the file
	if err := os.WriteFile(testFile, []byte("modified content"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(testFile); err != nil {
		t.Fatal(err)
	}

	// Stash the changes
	if err := core.Stash(); err != nil {
		t.Fatalf("Stash failed: %v", err)
	}

	// Verify working directory is clean
	isDirty, err := core.IsWorkDirDirty()
	if err != nil {
		t.Fatal(err)
	}
	if isDirty {
		t.Error("Working directory should be clean after stash")
	}

	// Verify stash reference exists
	stashRefPath := filepath.Join(tmpDir, ".kitcat", "refs", "stash")
	stack := readStashStack(t, stashRefPath)
	if len(stack) == 0 {
		t.Fatal("Stash reference should exist and contain at least one entry")
	}
	stashHash := stack[0]

	// Verify stash commit exists in commits.log
	stashCommit, err := storage.FindCommit(stashHash)
	if err != nil {
		t.Fatalf("Stash commit should exist: %v", err)
	}

	// Verify WIP message format
	if !strings.HasPrefix(stashCommit.Message, "WIP on ") {
		t.Errorf("Stash commit message should start with 'WIP on', got: %s", stashCommit.Message)
	}
}

// TestStash_CleanWorkingDirectory tests stashing with a clean working directory
func TestStash_CleanWorkingDirectory(t *testing.T) {
	_, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

	// Create and commit initial file
	testFile := "test.txt"
	if err := os.WriteFile(testFile, []byte("content"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(testFile); err != nil {
		t.Fatal(err)
	}
	if _, _, err := core.Commit("initial commit"); err != nil {
		t.Fatal(err)
	}

	// Try to stash with clean directory
	err := core.Stash()
	if err == nil {
		t.Fatal("Stash should fail with clean working directory")
	}
	if !strings.Contains(err.Error(), "nothing to stash") {
		t.Errorf("Expected 'nothing to stash' error, got: %v", err)
	}
}

// TestStash_StagedAndUnstagedChanges tests stashing with mixed changes
func TestStash_StagedAndUnstagedChanges(t *testing.T) {
	tmpDir, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

	// Create and commit initial files
	file1 := "file1.txt"
	file2 := "file2.txt"
	if err := os.WriteFile(file1, []byte("content1"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte("content2"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(file1); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(file2); err != nil {
		t.Fatal(err)
	}
	if _, _, err := core.Commit("initial commit"); err != nil {
		t.Fatal(err)
	}

	// Create staged changes
	if err := os.WriteFile(file1, []byte("staged content"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(file1); err != nil {
		t.Fatal(err)
	}

	// Create unstaged changes
	if err := os.WriteFile(file2, []byte("unstaged content"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Stash all changes
	if err := core.Stash(); err != nil {
		t.Fatalf("Stash failed: %v", err)
	}

	// Verify working directory matches HEAD
	content1, err := os.ReadFile(file1)
	if err != nil {
		t.Fatal(err)
	}
	if string(content1) != "content1" {
		t.Errorf("file1 should be reset to HEAD content, got: %s", string(content1))
	}

	content2, err := os.ReadFile(file2)
	if err != nil {
		t.Fatal(err)
	}
	if string(content2) != "content2" {
		t.Errorf("file2 should be reset to HEAD content, got: %s", string(content2))
	}

	// Verify stash reference exists
	stashRefPath := filepath.Join(tmpDir, ".kitcat", "refs", "stash")
	if stack := readStashStack(t, stashRefPath); len(stack) == 0 {
		t.Error("Stash reference should exist")
	}
}

// TestStashPop_Success tests successful stash pop
func TestStashPop_Success(t *testing.T) {
	tmpDir, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

	// Create and commit initial file
	testFile := "test.txt"
	originalContent := "original content"
	modifiedContent := "modified content"

	if err := os.WriteFile(testFile, []byte(originalContent), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(testFile); err != nil {
		t.Fatal(err)
	}
	if _, _, err := core.Commit("initial commit"); err != nil {
		t.Fatal(err)
	}

	// Modify and stash
	if err := os.WriteFile(testFile, []byte(modifiedContent), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(testFile); err != nil {
		t.Fatal(err)
	}
	if err := core.Stash(); err != nil {
		t.Fatal(err)
	}

	// Verify file is reset
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != originalContent {
		t.Errorf("File should be reset after stash, got: %s", string(content))
	}

	// Pop the stash
	if err := core.StashPop(); err != nil {
		t.Fatalf("StashPop failed: %v", err)
	}

	// Verify changes are restored
	content, err = os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != modifiedContent {
		t.Errorf("File should be restored after pop, got: %s, want: %s", string(content), modifiedContent)
	}

	// Verify stash reference is deleted
	stashRefPath := filepath.Join(tmpDir, ".kitcat", "refs", "stash")
	if stack := readStashStack(t, stashRefPath); len(stack) != 0 {
		t.Error("Stash reference should be deleted after pop")
	}
}

// TestStashPop_NoStash tests popping when no stash exists
func TestStashPop_NoStash(t *testing.T) {
	_, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

	// Create and commit initial file
	testFile := "test.txt"
	if err := os.WriteFile(testFile, []byte("content"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(testFile); err != nil {
		t.Fatal(err)
	}
	if _, _, err := core.Commit("initial commit"); err != nil {
		t.Fatal(err)
	}

	// Try to pop without stash
	err := core.StashPop()
	if err == nil {
		t.Fatal("StashPop should fail when no stash exists")
	}
	if !strings.Contains(err.Error(), "no stash found") {
		t.Errorf("Expected 'no stash found' error, got: %v", err)
	}
}

// TestStashPop_DirtyWorkingDirectory tests popping with dirty working directory
func TestStashPop_DirtyWorkingDirectory(t *testing.T) {
	_, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

	// Create and commit initial file
	testFile := "test.txt"
	if err := os.WriteFile(testFile, []byte("original"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(testFile); err != nil {
		t.Fatal(err)
	}
	if _, _, err := core.Commit("initial commit"); err != nil {
		t.Fatal(err)
	}

	// Modify and stash
	if err := os.WriteFile(testFile, []byte("stashed"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(testFile); err != nil {
		t.Fatal(err)
	}
	if err := core.Stash(); err != nil {
		t.Fatal(err)
	}

	// Make new changes
	if err := os.WriteFile(testFile, []byte("new changes"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Try to pop with dirty directory
	err := core.StashPop()
	if err == nil {
		t.Fatal("StashPop should fail with dirty working directory")
	}
	if !strings.Contains(err.Error(), "would be overwritten") {
		t.Errorf("Expected 'would be overwritten' error, got: %v", err)
	}
}

// TestStash_WIPCommitMessage tests the WIP commit message format
func TestStash_WIPCommitMessage(t *testing.T) {
	tmpDir, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

	// Create and commit initial file
	testFile := "test.txt"
	commitMessage := "test commit message"
	if err := os.WriteFile(testFile, []byte("content"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(testFile); err != nil {
		t.Fatal(err)
	}
	if _, _, err := core.Commit(commitMessage); err != nil {
		t.Fatal(err)
	}

	// Modify and stash
	if err := os.WriteFile(testFile, []byte("modified"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(testFile); err != nil {
		t.Fatal(err)
	}
	if err := core.Stash(); err != nil {
		t.Fatal(err)
	}

	// Read stash commit
	stashRefPath := filepath.Join(tmpDir, ".kitcat", "refs", "stash")
	stack := readStashStack(t, stashRefPath)
	if len(stack) == 0 {
		t.Fatalf("stash stack should not be empty")
	}
	stashHash := stack[0]
	stashCommit, err := storage.FindCommit(stashHash)
	if err != nil {
		t.Fatal(err)
	}

	// Verify message format: "WIP on <branch>: <commit_message>"
	expectedPrefix := "WIP on main: " + commitMessage
	if stashCommit.Message != expectedPrefix {
		t.Errorf("Expected stash message '%s', got: '%s'", expectedPrefix, stashCommit.Message)
	}
}

// TestStash_PreservesIndex tests that stash preserves the staged index
func TestStash_PreservesIndex(t *testing.T) {
	tmpDir, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

	// Create and commit initial files
	file1 := "file1.txt"
	file2 := "file2.txt"
	if err := os.WriteFile(file1, []byte("content1"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte("content2"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(file1); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(file2); err != nil {
		t.Fatal(err)
	}
	if _, _, err := core.Commit("initial commit"); err != nil {
		t.Fatal(err)
	}

	// Stage only file1
	if err := os.WriteFile(file1, []byte("modified1"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(file1); err != nil {
		t.Fatal(err)
	}

	// Stash
	if err := core.Stash(); err != nil {
		t.Fatal(err)
	}

	// Read stash commit and verify tree
	stashRefPath := filepath.Join(tmpDir, ".kitcat", "refs", "stash")
	stack := readStashStack(t, stashRefPath)
	if len(stack) == 0 {
		t.Fatalf("stash stack should not be empty")
	}
	stashHash := stack[0]
	stashCommit, err := storage.FindCommit(stashHash)
	if err != nil {
		t.Fatal(err)
	}

	// Parse stash tree
	stashTree, err := storage.ParseTree(stashCommit.TreeHash)
	if err != nil {
		t.Fatal(err)
	}

	// Verify both files are in the stash tree
	if _, ok := stashTree[file1]; !ok {
		t.Error("file1 should be in stash tree")
	}
	if _, ok := stashTree[file2]; !ok {
		t.Error("file2 should be in stash tree")
	}
}

// TestStash_MultipleFiles tests stashing multiple files
func TestStash_MultipleFiles(t *testing.T) {
	_, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

	// Create and commit multiple files
	files := []string{"file1.txt", "file2.txt", "file3.txt"}
	for _, file := range files {
		if err := os.WriteFile(file, []byte("original"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := core.AddFile(file); err != nil {
			t.Fatal(err)
		}
	}
	if _, _, err := core.Commit("initial commit"); err != nil {
		t.Fatal(err)
	}

	// Modify all files
	for _, file := range files {
		if err := os.WriteFile(file, []byte("modified"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := core.AddFile(file); err != nil {
			t.Fatal(err)
		}
	}

	// Stash
	if err := core.Stash(); err != nil {
		t.Fatalf("Stash failed: %v", err)
	}

	// Verify all files are reset
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		if string(content) != "original" {
			t.Errorf("File %s should be reset, got: %s", file, string(content))
		}
	}

	// Pop and verify all files are restored
	if err := core.StashPop(); err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		if string(content) != "modified" {
			t.Errorf("File %s should be restored, got: %s", file, string(content))
		}
	}
}

// TestStash_NoCommits tests stashing when there are no commits
func TestStash_NoCommits(t *testing.T) {
	_, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

	// Create a file without committing
	testFile := "test.txt"
	if err := os.WriteFile(testFile, []byte("content"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(testFile); err != nil {
		t.Fatal(err)
	}

	// Try to stash
	err := core.Stash()
	if err == nil {
		t.Fatal("Stash should fail when there are no commits")
	}
	if !strings.Contains(err.Error(), "no commits yet") {
		t.Errorf("Expected 'no commits yet' error, got: %v", err)
	}
}

// TestStash_UnstagedChanges integration test for stashing unstaged changes
func TestStash_UnstagedChanges(t *testing.T) {
	tmpDir, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

	// 1. Setup: Initialize and commit file.txt v1
	testFile := "file.txt"
	if err := os.WriteFile(testFile, []byte("v1"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(testFile); err != nil {
		t.Fatal(err)
	}
	if _, _, err := core.Commit("initial commit"); err != nil {
		t.Fatal(err)
	}

	// 2. Modify: Change content to v2 (unstaged)
	if err := os.WriteFile(testFile, []byte("v2"), 0o644); err != nil {
		t.Fatal(err)
	}

	// 3. Push: Run stash
	if err := core.Stash(); err != nil {
		t.Fatalf("Stash failed: %v", err)
	}

	// 4. Assert: Working directory is clean (content reverts to v1)
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "v1" {
		t.Errorf("File should be reset to v1 after stash, got: %s", string(content))
	}

	// 5. Pop: Run stash pop
	if err := core.StashPop(); err != nil {
		t.Fatalf("StashPop failed: %v", err)
	}

	// 6. Assert: Working directory contains v2
	content, err = os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "v2" {
		t.Errorf("File should be restored to v2 after pop, got: %s", string(content))
	}

	// 7. Verification: Ensure stash reference is gone
	stashRefPath := filepath.Join(tmpDir, ".kitcat", "refs", "stash")
	if stack := readStashStack(t, stashRefPath); len(stack) != 0 {
		t.Error("Stash reference should be deleted after pop")
	}
}

// TestStash_ListAndOrder verifies multiple stashes are stacked newest-first.
func TestStash_ListAndOrder(t *testing.T) {
	_, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

	file := "file.txt"
	if err := os.WriteFile(file, []byte("v1"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(file); err != nil {
		t.Fatal(err)
	}
	if _, _, err := core.Commit("initial"); err != nil {
		t.Fatal(err)
	}

	// First stash (v2)
	if err := os.WriteFile(file, []byte("v2"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(file); err != nil {
		t.Fatal(err)
	}
	if _, err := core.StashPush("first stash"); err != nil {
		t.Fatal(err)
	}

	// Second stash (v3)
	if err := os.WriteFile(file, []byte("v3"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(file); err != nil {
		t.Fatal(err)
	}
	if _, err := core.StashPush("second stash"); err != nil {
		t.Fatal(err)
	}

	commits, err := core.StashList()
	if err != nil {
		t.Fatal(err)
	}
	if len(commits) != 2 {
		t.Fatalf("expected 2 stashes, got %d", len(commits))
	}
	if commits[0].Message != "second stash" || commits[1].Message != "first stash" {
		t.Fatalf("unexpected stash order: %s, %s", commits[0].Message, commits[1].Message)
	}
}

// TestStash_ApplyByIndex ensures apply uses the requested index and keeps the stack.
func TestStash_ApplyByIndex(t *testing.T) {
	tmpDir, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

	file := "file.txt"
	if err := os.WriteFile(file, []byte("base"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(file); err != nil {
		t.Fatal(err)
	}
	if _, _, err := core.Commit("base"); err != nil {
		t.Fatal(err)
	}

	// Stash v2
	if err := os.WriteFile(file, []byte("v2"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(file); err != nil {
		t.Fatal(err)
	}
	if _, err := core.StashPush("stash v2"); err != nil {
		t.Fatal(err)
	}

	// Stash v3
	if err := os.WriteFile(file, []byte("v3"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(file); err != nil {
		t.Fatal(err)
	}
	if _, err := core.StashPush("stash v3"); err != nil {
		t.Fatal(err)
	}

	// Apply the older stash (index 1) -> expect v2
	if err := core.StashApply(1); err != nil {
		t.Fatalf("StashApply failed: %v", err)
	}

	content, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "v2" {
		t.Fatalf("expected file content 'v2', got '%s'", string(content))
	}

	// Stack should remain with 2 entries
	stack := readStashStack(t, filepath.Join(tmpDir, ".kitcat", "refs", "stash"))
	if len(stack) != 2 {
		t.Fatalf("stash stack should remain intact, expected 2 entries, got %d", len(stack))
	}
}

// TestStash_DropAndClear validates drop by index and clear.
func TestStash_DropAndClear(t *testing.T) {
	tmpDir, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

	file := "file.txt"
	if err := os.WriteFile(file, []byte("base"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(file); err != nil {
		t.Fatal(err)
	}
	if _, _, err := core.Commit("base"); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(file, []byte("v2"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(file); err != nil {
		t.Fatal(err)
	}
	if _, err := core.StashPush("stash v2"); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(file, []byte("v3"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(file); err != nil {
		t.Fatal(err)
	}
	if _, err := core.StashPush("stash v3"); err != nil {
		t.Fatal(err)
	}

	// Drop older stash (index 1)
	if err := core.StashDrop(1); err != nil {
		t.Fatalf("StashDrop failed: %v", err)
	}
	stack := readStashStack(t, filepath.Join(tmpDir, ".kitcat", "refs", "stash"))
	if len(stack) != 1 {
		t.Fatalf("expected 1 stash after drop, got %d", len(stack))
	}

	// Clear should remove all
	if err := core.StashClear(); err != nil {
		t.Fatalf("StashClear failed: %v", err)
	}
	stack = readStashStack(t, filepath.Join(tmpDir, ".kitcat", "refs", "stash"))
	if len(stack) != 0 {
		t.Fatalf("expected 0 stashes after clear, got %d", len(stack))
	}
}
