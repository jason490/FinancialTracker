package stripebilling

import (
	"context"
	"errors"
	"fmt"

	"FinancialTracker/internal/models"

	"github.com/stripe/stripe-go/v86"
)

// ensureCustomer returns an existing Stripe customer ID or creates one for billing.
func (s *Service) ensureCustomer(ctx context.Context, user *models.User) (string, error) {
	if user.StripeCustomerID != "" {
		return user.StripeCustomerID, nil
	}

	email := user.Email
	name := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	params := &stripe.CustomerCreateParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
		Metadata: map[string]string{
			"user_id": fmt.Sprintf("%d", user.ID),
		},
	}
	customer, err := s.client.V1Customers.Create(ctx, params)
	if err != nil {
		return "", err
	}
	if err := s.store.UpdateUserStripeCustomerID(user.ID, customer.ID); err != nil {
		return "", err
	}
	return customer.ID, nil
}

// resolveUserIDFromCustomer looks up the app user for a Stripe customer ID.
func (s *Service) resolveUserIDFromCustomer(customerID string) (int64, error) {
	if customerID == "" {
		return 0, errors.New("missing stripe customer id")
	}
	user, err := s.store.GetUserByStripeCustomerID(customerID)
	if err != nil {
		return 0, err
	}
	if user == nil {
		return 0, errors.New("no user for stripe customer")
	}
	return user.ID, nil
}
