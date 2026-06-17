package legacyHandler

import (
	"FinancialTracker/internal/models"
	"FinancialTracker/internal/storage"
	"FinancialTracker/internal/utils"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
)

// ReauthTimeout is the time window for high-privilege actions (5 minutes)
const ReauthTimeout = 300

// sanitize removes leading/trailing whitespace
func sanitize(s string) string {
    return utils.Sanitize(s)
}

// createCookie creates a new authentication cookie
func createCookie(sessionid string, remember bool) *http.Cookie {
	env := os.Getenv("ENV")
    var cookie = new(http.Cookie)
    cookie.Name = "Session"
    cookie.Value = sessionid
    if remember {
        cookie.Expires = time.Now().AddDate(1, 0, 0)
    } else {
        cookie.Expires = time.Now().AddDate(0, 0, 1)
    }
    cookie.HttpOnly = true
	if env == "development" {
    	cookie.Secure = false
	} else {
    	cookie.Secure = true
	}
    cookie.Path = "/"
    return cookie
}

// Render renders a templ component
func Render(ctx *echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

// ErrorMessages maps numeric error codes to string messages
var ErrorMessages = map[string]string{
	"1": "Authentication failed. Please try again.",
	"2": "This SSO account is already linked to another user.",
	"3": "Failed to link SSO account.",
	"4": "Identity mismatch. Please log in with the correct account.",
	"5": "A database error occurred. Please try again later.",
	"6": "Your session has expired. Please log in again.",
	"7": "Invalid email or password.",
	"8": "User already exists.",
}

// getHXTriggers parses the current HX-Trigger header into a map for merging.
func getHXTriggers(c *echo.Context) map[string]any {
	existing := c.Response().Header().Get("HX-Trigger")
	triggers := make(map[string]any)

	if existing == "" {
		return triggers
	}

	if err := json.Unmarshal([]byte(existing), &triggers); err != nil {
		for part := range strings.SplitSeq(existing, ",") {
			key := strings.TrimSpace(part)
			if key != "" {
				triggers[key] = true
			}
		}
	}
	return triggers
}

// setHXTriggers writes the merged HX-Trigger header.
func setHXTriggers(c *echo.Context, triggers map[string]any) {
	data, _ := json.Marshal(triggers)
	c.Response().Header().Set("HX-Trigger", string(data))
}

// AddHXTriggerEvent adds a simple named HTMX event to the response triggers.
func AddHXTriggerEvent(c *echo.Context, name string) {
	if name == "" {
		return
	}
	triggers := getHXTriggers(c)
	triggers[name] = true
	setHXTriggers(c, triggers)
}

// AddNotification sets the HX-Trigger header for toast notifications.
func AddNotification(c *echo.Context, msg, nType string) {
	notification := map[string]string{
		"msg":  msg,
		"type": nType,
	}

	triggers := getHXTriggers(c)

	notifs, ok := triggers["notification"].([]any)
	if !ok {
		if single, ok := triggers["notification"].(map[string]any); ok {
			notifs = []any{single}
		} else {
			notifs = []any{}
		}
	}

	triggers["notification"] = append(notifs, notification)
	setHXTriggers(c, triggers)
}

// QueuePageNotification appends a notification to PageData for display on full page load.
func QueuePageNotification(pageData *models.PageData, msg, nType string) {
	if pageData == nil {
		return
	}
	pageData.Notifications = append(pageData.Notifications, models.Notification{
		Message: msg,
		Type:    nType,
	})
}

// GetPageData populates the PageData struct with user info
func GetPageData(c *echo.Context, store *storage.Storage, title string) *models.PageData {
	user := getCurrentUser(c, store)

	return &models.PageData{
		User:  user,
		Title: title,
	}
}

// getCurrentUser retrieves the current user based on session
func getCurrentUser(c *echo.Context, store *storage.Storage) *models.User {
	userID := c.Get("user_id")
	if userID == nil {
		return nil
	}
	user, err := store.GetUserByID(userID.(int64))
	if err != nil || user == nil {
		return nil
	}

	return user
}
