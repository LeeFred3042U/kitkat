package models

import "time"


type Commit struct {
	ID        string
	Message   string
	Timestamp time.Time
	TreeHash  string
}