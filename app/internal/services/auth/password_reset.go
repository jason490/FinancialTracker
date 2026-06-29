package auth

import (
	"FinancialTracker/internal/models"
	"FinancialTracker/internal/utils"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	resetCodeTTL        = 15 * time.Minute
	resetMaxAttempts    = 5
	resetResendCooldown = 60 * time.Second
)

// ResetCodeTTLSeconds is the password reset code lifetime exposed to clients.
const ResetCodeTTLSeconds = int64(resetCodeTTL / time.Second)

var (
	// ErrInvalidResetCode is returned when a reset code cannot be verified.
	ErrInvalidResetCode = errors.New("invalid or expired reset code")
	resetCodePattern    = regexp.MustCompile(`^\d{6}$`)
)

// RequestPasswordReset issues a temporary reset code when the account has a password.
func (s *AuthService) RequestPasswordReset(email string) error {
	email = strings.ToLower(utils.Sanitize(email))
	if email == "" {
		return errors.New("email is required")
	}
	if msg, invalid := utils.ValidateEmail(email); invalid {
		return errors.New(msg)
	}

	user, err := s.store.GetUserByEmail(email)
	if err != nil {
		log.Printf("[auth] password reset lookup failed for %s: %v", email, err)
		return nil
	}
	if user == nil || user.PasswordHash == "" {
		return nil
	}

	since := time.Now().Unix() - int64(resetResendCooldown.Seconds())
	issuedRecently, err := s.store.CountPasswordResetCodesSince(user.ID, since)
	if err != nil {
		log.Printf("[auth] password reset cooldown check failed for user %d: %v", user.ID, err)
		return nil
	}
	if issuedRecently >= 2 {
		return nil
	}

	code, err := generateResetCode()
	if err != nil {
		log.Printf("[auth] password reset code generation failed for user %d: %v", user.ID, err)
		return nil
	}

	codeHash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[auth] password reset code hash failed for user %d: %v", user.ID, err)
		return nil
	}

	if err := s.store.InvalidatePasswordResetCodes(user.ID); err != nil {
		log.Printf("[auth] password reset invalidate failed for user %d: %v", user.ID, err)
		return nil
	}

	expiresAt := time.Now().Add(resetCodeTTL).Unix()
	if err := s.store.CreatePasswordResetCode(user.ID, string(codeHash), expiresAt); err != nil {
		log.Printf("[auth] password reset store failed for user %d: %v", user.ID, err)
		return nil
	}

	if err := s.mailSender.SendPasswordResetCode(user.Email, user.FirstName, code); err != nil {
		log.Printf("[auth] password reset mail failed for %s: %v", user.Email, err)
	}

	return nil
}

// VerifyPasswordResetCode checks a reset code without changing the password.
func (s *AuthService) VerifyPasswordResetCode(email, code string) (int64, error) {
	_, resetRow, err := s.validateResetCode(email, code)
	if err != nil {
		return 0, err
	}
	return resetRow.ExpiresAt, nil
}

// ConfirmPasswordReset verifies a reset code and sets a new password.
func (s *AuthService) ConfirmPasswordReset(email, code, newPassword, confirmPassword string) error {
	email = strings.ToLower(utils.Sanitize(email))
	code = strings.TrimSpace(code)

	if email == "" {
		return errors.New("email is required")
	}
	if msg, invalid := utils.ValidateEmail(email); invalid {
		return errors.New(msg)
	}
	if !resetCodePattern.MatchString(code) {
		return errors.New("reset code must be 6 digits")
	}
	if msg, invalid := utils.ValidatePassword(newPassword); invalid {
		return errors.New(msg)
	}
	if newPassword != confirmPassword {
		return errors.New("passwords do not match")
	}

	user, resetRow, err := s.validateResetCode(email, code)
	if err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.store.UpdateUserPassword(user.ID, string(hash)); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	if err := s.store.MarkPasswordResetCodeUsed(resetRow.ID); err != nil {
		return fmt.Errorf("failed to mark reset code used: %w", err)
	}
	if err := s.store.DeleteUserSessions(user.ID); err != nil {
		return fmt.Errorf("failed to invalidate sessions: %w", err)
	}

	return nil
}

func (s *AuthService) validateResetCode(email, code string) (*models.User, *models.PasswordResetCode, error) {
	email = strings.ToLower(utils.Sanitize(email))
	code = strings.TrimSpace(code)

	if email == "" {
		return nil, nil, errors.New("email is required")
	}
	if msg, invalid := utils.ValidateEmail(email); invalid {
		return nil, nil, errors.New(msg)
	}
	if !resetCodePattern.MatchString(code) {
		return nil, nil, errors.New("reset code must be 6 digits")
	}

	user, err := s.store.GetUserByEmail(email)
	if err != nil || user == nil {
		return nil, nil, ErrInvalidResetCode
	}

	resetRow, err := s.store.GetActivePasswordResetCode(user.ID)
	if err != nil || resetRow == nil {
		return nil, nil, ErrInvalidResetCode
	}
	if resetRow.Attempts >= resetMaxAttempts {
		return nil, nil, ErrInvalidResetCode
	}

	if err := bcrypt.CompareHashAndPassword([]byte(resetRow.CodeHash), []byte(code)); err != nil {
		_ = s.store.IncrementPasswordResetAttempts(resetRow.ID)
		return nil, nil, ErrInvalidResetCode
	}

	return user, resetRow, nil
}

// generateResetCode returns a cryptographically random 6-digit code.
func generateResetCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}
