package storage

import (
	"bufio"
	"os"
	"strings"
)

const indexPath = ".kitkat/index"

// reads index into map[path]hash
func LoadIndex() (map[string]string, error) {
	index := make(map[string]string)

	f, err := os.Open(indexPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 {
			index[parts[0]] = parts[1]
		}
	}

	return index, scanner.Err()
}

// writes full index (de-duplicated)
func WriteIndex(index map[string]string) error {
	f, err := os.Create(indexPath) // overwrite
	if err != nil {
		return err
	}
	defer f.Close()

	for path, hash := range index {
		_, err := f.WriteString(path + " " + hash + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
