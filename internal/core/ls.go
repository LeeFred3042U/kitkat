package core

import (
	"fmt"

	"github.com/LeeFred3042U/kitcat/internal/storage"
)

// ListFiles prints all tracked file paths from the index
func ListFiles() error {
	index, err := storage.LoadIndex()
	if err != nil {
		return err
	}

	for path := range index {
		fmt.Println(path)
	}
	return nil
}
