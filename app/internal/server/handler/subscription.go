package handler

import (
	"FinancialTracker/internal/config"
	"FinancialTracker/internal/models"
	"FinancialTracker/internal/models/external"
	"FinancialTracker/internal/services/financial"
	stripebilling "FinancialTracker/internal/services/stripebilling"
	subscriptionService "FinancialTracker/internal/services/subscription"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
)

// SubscriptionHandler serves JSON subscription endpoints for the SPA.
type SubscriptionHandler struct {
	subscriptionService *subscriptionService.Service
	billingService      *stripebilling.Service
}

// NewSubscriptionHandler creates a SubscriptionHandler.
func NewSubscriptionHandler(
	subscriptionService *subscriptionService.Service,
	billingService *stripebilling.Service,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
		billingService:      billingService,
	}
}

// HandleGetSubscription returns the user's plan, billing period, and available tiers.
func (h *SubscriptionHandler) HandleGetSubscription(c *echo.Context) error {
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	payload, err := h.buildPayload(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("subscription_error", "Failed to load subscription"))
	}

	return c.JSON(http.StatusOK, payload)
}

// HandleChangeSubscription applies a skeleton tier change (development only until Stripe is integrated).
func (h *SubscriptionHandler) HandleChangeSubscription(c *echo.Context) error {
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	var req external.ChangeSubscriptionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid request body"))
	}
	if req.Tier == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Tier is required"))
	}

	if err := h.subscriptionService.ChangeTier(userID, req.Tier); err != nil {
		switch {
		case errors.Is(err, subscriptionService.ErrSubscriptionsDisabled):
			return c.JSON(http.StatusForbidden, ErrorResponse("subscriptions_disabled", "Subscriptions are disabled on this server"))
		case errors.Is(err, subscriptionService.ErrInvalidTier):
			return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_tier", err.Error()))
		case errors.Is(err, subscriptionService.ErrAlreadyOnTier):
			return c.JSON(http.StatusBadRequest, ErrorResponse("already_on_tier", err.Error()))
		case errors.Is(err, subscriptionService.ErrStripeRequired):
			return c.JSON(http.StatusNotImplemented, ErrorResponse("stripe_required", "Use Stripe checkout to change plans in production"))
		default:
			return c.JSON(http.StatusInternalServerError, ErrorResponse("subscription_error", err.Error()))
		}
	}

	payload, err := h.buildPayload(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("subscription_error", "Plan updated but failed to reload subscription"))
	}

	return c.JSON(http.StatusOK, external.ChangeSubscriptionResponse{Subscription: *payload})
}

// HandleCreateCheckoutSession starts a Stripe Checkout subscription flow.
func (h *SubscriptionHandler) HandleCreateCheckoutSession(c *echo.Context) error {
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	var req external.CheckoutSessionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid request body"))
	}
	if req.Tier == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Tier is required"))
	}

	successURL := billingReturnURL("success")
	cancelURL := billingReturnURL("cancelled")
	url, err := h.billingService.CreateCheckoutSession(c.Request().Context(), userID, req.Tier, successURL, cancelURL)
	if err != nil {
		switch {
		case errors.Is(err, stripebilling.ErrBillingNotConfigured):
			return c.JSON(http.StatusServiceUnavailable, ErrorResponse("billing_not_configured", "Stripe billing is not configured yet"))
		case errors.Is(err, stripebilling.ErrPaidTierRequired):
			return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_tier", "Checkout requires a paid plan"))
		case errors.Is(err, subscriptionService.ErrInvalidTier):
			return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_tier", err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, ErrorResponse("checkout_error", "Failed to start checkout"))
		}
	}

	return c.JSON(http.StatusOK, external.CheckoutSessionResponse{URL: url})
}

// HandleCreateBillingPortal opens the Stripe Customer Portal for the signed-in user.
func (h *SubscriptionHandler) HandleCreateBillingPortal(c *echo.Context) error {
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	returnURL := billingReturnURL("")
	url, err := h.billingService.CreatePortalSession(c.Request().Context(), userID, returnURL)
	if err != nil {
		switch {
		case errors.Is(err, stripebilling.ErrBillingNotConfigured):
			return c.JSON(http.StatusServiceUnavailable, ErrorResponse("billing_not_configured", "Stripe billing is not configured yet"))
		case errors.Is(err, stripebilling.ErrNoStripeCustomer):
			return c.JSON(http.StatusBadRequest, ErrorResponse("no_billing_account", "No billing account found. Upgrade through checkout first."))
		default:
			return c.JSON(http.StatusInternalServerError, ErrorResponse("portal_error", "Failed to open billing portal"))
		}
	}

	return c.JSON(http.StatusOK, external.BillingPortalResponse{URL: url})
}

// HandleStripeWebhook receives Stripe billing events and updates local subscription state.
func (h *SubscriptionHandler) HandleStripeWebhook(c *echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil || len(body) == 0 {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Missing webhook body"))
	}

	signature := c.Request().Header.Get("Stripe-Signature")
	if err := h.billingService.HandleWebhook(body, signature); err != nil {
		switch {
		case errors.Is(err, stripebilling.ErrWebhookNotConfigured):
			return c.JSON(http.StatusServiceUnavailable, ErrorResponse("webhook_not_configured", "Stripe webhook secret is not configured"))
		case errors.Is(err, stripebilling.ErrWebhookVerification):
			return c.JSON(http.StatusBadRequest, ErrorResponse("webhook_verification_failed", "Webhook verification failed"))
		default:
			log.Errorf("Stripe webhook handling failed: %v", err)
		}
	}

	return c.NoContent(http.StatusOK)
}

func (h *SubscriptionHandler) buildPayload(userID int64) (*external.SubscriptionPayload, error) {
	sub, err := h.subscriptionService.GetUserSubscription(userID)
	if err != nil {
		return nil, err
	}
	period, err := h.subscriptionService.CurrentBillingPeriod(userID)
	if err != nil {
		return nil, err
	}
	limits, err := h.subscriptionService.EffectiveLimits(userID)
	if err != nil {
		return nil, err
	}
	privileges, err := h.subscriptionService.PrivilegesForUser(userID)
	if err != nil {
		return nil, err
	}

	tier := sub.Tier
	if tier == "" {
		tier = models.SubscriptionTierFree
	}

	return &external.SubscriptionPayload{
		Tier:                  tier,
		Billing:               *period,
		Limits:                limits,
		Plans:                 subscriptionService.TierCatalog,
		StripeConfigured:      config.SubscriptionsEnabled() && financial.StripeConfigured(),
		BillingEnabled:        stripebilling.BillingReady(),
		SubscriptionsEnabled:  config.SubscriptionsEnabled(),
		HasActiveSubscription: sub.StripeSubscriptionID != "",
		CanChangePlan:         h.subscriptionService.CanChangePlanDirectly(),
		Privileges:            privileges,
	}, nil
}

func billingReturnURL(checkoutStatus string) string {
	base := strings.TrimRight(os.Getenv("FRONTEND_URL"), "/")
	if base == "" {
		base = "http://localhost"
	}
	url := base + "/settings?tab=plan"
	if checkoutStatus != "" {
		url += "&checkout=" + checkoutStatus
	}
	return url
}
