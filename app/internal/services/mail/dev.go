package mail

import (
	"log"
)

// DevSender logs password reset codes to stdout for local development.
type DevSender struct{}

// SendPasswordResetCode logs the reset code instead of sending email.
func (DevSender) SendPasswordResetCode(to, firstName, code string) error {
	name := firstName
	if name == "" {
		name = "there"
	}
	log.Printf("[mail:dev] password reset code for %s (%s): %s", to, name, code)
	return nil
}
