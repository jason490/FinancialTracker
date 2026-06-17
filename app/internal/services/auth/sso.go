package auth

import (
	"FinancialTracker/internal/models"
	"FinancialTracker/internal/storage"
	"FinancialTracker/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

const apiStatePrefix = "api:"
const defaultAPIReturnTo = "/auth/sso/complete"

type SSOService struct {
	store *storage.Storage
}

func NewSSOService(store *storage.Storage) *SSOService {
	return &SSOService{
		store: store,
	}
}

// GoogleUserInfo represents data returned from Google UserInfo API
type GoogleUserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// HandleReauth handles re-authentication for sensitive actions
func (s *SSOService) HandleReauth(userID int64, sessionID string, info GoogleUserInfo) error {
	user, err := s.store.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	if user.Email != info.Email {
		return errors.New("identity mismatch")
	}

	return s.store.UpdateSessionReauth(sessionID, time.Now().Unix())
}

// LinkGoogleAccount links a Google SSO account to an existing user
func (s *SSOService) LinkGoogleAccount(userID int64, info GoogleUserInfo) error {
	existing, _ := s.store.GetUserBySSO("google", info.ID)
	if existing != nil {
		return errors.New("SSO account already linked to another user")
	}

	return s.store.LinkSSO(userID, "google", info.ID)
}

// AuthenticateViaGoogle handles login/registration via Google SSO
func (s *SSOService) AuthenticateViaGoogle(info GoogleUserInfo) (*models.Session, error) {
	// 1. Try finding user by this SSO
	user, err := s.store.GetUserBySSO("google", info.ID)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// 2. If not found, try finding by email and link
	if user == nil {
		user, err = s.store.GetUserByEmail(info.Email)
		if err != nil {
			return nil, fmt.Errorf("database error: %w", err)
		}

		if user != nil {
			if err := s.store.LinkSSO(user.ID, "google", info.ID); err != nil {
				return nil, fmt.Errorf("failed to link SSO: %w", err)
			}
		}
	}

	// 3. If still not found, create new user
	if user == nil {
		parts := strings.Split(info.Name, " ")
		firstName := parts[0]
		lastName := ""
		if len(parts) > 1 {
			lastName = strings.Join(parts[1:], " ")
		}

		user = &models.User{
			Email:           info.Email,
			FirstName:       firstName,
			LastName:        lastName,
			ThemePreference: "system",
		}
		if err := s.store.CreateUser(user); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
		if err := s.store.LinkSSO(user.ID, "google", info.ID); err != nil {
			return nil, fmt.Errorf("failed to link SSO after registration: %w", err)
		}
	}

	// Create session
	sessionID := uuid.New().String()
	session := &models.Session{
		ID:                sessionID,
		UserID:            user.ID,
		ExpiresAt:         time.Now().AddDate(0, 0, 1).Unix(),
		ReauthenticatedAt: time.Now().Unix(),
	}
	if err := s.store.CreateSession(session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

// CompleteGoogleLogin exchanges an OAuth code and returns a new session.
func (s *SSOService) CompleteGoogleLogin(ctx context.Context, config *oauth2.Config, code string) (*models.Session, error) {
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, errors.New("authentication failed")
	}

	userInfo, err := s.FetchGoogleUserInfo(token.AccessToken)
	if err != nil {
		return nil, err
	}

	return s.AuthenticateViaGoogle(userInfo)
}

// FetchGoogleUserInfo retrieves profile data from Google's userinfo endpoint.
func (s *SSOService) FetchGoogleUserInfo(accessToken string) (GoogleUserInfo, error) {
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return GoogleUserInfo{}, errors.New("authentication failed")
	}
	defer response.Body.Close()

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		return GoogleUserInfo{}, errors.New("authentication failed")
	}

	return userInfo, nil
}

// GetUserForSession returns the user associated with a session ID.
func (s *SSOService) GetUserForSession(sessionID string) (*models.User, error) {
	session, err := s.store.GetSession(sessionID)
	if err != nil || session == nil {
		return nil, errors.New("invalid or expired session")
	}

	user, err := s.store.GetUserByID(session.UserID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("failed to load user profile: %w", err)
	}

	return user, nil
}

// APIOAuthState contains the decoded OAuth state for API clients.
type APIOAuthState struct {
	ReturnTo string
	Mode     string
}

// BuildAPIOAuthState encodes the frontend return URL and OAuth mode into a secure state value.
func (s *SSOService) BuildAPIOAuthState(returnTo, mode string) string {
	nonce := uuid.New().String()
	if mode == "" {
		mode = "login"
	}
	payload := fmt.Sprintf("%s|%s|%d|%s", returnTo, nonce, time.Now().Unix(), mode)

	encrypted, err := utils.Encrypt(payload)
	if err != nil {
		return apiStatePrefix + "error"
	}

	return apiStatePrefix + encrypted
}

// ParseAPIState decodes and verifies the frontend OAuth state for API clients.
func (s *SSOService) ParseAPIState(state string) (*APIOAuthState, error) {
	if !strings.HasPrefix(state, apiStatePrefix) {
		return nil, ErrInvalidOAuthState
	}

	encryptedPayload := strings.TrimPrefix(state, apiStatePrefix)
	payload, err := utils.Decrypt(encryptedPayload)
	if err != nil {
		return nil, ErrInvalidOAuthState
	}

	parts := strings.Split(payload, "|")
	if len(parts) < 3 {
		return nil, ErrInvalidOAuthState
	}

	returnTo := parts[0]
	mode := "login"
	if len(parts) >= 4 && parts[3] != "" {
		mode = parts[3]
	}

	if returnTo == "" {
		returnTo = defaultAPIReturnTo
	}

	return &APIOAuthState{
		ReturnTo: returnTo,
		Mode:     mode,
	}, nil
}

// ParseAPIReturnTo decodes and verifies the frontend return URL from an encrypted OAuth state value.
func (s *SSOService) ParseAPIReturnTo(state string) (string, error) {
	oauthState, err := s.ParseAPIState(state)
	if err != nil {
		return "", err
	}
	return oauthState.ReturnTo, nil
}

// APICallbackResult contains the redirect target and session info after Google SSO.
type APICallbackResult struct {
	ReturnTo  string
	SessionID string
}

// CompleteGoogleLink exchanges a Google OAuth code and links it to an existing user.
func (s *SSOService) CompleteGoogleLink(ctx context.Context, config *oauth2.Config, code string, userID int64) error {
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return errors.New("authentication failed")
	}

	userInfo, err := s.FetchGoogleUserInfo(token.AccessToken)
	if err != nil {
		return err
	}

	return s.LinkGoogleAccount(userID, userInfo)
}

// CompleteAPICallback exchanges a Google OAuth code and issues session info.
func (s *SSOService) CompleteAPICallback(
	ctx context.Context,
	config *oauth2.Config,
	code string,
	state string,
) (*APICallbackResult, error) {
	oauthState, err := s.ParseAPIState(state)
	if err != nil {
		return nil, err
	}

	session, err := s.CompleteGoogleLogin(ctx, config, code)
	if err != nil {
		return &APICallbackResult{ReturnTo: oauthState.ReturnTo}, err
	}

	return &APICallbackResult{
		ReturnTo:  oauthState.ReturnTo,
		SessionID: session.ID,
	}, nil
}
