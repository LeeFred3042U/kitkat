package core

import (
	"os"
	"testing"
)

func Test_CheckoutFile_OverwritesDirtyFile(t *testing.T) {
	// Restores the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(cwd)

	// Create temp directory
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Initialize kitkat repository
	if err := InitRepo(); err != nil {
		t.Fatal(err)
	}

	fileName := "foo.txt"
	v1Content := []byte("v1")
	v2Content := []byte("v2")

	// 1. Create file with v1
	if err := os.WriteFile(fileName, v1Content, 0644); err != nil {
		t.Fatal(err)
	}

	// 2. Stage and Commit v1
	if err := AddFile(fileName); err != nil {
		t.Fatalf("failed to stage file: %v", err)
	}
	if _, _, err := Commit("initial commit"); err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	// 3. Modify file to v2 (dirty working copy)
	if err := os.WriteFile(fileName, v2Content, 0644); err != nil {
		t.Fatal(err)
	}

	// Verify it is v2 before checkout
	content, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "v2" {
		t.Fatalf("setup failed: expected v2, got %s", string(content))
	}

	// 4. Checkout file
	// This should overwrite v2 with v1 from the commit
	if err := CheckoutFile(fileName); err != nil {
		t.Fatalf("CheckoutFile failed: %v", err)
	}

	// 5. Assert destructive behavior
	// If the bug exists, content should now be v1
	content, err = os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != "v1" {
		t.Errorf("Expected CheckoutFile to overwrite dirty file (reproducing bug), but got content: %s", string(content))
	} else {
		t.Logf("Confirmed: CheckoutFile overwrote dirty file as expected (bug reproduction success)")
	}
}
