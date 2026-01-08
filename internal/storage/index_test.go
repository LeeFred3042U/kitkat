package storage

import (
	"os"
	"strings"
	"testing"
)

func Test_LoadIndex_TruncatedFile(t *testing.T) {
	// 1. Create a temporary directory for the test
	tmpDir := t.TempDir()

	// 2. Save current working directory so we can return to it later
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current wd: %v", err)
	}

	// 3. Switch to the temp directory
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change wd: %v", err)
	}
	// Ensure we switch back to the original directory when the test finishes
	defer os.Chdir(oldWd)

	// 4. Setup the .kitkat directory
	err = os.Mkdir(".kitkat", 0755)
	if err != nil {
		t.Fatalf("failed to create .kitkat dir: %v", err)
	}

	// 5. Create a truncated JSON file in .kitkat/index
	// This JSON is missing its closing brace
	truncatedJSON := `{"file.txt": "sha1hashcode"`
	err = os.WriteFile(indexPath, []byte(truncatedJSON), 0644)
	if err != nil {
		t.Fatalf("failed to write truncated index: %v", err)
	}

	// 6. Call the production loader
	entries, err := LoadIndex()

	// 7. Assertions: We expect an error and NO entries
	if err == nil {
		t.Error("expected an error when loading a truncated index, but got nil")
	}

	if entries != nil && len(entries) > 0 {
		t.Errorf("expected entries to be nil or empty on failure, got: %v", entries)
	}

	// 8. Verify error message contains 'parse' or 'index'
	errMsg := strings.ToLower(err.Error())
	if !strings.Contains(errMsg, "parse") && !strings.Contains(errMsg, "index") {
		t.Errorf("error message should mention parse failure, got: %s", err.Error())
	}
}
