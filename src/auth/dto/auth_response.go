package dto

import "gobackend/src/auth/dao"

// AuthResponse is delivered back to the client after a successful OAuth2 flow.
type AuthResponse struct {
	Token string   `json:"token"`
	User  dao.User `json:"user"`
}
