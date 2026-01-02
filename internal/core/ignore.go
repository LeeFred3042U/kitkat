package core

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// IgnorePattern represents a single pattern from .kitignore
type IgnorePattern struct {
	Original    string // The original pattern line from .kitignore
	Pattern     string // The processed pattern (without comments/whitespace)
	IsDirectory bool   // True if pattern ends with '/' (directory-only pattern)
	LineNumber  int    // Line number in .kitignore for error reporting
}

// Global cache for ignore patterns
var (
	ignoreCache     []IgnorePattern
	ignoreCacheMu   sync.RWMutex
	ignoreCacheInit bool
)

// LoadIgnorePatterns reads and parses the .kitignore file
// Returns an empty slice if .kitignore doesn't exist (not an error)
// Skips invalid patterns with a warning to stderr
func LoadIgnorePatterns() ([]IgnorePattern, error) {
	// Check cache first
	ignoreCacheMu.RLock()
	if ignoreCacheInit {
		patterns := ignoreCache
		ignoreCacheMu.RUnlock()
		return patterns, nil
	}
	ignoreCacheMu.RUnlock()

	// Acquire write lock to populate cache
	ignoreCacheMu.Lock()
	defer ignoreCacheMu.Unlock()

	// Double-check after acquiring write lock
	if ignoreCacheInit {
		return ignoreCache, nil
	}

	patterns := []IgnorePattern{}

	// Open .kitignore file
	file, err := os.Open(".kitignore")
	if err != nil {
		if os.IsNotExist(err) {
			// No .kitignore file is not an error, just return empty patterns
			ignoreCache = patterns
			ignoreCacheInit = true
			return patterns, nil
		}
		return nil, fmt.Errorf("error reading .kitignore: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// Trim whitespace
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check if this is a directory pattern (ends with /)
		isDirectory := strings.HasSuffix(line, "/")
		pattern := line

		// Remove trailing slash for processing, we'll handle it separately
		if isDirectory {
			pattern = strings.TrimSuffix(pattern, "/")
		}

		// Validate the pattern
		if !isValidPattern(pattern) {
			fmt.Fprintf(os.Stderr, "warning: .kitignore line %d: invalid pattern '%s' (skipping)\n", lineNumber, line)
			continue
		}

		patterns = append(patterns, IgnorePattern{
			Original:    line,
			Pattern:     pattern,
			IsDirectory: isDirectory,
			LineNumber:  lineNumber,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading .kitignore: %w", err)
	}

	// Cache the results
	ignoreCache = patterns
	ignoreCacheInit = true

	return patterns, nil
}

// ShouldIgnore checks if a path should be ignored based on patterns
// Returns false if the path is already tracked (tracked files are never ignored)
// Returns true if the path matches any ignore pattern
func ShouldIgnore(path string, patterns []IgnorePattern, trackedFiles map[string]string) bool {
	// Already tracked files are never ignored
	if _, isTracked := trackedFiles[path]; isTracked {
		return false
	}

	// Check if path matches any pattern
	for _, pattern := range patterns {
		if matchesPattern(path, pattern) {
			return true
		}
	}

	return false
}

// matchesPattern checks if a path matches a specific ignore pattern
// Handles glob patterns, directory patterns, and recursive patterns (**)
func matchesPattern(path string, pattern IgnorePattern) bool {
	// Normalize path separators for cross-platform compatibility
	path = filepath.ToSlash(path)
	patternStr := filepath.ToSlash(pattern.Pattern)

	// If it's a directory pattern, check if path is under that directory
	if pattern.IsDirectory {
		// Match the directory itself or any files/subdirectories under it
		if path == patternStr {
			return true
		}
		if strings.HasPrefix(path, patternStr+"/") {
			return true
		}
		return false
	}

	// Handle ** (recursive directory matching)
	if strings.Contains(patternStr, "**") {
		return matchesRecursivePattern(path, patternStr)
	}

	// Basic glob matching using filepath.Match
	matched, err := filepath.Match(patternStr, filepath.Base(path))
	if err != nil {
		// This shouldn't happen if pattern validation worked, but handle it anyway
		return false
	}

	if matched {
		return true
	}

	// Also try matching the full path for patterns that might include directories
	matched, err = filepath.Match(patternStr, path)
	if err != nil {
		return false
	}

	return matched
}

// matchesRecursivePattern handles patterns with ** (matches any number of directories)
func matchesRecursivePattern(path, pattern string) bool {
	// Split pattern by **
	parts := strings.Split(pattern, "**")

	// Simple case: pattern is just "**" (matches everything)
	if len(parts) == 2 && parts[0] == "" && parts[1] == "" {
		return true
	}

	// Handle pattern like "**/*.txt" or "**/temp"
	if len(parts) == 2 {
		prefix := strings.TrimSuffix(parts[0], "/")
		suffix := strings.TrimPrefix(parts[1], "/")

		// Check prefix
		if prefix != "" && !strings.HasPrefix(path, prefix) {
			return false
		}

		// Check suffix
		if suffix != "" {
			// For suffix matching, check if any part of the path matches
			pathParts := strings.Split(path, "/")
			for i := range pathParts {
				subPath := strings.Join(pathParts[i:], "/")
				matched, err := filepath.Match(suffix, subPath)
				if err == nil && matched {
					return true
				}
				// Also try matching just the file name
				matched, err = filepath.Match(suffix, pathParts[i])
				if err == nil && matched {
					return true
				}
			}
			return false
		}

		return true
	}

	// For more complex patterns, fall back to simpler matching
	// This is a simplified implementation; a full implementation would be more complex
	return false
}

// isValidPattern validates a glob pattern
// Returns true if the pattern is valid, false otherwise
func isValidPattern(pattern string) bool {
	// Empty pattern is invalid
	if pattern == "" {
		return false
	}

	// Try to match against a dummy path to validate syntax
	_, err := filepath.Match(pattern, "test")
	if err != nil {
		return false
	}

	// Additional validation: check for invalid escape sequences on Windows
	// filepath.Match should catch most issues, but we can add more checks if needed

	return true
}

// ClearIgnoreCache clears the cached ignore patterns
// This is useful for testing or when .kitignore is modified
func ClearIgnoreCache() {
	ignoreCacheMu.Lock()
	defer ignoreCacheMu.Unlock()
	ignoreCache = nil
	ignoreCacheInit = false
}
