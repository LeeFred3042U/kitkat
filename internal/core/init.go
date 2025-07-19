package core

import (
	"os"
)

const (
	repoDir    = ".kitkat"
	indexPath  = ".kitkat/index"
	objectsDir = ".kitkat/objects"
)


func InitRepo() error {
	err := os.Mkdir(".kitkat", 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	dirs := []string{".kitkat/objects"}
	for _, dir := range dirs {
		if err := os.Mkdir(dir, 0755); err != nil && !os.IsExist(err) {
			return err
		}
	}

	files := []string{".kitkat/index", ".kitkat/commits.log"}
	for _, file := range files {
		f, err := os.Create(file)
		if err != nil {
			return err
		}
		f.Close()
	}

	return nil
}