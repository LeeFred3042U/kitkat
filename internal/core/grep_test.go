package core

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// helper to run commands
func runCmd(t *testing.T, dir string, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("command failed: %s %v", name, args)
	}
}

func TestGrepBasic(t *testing.T) {
	tmpDir := t.TempDir()

	// IMPORTANT: change working directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	// Init git repo
	runCmd(t, tmpDir, "git", "init")

	// Create tracked file
	fileContent := "package main\n\nfunc main() {}\n"
	err = os.WriteFile("main.go", []byte(fileContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Create untracked file (should be ignored)
	_ = os.WriteFile("ignore.go", []byte("func ignored() {}\n"), 0644)

	// Track only main.go
	runCmd(t, tmpDir, "git", "add", "main.go")

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

	// Assertions
	if !strings.Contains(output, "main.go:3:func main()") {
		t.Fatalf("expected match not found, got:\n%s", output)
	}

	if strings.Contains(output, "ignore.go") {
		t.Fatalf("untracked file was searched")
	}
}
