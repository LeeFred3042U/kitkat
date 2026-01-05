package test

import (
	"strings"
	"testing"
)

// runUninitTest runs standard uninitialized repo test
func runUninitTest(t *testing.T, name string, fn func() error) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		NewTestRepo(t, false)
		err := fn()
		if err == nil {
			t.Error("Expected error in uninitialized repo")
		}
	})
}

// runHappyTest runs success path with repo
func runHappyTest(t *testing.T, name string, fn func(*TestRepo)) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		repo := NewTestRepo(t, true)
		fn(repo)
	})
}

// runErrorTest runs expected error cases
func runErrorTest(t *testing.T, name string, wantMsg string, fn func(*TestRepo) error) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		repo := NewTestRepo(t, true)
		err := fn(repo)
		if err == nil {
			t.Errorf("Expected error containing %q", wantMsg)
			return
		}
		if !strings.Contains(err.Error(), wantMsg) {
			t.Errorf("Expected %q, got: %v", wantMsg, err)
		}
	})
}
