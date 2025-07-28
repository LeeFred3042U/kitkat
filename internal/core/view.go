package core

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/LeeFred3042U/kitkat/internal/models"
)

// ViewLogs loads logs.json and prints all entries using the custom format ^-^
// Logs are shown in the order stored 
// No sorting or filtering yet
func ViewLogs() error {
	data, err := os.ReadFile(".kitkat/logs.json")
	if err != nil {
		return err
	}

	var logs []models.LogEntry
	if err := json.Unmarshal(data, &logs); err != nil {
		return err
	}

	for _, log := range logs {
		printFormattedLog(log)
	}

	return nil
}

// printFormattedLog renders a single log entry in a readable, styled format somehat!
// Timestamps are local 
// tag is hidden if empty.
func printFormattedLog(entry models.LogEntry) {
	ts := entry.Timestamp.Format("2006-01-02 15:04")

	// Omit tag display if not set
	if entry.Tag != "" {
		fmt.Printf("🔹 [%s] (%s) %s\n", ts, entry.Tag, entry.Message)
	} else {
		fmt.Printf("🔹 [%s] %s\n", ts, entry.Message)
	}

	fmt.Printf("    ID: %s\n\n", entry.ID)
}