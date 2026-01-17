package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/LeeFred3042U/kitcat/internal/testutil"
)

// TestCreateBranch_InvalidName tests that CreateBranch rejects invalid branch names
func TestCreateBranch_InvalidName(t *testing.T) {
	_, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

	err := CreateBranch("../HEAD")
	if err == nil {
		t.Error("CreateBranch should reject '../HEAD' but it succeeded")
	}

	// Verify that .kitcat/HEAD was not altered
	headContent, err := os.ReadFile(".kitcat/HEAD")
	if err != nil {
		t.Fatalf("failed to read HEAD: %v", err)
	}
	expected := "ref: refs/heads/main\n"
	if string(headContent) != expected {
		t.Errorf("HEAD was altered! Expected %q, got %q", expected, string(headContent))
	}
}

// TestCreateBranch_InvalidName_ValidName tests that valid names are accepted
func TestCreateBranch_InvalidName_ValidName(t *testing.T) {
	_, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

	err := CreateBranch("feature-branch")
	if err != nil {
		t.Errorf("CreateBranch should accept 'feature-branch' but got error: %v", err)
	}

	// Verify branch was created
	branchPath := filepath.Join(".kitcat", "refs", "heads", "feature-branch")
	if _, err := os.Stat(branchPath); os.IsNotExist(err) {
		t.Error("Branch file was not created")
	}
}

// TestCreateBranch_InvalidName_ParentTraversal tests parent directory traversal attempts
func TestCreateBranch_InvalidName_ParentTraversal(t *testing.T) {
	_, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

	err := CreateBranch("../../etc/passwd")
	if err == nil {
		t.Error("CreateBranch should reject '../../etc/passwd' but it succeeded")
	}

	// Verify no file was created outside .kitcat
	if _, err := os.Stat("../../etc/passwd"); err == nil {
		t.Error("File was created outside repository!")
	}
}

// TestCreateBranch_InvalidName_Backslash tests backslash path separator
func TestCreateBranch_InvalidName_Backslash(t *testing.T) {
	_, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

	err := CreateBranch("..\\HEAD")
	if err == nil {
		t.Error("CreateBranch should reject '..\\HEAD' but it succeeded")
	}
}

// TestCreateBranch_InvalidName_ForwardSlash tests forward slash path separator
func TestCreateBranch_InvalidName_ForwardSlash(t *testing.T) {
	_, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

	err := CreateBranch("../refs/heads/malicious")
	if err == nil {
		t.Error("CreateBranch should reject '../refs/heads/malicious' but it succeeded")
	}
}

// TestCreateBranch_InvalidName_ControlChar tests control character injection
func TestCreateBranch_InvalidName_ControlChar(t *testing.T) {
	_, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

	// Test with null byte
	err := CreateBranch("branch\x00name")
	if err == nil {
		t.Error("CreateBranch should reject branch name with null byte but it succeeded")
	}

	// Test with newline
	err = CreateBranch("branch\nname")
	if err == nil {
		t.Error("CreateBranch should reject branch name with newline but it succeeded")
	}

	// Test with tab
	err = CreateBranch("branch\tname")
	if err == nil {
		t.Error("CreateBranch should reject branch name with tab but it succeeded")
	}
}

// TestRenameCurrentBranch_InvalidName tests that RenameCurrentBranch rejects invalid names
func TestRenameCurrentBranch_InvalidName(t *testing.T) {
	_, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

	err := RenameCurrentBranch("../HEAD")
	if err == nil {
		t.Error("RenameCurrentBranch should reject '../HEAD' but it succeeded")
	}

	// Verify that HEAD still points to main
	headContent, err := os.ReadFile(".kitcat/HEAD")
	if err != nil {
		t.Fatalf("failed to read HEAD: %v", err)
	}
	expected := "ref: refs/heads/main\n"
	if string(headContent) != expected {
		t.Errorf("HEAD was altered! Expected %q, got %q", expected, string(headContent))
	}

	// Verify main branch still exists
	mainPath := filepath.Join(".kitcat", "refs", "heads", "main")
	if _, err := os.Stat(mainPath); os.IsNotExist(err) {
		t.Error("Main branch was deleted!")
	}
}

// TestRenameCurrentBranch_InvalidName_ValidRename tests that valid renames work
func TestRenameCurrentBranch_InvalidName_ValidRename(t *testing.T) {
	_, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

	err := RenameCurrentBranch("develop")
	if err != nil {
		t.Errorf("RenameCurrentBranch should accept 'develop' but got error: %v", err)
	}

	// Verify HEAD points to develop
	headContent, err := os.ReadFile(".kitcat/HEAD")
	if err != nil {
		t.Fatalf("failed to read HEAD: %v", err)
	}
	expected := "ref: refs/heads/develop\n"
	if string(headContent) != expected {
		t.Errorf("HEAD not updated correctly. Expected %q, got %q", expected, string(headContent))
	}

	// Verify develop branch exists and main doesn't
	developPath := filepath.Join(".kitcat", "refs", "heads", "develop")
	if _, err := os.Stat(developPath); os.IsNotExist(err) {
		t.Error("Develop branch was not created")
	}

	mainPath := filepath.Join(".kitcat", "refs", "heads", "main")
	if _, err := os.Stat(mainPath); err == nil {
		t.Error("Main branch still exists after rename")
	}
}

// TestRenameCurrentBranch_InvalidName_InvalidRename tests invalid rename attempts
func TestRenameCurrentBranch_InvalidName_InvalidRename(t *testing.T) {
	_, cleanup := testutil.SetupTestRepo(t)
	defer cleanup()

	testCases := []struct {
		name     string
		newName  string
		expected string
	}{
		{"parent traversal", "../../malicious", "should reject parent traversal"},
		{"backslash", "..\\HEAD", "should reject backslash"},
		{"forward slash", "../refs/heads/bad", "should reject forward slash with .."},
		{"null byte", "branch\x00name", "should reject null byte"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := RenameCurrentBranch(tc.newName)
			if err == nil {
				t.Errorf("RenameCurrentBranch %s (tried %q)", tc.expected, tc.newName)
			}

			// Verify HEAD still points to main
			headContent, err := os.ReadFile(".kitcat/HEAD")
			if err != nil {
				t.Fatalf("failed to read HEAD: %v", err)
			}
			expected := "ref: refs/heads/main\n"
			if string(headContent) != expected {
				t.Errorf("HEAD was altered! Expected %q, got %q", expected, string(headContent))
			}
		})
	}
}
