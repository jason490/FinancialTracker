package auth

import (
	"FinancialTracker/internal/models"
	"FinancialTracker/internal/services/mail"
	"FinancialTracker/internal/storage"
	"FinancialTracker/internal/utils"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	store      *storage.Storage
	mailSender mail.Sender
}

func NewAuthService(store *storage.Storage, mailSender mail.Sender) *AuthService {
	return &AuthService{
		store:      store,
		mailSender: mailSender,
	}
}

// Authenticate verifies user credentials and returns a new session
func (s *AuthService) Authenticate(email, password string, remember bool) (*models.Session, error) {
	email = utils.Sanitize(email)
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	user, err := s.store.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	sessionID := uuid.New().String()
	expiresAt := time.Now().AddDate(0, 0, 1).Unix()
	if remember {
		expiresAt = time.Now().AddDate(1, 0, 0).Unix()
	}

	session := &models.Session{
		ID:                sessionID,
		UserID:            user.ID,
		ExpiresAt:         expiresAt,
		ReauthenticatedAt: time.Now().Unix(),
	}

	if err := s.store.CreateSession(session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

// Register creates a new user and returns a new session
func (s *AuthService) Register(firstName, lastName, email, password, confirmPassword string) (*models.Session, error) {
	firstName = utils.Sanitize(firstName)
	lastName = utils.Sanitize(lastName)
	email = strings.ToLower(utils.Sanitize(email))

	if msg, err := utils.ValidateName(firstName, "First name"); err {
		return nil, errors.New(msg)
	}
	if msg, err := utils.ValidateName(lastName, "Last name"); err {
		return nil, errors.New(msg)
	}
	if msg, err := utils.ValidateEmail(email); err {
		return nil, errors.New(msg)
	}
	if password != confirmPassword {
		return nil, errors.New("passwords do not match")
	}
	if msg, err := utils.ValidatePassword(password); err {
		return nil, errors.New(msg)
	}

	existing, _ := s.store.GetUserByEmail(email)
	if existing != nil {
		return nil, errors.New("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Email:           email,
		FirstName:       firstName,
		LastName:        lastName,
		PasswordHash:    string(hash),
		ThemePreference: "system",
	}

	if err := s.store.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	sessionID := uuid.New().String()
	session := &models.Session{
		ID:                sessionID,
		UserID:            user.ID,
		ExpiresAt:         time.Now().AddDate(0, 0, 1).Unix(),
		ReauthenticatedAt: time.Now().Unix(),
	}

	if err := s.store.CreateSession(session); err != nil {
		return nil, fmt.Errorf("failed to create session after registration: %w", err)
	}

	return session, nil
}

// CompleteOnboarding marks the user's onboarding wizard as finished.
func (s *AuthService) CompleteOnboarding(userID int64) error {
	return s.store.CompleteOnboarding(userID)
}

// Logout invalidates a user session
func (s *AuthService) Logout(sessionID string) error {
	if sessionID == "" {
		return nil
	}
	return s.store.DeleteSession(sessionID)
}

