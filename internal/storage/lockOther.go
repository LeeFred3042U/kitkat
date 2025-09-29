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

func unlock(f *os.File) {
	f.Close()
}