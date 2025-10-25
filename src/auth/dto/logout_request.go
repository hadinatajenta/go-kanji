package dto

// LogoutRequest represents the payload to record a logout event.
type LogoutRequest struct {
	Detail string `json:"detail"`
}
