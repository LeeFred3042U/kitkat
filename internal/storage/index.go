package storage

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

const indexPath = ".kitkat/index"

// small helper lock using a .lock file (Unix-only)
func lock(path string) (*os.File, error) {
	lockFile := path + ".lock"
	f, err := os.OpenFile(lockFile, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		f.Close()
		return nil, err
	}
	return f, nil
}

func unlock(f *os.File) {
	_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	f.Close()
}

// LoadIndex reads .kitkat/index and returns a map of path -> hash
// If file doesn't exist return empty map (convenience for fresh repo)
func LoadIndex() (map[string]string, error) {
	index := make(map[string]string)

	// tolerate missing index (fresh repo)
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return index, nil
	}

	f, err := os.Open(indexPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 {
			index[parts[0]] = parts[1]
		}
	}
	return index, scanner.Err()
}

// WriteIndex writes index atomically (temp-file + rename) while holding a lock
func WriteIndex(index map[string]string) error {
	// ensure parent dir
	if err := os.MkdirAll(filepath.Dir(indexPath), 0755); err != nil {
		return err
	}

	l, err := lock(indexPath)
	if err != nil {
		return err
	}
	defer unlock(l)

	tmp := indexPath + ".tmp"
	var b []byte
	for path, hash := range index {
		b = append(b, []byte(fmt.Sprintf("%s %s\n", path, hash))...)
	}
	if err := ioutil.WriteFile(tmp, b, 0644); err != nil {
		return err
	}
	// rename is atomic on same FS
	if err := os.Rename(tmp, indexPath); err != nil {
		return err
	}
	return nil
}
