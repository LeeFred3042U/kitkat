package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/LeeFred3042U/kitkat/internal/core"
)

func TestConfigCommand(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	
	// Mock HOME so UserHomeDir() points to our temp directory
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	
	os.Chdir(tempDir)
	defer func() {
		os.Chdir(origDir)
		os.Setenv("HOME", origHome)
	}()

	t.Run("Standard Get and Set with Persistence", func(t *testing.T) {
		key := "user.email"
		val := "bot@test.com"

		// Write config
		if err := core.SetConfig(key, val); err != nil {
			t.Fatalf("SetConfig failed: %v", err)
		}

		// Verify the physical file exists in the mocked home
		configPath := filepath.Join(tempDir, ".kitkatconfig")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Error("Config file was not physically created in home directory")
		}

		// Verify GetConfig retrieves correctly
		got, ok, err := core.GetConfig(key)
		if err != nil || !ok || got != val {
			t.Errorf("GetConfig = %v, %v; want %v, true", got, ok, val)
		}
	})

	t.Run("Overwrite and Special Characters", func(t *testing.T) {
		key := "complex.key-name"
		val1 := "first_value"
		val2 := "second value with spaces"

		core.SetConfig(key, val1)
		core.SetConfig(key, val2) // Overwrite

		got, _, _ := core.GetConfig(key)
		if got != val2 {
			t.Errorf("Expected %s, got %s", val2, got)
		}
	})

	t.Run("Empty Values", func(t *testing.T) {
		key := "empty.val"
		if err := core.SetConfig(key, ""); err != nil {
			t.Errorf("Should allow setting empty values: %v", err)
		}
		got, ok, _ := core.GetConfig(key)
		if !ok || got != "" {
			t.Errorf("Failed to retrieve empty value: got %s, ok %v", got, ok)
		}
	})
}