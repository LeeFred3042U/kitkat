package storage

import (
	"encoding/json"
	"os"
	"github.com/LeeFred3042U/kitkat/internal/models"
)

const logPath = ".kitkat/logs.json"

func AppendLog(entry models.LogEntry) error {
	var logs []models.LogEntry

	// load existing
	data, err := os.ReadFile(logPath)
	if err == nil {
		_ = json.Unmarshal(data, &logs)
	}

	logs = append(logs, entry)

	// write updated
	out, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(logPath, out, 0644)
}