package dto

import "time"

// Log represents the log information exposed via the API.
type Log struct {
	UserName  string    `json:"user_name"`
	Action    string    `json:"action"`
	Detail    string    `json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}

// NewLog describes payload required to create a log entry.
type NewLog struct {
	UserID int64  `json:"user_id"`
	Action string `json:"action"`
	Detail string `json:"detail"`
}
