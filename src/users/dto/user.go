package dto

import "time"

// User represents the payload returned to API consumers.
type User struct {
	Reference   string    `json:"reference"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	PictureURL  string    `json:"picture_url,omitempty"`
	LastLoginAt time.Time `json:"last_login_at"`
	CreatedAt   time.Time `json:"created_at"`
	Provider    string    `json:"provider"`
}
