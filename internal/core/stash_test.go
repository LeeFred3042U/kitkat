package core

import (
	"os"
	"strings"
	"testing"
)

func TestStash(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(cwd)

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	if err := InitRepo(); err != nil {
		t.Fatal(err)
	}

	// 1. Setup initial state
	if err := os.WriteFile("foo.txt", []byte("v1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := AddFile("foo.txt"); err != nil {
		t.Fatal(err)
	}
	if _, _, err := Commit("initial commit"); err != nil {
		t.Fatal(err)
	}

	// 2. Create changes to stash
	// Modify foo.txt (staged)
	if err := os.WriteFile("foo.txt", []byte("v2"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := AddFile("foo.txt"); err != nil {
		t.Fatal(err)
	}

	// Create bar.txt (untracked)
	if err := os.WriteFile("bar.txt", []byte("bar"), 0644); err != nil {
		t.Fatal(err)
	}

	// 3. Stash
	if err := StashSave("test stash"); err != nil {
		t.Fatalf("StashSave failed: %v", err)
	}

	// 4. Verify working directory is clean (reset to HEAD)
	fooContent, err := os.ReadFile("foo.txt")
	if err != nil {
		t.Fatal(err)
	}
	if string(fooContent) != "v1" {
		t.Errorf("expected foo.txt to be 'v1', got '%s'", string(fooContent))
	}

	if _, err := os.Stat("bar.txt"); !os.IsNotExist(err) {
		t.Errorf("expected bar.txt to be removed (untracked)")
	}

	// 5. Verify Stash List
	logContent, err := os.ReadFile(".kitkat/logs/refs/stash")
	if err != nil {
		t.Fatal(err)
	}
	if len(logContent) == 0 {
		t.Errorf("expected stash log to be non-empty")
	}

	// 6. Pop Stash
	if err := StashPop(0); err != nil {
		t.Fatalf("StashPop failed: %v", err)
	}

	// 7. Verify changes restored
	fooContent, err = os.ReadFile("foo.txt")
	if err != nil {
		t.Fatal(err)
	}
	if string(fooContent) != "v2" {
		t.Errorf("expected foo.txt to be 'v2', got '%s'", string(fooContent))
	}

	barContent, err := os.ReadFile("bar.txt")
	if err != nil {
		t.Fatal(err)
	}
	if string(barContent) != "bar" {
		t.Errorf("expected bar.txt to be 'bar', got '%s'", string(barContent))
	}

	// 8. Verify Stash List is empty
	if _, err := os.Stat(".kitkat/logs/refs/stash"); !os.IsNotExist(err) {
		content, _ := os.ReadFile(".kitkat/logs/refs/stash")
		if len(strings.TrimSpace(string(content))) > 0 {
			t.Errorf("expected stash log to be empty, got: %s", string(content))
		}
	}
}
