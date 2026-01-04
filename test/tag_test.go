package test

import (
	"os"
	"reflect"
	"testing"

	"github.com/LeeFred3042U/kitkat/internal/core"
)

func TestTagCommand(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(origDir)

	t.Run("Tagging before init", func(t *testing.T) {
		err := core.CreateTag("v1", "hash123")
		if err == nil {
			t.Error("Expected error when tagging in uninitialized repo")
		}
	})

	core.InitRepo()

	t.Run("Create Tag and Verify Content", func(t *testing.T) {
		tagName := "v1.0"
		commitID := "abc1234567890"
		
		if err := core.CreateTag(tagName, commitID); err != nil {
			t.Fatalf("Failed to create tag: %v", err)
		}

		// Verify physical file content
		tagPath := ".kitkat/refs/tags/" + tagName
		data, err := os.ReadFile(tagPath)
		if err != nil {
			t.Fatalf("Could not read tag file: %v", err)
		}
		if string(data) != commitID {
			t.Errorf("Tag content mismatch: got %s, want %s", string(data), commitID)
		}
	})

	t.Run("List Multiple Tags", func(t *testing.T) {
		tags := []string{"alpha", "beta", "gamma"}
		for _, name := range tags {
			core.CreateTag(name, "some-hash")
		}

		got, err := core.ListTags() //
		if err != nil {
			t.Fatalf("ListTags failed: %v", err)
		}

		// ListTags sorts alphabetically
		// Existing tags from previous test: v1.0
		expected := []string{"alpha", "beta", "gamma", "v1.0"}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("ListTags mismatch.\nGot: %v\nWant: %v", got, expected)
		}
	})

	t.Run("Prevent Duplicate Tags", func(t *testing.T) {
		tagName := "duplicate"
		core.CreateTag(tagName, "hash1")
		err := core.CreateTag(tagName, "hash2") //
		if err == nil {
			t.Error("Expected error when creating a tag that already exists")
		}
	})
}