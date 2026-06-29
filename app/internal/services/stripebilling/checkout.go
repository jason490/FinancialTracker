package stripebilling

import (
	"context"
	"errors"
	"fmt"
	"time"

	"FinancialTracker/internal/models"
	"FinancialTracker/internal/services/financial"
	subscriptionService "FinancialTracker/internal/services/subscription"

	"github.com/stripe/stripe-go/v86"
)

// CreateCheckoutSession starts a Stripe Checkout subscription flow for the given tier.
func (s *Service) CreateCheckoutSession(ctx context.Context, userID int64, tier, successURL, cancelURL string) (string, error) {
	if !BillingReady() {
		return "", ErrBillingNotConfigured
	}
	if tier != models.SubscriptionTierPlus && tier != models.SubscriptionTierPremium {
		return "", subscriptionService.ErrInvalidTier
	}

	priceID := subscriptionService.PriceIDForTier(tier)
	if priceID == "" {
		return "", ErrBillingNotConfigured
	}

	user, err := s.store.GetUserByID(userID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	customerID, err := s.ensureCustomer(ctx, user)
	if err != nil {
		return "", err
	}

	params := &stripe.CheckoutSessionCreateParams{
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		Customer:   stripe.String(customerID),
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		LineItems: []*stripe.CheckoutSessionCreateLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		ClientReferenceID: stripe.String(fmt.Sprintf("%d", userID)),
		Metadata: map[string]string{
			"user_id": fmt.Sprintf("%d", userID),
			"tier":    tier,
		},
		SubscriptionData: &stripe.CheckoutSessionCreateSubscriptionDataParams{
			Metadata: map[string]string{
				"user_id": fmt.Sprintf("%d", userID),
				"tier":    tier,
			},
		},
	}

	couponID, err := s.subscription.StripeCouponForUser(userID)
	if err != nil {
		return "", err
	}
	if couponID != "" {
		params.Discounts = []*stripe.CheckoutSessionCreateDiscountParams{
			{Coupon: stripe.String(couponID)},
		}
	}

	session, err := s.client.V1CheckoutSessions.Create(ctx, params)
	if err != nil {
		return "", err
	}
	if session.URL == "" {
		return "", errors.New("stripe checkout session missing redirect URL")
	}
	return session.URL, nil
}

// CreatePortalSession opens the Stripe Customer Portal for subscription management.
func (s *Service) CreatePortalSession(ctx context.Context, userID int64, returnURL string) (string, error) {
	if !financial.StripeConfigured() {
		return "", ErrBillingNotConfigured
	}

	user, err := s.store.GetUserByID(userID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}
	if user.StripeCustomerID == "" {
		return "", ErrNoStripeCustomer
	}

	params := &stripe.BillingPortalSessionCreateParams{
		Customer:  stripe.String(user.StripeCustomerID),
		ReturnURL: stripe.String(returnURL),
	}
	session, err := s.client.V1BillingPortalSessions.Create(ctx, params)
	if err != nil {
		return "", err
	}
	if session.URL == "" {
		return "", errors.New("stripe portal session missing redirect URL")
	}
	return session.URL, nil
}

// subscriptionStartedAt picks a billing anchor from a Stripe subscription.
func subscriptionStartedAt(sub *stripe.Subscription) int64 {
	if sub == nil {
		return time.Now().Unix()
	}
	if sub.StartDate > 0 {
		return sub.StartDate
	}
	if sub.BillingCycleAnchor > 0 {
		return sub.BillingCycleAnchor
	}
	if sub.Created > 0 {
		return sub.Created
	}
	return time.Now().Unix()
}
