package auth

import (
	"FinancialTracker/internal/config"
	"FinancialTracker/internal/models"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	registrationCodeTTL     = 48 * time.Hour
	registrationMaxAttempts = 5
)

// RegistrationCodeTTLSeconds is the invite code lifetime exposed to clients.
const RegistrationCodeTTLSeconds = int64(registrationCodeTTL / time.Second)

var (
	// ErrInvalidRegistrationCode is returned when an invite code cannot be verified.
	ErrInvalidRegistrationCode = errors.New("invalid or expired registration code")
	// ErrRegistrationCodeRequired is returned when invite-only registration lacks a code.
	ErrRegistrationCodeRequired = errors.New("registration code is required")
	registrationCodePattern     = regexp.MustCompile(`^[A-Z0-9]{8}$`)
)

// RegistrationGateEnabled reports whether new accounts require an admin invite code.
func RegistrationGateEnabled() bool {
	return config.RegistrationGateEnabled()
}

const developmentRegistrationAdminEmail = "test@test.com"

// IsRegistrationAdmin reports whether the email may issue registration invite codes.
func IsRegistrationAdmin(email string) bool {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return false
	}
	if config.IsDevelopment() && email == developmentRegistrationAdminEmail {
		return true
	}
	for _, admin := range config.RegistrationAdminEmails() {
		if admin == email {
			return true
		}
	}
	return false
}

// GenerateRegistrationCode creates a single-use invite code and returns the plaintext value once.
func (s *AuthService) GenerateRegistrationCode(createdByUserID *int64) (string, int64, error) {
	code, err := generateRegistrationCodeValue()
	if err != nil {
		return "", 0, fmt.Errorf("failed to generate registration code: %w", err)
	}

	codeHash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return "", 0, fmt.Errorf("failed to hash registration code: %w", err)
	}

	expiresAt := time.Now().Add(registrationCodeTTL).Unix()
	if err := s.store.CreateRegistrationCode(createdByUserID, string(codeHash), expiresAt); err != nil {
		return "", 0, fmt.Errorf("failed to store registration code: %w", err)
	}

	return code, expiresAt, nil
}

// ValidateAndConsumeRegistrationCode verifies an invite code and marks it used for the new user.
func (s *AuthService) ValidateAndConsumeRegistrationCode(code string, usedByUserID int64) error {
	if !RegistrationGateEnabled() {
		return nil
	}

	row, err := s.findRegistrationCode(code)
	if err != nil {
		return err
	}
	if row == nil {
		return ErrInvalidRegistrationCode
	}

	return s.store.MarkRegistrationCodeUsed(row.ID, usedByUserID)
}

// RequireRegistrationCode validates an invite code before account creation.
func (s *AuthService) RequireRegistrationCode(code string) error {
	if !RegistrationGateEnabled() {
		return nil
	}

	code = strings.ToUpper(strings.TrimSpace(code))
	if code == "" {
		return ErrRegistrationCodeRequired
	}
	if !registrationCodePattern.MatchString(code) {
		return ErrInvalidRegistrationCode
	}

	_, err := s.findRegistrationCode(code)
	return err
}

func (s *AuthService) findRegistrationCode(code string) (*models.RegistrationCode, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	if code == "" {
		return nil, ErrRegistrationCodeRequired
	}
	if !registrationCodePattern.MatchString(code) {
		return nil, ErrInvalidRegistrationCode
	}

	codes, err := s.store.ListActiveRegistrationCodes()
	if err != nil {
		return nil, fmt.Errorf("failed to list registration codes: %w", err)
	}

	for _, row := range codes {
		if row.Attempts >= registrationMaxAttempts {
			continue
		}
		if err := bcrypt.CompareHashAndPassword([]byte(row.CodeHash), []byte(code)); err == nil {
			return row, nil
		}
	}

	return nil, ErrInvalidRegistrationCode
}

// generateRegistrationCodeValue returns a cryptographically random 8-character uppercase code.
func generateRegistrationCodeValue() (string, error) {
	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	out := make([]byte, 8)
	for i := range out {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		out[i] = alphabet[n.Int64()]
	}
	return string(out), nil
}
