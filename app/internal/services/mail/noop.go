package mail

import (
	"fmt"
	"log"
)

// NoopSender discards outbound mail and logs a warning (production without SMTP configured).
type NoopSender struct{}

// SendPasswordResetCode logs that mail is not configured and does not deliver the code.
func (NoopSender) SendPasswordResetCode(to, firstName, code string) error {
	log.Printf("[mail:noop] password reset requested for %s but SMTP is not configured; code not sent", to)
	return fmt.Errorf("mail sender not configured")
}
