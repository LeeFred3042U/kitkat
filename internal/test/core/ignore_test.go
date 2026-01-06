package core_test

import (
	"testing"

	"github.com/LeeFred3042U/kitkat/internal/core"
)

func TestShouldIgnore(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		patterns     []core.IgnorePattern
		trackedFiles map[string]string // simulate tracked files
		want         bool
	}{
		// 1. Literal Matches
		{
			name: "Literal file match",
			path: "file.txt",
			patterns: []core.IgnorePattern{
				{Pattern: "file.txt", Original: "file.txt"},
			},
			want: true,
		},
		{
			name: "Literal file mismatch",
			path: "other.txt",
			patterns: []core.IgnorePattern{
				{Pattern: "file.txt", Original: "file.txt"},
			},
			want: false,
		},
		{
			name: "Literal file in subdir match",
			path: "subdir/file.txt",
			patterns: []core.IgnorePattern{
				{Pattern: "subdir/file.txt", Original: "subdir/file.txt"},
			},
			want: true,
		},

		// 2. Directory Matches
		{
			name: "Directory match",
			path: "build/output.log",
			patterns: []core.IgnorePattern{
				{Pattern: "build", Original: "build/", IsDirectory: true},
			},
			want: true,
		},
		{
			name: "Directory exact match",
			path: "build",
			patterns: []core.IgnorePattern{
				{Pattern: "build", Original: "build/", IsDirectory: true},
			},
			want: true,
		},
		{
			name: "Directory mismatch (matches prefix but not dir)",
			path: "builder.go",
			patterns: []core.IgnorePattern{
				{Pattern: "build", Original: "build/", IsDirectory: true},
			},
			want: false,
		},

		// 3. Wildcard Matches
		{
			name: "Extension wildcard",
			path: "error.log",
			patterns: []core.IgnorePattern{
				{Pattern: "*.log", Original: "*.log"},
			},
			want: true,
		},
		{
			name: "Extension wildcard in subdir",
			path: "logs/error.log",
			patterns: []core.IgnorePattern{
				{Pattern: "*.log", Original: "*.log"},
			},
			want: true,
		},
		{
			name: "Prefix wildcard",
			path: "temp_123",
			patterns: []core.IgnorePattern{
				{Pattern: "temp*", Original: "temp*"},
			},
			want: true,
		},
		{
			name: "Question mark wildcard",
			path: "image.jzg", // intentional typo to test ?
			patterns: []core.IgnorePattern{
				{Pattern: "image.j?g", Original: "image.j?g"},
			},
			want: true,
		},

		// 4. Recursive Wildcard (**)
		{
			name: "Recursive logs match",
			path: "logs/mw/error.log",
			patterns: []core.IgnorePattern{
				{Pattern: "logs/**/*.log", Original: "logs/**/*.log"},
			},
			want: true,
		},
		{
			name: "Double star prefix",
			path: "foo/bar/baz.txt",
			patterns: []core.IgnorePattern{
				{Pattern: "**/baz.txt", Original: "**/baz.txt"},
			},
			want: true,
		},
		{
			name: "Double star suffix",
			path: "foo/bar/node_modules/cache",
			patterns: []core.IgnorePattern{
				{Pattern: "foo/**", Original: "foo/**"},
			},
			want: true,
		},

		// 5. Tracked Files (Priority)
		{
			name: "Tracked file ignored in patterns but should not be ignored",
			path: "ignored.txt",
			patterns: []core.IgnorePattern{
				{Pattern: "ignored.txt", Original: "ignored.txt"},
			},
			trackedFiles: map[string]string{
				"ignored.txt": "hash",
			},
			want: false,
		},

		// 6. Edge Cases
		{
			name:     "Empty list of patterns",
			path:     "file.txt",
			patterns: []core.IgnorePattern{},
			want:     false,
		},
		{
			name: "Complex combination",
			path: "src/vendor/lib.o",
			patterns: []core.IgnorePattern{
				{Pattern: "*.o", Original: "*.o"},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := core.ShouldIgnore(tt.path, tt.patterns, tt.trackedFiles)
			if got != tt.want {
				t.Errorf("ShouldIgnore(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
