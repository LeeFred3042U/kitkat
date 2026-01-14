package core

import (
	"os"
	"testing"

	"github.com/LeeFred3042U/kitcat/internal/storage"
)

func TestMoveFile(t *testing.T) {
	// Restores the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Chdir(cwd)
	}()

	// Create temp directory
	tmpDir := t.TempDir()

	// Change working directory into temp repo
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Initialize kitcat repository
	if err := InitRepo(); err != nil {
		t.Fatal(err)
	}

	oldPath := "old_test.txt"
	newPath := "new_test.txt"

	// Create old file
	if err := os.WriteFile(oldPath, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Stage old file
	if err := AddFile(oldPath); err != nil {
		t.Fatal(err)
	}

	// Move file
	if err := MoveFile(oldPath, newPath, false); err != nil {
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
	idx, err := storage.LoadIndex()
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := idx[newPath]; !ok {
		t.Fatalf("expected new file to be staged")
	}
}

func TestMoveFile_DestinationExists(t *testing.T) {
	// Restores the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Chdir(cwd)
	}()

	// Create temp directory
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Initialize kitcat repository
	if err := InitRepo(); err != nil {
		t.Fatal(err)
	}

	// Create source and destination files
	src := "source.txt"
	dst := "destination.txt"
	if err := os.WriteFile(src, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dst, []byte("hola"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Stage source and destination files
	if err := AddFile(src); err != nil {
		t.Fatal(err)
	}
	if err := AddFile(dst); err != nil {
		t.Fatal(err)
	}

	// Moves the file without the flag
	if err := MoveFile(src, dst, false); err == nil {
		t.Fatalf("expected error when destination exists")
	}

	// Moves the file with the flag
	if err := MoveFile(src, dst, true); err != nil {
		t.Fatal(err)
	}

	// Checks for the files to be moved
	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Fatalf("expected src to be moved")
	}

	if _, err := os.Stat(dst); err != nil {
		t.Fatalf("expected dst to exist")
	}

	// Index should contain the destination file
	idx, err := storage.LoadIndex()
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := idx[dst]; !ok {
		t.Fatalf("expected new file to be staged")
	}
}

func TestMoveFile_SamePath(t *testing.T) {
	// Restores the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Chdir(cwd)
	}()

	// Create a temp directory
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Initialize kitcat repository
	if err := InitRepo(); err != nil {
		t.Fatal(err)
	}

	// Create source and destination files
	f := "file"
	if err := os.WriteFile(f, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Stage source and destination files
	if err := AddFile(f); err != nil {
		t.Fatal(err)
	}

	if err := MoveFile(f, f, false); err == nil {
		t.Fatalf("expected error for same source and destination")
	}

	if err := MoveFile(f, f, true); err == nil {
		t.Fatalf("expected error for same source and destination")
	}
}

func TestMoveFile_LocalChanges_WithoutForce(t *testing.T) {
	// Restores the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Chdir(cwd)
	}()

	// Create temp directory
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Initialize kitcat repository
	if err := InitRepo(); err != nil {
		t.Fatal(err)
	}

	oldPath := "test.txt"
	newPath := "renamed.txt"

	// Create and stage file
	if err := os.WriteFile(oldPath, []byte("original content"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := AddFile(oldPath); err != nil {
		t.Fatal(err)
	}

	// Modify the file after staging
	if err := os.WriteFile(oldPath, []byte("modified content"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Attempt to move without force flag
	err = MoveFile(oldPath, newPath, false)
	if err == nil {
		t.Fatalf("expected error when moving modified file without force")
	}
	if err.Error() != "local changes present, use -f to force" {
		t.Fatalf("expected 'local changes present' error, got: %v", err)
	}

	// Verify file remains at oldPath
	if _, err := os.Stat(oldPath); err != nil {
		t.Fatalf("expected file to remain at oldPath after failed move")
	}

	// Verify newPath does not exist
	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		t.Fatalf("expected newPath to not exist after failed move")
	}

	// Verify index still contains oldPath
	idx, err := storage.LoadIndex()
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := idx[oldPath]; !ok {
		t.Fatalf("expected oldPath to remain in index after failed move")
	}
}

func TestMoveFile_LocalChanges_WithForce(t *testing.T) {
	// Restores the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Chdir(cwd)
	}()

	// Create temp directory
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Initialize kitcat repository
	if err := InitRepo(); err != nil {
		t.Fatal(err)
	}

	oldPath := "test.txt"
	newPath := "renamed.txt"

	// Create and stage file
	if err := os.WriteFile(oldPath, []byte("original content"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := AddFile(oldPath); err != nil {
		t.Fatal(err)
	}

	// Modify the file after staging
	if err := os.WriteFile(oldPath, []byte("modified content"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Move with force flag
	if err := MoveFile(oldPath, newPath, true); err != nil {
		t.Fatalf("expected no error when moving with force flag, got: %v", err)
	}

	// Verify file moved to newPath
	if _, err := os.Stat(newPath); err != nil {
		t.Fatalf("expected file to exist at newPath")
	}

	// Verify oldPath no longer exists
	if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
		t.Fatalf("expected oldPath to be removed after move")
	}

	// Verify index contains newPath
	idx, err := storage.LoadIndex()
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := idx[newPath]; !ok {
		t.Fatalf("expected newPath to be in index after move")
	}
	if _, ok := idx[oldPath]; ok {
		t.Fatalf("expected oldPath to be removed from index after move")
	}
}

func TestMoveFile_UntrackedFile(t *testing.T) {
	// Restores the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Chdir(cwd)
	}()

	// Create temp directory
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Initialize kitcat repository
	if err := InitRepo(); err != nil {
		t.Fatal(err)
	}

	oldPath := "untracked.txt"
	newPath := "moved.txt"

	// Create file but do NOT stage it
	if err := os.WriteFile(oldPath, []byte("untracked content"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Move untracked file without force flag (should succeed)
	if err := MoveFile(oldPath, newPath, false); err != nil {
		t.Fatalf("expected no error when moving untracked file, got: %v", err)
	}

	// Verify file moved to newPath
	if _, err := os.Stat(newPath); err != nil {
		t.Fatalf("expected file to exist at newPath")
	}

	// Verify oldPath no longer exists
	if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
		t.Fatalf("expected oldPath to be removed after move")
	}

	// Verify index contains newPath (it gets staged by MoveFile)
	idx, err := storage.LoadIndex()
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := idx[newPath]; !ok {
		t.Fatalf("expected newPath to be staged after move")
	}
}
