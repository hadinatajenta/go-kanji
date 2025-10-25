package delivery

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"

	"gobackend/src/auth/dto"
	authinterfaces "gobackend/src/auth/interfaces"
	authservice "gobackend/src/auth/service"
	"gobackend/src/auth/validation"
)

const (
	stateCookieName = "oauthstate"
	stateTTL        = 10 * time.Minute
)

// Handler wires HTTP requests to the auth service layer.
type Handler struct {
	service            authinterfaces.AuthService
	successRedirectURL string
	failureRedirectURL string
}

// NewHandler instantiates an auth HTTP handler.
func NewHandler(service authinterfaces.AuthService, successRedirectURL, failureRedirectURL string) *Handler {
	return &Handler{
		service:            service,
		successRedirectURL: successRedirectURL,
		failureRedirectURL: failureRedirectURL,
	}
}

// GoogleLogin initiates the OAuth2 login by redirecting to Google.
func (h *Handler) GoogleLogin(ctx *gin.Context) {
	state, err := generateState()
	if err != nil {
		writeError(ctx, http.StatusInternalServerError, "failed to generate oauth state")
		return
	}

	ctx.SetCookie(
		stateCookieName,
		state,
		int(stateTTL.Seconds()),
		"/",
		"",
		ctx.Request.TLS != nil,
		true,
	)

	loginURL := h.service.GetGoogleLoginURL(state)
	ctx.Redirect(http.StatusTemporaryRedirect, loginURL)
}

// GoogleCallback handles Google's OAuth2 callback.
func (h *Handler) GoogleCallback(ctx *gin.Context) {
	query := ctx.Request.URL.Query()
	req := dto.GoogleCallbackRequest{
		Code:  query.Get("code"),
		State: query.Get("state"),
	}

	if err := validation.ValidateGoogleCallback(req); err != nil {
		writeValidationError(ctx, err)
		return
	}

	if cookie, err := ctx.Request.Cookie(stateCookieName); err == nil && cookie.Value != "" {
		if cookie.Value != req.State {
			writeError(ctx, http.StatusBadRequest, "state mismatch")
			return
		}

		ctx.SetCookie(stateCookieName, "", -1, "/", "", false, true)
	}

	result, err := h.service.HandleGoogleCallback(ctx.Request.Context(), req)
	if err != nil {
		if errors.Is(err, authservice.ErrUnauthorized) {
			h.handleUnauthorized(ctx)
			return
		}

		writeError(ctx, http.StatusBadGateway, err.Error())
		return
	}

	if h.successRedirectURL != "" {
		redirectURL, parseErr := url.Parse(h.successRedirectURL)
		if parseErr != nil {
			writeError(ctx, http.StatusInternalServerError, "invalid success redirect url")
			return
		}

		query := redirectURL.Query()
		query.Set("token", result.Token)
		query.Set("email", result.User.Email)
		query.Set("name", result.User.Name)
		if result.User.PictureURL != "" {
			query.Set("picture", result.User.PictureURL)
		}
		redirectURL.RawQuery = query.Encode()

		ctx.Redirect(http.StatusTemporaryRedirect, redirectURL.String())
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func generateState() (string, error) {
	random := make([]byte, 32)
	if _, err := rand.Read(random); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(random), nil
}

func writeValidationError(ctx *gin.Context, err error) {
	status := http.StatusBadRequest
	switch {
	case errors.Is(err, validation.ErrMissingCode):
		status = http.StatusBadRequest
	}

	writeError(ctx, status, err.Error())
}

func writeError(ctx *gin.Context, status int, message string) {
	ctx.JSON(status, gin.H{"error": message})
}

func (h *Handler) handleUnauthorized(ctx *gin.Context) {
	ctx.SetCookie(stateCookieName, "", -1, "/", "", false, true)

	if h.failureRedirectURL != "" {
		redirectURL, err := url.Parse(h.failureRedirectURL)
		if err == nil {
			query := redirectURL.Query()
			query.Set("error", authservice.ErrUnauthorized.Error())
			redirectURL.RawQuery = query.Encode()
			ctx.Redirect(http.StatusTemporaryRedirect, redirectURL.String())
			return
		}
	}

	writeError(ctx, http.StatusUnauthorized, authservice.ErrUnauthorized.Error())
}
