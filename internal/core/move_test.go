package core_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/LeeFred3042U/kitcat/internal/core"
	"github.com/LeeFred3042U/kitcat/internal/testutil"
)

func TestMoveFile(t *testing.T) {
	repoDir, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

	oldPath := "old_test.txt"
	newPath := "new_test.txt"

	if err := os.WriteFile(oldPath, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := core.AddFile(oldPath); err != nil {
		t.Fatal(err)
	}

	if err := core.MoveFile(oldPath, newPath, false); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
		t.Fatalf("expected old file to be removed")
	}

	if _, err := os.Stat(newPath); err != nil {
		t.Fatalf("expected new file to exist")
	}

	idx, err := loadIndexForTest(filepath.Join(repoDir, ".kitcat", "index"))
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := idx[newPath]; !ok {
		t.Fatalf("expected new file to be staged")
	}
}

func TestMoveFile_DestinationExists(t *testing.T) {
	repoDir, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

	src := "source.txt"
	dst := "destination.txt"

	if err := os.WriteFile(src, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dst, []byte("hola"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := core.AddFile(src); err != nil {
		t.Fatal(err)
	}
	if err := core.AddFile(dst); err != nil {
		t.Fatal(err)
	}

	if err := core.MoveFile(src, dst, false); err == nil {
		t.Fatalf("expected error when destination exists")
	}

	if err := core.MoveFile(src, dst, true); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Fatalf("expected src to be moved")
	}

	if _, err := os.Stat(dst); err != nil {
		t.Fatalf("expected dst to exist")
	}

	idx, err := loadIndexForTest(filepath.Join(repoDir, ".kitcat", "index"))
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := idx[dst]; !ok {
		t.Fatalf("expected new file to be staged")
	}
}

func TestMoveFile_SamePath(t *testing.T) {
	_, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

	f := "file"

	if err := os.WriteFile(f, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := core.AddFile(f); err != nil {
		t.Fatal(err)
	}

	if err := core.MoveFile(f, f, false); err == nil {
		t.Fatalf("expected error for same source and destination")
	}

	if err := core.MoveFile(f, f, true); err == nil {
		t.Fatalf("expected error for same source and destination")
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
