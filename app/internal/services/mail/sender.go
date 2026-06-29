package mail

// Sender delivers outbound email for authentication flows.
type Sender interface {
	SendPasswordResetCode(to, firstName, code string) error
}
