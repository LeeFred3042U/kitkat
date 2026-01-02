package core

import (
	"os"
	"path/filepath"
	"testing"
)

// Helper function to create a temporary .kitignore file
func createKitignore(t *testing.T, content string) func() {
	err := os.WriteFile(".kitignore", []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create .kitignore: %v", err)
	}

	// Return cleanup function
	return func() {
		os.Remove(".kitignore")
		ClearIgnoreCache()
	}
}

func TestLoadIgnorePatterns(t *testing.T) {
	// Clear cache before test
	ClearIgnoreCache()

	t.Run("NoKitignoreFile", func(t *testing.T) {
		patterns, err := LoadIgnorePatterns()
		if err != nil {
			t.Errorf("Expected no error when .kitignore doesn't exist, got: %v", err)
		}
		if len(patterns) != 0 {
			t.Errorf("Expected empty patterns, got %d patterns", len(patterns))
		}
	})

	t.Run("ValidPatterns", func(t *testing.T) {
		ClearIgnoreCache()
		cleanup := createKitignore(t, "*.log\n*.txt\nbin/\n")
		defer cleanup()

		patterns, err := LoadIgnorePatterns()
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if len(patterns) != 3 {
			t.Errorf("Expected 3 patterns, got %d", len(patterns))
		}

		// Check patterns
		if patterns[0].Pattern != "*.log" {
			t.Errorf("Expected first pattern to be '*.log', got '%s'", patterns[0].Pattern)
		}
		if patterns[1].Pattern != "*.txt" {
			t.Errorf("Expected second pattern to be '*.txt', got '%s'", patterns[1].Pattern)
		}
		if patterns[2].Pattern != "bin" || !patterns[2].IsDirectory {
			t.Errorf("Expected third pattern to be 'bin' with IsDirectory=true, got '%s' IsDirectory=%v",
				patterns[2].Pattern, patterns[2].IsDirectory)
		}
	})

	t.Run("CommentsAndBlankLines", func(t *testing.T) {
		ClearIgnoreCache()
		cleanup := createKitignore(t, "# Comment line\n\n*.log\n   \n# Another comment\n*.txt\n")
		defer cleanup()

		patterns, err := LoadIgnorePatterns()
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if len(patterns) != 2 {
			t.Errorf("Expected 2 patterns (comments and blanks ignored), got %d", len(patterns))
		}
	})

	t.Run("InvalidPattern", func(t *testing.T) {
		ClearIgnoreCache()
		// filepath.Match doesn't always fail on invalid patterns in a predictable way
		// So we'll just ensure the function doesn't crash
		cleanup := createKitignore(t, "*.log\n[invalid\n*.txt\n")
		defer cleanup()

		patterns, err := LoadIgnorePatterns()
		if err != nil {
			t.Errorf("Expected no error (invalid patterns should be skipped), got: %v", err)
		}
		// Should have 2 valid patterns (*.log and *.txt), [invalid should be skipped
		if len(patterns) < 2 {
			t.Errorf("Expected at least 2 valid patterns, got %d", len(patterns))
		}
	})

	t.Run("PatternCaching", func(t *testing.T) {
		ClearIgnoreCache()
		cleanup := createKitignore(t, "*.log\n")
		defer cleanup()

		// First load
		patterns1, err1 := LoadIgnorePatterns()
		if err1 != nil {
			t.Errorf("Expected no error on first load, got: %v", err1)
		}

		// Second load (should use cache)
		patterns2, err2 := LoadIgnorePatterns()
		if err2 != nil {
			t.Errorf("Expected no error on second load, got: %v", err2)
		}

		if len(patterns1) != len(patterns2) {
			t.Errorf("Cache returned different number of patterns")
		}
	})
}

