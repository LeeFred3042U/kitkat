//go:build !linux && !darwin && !freebsd

package storage

import "os"

func lock(path string) (*os.File, error) {
	lockFile := path + ".lock"
	f, err := os.OpenFile(lockFile, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// LockFile is a no-op on non-Unix systems (unless we add valid windows locking later)
func LockFile(f *os.File) error {
	return nil
}

func unlock(f *os.File) {
	f.Close()
}
