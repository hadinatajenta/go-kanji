package dao

import "time"

// User represents a persisted user record.
type User struct {
	ID          int64     `json:"id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	Provider    string    `json:"provider"`
	ProviderID  string    `json:"provider_id"`
	PictureURL  string    `json:"picture_url"`
	CreatedAt   time.Time `json:"created_at"`
	LastLoginAt time.Time `json:"last_login_at"`
}
