package testutil

import (
	"os"
	"testing"

	"github.com/LeeFred3042U/kitcat/internal/core"
)

// SetupTestRepo creates a temporary kitcat repository for tests.
// Returns the repo path and a cleanup function.
func SetupTestRepo(t *testing.T) (string, func()) {
	t.Helper()

	// Create temp directory
	tempDir := t.TempDir()

	// Save current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	// Change to temp repo
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	// Initialize repo
	if err := core.InitRepo(); err != nil {
		t.Fatalf("failed to init repo: %v", err)
	}

	cleanup := func() {
		_ = os.Chdir(cwd)
	}

	return tempDir, cleanup
}
