package main

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCLIExitCodeOnInvalidCommand(t *testing.T) {
	// Setup isolated environment
	tmpDir := t.TempDir()
	binName := "kitkat"
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	binPath := filepath.Join(tmpDir, binName)

	// Build the binary from current source
	buildCmd := exec.Command("go", "build", "-o", binPath, "main.go")
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build kitkat binary: %v\nOutput: %s", err, output)
	}

	// Run an invalid command
	// Note: Currently, cmd/main.go has logic for unknown commands:
	// if !ok { fmt.Println("Unknown command..."); os.Exit(1) }
	// This specific path should pass.
	cmd := exec.Command(binPath, "definitely-not-a-command")
	err := cmd.Run()

	if err == nil {
		t.Error("Expected error (non-zero exit code) for unknown command, but got exit code 0")
	}

	// Verification of the "Silent Failure" risk mentioned in audits:
	// Commands that return normally after printing an error (e.g. 'add' with missing file)
	// will fail this test if the exit code is 0.
	t.Run("CommandLogicErrorExitCode", func(t *testing.T) {
		// Initialize repo first to trigger porcelain logic
		initCmd := exec.Command(binPath, "init")
		initCmd.Dir = tmpDir
		_ = initCmd.Run()

		// Attempt to add a non-existent file
		addCmd := exec.Command(binPath, "add", "non-existent-file.txt")
		addCmd.Dir = tmpDir
		addErr := addCmd.Run()

		if addErr == nil {
			t.Errorf("FAIL: 'kitkat add' on missing file exited with code 0, expected non-zero")
		}
	})
}
