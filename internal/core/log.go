package core

import (
	"time"
	"github.com/google/uuid"
	"github.com/LeeFred3042U/kitkat/internal/models"
	"github.com/LeeFred3042U/kitkat/internal/storage"
)

func LogMessage(msg string) error {
	entry := models.LogEntry{
		ID:        uuid.New().String(),
		Message:   msg,
		Timestamp: time.Now(),
	}
	return storage.AppendLog(entry)
}