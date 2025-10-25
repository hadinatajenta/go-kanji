package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"gobackend/src/auth/dao"
	"gobackend/src/auth/dto"
	authinterfaces "gobackend/src/auth/interfaces"
	logdto "gobackend/src/logs/dto"
	loginterfaces "gobackend/src/logs/interfaces"
)

const (
	googleProvider     = "google"
	userInfoEndpoint   = "https://www.googleapis.com/oauth2/v2/userinfo"
	defaultHTTPTimeout = 10 * time.Second
	allowedDisplayName = "Hadinata Jenta"
)

var (
	// ErrInvalidConfig signals that the service was initialised with an incomplete configuration.
	ErrInvalidConfig = errors.New("invalid google oauth config")
	// ErrUnauthorized indicates the authenticated Google account is not permitted.
	ErrUnauthorized = errors.New("unauthorize")
)

// GoogleAuthConfig contains all configuration values required by GoogleAuthService.
type GoogleAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	JWTSecret    string
	TokenTTL     time.Duration
	HTTPClient   *http.Client
	LogService   loginterfaces.Service
}

type googleUserInfo struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	GivenName string `json:"given_name"`
	Picture   string `json:"picture"`
}

type authClaims struct {
	Email    string `json:"email"`
	Provider string `json:"provider"`
	jwt.RegisteredClaims
}

// GoogleAuthService implements the OAuth2 flow against Google.
type GoogleAuthService struct {
	repo        authinterfaces.UserRepository
	oauthConfig *oauth2.Config
	jwtSecret   []byte
	tokenTTL    time.Duration
	httpClient  *http.Client
	logService  loginterfaces.Service
}

var _ authinterfaces.AuthService = (*GoogleAuthService)(nil)

// NewGoogleAuthService constructs a new GoogleAuthService.
func NewGoogleAuthService(repo authinterfaces.UserRepository, cfg GoogleAuthConfig) (*GoogleAuthService, error) {
	if repo == nil ||
		cfg.ClientID == "" ||
		cfg.ClientSecret == "" ||
		cfg.RedirectURL == "" ||
		cfg.JWTSecret == "" ||
		cfg.LogService == nil {
		return nil, ErrInvalidConfig
	}

	if len(cfg.Scopes) == 0 {
		cfg.Scopes = []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		}
	}

	if cfg.TokenTTL <= 0 {
		cfg.TokenTTL = time.Hour
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultHTTPTimeout}
	}

	oauthConfig := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes:       cfg.Scopes,
		Endpoint:     google.Endpoint,
	}

	return &GoogleAuthService{
		repo:        repo,
		oauthConfig: oauthConfig,
		jwtSecret:   []byte(cfg.JWTSecret),
		tokenTTL:    cfg.TokenTTL,
		httpClient:  httpClient,
		logService:  cfg.LogService,
	}, nil
}

// GetGoogleLoginURL produces the Google consent screen URL.
func (s *GoogleAuthService) GetGoogleLoginURL(state string) string {
	return s.oauthConfig.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "select_account"),
		oauth2.SetAuthURLParam("include_granted_scopes", "true"),
	)
}

// HandleGoogleCallback completes the OAuth2 flow once Google redirects back to the application.
func (s *GoogleAuthService) HandleGoogleCallback(ctx context.Context, req dto.GoogleCallbackRequest) (*dto.AuthResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultHTTPTimeout)
	defer cancel()

	ctx = context.WithValue(ctx, oauth2.HTTPClient, s.httpClient)

	token, err := s.oauthConfig.Exchange(ctx, req.Code)
	if err != nil {
		return nil, fmt.Errorf("exchange authorization code: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, userInfoEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("build user info request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := s.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("fetch google user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google userinfo responded with status %d", resp.StatusCode)
	}

	var info googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("decode google user info: %w", err)
	}

	user, err := s.ensureUser(ctx, info)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateLoginTimestamp(ctx, user.ID); err != nil {
		return nil, fmt.Errorf("update login timestamp: %w", err)
	}

	user.LastLoginAt = time.Now()

	if err := s.logService.Record(ctx, logdto.NewLog{
		UserID: user.ID,
		Action: "login",
		Detail: "authenticated via Google OAuth",
	}); err != nil {
		return nil, fmt.Errorf("record login log: %w", err)
	}

	tokenString, err := s.generateJWT(*user)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token: tokenString,
		User:  *user,
	}, nil
}

func (s *GoogleAuthService) ensureUser(ctx context.Context, info googleUserInfo) (*dao.User, error) {
	displayName := info.Name
	if displayName == "" {
		displayName = info.GivenName
	}

	if displayName == "" {
		displayName = info.Email
	}

	if displayName != allowedDisplayName {
		return nil, ErrUnauthorized
	}

	existing, err := s.repo.FindByProvider(ctx, googleProvider, info.ID)
	if err != nil {
		return nil, fmt.Errorf("find user by provider: %w", err)
	}

	if existing != nil {
		if existing.Name != allowedDisplayName {
			return nil, ErrUnauthorized
		}
		// Keep latest picture if previously missing.
		if existing.PictureURL == "" && info.Picture != "" {
			existing.PictureURL = info.Picture
		}
		return existing, nil
	}

	newUser := dao.User{
		Email:      info.Email,
		Name:       displayName,
		Provider:   googleProvider,
		ProviderID: info.ID,
		PictureURL: info.Picture,
	}

	created, err := s.repo.Create(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return created, nil
}

func (s *GoogleAuthService) generateJWT(user dao.User) (string, error) {
	now := time.Now()
	claims := authClaims{
		Email:    user.Email,
		Provider: user.Provider,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(user.ID, 10),
			Issuer:    "gobackend",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("sign jwt token: %w", err)
	}

	return signed, nil
}

// ExtractUserID parses the JWT token and returns the embedded user ID.
func (s *GoogleAuthService) ExtractUserID(token string) (int64, error) {
	if token == "" {
		return 0, fmt.Errorf("token is required")
	}

	claims := &authClaims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil {
		return 0, fmt.Errorf("parse token: %w", err)
	}

	if !parsed.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("decode subject: %w", err)
	}

	return userID, nil
}
