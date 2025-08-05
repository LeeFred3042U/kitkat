package core

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/LeeFred3042U/kitkat/internal/models"
)

// GrepLogs prints logs
func GrepLogs(term string) error {
	data, err := os.ReadFile(".kitkat/logs.json")
	if err != nil {
		return err
	}

	var logs []models.LogEntry
	if err := json.Unmarshal(data, &logs); err != nil {
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