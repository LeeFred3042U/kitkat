package core

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/LeeFred3042U/kitkat/internal/models"
)

// List Logs By Tag 
func ListLogsByTag(tag string) error {
	data, err := os.ReadFile(".kitkat/logs.json")
	if err != nil {
		return err
	}

	var logs []models.LogEntry
	if err := json.Unmarshal(data, &logs); err != nil {
		return err
	}

	for _, entry := range logs {
		if entry.Tag == tag {
			fmt.Printf("[%s] (%s) %s\n", entry.ID, entry.Tag, entry.Message)
		}
	}

	return nil
}
