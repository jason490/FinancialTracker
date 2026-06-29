package mail

import (
	"fmt"
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

// TODO: configure SMTP for production password reset emails.
//
// Set these environment variables to enable automatic email delivery:
//   - MAIL_SMTP_HOST     (required) e.g. smtp.sendgrid.net
//   - MAIL_SMTP_PORT     (required) e.g. 587
//   - MAIL_SMTP_USER     (required) SMTP username
//   - MAIL_SMTP_PASSWORD (required) SMTP password or API key
//   - MAIL_FROM          (required) sender address, e.g. noreply@yourdomain.com
//   - MAIL_FROM_NAME     (optional) display name, defaults to "Financial Tracker"
//
// When MAIL_SMTP_HOST is set, NewSenderFromEnv selects SMTPSender instead of DevSender/NoopSender.

// SMTPSender sends password reset codes via SMTP.
type SMTPSender struct {
	host     string
	port     int
	user     string
	password string
	from     string
	fromName string
}

// NewSMTPSenderFromEnv builds an SMTPSender from MAIL_* environment variables.
func NewSMTPSenderFromEnv() (*SMTPSender, error) {
	host := strings.TrimSpace(os.Getenv("MAIL_SMTP_HOST"))
	portStr := strings.TrimSpace(os.Getenv("MAIL_SMTP_PORT"))
	user := strings.TrimSpace(os.Getenv("MAIL_SMTP_USER"))
	password := os.Getenv("MAIL_SMTP_PASSWORD")
	from := strings.TrimSpace(os.Getenv("MAIL_FROM"))
	fromName := strings.TrimSpace(os.Getenv("MAIL_FROM_NAME"))

	missing := []string{}
	if host == "" {
		missing = append(missing, "MAIL_SMTP_HOST")
	}
	if portStr == "" {
		missing = append(missing, "MAIL_SMTP_PORT")
	}
	if user == "" {
		missing = append(missing, "MAIL_SMTP_USER")
	}
	if password == "" {
		missing = append(missing, "MAIL_SMTP_PASSWORD")
	}
	if from == "" {
		missing = append(missing, "MAIL_FROM")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("incomplete SMTP configuration: missing %s", strings.Join(missing, ", "))
	}

	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 {
		return nil, fmt.Errorf("invalid MAIL_SMTP_PORT: %q", portStr)
	}
	if fromName == "" {
		fromName = "Financial Tracker"
	}

	return &SMTPSender{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		from:     from,
		fromName: fromName,
	}, nil
}

// SendPasswordResetCode emails a plain-text reset code to the user.
func (s *SMTPSender) SendPasswordResetCode(to, firstName, code string) error {
	name := firstName
	if name == "" {
		name = "there"
	}

	subject := "Your password reset code"
	body := fmt.Sprintf(`Hi %s,

You requested a password reset for your Financial Tracker account.

Your reset code is: %s

This code expires in 15 minutes. If you did not request a reset, you can ignore this email.

— Financial Tracker
`, name, code)

	fromHeader := fmt.Sprintf("%s <%s>", s.fromName, s.from)
	msg := strings.Join([]string{
		"From: " + fromHeader,
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	auth := smtp.PlainAuth("", s.user, s.password, s.host)
	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg))
}
