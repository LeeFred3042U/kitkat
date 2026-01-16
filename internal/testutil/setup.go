package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/LeeFred3042U/kitcat/internal/core"
)

// SetupTestRepo creates a temporary kitcat repository for tests.
// It returns the temp directory path and a cleanup function.
func SetupTestRepo(t *testing.T) (string, func()) {
	t.Helper()

	// Save current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Create temp directory
	tempDir := t.TempDir()

	// Change to temp directory
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}

	// Initialize repository
	if err := core.InitRepo(); err != nil {
		t.Fatal(err)
	}

	// Create dummy HEAD if needed
	headPath := filepath.Join(tempDir, ".kitcat", "HEAD")
	if err := os.WriteFile(headPath, []byte("ref: refs/heads/main\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Cleanup function
	cleanup := func() {
		_ = os.Chdir(cwd)
		_ = os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}
