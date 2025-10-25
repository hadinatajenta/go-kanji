package delivery

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"gobackend/shared/response"
	"gobackend/src/auth/dto"
	authinterfaces "gobackend/src/auth/interfaces"
	authservice "gobackend/src/auth/service"
	"gobackend/src/auth/validation"
	logdto "gobackend/src/logs/dto"
	loginterfaces "gobackend/src/logs/interfaces"
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
	logService         loginterfaces.Service
}

// NewHandler instantiates an auth HTTP handler.
func NewHandler(service authinterfaces.AuthService, successRedirectURL, failureRedirectURL string, logService loginterfaces.Service) *Handler {
	return &Handler{
		service:            service,
		successRedirectURL: successRedirectURL,
		failureRedirectURL: failureRedirectURL,
		logService:         logService,
	}
}

// GoogleLogin initiates the OAuth2 login by redirecting to Google.
func (h *Handler) GoogleLogin(ctx *gin.Context) {
	state, err := generateState()
	if err != nil {
		response.InternalError(ctx, "failed to generate oauth state", err.Error())
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
		response.BadRequest(ctx, err.Error(), nil)
		return
	}

	if cookie, err := ctx.Request.Cookie(stateCookieName); err == nil && cookie.Value != "" {
		if cookie.Value != req.State {
			response.BadRequest(ctx, "state mismatch", nil)
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

		response.InternalError(ctx, "failed to complete login", err.Error())
		return
	}

	if h.successRedirectURL != "" {
		redirectURL, parseErr := url.Parse(h.successRedirectURL)
		if parseErr != nil {
			response.InternalError(ctx, "invalid success redirect url", parseErr.Error())
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

	response.OK(ctx, "login successful", result)
}

// Logout registers a logout activity in the audit logs.
func (h *Handler) Logout(ctx *gin.Context) {
	var req dto.LogoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "invalid payload", err.Error())
		return
	}

	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		response.Unauthorized(ctx, "missing authorization header")
		return
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
	if token == "" {
		response.Unauthorized(ctx, "invalid authorization header")
		return
	}

	userID, err := h.service.ExtractUserID(token)
	if err != nil {
		response.Unauthorized(ctx, "invalid or expired token")
		return
	}

	detail := req.Detail
	if detail == "" {
		detail = "user initiated logout"
	}

	entry := logdto.NewLog{
		UserID: userID,
		Action: "logout",
		Detail: detail,
	}

	if err := h.logService.Record(ctx.Request.Context(), entry); err != nil {
		response.InternalError(ctx, "failed to record logout", err.Error())
		return
	}

	response.OK(ctx, "logout recorded", gin.H{"status": "ok"})
}

func generateState() (string, error) {
	random := make([]byte, 32)
	if _, err := rand.Read(random); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(random), nil
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

	response.Unauthorized(ctx, authservice.ErrUnauthorized.Error())
}
