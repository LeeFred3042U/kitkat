package storage

import (
	"os"
	"fmt"
	"time"
	"bufio"
	"strings"

	"github.com/LeeFred3042U/kitkat/internal/models"
)

const logPath = ".kitkat/logs.txt"

// Appends a new log entry
func AppendLog(entry models.LogEntry) error {
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	timestamp := entry.Timestamp.Format(time.RFC3339)
	logLine := fmt.Sprintf("id:%s tag:%s message:\"%s\" timestamp:%s\n", entry.ID, entry.Tag, entry.Message, timestamp)

	if _, err := f.WriteString(logLine); err != nil {
		return err
	}

	return nil
}

// Reads all log entries
func ReadLogs() ([]models.LogEntry, error) {
	var logs []models.LogEntry
	f, err := os.Open(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return logs, nil
		}
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		entry := models.LogEntry{}
		for _, part := range parts {
			kv := strings.SplitN(part, ":", 2)
			if len(kv) != 2 {
				continue
			}
			key, value := kv[0], kv[1]
			switch key {
			case "id":
				entry.ID = value
			case "tag":
				entry.Tag = value
			case "message":
				entry.Message = strings.Trim(value, "\"")
			case "timestamp":
				entry.Timestamp, _ = time.Parse(time.RFC3339, value)
			}
		}
		logs = append(logs, entry)
	}
	return logs, scanner.Err()
}

// Writes a slice of log entries, does overwrite
func WriteLogs(logs []models.LogEntry) error {
	f, err := os.Create(logPath)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, entry := range logs {
		timestamp := entry.Timestamp.Format(time.RFC3339)
		logLine := fmt.Sprintf("id:%s tag:%s message:\"%s\" timestamp:%s\n", entry.ID, entry.Tag, entry.Message, timestamp)
		if _, err := f.WriteString(logLine); err != nil {
			return err
		}
	}
	return nil
}