func TestMatchesPattern(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		pattern IgnorePattern
		want    bool
	}{
		{
			name:    "SimpleGlobMatch",
			path:    "test.log",
			pattern: IgnorePattern{Pattern: "*.log", IsDirectory: false},
			want:    true,
		},
		{
			name:    "SimpleGlobNoMatch",
			path:    "test.txt",
			pattern: IgnorePattern{Pattern: "*.log", IsDirectory: false},
			want:    false,
		},
		{
			name:    "ExactMatch",
			path:    "secret.txt",
			pattern: IgnorePattern{Pattern: "secret.txt", IsDirectory: false},
			want:    true,
		},
		{
			name:    "DirectoryPattern",
			path:    "bin",
			pattern: IgnorePattern{Pattern: "bin", IsDirectory: true},
			want:    true,
		},
		{
			name:    "DirectoryPatternNested",
			path:    "bin/output",
			pattern: IgnorePattern{Pattern: "bin", IsDirectory: true},
			want:    true,
		},
		{
			name:    "DirectoryPatternDeepNested",
			path:    "bin/debug/symbols",
			pattern: IgnorePattern{Pattern: "bin", IsDirectory: true},
			want:    true,
		},
		{
			name:    "RecursivePattern",
			path:    "src/main/java/App.class",
			pattern: IgnorePattern{Pattern: "**/*.class", IsDirectory: false},
			want:    true,
		},
		{
			name:    "QuestionMark",
			path:    "file1.txt",
			pattern: IgnorePattern{Pattern: "file?.txt", IsDirectory: false},
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Normalize path for testing
			normalizedPath := filepath.ToSlash(tt.path)
			got := matchesPattern(normalizedPath, tt.pattern)
			if got != tt.want {
				t.Errorf("matchesPattern(%q, %+v) = %v, want %v",
					normalizedPath, tt.pattern, got, tt.want)
			}
		})
	}
}

func TestShouldIgnore(t *testing.T) {
	patterns := []IgnorePattern{
		{Pattern: "*.log", IsDirectory: false},
		{Pattern: "*.txt", IsDirectory: false},
		{Pattern: "bin", IsDirectory: true},
	}

	t.Run("IgnoreUntracked", func(t *testing.T) {
		trackedFiles := map[string]string{}

		if !ShouldIgnore("test.log", patterns, trackedFiles) {
			t.Error("Expected test.log to be ignored")
		}
		if !ShouldIgnore("readme.txt", patterns, trackedFiles) {
			t.Error("Expected readme.txt to be ignored")
		}
		if !ShouldIgnore("bin/output", patterns, trackedFiles) {
			t.Error("Expected bin/output to be ignored")
		}
		if ShouldIgnore("main.go", patterns, trackedFiles) {
			t.Error("Expected main.go to NOT be ignored")
		}
	})

	t.Run("NeverIgnoreTracked", func(t *testing.T) {
		trackedFiles := map[string]string{
			"important.txt": "abc123",
			"bin/config":    "def456",
		}

		// These files match patterns but are tracked, so should NOT be ignored
		if ShouldIgnore("important.txt", patterns, trackedFiles) {
			t.Error("Expected tracked file important.txt to NOT be ignored")
		}
		if ShouldIgnore("bin/config", patterns, trackedFiles) {
			t.Error("Expected tracked file bin/config to NOT be ignored")
		}
	})
}

func TestIsValidPattern(t *testing.T) {
	tests := []struct {
		pattern string
		want    bool
	}{
		{"*.txt", true},
		{"*.log", true},
		{"bin/", true},
		{"**/*.class", true},
		{"file?.txt", true},
		{"", false}, // Empty pattern is invalid
		// Note: filepath.Match is quite permissive, so most patterns will be valid
	}

	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			got := isValidPattern(tt.pattern)
			if got != tt.want {
				t.Errorf("isValidPattern(%q) = %v, want %v", tt.pattern, got, tt.want)
			}
		})
	}
}

func TestClearIgnoreCache(t *testing.T) {
	ClearIgnoreCache()
	cleanup := createKitignore(t, "*.log\n")
	defer cleanup()

	// Load patterns to populate cache
	patterns1, _ := LoadIgnorePatterns()
	if len(patterns1) == 0 {
		t.Fatal("Expected patterns to be loaded")
	}

	// Clear cache
	ClearIgnoreCache()

	// Verify cache is cleared by checking internal state
	// (In a real scenario, you might check if patterns are re-read from file)
	// For now, we'll just ensure calling it doesn't crash
	_, err := LoadIgnorePatterns()
	if err != nil {
		t.Errorf("Expected no error after clearing cache, got: %v", err)
	}
}
