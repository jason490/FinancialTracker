package external

import "FinancialTracker/internal/models"

// ConnectionAccountView is a sanitized bank account for connection management.
type ConnectionAccountView struct {
	AccountID string  `json:"account_id"`
	Name      string  `json:"name"`
	Mask      string  `json:"mask"`
	Subtype   string  `json:"subtype"`
	Balance   float64 `json:"balance"`
	Currency  string  `json:"currency"`
	Status    string  `json:"status"`
	IsHidden  bool    `json:"is_hidden"`
}

// ConnectionView is an institution connection with its accounts.
type ConnectionView struct {
	RowID           string                  `json:"row_id"`
	InstitutionName string                  `json:"institution_name"`
	Status          string                  `json:"status"`
	CreatedAt       int64                   `json:"created_at"`
	LastSynced      int64                   `json:"last_synced"`
	Accounts        []ConnectionAccountView `json:"accounts"`
}

// ConnectionsPayload lists all bank connections for the user.
type ConnectionsPayload struct {
	Provider    string             `json:"provider"`
	Connections []ConnectionView   `json:"connections"`
	Usage       models.PlaidUsage  `json:"usage"`
}

// ProviderInfoResponse exposes the active financial provider to the client.
type ProviderInfoResponse struct {
	Provider       string `json:"provider"`
	PublishableKey string `json:"publishable_key,omitempty"`
}

// CreateSessionResponse returns credentials to launch the provider link flow.
type CreateSessionResponse struct {
	LinkToken    string `json:"link_token,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
}

// CompleteConnectionRequest completes a bank link flow on the server.
type CompleteConnectionRequest struct {
	PublicToken string `json:"public_token,omitempty"`
	SessionID   string `json:"session_id,omitempty"`
}
