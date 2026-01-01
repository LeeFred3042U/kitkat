package core

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// getConfigPath returns the absolute path to the global kitkat config file
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".kitkatconfig"), nil
}

// readConfig loads the config file into a map
func readConfig() (map[string]string, error) {
	config := make(map[string]string)
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	// It's okay if the file doesn't exist yet, just return an empty map
	if os.IsNotExist(err) {
		return config, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// The file format is simple: key=value
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			config[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return config, scanner.Err()
}

// SetConfig sets a key-value pair in the global config file
func SetConfig(key, value string) error {
	config, err := readConfig()
	if err != nil {
		return fmt.Errorf("could not read existing config: %w", err)
	}

	// Update the map with the new key-value pair
	config[key] = value

	// Write the entire updated map back to the file
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	for k, v := range config {
		if _, err := fmt.Fprintf(file, "%s = %s\n", k, v); err != nil {
			return err
		}
	}

	return nil
}

// GetConfig reads a key from the global config file
func GetConfig(key string) (string, bool, error) {
	config, err := readConfig()
	if err != nil {
		return "", false, err
	}
	value, ok := config[key]
	return value, ok, nil
}

// PrintAllConfig prints all key-value pairs in the config file
func PrintAllConfig() error {
	config, err := readConfig()
	if err != nil {
		return err
	}
	for k, v := range config {
		fmt.Printf("%s = %s\n", k, v)
	}
	return nil
}
