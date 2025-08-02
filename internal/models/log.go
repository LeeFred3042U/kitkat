package models

import "time"

// LogEntry represents a single user log event.
// Used for productivity tracking, tagging, filtering and so on....
type LogEntry struct {
	ID        string    `json:"id"`
	Tag       string    `json:"tag"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}