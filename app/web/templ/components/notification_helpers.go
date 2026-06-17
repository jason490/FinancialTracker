package components

import (
	"encoding/json"

	"FinancialTracker/internal/models"
)

// NotificationsJSON serializes page-load notifications for client-side bootstrap.
func NotificationsJSON(notifications []models.Notification) string {
	type payload struct {
		Msg  string `json:"msg"`
		Type string `json:"type"`
	}
	items := make([]payload, len(notifications))
	for i, n := range notifications {
		items[i] = payload{Msg: n.Message, Type: n.Type}
	}
	b, err := json.Marshal(items)
	if err != nil {
		return "[]"
	}
	return string(b)
}
