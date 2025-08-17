package core

import (
	"fmt"

	"github.com/LeeFred3042U/kitkat/internal/models"
	"github.com/LeeFred3042U/kitkat/internal/storage"
)

// View Logs
func ViewLogs() error {
	logs, err := storage.ReadLogs()
	if err != nil {
		return err
	}

	for _, log := range logs {
		printFormattedLog(log)
	}

	return nil
}

// Renders a single log entry in a readable, styled format somewhat!
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