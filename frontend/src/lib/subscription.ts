import { clientApiRequest } from "./api";
import type {
  BillingPortalResponse,
  ChangeSubscriptionResponse,
  CheckoutSessionResponse,
  SubscriptionPayload,
} from "./types";

// getSubscription loads the user's plan and billing period.
export async function getSubscription(): Promise<SubscriptionPayload> {
  return clientApiRequest<SubscriptionPayload>("/api/v1/subscription");
}

// changeSubscription applies a skeleton tier change (development only).
export async function changeSubscription(tier: string): Promise<ChangeSubscriptionResponse> {
  return clientApiRequest<ChangeSubscriptionResponse>("/api/v1/subscription/change", {
    method: "POST",
    body: { tier },
  });
}

// createCheckoutSession starts Stripe Checkout for a paid tier.
export async function createCheckoutSession(tier: string): Promise<CheckoutSessionResponse> {
  return clientApiRequest<CheckoutSessionResponse>("/api/v1/subscription/checkout", {
    method: "POST",
    body: { tier },
  });
}

// createBillingPortal opens the Stripe Customer Portal for the signed-in user.
export async function createBillingPortal(): Promise<BillingPortalResponse> {
  return clientApiRequest<BillingPortalResponse>("/api/v1/subscription/portal", {
    method: "POST",
  });
}
