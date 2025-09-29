package models

import "time"

type Commit struct {
	ID        		string
	Parent    		string
	Message   		string
	Timestamp 		time.Time
	TreeHash  		string
	AuthorName 		string
	AuthorEmail		string
}
