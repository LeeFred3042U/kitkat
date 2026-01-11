package core

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestGrepBasic(t *testing.T) {
	tmpDir := t.TempDir()

	// Change working directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Create test files
	mainContent := "package main\n\nfunc main() {}\n"
	err = os.WriteFile("main.go", []byte(mainContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	_ = os.WriteFile("ignore.go", []byte("func ignored() {}\n"), 0644)

	// Create kitcat index
	err = os.Mkdir(".kitcat", 0755)
	if err != nil {
		t.Fatal(err)
	}

	entries := []IndexEntry{
		{Path: "main.go", Hash: "dummy"},
	}

	err = SaveIndex(entries)
	if err != nil {
		t.Fatal(err)
	}

	// Capture stdout
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run grep
	err = Grep([]string{"--line-number", "func"})
	if err != nil {
		t.Fatalf("grep returned error: %v", err)
	}

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	buf.ReadFrom(r)

	output := buf.String()

	// Assert expected match
	if !strings.Contains(output, "main.go:3:func main()") {
		t.Fatalf("expected match not found, got:\n%s", output)
	}

	// Ensure untracked file is ignored
	if strings.Contains(output, "ignore.go") {
		t.Fatalf("untracked file should not be searched")
	}
}
