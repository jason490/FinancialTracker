package models

// Notification represents a message to be displayed to the user
type Notification struct {
	Message string `json:"message"`
	Type    string `json:"type"` // "error", "success", "info"
}

// PageData represents the data passed to every page template
type PageData struct {
	User          *User          `json:"user"`
	Title         string         `json:"title"`
	Data          any            `json:"data"`
	ReauthSuccess bool           `json:"reauth_success"`
	Notifications []Notification `json:"notifications,omitempty"`
}
