package core

import (
	"fmt"
	
	"github.com/LeeFred3042U/kitkat/internal/storage"
)

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