package stripebilling

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"FinancialTracker/internal/models"
	subscriptionService "FinancialTracker/internal/services/subscription"

	"github.com/labstack/gommon/log"
	"github.com/stripe/stripe-go/v86"
	"github.com/stripe/stripe-go/v86/webhook"
)

// HandleWebhook verifies and processes Stripe billing webhook events.
func (s *Service) HandleWebhook(body []byte, signature string) error {
	secret := webhookSecret()
	if secret == "" {
		return ErrWebhookNotConfigured
	}

	event, err := webhook.ConstructEvent(body, signature, secret)
	if err != nil {
		return ErrWebhookVerification
	}

	ctx := context.Background()
	switch event.Type {
	case "checkout.session.completed":
		return s.handleCheckoutCompleted(ctx, event.Data.Raw)
	case "customer.subscription.created", "customer.subscription.updated":
		return s.handleSubscriptionChanged(ctx, event.Data.Raw)
	case "customer.subscription.deleted":
		return s.handleSubscriptionDeleted(ctx, event.Data.Raw)
	default:
		return nil
	}
}

func (s *Service) handleCheckoutCompleted(ctx context.Context, raw json.RawMessage) error {
	var session stripe.CheckoutSession
	if err := json.Unmarshal(raw, &session); err != nil {
		return err
	}
	if session.Mode != stripe.CheckoutSessionModeSubscription {
		return nil
	}

	userID, err := s.userIDFromCheckoutSession(&session)
	if err == nil && session.Customer != nil && session.Customer.ID != "" {
		_ = s.store.UpdateUserStripeCustomerID(userID, session.Customer.ID)
	}

	if session.Subscription == nil || session.Subscription.ID == "" {
		return nil
	}

	sub, err := s.client.V1Subscriptions.Retrieve(ctx, session.Subscription.ID, &stripe.SubscriptionRetrieveParams{})
	if err != nil {
		return err
	}
	return s.applySubscription(ctx, sub)
}

func (s *Service) handleSubscriptionChanged(ctx context.Context, raw json.RawMessage) error {
	var sub stripe.Subscription
	if err := json.Unmarshal(raw, &sub); err != nil {
		return err
	}
	return s.applySubscription(ctx, &sub)
}

func (s *Service) handleSubscriptionDeleted(ctx context.Context, raw json.RawMessage) error {
	var sub stripe.Subscription
	if err := json.Unmarshal(raw, &sub); err != nil {
		return err
	}

	userID, err := s.userIDFromSubscription(&sub)
	if err != nil {
		log.Warnf("Stripe subscription deleted without resolvable user: %v", err)
		return nil
	}
	return s.subscription.ApplyStripeSubscription(userID, models.SubscriptionTierFree, 0, "")
}

func (s *Service) applySubscription(ctx context.Context, sub *stripe.Subscription) error {
	if sub == nil {
		return nil
	}
	if sub.Status != stripe.SubscriptionStatusActive && sub.Status != stripe.SubscriptionStatusTrialing {
		return nil
	}

	tier, err := s.tierFromSubscription(sub)
	if err != nil {
		return err
	}
	userID, err := s.userIDFromSubscription(sub)
	if err != nil {
		return err
	}

	startedAt := subscriptionStartedAt(sub)
	return s.subscription.ApplyStripeSubscription(userID, tier, startedAt, sub.ID)
}

func (s *Service) tierFromSubscription(sub *stripe.Subscription) (string, error) {
	if sub == nil || len(sub.Items.Data) == 0 {
		return "", errors.New("subscription has no line items")
	}
	priceID := ""
	if sub.Items.Data[0].Price != nil {
		priceID = sub.Items.Data[0].Price.ID
	}
	tier := subscriptionService.TierForPriceID(priceID)
	if tier == "" {
		return "", errors.New("unknown stripe price id")
	}
	return tier, nil
}

func (s *Service) userIDFromSubscription(sub *stripe.Subscription) (int64, error) {
	if sub == nil {
		return 0, errors.New("missing subscription")
	}
	if sub.Metadata != nil {
		if raw, ok := sub.Metadata["user_id"]; ok && raw != "" {
			var userID int64
			if _, err := fmt.Sscanf(raw, "%d", &userID); err == nil && userID > 0 {
				return userID, nil
			}
		}
	}
	if sub.Customer == nil || sub.Customer.ID == "" {
		return 0, errors.New("subscription missing customer")
	}
	return s.resolveUserIDFromCustomer(sub.Customer.ID)
}

func (s *Service) userIDFromCheckoutSession(session *stripe.CheckoutSession) (int64, error) {
	if session == nil {
		return 0, errors.New("missing checkout session")
	}
	if session.Metadata != nil {
		if raw, ok := session.Metadata["user_id"]; ok && raw != "" {
			var userID int64
			if _, err := fmt.Sscanf(raw, "%d", &userID); err == nil && userID > 0 {
				return userID, nil
			}
		}
	}
	if session.ClientReferenceID != "" {
		var userID int64
		if _, err := fmt.Sscanf(session.ClientReferenceID, "%d", &userID); err == nil && userID > 0 {
			return userID, nil
		}
	}
	if session.Customer != nil && session.Customer.ID != "" {
		return s.resolveUserIDFromCustomer(session.Customer.ID)
	}
	return 0, errors.New("checkout session missing user reference")
}
