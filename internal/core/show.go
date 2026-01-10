package core

import (
	"fmt"

	"github.com/LeeFred3042U/kitcat/internal/storage"
)

// Displays the contents of a kitkat object
func ShowObject(hash string) error {
	data, err := storage.ReadObject(hash)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
