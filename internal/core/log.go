package core

import (
	"crypto/sha1"
	"encoding/hex"
	"time"

	"github.com/LeeFred3042U/kitkat/internal/models"
	"github.com/LeeFred3042U/kitkat/internal/storage"
)

// Computes a SHA-1 over log content
func hashLog(l models.LogEntry) string {
	h := sha1.New()
	h.Write([]byte(l.Message))
	h.Write([]byte(l.Tag))
	h.Write([]byte(l.Timestamp.UTC().String()))
	return hex.EncodeToString(h.Sum(nil))
}

func LogMessage(msg string) error {
	entry := models.LogEntry{
		Message:   msg,
		Timestamp: time.Now().UTC(),
	}
	entry.ID = hashLog(entry)
	return storage.AppendLog(entry)
}