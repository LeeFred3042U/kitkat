package test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/LeeFred3042U/kitkat/internal/core"
	"github.com/LeeFred3042U/kitkat/internal/models"
	"github.com/LeeFred3042U/kitkat/internal/storage"
)

// TestRepo represents an isolated temporary KitKat repository for integration testing.
type TestRepo struct {
	t        *testing.T
	dir      string
	commits  []models.Commit
	branches map[string]string
	files    map[string]string
}

// NewTestRepo creates and initializes a test repository with optional repo setup.
func NewTestRepo(t *testing.T, initRepo bool) *TestRepo {
	t.Helper()
	tr := &TestRepo{
		t:        t,
		dir:      "",
		commits:  []models.Commit{},
		branches: make(map[string]string),
		files:    make(map[string]string),
	}
	tr.dir = tr.setupSandbox(initRepo)
	return tr
}

func logFatal(t *testing.T, msg string, err error) {
	if err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
}

func logError(t *testing.T, msg string, err error) {
	if err != nil {
		t.Errorf("%s: %v", msg, err)
	}
}

// setupSandbox creates an isolated temporary environment and optionally initializes a repo.
func (tr *TestRepo) setupSandbox(initRepo bool) string {
	tr.t.Helper()
	tempDir := tr.t.TempDir()
	origDir, err := os.Getwd()
	logFatal(tr.t, "Getwd failed", err)

	origHome := os.Getenv("HOME")
	logFatal(tr.t, "Setenv HOME", os.Setenv("HOME", tempDir))
	logFatal(tr.t, "Chdir tempDir", os.Chdir(tempDir))

	tr.t.Cleanup(func() {
		logError(tr.t, "Cleanup Chdir", os.Chdir(origDir))
		os.Setenv("HOME", origHome)
	})

	if initRepo {
		logFatal(tr.t, "InitRepo", core.InitRepo())
	}
	return tempDir
}

// AddFile creates a file in the working directory. Returns tr for fluent chaining.
func (tr *TestRepo) AddFile(path, content string) *TestRepo {
	tr.t.Helper()
	logFatal(tr.t, fmt.Sprintf("WriteFile %s", path), os.WriteFile(path, []byte(content), 0o644))
	tr.files[path] = content
	return tr
}

// modifyFile modifies an existing file without committing (dirty working directory).
func (tr *TestRepo) modifyFile(path, content string) error {
	tr.t.Helper()
	return os.WriteFile(path, []byte(content), 0o644)
}

// CommitSingle commits the current state. Returns tr for fluent chaining.
func (tr *TestRepo) CommitSingle(msg string) *TestRepo {
	tr.t.Helper()
	last, err := storage.GetLastCommit()
	logError(tr.t, "GetLastCommit after CommitSingle", err)
	if last.ID != "" || last.Parent != "" {
		tr.commits = append(tr.commits, last)
	}
	return tr
}

// CommitAll stages all files and commits with given message. Returns tr for fluent chaining.
func (tr *TestRepo) CommitAll(msg string) *TestRepo {
	tr.t.Helper()
	core.CommitAll(msg)
	last, err := storage.GetLastCommit()
	logError(tr.t, "GetLastCommit after CommitAll", err)
	if last.ID != "" || last.Parent != "" {
		tr.commits = append(tr.commits, last)
	}
	return tr
}

// CreateBranch creates a new branch at current HEAD. Returns tr for fluent chaining.
func (tr *TestRepo) CreateBranch(name string) *TestRepo {
	tr.t.Helper()
	logFatal(tr.t, fmt.Sprintf("CreateBranch %s", name), core.CreateBranch(name))
	last, err := storage.GetLastCommit()
	logError(tr.t, "GetLastCommit for branch", err)
	if last.ID != "" {
		tr.branches[name] = last.ID
	}
	return tr
}

// CheckoutFile restores a file from the last commit. Returns tr for fluent chaining.
func (tr *TestRepo) CheckoutFile(path string) *TestRepo {
	tr.t.Helper()
	logFatal(tr.t, fmt.Sprintf("CheckoutFile %s", path), core.CheckoutFile(path))
	return tr
}

// CheckoutBranch switches to the named branch. Returns tr for fluent chaining.
func (tr *TestRepo) CheckoutBranch(name string) *TestRepo {
	tr.t.Helper()
	logFatal(tr.t, fmt.Sprintf("CheckoutBranch %s", name), core.CheckoutBranch(name))
	return tr
}

// CheckoutCommit moves HEAD to a specific commit (detached state). Returns tr for fluent chaining.
func (tr *TestRepo) CheckoutCommit(hash string) *TestRepo {
	tr.t.Helper()
	logFatal(tr.t, fmt.Sprintf("CheckoutCommit %s", hash), core.CheckoutCommit(hash))
	return tr
}

// AssertFileContent verifies file content matches expected value.
func (tr *TestRepo) AssertFileContent(path, want string) {
	tr.t.Helper()
	content, err := os.ReadFile(path)
	if err != nil || string(content) != want {
		tr.t.Errorf("AssertFile %s: got %q (err: %v), want %q", path, content, err, want)
	}
}

// AssertHEAD verifies HEAD points to expected reference or commit.
func (tr *TestRepo) AssertHEAD(want string) {
	tr.t.Helper()
	headPath := filepath.Join(tr.dir, ".kitkat/HEAD")
	head, err := os.ReadFile(headPath)
	if err != nil || strings.TrimSpace(string(head)) != want {
		tr.t.Errorf("AssertHEAD: got %q (err: %v), want %q", head, err, want)
	}
}

// AssertBranchExists verifies a branch ref file exists.
func (tr *TestRepo) AssertBranchExists(name string) {
	tr.t.Helper()
	refPath := filepath.Join(tr.dir, ".kitkat/refs/heads", name)
	if _, err := os.Stat(refPath); err != nil {
		tr.t.Errorf("Branch %s ref missing: %v", name, err)
	}
}
