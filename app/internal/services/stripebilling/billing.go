package stripebilling

import (
	"errors"
)

var (
	ErrBillingNotConfigured = errors.New("stripe billing is not configured")
	ErrPaidTierRequired     = errors.New("checkout requires a paid plan")
	ErrNoStripeCustomer     = errors.New("no stripe customer on file")
	ErrWebhookNotConfigured = errors.New("stripe webhook secret is not configured")
	ErrWebhookVerification  = errors.New("stripe webhook verification failed")
)
