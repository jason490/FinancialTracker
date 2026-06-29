package subscription

import (
	"FinancialTracker/internal/config"
	"FinancialTracker/internal/models"
	"FinancialTracker/internal/storage"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	unlimitedItems    = 1_000_000
	unlimitedAPICalls = 1_000_000
)

// EffectiveLimits returns the quota limits for a user, honoring operator overrides.
func (s *Service) EffectiveLimits(userID int64) (models.PlaidLimits, error) {
	if !config.SubscriptionsEnabled() {
		return models.PlaidLimits{
			MaxItems:         unlimitedItems,
			MaxAPICallsMonth: unlimitedAPICalls,
		}, nil
	}

	sub, err := s.GetUserSubscription(userID)
	if err != nil {
		return models.PlaidLimits{}, err
	}
	limits := LimitsForTier(sub.Tier)

	priv, err := s.store.GetUserPrivileges(userID)
	if err != nil {
		return models.PlaidLimits{}, err
	}
	if priv != nil && priv.UnlimitedLimits {
		return models.PlaidLimits{
			MaxItems:         unlimitedItems,
			MaxAPICallsMonth: unlimitedAPICalls,
		}, nil
	}
	return limits, nil
}

// PrivilegesForUser returns the client-facing privilege summary for a user.
func (s *Service) PrivilegesForUser(userID int64) (models.SubscriptionPrivileges, error) {
	priv, err := s.store.GetUserPrivileges(userID)
	if err != nil {
		return models.SubscriptionPrivileges{}, err
	}
	if priv == nil {
		return models.SubscriptionPrivileges{}, nil
	}
	return models.SubscriptionPrivileges{
		UnlimitedLimits: priv.UnlimitedLimits,
		HasDiscount:     strings.TrimSpace(priv.StripeCouponID) != "",
	}, nil
}

// StripeCouponForUser returns a Stripe coupon ID to apply at checkout, if configured.
func (s *Service) StripeCouponForUser(userID int64) (string, error) {
	priv, err := s.store.GetUserPrivileges(userID)
	if err != nil {
		return "", err
	}
	if priv == nil {
		return "", nil
	}
	return strings.TrimSpace(priv.StripeCouponID), nil
}

// HasUnlimitedLimits reports whether a user bypasses subscription quotas.
func (s *Service) HasUnlimitedLimits(userID int64) (bool, error) {
	priv, err := s.store.GetUserPrivileges(userID)
	if err != nil {
		return false, err
	}
	return priv != nil && priv.UnlimitedLimits, nil
}

// SyncPrivilegeOverridesFromEnv applies SUBSCRIPTION_OVERRIDES to the database on startup.
func (s *Service) SyncPrivilegeOverridesFromEnv() {
	raw := strings.TrimSpace(os.Getenv("SUBSCRIPTION_OVERRIDES"))
	if raw == "" {
		return
	}

	entries := strings.Split(raw, ";")
	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		email, attrs, ok := strings.Cut(entry, "=")
		if !ok {
			log.Printf("[subscription] skipping invalid override entry %q", entry)
			continue
		}
		email = strings.TrimSpace(strings.ToLower(email))
		if email == "" {
			continue
		}

		user, err := s.store.GetUserByEmail(email)
		if err != nil {
			log.Printf("[subscription] override lookup failed for %s: %v", email, err)
			continue
		}
		if user == nil {
			log.Printf("[subscription] override skipped; no user for %s", email)
			continue
		}

		priv := &models.UserPrivileges{UserID: user.ID}
		for _, attr := range strings.Split(attrs, ",") {
			attr = strings.TrimSpace(attr)
			if attr == "" {
				continue
			}
			key, value, hasValue := strings.Cut(attr, "=")
			key = strings.TrimSpace(strings.ToLower(key))
			value = strings.TrimSpace(value)
			switch key {
			case "unlimited":
				priv.UnlimitedLimits = true
			case "coupon":
				if hasValue && value != "" {
					priv.StripeCouponID = value
				}
			case "notes":
				if hasValue {
					priv.Notes = value
				}
			default:
				log.Printf("[subscription] unknown override key %q for %s", key, email)
			}
		}

		if err := s.store.UpsertUserPrivileges(priv); err != nil {
			log.Printf("[subscription] failed to save overrides for %s: %v", email, err)
			continue
		}
		log.Printf("[subscription] applied overrides for %s", email)
	}
}

