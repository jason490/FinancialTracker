package auth

import (
	"FinancialTracker/internal/models/external"
	"fmt"
)

// GetSessionProfile returns the minimal user profile for authenticated clients.
func (s *AuthService) GetSessionProfile(userID int64) (*external.SessionProfile, error) {
	user, err := s.store.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to load user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return &external.SessionProfile{
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		Email:               user.Email,
		ThemePreference:     user.ThemePreference,
		OnboardingCompleted: user.OnboardingCompleted,
	}, nil
}
