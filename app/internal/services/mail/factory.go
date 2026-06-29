package mail

import (
	"FinancialTracker/internal/config"
	"log"
	"os"
	"strings"
)

// NewSenderFromEnv returns the configured mail sender for the current environment.
func NewSenderFromEnv() Sender {
	if strings.TrimSpace(os.Getenv("MAIL_SMTP_HOST")) != "" {
		sender, err := NewSMTPSenderFromEnv()
		if err != nil {
			log.Printf("[mail] SMTP configuration error: %v; falling back to noop sender", err)
			return NoopSender{}
		}
		return sender
	}
	if config.IsDevelopment() {
		return DevSender{}
	}
	return NoopSender{}
}
