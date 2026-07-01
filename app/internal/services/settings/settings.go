package settings

import (
	"FinancialTracker/internal/models"
	"FinancialTracker/internal/models/external"
	"FinancialTracker/internal/services/auth"
	"FinancialTracker/internal/storage"
	"FinancialTracker/internal/utils"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const reauthTimeoutSeconds = 300

type SettingsService struct {
	store *storage.Storage
}

func NewSettingsService(store *storage.Storage) *SettingsService {
	return &SettingsService{
		store: store,
	}
}

// GetProfile returns the settings profile for an authenticated user.
func (s *SettingsService) GetProfile(userID int64) (*external.SettingsProfile, error) {
	user, err := s.store.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	providers := make([]string, 0, len(user.SSOs))
	for _, sso := range user.SSOs {
		providers = append(providers, sso.Provider)
	}

	return &external.SettingsProfile{
		FirstName:                user.FirstName,
		LastName:                 user.LastName,
		Email:                    user.Email,
		ThemePreference:          user.ThemePreference,
		HasPassword:              user.PasswordHash != "",
		SSOProviders:             providers,
		IsRegistrationAdmin:      auth.IsRegistrationAdmin(user.Email),
		RegistrationCodeRequired: auth.RegistrationGateEnabled(),
	}, nil
}

// UpdateNames updates the user's first and last name.
func (s *SettingsService) UpdateNames(userID int64, firstName, lastName string) error {
	user, err := s.store.GetUserByID(userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}
	return s.UpdateProfile(userID, firstName, lastName, user.Email)
}

// UnlinkSSO removes an SSO provider when the user retains another login method.
func (s *SettingsService) UnlinkSSO(userID int64, provider string) error {
	user, err := s.store.GetUserByID(userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	if user.PasswordHash == "" && len(user.SSOs) <= 1 {
		return errors.New("cannot remove your only login method")
	}

	return s.store.UnlinkSSO(userID, provider)
}

// UpdateTheme updates the user's theme preference
func (s *SettingsService) UpdateTheme(userID int64, theme string) error {
	if !models.IsValidThemePreference(theme) {
		return errors.New("invalid theme preference")
	}
	return s.store.UpdateUserTheme(userID, theme)
}

// UpdateProfile updates user profile information
func (s *SettingsService) UpdateProfile(userID int64, firstName, lastName, email string) error {
	firstName = utils.Sanitize(firstName)
	lastName = utils.Sanitize(lastName)
	email = utils.Sanitize(strings.ToLower(email))

	if msg, err := utils.ValidateName(firstName, "First name"); err {
		return errors.New(msg)
	}
	if msg, err := utils.ValidateName(lastName, "Last name"); err {
		return errors.New(msg)
	}
	if msg, err := utils.ValidateEmail(email); err {
		return errors.New(msg)
	}

	return s.store.UpdateUserInfo(userID, firstName, lastName, email)
}

// UpdatePassword verifies current password and updates to new password
func (s *SettingsService) UpdatePassword(userID int64, currentPassword, newPassword, confirmPassword string) error {
	user, err := s.store.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if user.PasswordHash != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
			return errors.New("current password incorrect")
		}
	}

	if msg, err := utils.ValidatePassword(newPassword); err {
		return errors.New(msg)
	}

	if newPassword != confirmPassword {
		return errors.New("new passwords do not match")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	return s.store.UpdateUserPassword(userID, string(hash))
}

// VerifyReauth verifies credentials for sensitive actions (like account deletion)
func (s *SettingsService) VerifyReauth(userID int64, email, password string) error {
	user, err := s.store.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if user.PasswordHash == "" {
		return errors.New("please use SSO to re-authenticate (no password set)")
	}

	email = utils.Sanitize(email)
	if email != user.Email {
		return errors.New("email does not match our records")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return errors.New("invalid password")
	}

	return nil
}

// DeleteAccount removes the user and all their data
func (s *SettingsService) DeleteAccount(userID int64) error {
	return s.store.DeleteUser(userID)
}

// IsSessionReauthenticated reports whether the session has a recent re-auth timestamp.
func (s *SettingsService) IsSessionReauthenticated(session *models.Session) bool {
	if session == nil {
		return false
	}
	return time.Now().Unix()-session.ReauthenticatedAt < reauthTimeoutSeconds
}

// CompletePasswordReauth verifies credentials and records re-auth on the session.
func (s *SettingsService) CompletePasswordReauth(userID int64, sessionID, email, password string) error {
	if err := s.VerifyReauth(userID, email, password); err != nil {
		return err
	}
	return s.store.UpdateSessionReauth(sessionID, time.Now().Unix())
}

// DeleteAccountAfterReauth removes the user when re-auth is still valid.
func (s *SettingsService) DeleteAccountAfterReauth(userID int64, session *models.Session) error {
	if !s.IsSessionReauthenticated(session) {
		return errors.New("verification expired. Please re-authenticate")
	}
	return s.DeleteAccount(userID)
}

// InvalidateSession removes an active session after sensitive account changes.
func (s *SettingsService) InvalidateSession(sessionID string) error {
	return s.store.DeleteSession(sessionID)
}