// SetUserPrivileges stores operator overrides for a user by email.
func (s *Service) SetUserPrivileges(email string, unlimited bool, couponID, notes string) error {
	user, err := s.store.GetUserByEmail(strings.TrimSpace(strings.ToLower(email)))
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}
	return s.store.UpsertUserPrivileges(&models.UserPrivileges{
		UserID:          user.ID,
		UnlimitedLimits: unlimited,
		StripeCouponID:  strings.TrimSpace(couponID),
		Notes:           strings.TrimSpace(notes),
	})
}

// ApplyStripeSubscription persists tier changes received from Stripe webhooks.
func (s *Service) ApplyStripeSubscription(userID int64, tier string, startedAt int64, subscriptionID string) error {
	if !isValidTier(tier) {
		return ErrInvalidTier
	}
	var startedPtr *int64
	if tier != models.SubscriptionTierFree && startedAt > 0 {
		startedPtr = &startedAt
	}
	var subPtr *string
	if subscriptionID != "" {
		subPtr = &subscriptionID
	}
	return s.store.ApplyStripeSubscription(userID, tier, startedPtr, subPtr)
}

// BillingEnabled reports whether Stripe price IDs are configured for checkout.
func BillingEnabled() bool {
	return PriceIDForTier(models.SubscriptionTierPlus) != "" &&
		PriceIDForTier(models.SubscriptionTierPremium) != ""
}

// PriceIDForTier returns the configured Stripe price ID for a subscription tier.
func PriceIDForTier(tier string) string {
	switch tier {
	case models.SubscriptionTierPlus:
		if config.IsDevelopment() {
			if id := strings.TrimSpace(os.Getenv("STRIPE_SANDBOX_PRICE_PLUS")); id != "" {
				return id
			}
		}
		return strings.TrimSpace(os.Getenv("STRIPE_PRICE_PLUS"))
	case models.SubscriptionTierPremium:
		if config.IsDevelopment() {
			if id := strings.TrimSpace(os.Getenv("STRIPE_SANDBOX_PRICE_PREMIUM")); id != "" {
				return id
			}
		}
		return strings.TrimSpace(os.Getenv("STRIPE_PRICE_PREMIUM"))
	default:
		return ""
	}
}

// TierForPriceID maps a Stripe price ID back to an internal tier.
func TierForPriceID(priceID string) string {
	priceID = strings.TrimSpace(priceID)
	if priceID == "" {
		return ""
	}
	if priceID == PriceIDForTier(models.SubscriptionTierPlus) {
		return models.SubscriptionTierPlus
	}
	if priceID == PriceIDForTier(models.SubscriptionTierPremium) {
		return models.SubscriptionTierPremium
	}
	return ""
}

// ParseUserIDOverride reads SUBSCRIPTION_OVERRIDE_USER_IDS for numeric user IDs.
// Format: "1=unlimited,coupon=abc;2=unlimited"
func ParseUserIDOverrides(store *storage.Storage) {
	raw := strings.TrimSpace(os.Getenv("SUBSCRIPTION_OVERRIDE_USER_IDS"))
	if raw == "" {
		return
	}
	for _, entry := range strings.Split(raw, ";") {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		idPart, attrs, ok := strings.Cut(entry, "=")
		if !ok {
			continue
		}
		userID, err := strconv.ParseInt(strings.TrimSpace(idPart), 10, 64)
		if err != nil || userID <= 0 {
			continue
		}
		priv := &models.UserPrivileges{UserID: userID}
		for _, attr := range strings.Split(attrs, ",") {
			key, value, hasValue := strings.Cut(strings.TrimSpace(attr), "=")
			switch strings.ToLower(strings.TrimSpace(key)) {
			case "unlimited":
				priv.UnlimitedLimits = true
			case "coupon":
				if hasValue {
					priv.StripeCouponID = strings.TrimSpace(value)
				}
			case "notes":
				if hasValue {
					priv.Notes = strings.TrimSpace(value)
				}
			}
		}
		if err := store.UpsertUserPrivileges(priv); err != nil {
			log.Printf("[subscription] failed user_id override for %d: %v", userID, err)
		}
	}
}
