package core

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/LeeFred3042U/kitkat/internal/models"
)

// TagLog updates the tag field of a log entry when given it's ID
// format: ./kitkat tag log-id tag-name
func TagLog(id string, newTag string) error {
	data, err := os.ReadFile(".kitkat/logs.json") // read existing log file
	if err != nil {
		return err
	}

	var logs []models.LogEntry
	if err := json.Unmarshal(data, &logs); err != nil {
		return err
	}

	found := false
	for i := range logs {
		if logs[i].ID == id {
			logs[i].Tag = newTag // mutate tag field
			found = true
			break
		}
	}

	if !found {
		return errors.New("log ID not found")
	}

	updated, err := json.MarshalIndent(logs, "", "  ") // re-encode logs
	if err != nil {
		return err
	}

	return os.WriteFile(".kitkat/logs.json", updated, 0644) // overwrite the file
}