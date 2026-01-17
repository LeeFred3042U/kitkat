package testutil

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const defaultBranch = "main"

// SetupTestRepo creates a temporary repository ready for tests and returns its path plus a cleanup func.
func SetupTestRepo(t *testing.T) (string, func()) {
	t.Helper()

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	repoDir := t.TempDir()
	if err := os.Chdir(repoDir); err != nil {
		t.Fatalf("failed to chdir to temp repo: %v", err)
	}

	cleanup := func() {
		_ = os.Chdir(cwd)
	}

	dirs := []string{
		".kitcat",
		".kitcat/objects",
		".kitcat/refs/heads",
		".kitcat/refs/tags",
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			cleanup()
			t.Fatalf("failed to create dir %s: %v", dir, err)
		}
	}

	files := []string{".kitcat/index", ".kitcat/commits.log"}
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			f, err := os.Create(file)
			if err != nil {
				cleanup()
				t.Fatalf("failed to create %s: %v", file, err)
			}
			_ = f.Close()
		} else if err != nil {
			cleanup()
			t.Fatalf("failed to stat %s: %v", file, err)
		}
	}

	headPath := filepath.Join(".kitcat", "HEAD")
	headContent := []byte("ref: refs/heads/" + defaultBranch + "\n")
	if err := os.WriteFile(headPath, headContent, 0o644); err != nil {
		cleanup()
		t.Fatalf("failed to write HEAD: %v", err)
	}

	mainBranchPath := filepath.Join(".kitcat", "refs", "heads", defaultBranch)
	if err := os.MkdirAll(filepath.Dir(mainBranchPath), 0o755); err != nil {
		cleanup()
		t.Fatalf("failed to ensure refs dir: %v", err)
	}
	if _, err := os.Stat(mainBranchPath); os.IsNotExist(err) {
		if err := os.WriteFile(mainBranchPath, []byte(""), 0o644); err != nil {
			cleanup()
			t.Fatalf("failed to create default branch ref: %v", err)
		}
	} else if err != nil {
		cleanup()
		t.Fatalf("failed to stat default branch ref: %v", err)
	}

	if err := setTestConfig("user.name", "Test User"); err != nil {
		cleanup()
		t.Fatalf("failed to set user.name: %v", err)
	}
	if err := setTestConfig("user.email", "test@example.com"); err != nil {
		cleanup()
		t.Fatalf("failed to set user.email: %v", err)
	}

	return repoDir, cleanup
}

func setTestConfig(key, value string) error {
	configPath, err := userConfigPath()
	if err != nil {
		return err
	}

	config, err := readConfigFile(configPath)
	if err != nil {
		return err
	}

	config[key] = value

	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		return err
	}

	f, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer f.Close()

	for k, v := range config {
		if _, err := fmt.Fprintf(f, "%s = %s\n", k, v); err != nil {
			return err
		}
	}

	return nil
}

func readConfigFile(path string) (map[string]string, error) {
	config := make(map[string]string)

	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return config, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			config[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	return config, scanner.Err()
}

func userConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".kitcatconfig"), nil
}
