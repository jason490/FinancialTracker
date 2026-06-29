package stripefc

import (
	"context"
	"errors"
	"fmt"

	"github.com/stripe/stripe-go/v86"
)

// ensureStripeCustomer returns an existing Stripe customer ID or creates one for the user.
func (s *Service) ensureStripeCustomer(ctx context.Context, userID int64) (string, error) {
	user, err := s.store.GetUserByID(userID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}
	if user.StripeCustomerID != "" {
		return user.StripeCustomerID, nil
	}

	if err := s.reserveStripeAPICall(userID); err != nil {
		return "", err
	}

	email := user.Email
	name := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	params := &stripe.CustomerCreateParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
	}
	customer, err := s.client.V1Customers.Create(ctx, params)
	if err != nil {
		return "", err
	}

	if err := s.store.UpdateUserStripeCustomerID(userID, customer.ID); err != nil {
		return "", err
	}
	return customer.ID, nil
}
