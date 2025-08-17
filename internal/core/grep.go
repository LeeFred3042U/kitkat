package core

import (
	"fmt"
	"strings"

	"github.com/LeeFred3042U/kitkat/internal/storage"
)

// Prints logs
func GrepLogs(term string) error {
	logs, err := storage.ReadLogs()
	if err != nil {
		return err
	}

	term = strings.ToLower(term)
	for _, entry := range logs {
		msg := strings.ToLower(entry.Message)
		tag := strings.ToLower(entry.Tag)

		if strings.Contains(msg, term) || strings.Contains(tag, term) {
			fmt.Printf("[%s] (%s) %s\n", entry.ID, entry.Tag, entry.Message)
		}
	}
	return nil
}