package external

import "FinancialTracker/internal/models"

// LoginRequest is the payload for email/password sign-in.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

// RegisterRequest is the payload for creating a new account.
type RegisterRequest struct {
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	ConfirmPassword  string `json:"confirm_password"`
	RegistrationCode string `json:"registration_code"`
}

// ForgotPasswordRequest requests a password reset code.
type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

// VerifyResetCodeRequest validates a password reset code.
type VerifyResetCodeRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

// ResetPasswordRequest sets a new password after code verification.
type ResetPasswordRequest struct {
	Email           string `json:"email"`
	Code            string `json:"code"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

// SaveDashboardLayoutRequest persists a dashboard widget layout from the SPA.
type SaveDashboardLayoutRequest struct {
	DeviceType string                   `json:"device_type"`
	Widgets    []models.DashboardWidget `json:"widgets"`
}

// BulkTagRequest is the payload for adding or removing a tag from multiple transactions.
type BulkTagRequest struct {
	TransactionIDs []string `json:"transaction_ids"`
	TagID          int64    `json:"tag_id"`
}

// TagFilterInput is a single auto-tagging rule for tag create/update.
type TagFilterInput struct {
	Pattern    string `json:"pattern"`
	FilterType string `json:"filter_type"`
}

// CreateTagRequest is the payload for creating a new tag.
type CreateTagRequest struct {
	CategoryID int64            `json:"category_id"`
	Name       string           `json:"name"`
	Color      string           `json:"color"`
	Filters    []TagFilterInput `json:"filters"`
	Apply      bool             `json:"apply"`
}

// UpdateTagRequest is the payload for updating an existing tag.
type UpdateTagRequest struct {
	Name       string           `json:"name"`
	Color      string           `json:"color"`
	CategoryID int64            `json:"category_id"`
	Filters    []TagFilterInput `json:"filters"`
	Apply      bool             `json:"apply"`
}

// MoveTagRequest moves a tag to a different category.
type MoveTagRequest struct {
	CategoryID int64 `json:"category_id"`
}

// CreateCategoryRequest is the payload for creating a category.
type CreateCategoryRequest struct {
	Name string `json:"name"`
}

// UpdateCategoryRequest is the payload for renaming a category.
type UpdateCategoryRequest struct {
	Name string `json:"name"`
}

// DeleteCategoryRequest describes how to handle tags when deleting a category.
type DeleteCategoryRequest struct {
	Action           string `json:"action"`
	TargetCategoryID int64  `json:"target_category_id,omitempty"`
}
