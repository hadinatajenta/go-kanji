package validation

import (
	"errors"
	"strings"

	"gobackend/src/auth/dto"
)

var (
	// ErrMissingCode indicates the OAuth2 callback payload does not contain the authorization code.
	ErrMissingCode = errors.New("authorization code is required")
)

// ValidateGoogleCallback ensures the mandatory fields are present in the callback request.
func ValidateGoogleCallback(req dto.GoogleCallbackRequest) error {
	if strings.TrimSpace(req.Code) == "" {
		return ErrMissingCode
	}

	return nil
}
