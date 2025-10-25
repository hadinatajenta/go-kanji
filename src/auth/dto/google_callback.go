package dto

// GoogleCallbackRequest represents the data received from Google on the OAuth2 callback flow.
type GoogleCallbackRequest struct {
	Code  string
	State string
}
