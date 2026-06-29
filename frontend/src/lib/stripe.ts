import type { Stripe } from "@stripe/stripe-js";
import { clientApiRequest } from "./api";
import type { CompleteConnectionRequest } from "./types";

let stripePromise: Promise<Stripe | null> | undefined;

// loadStripeClient initializes Stripe.js on demand for Financial Connections.
async function loadStripeClient(publishableKey: string): Promise<Stripe> {
  if (!publishableKey) {
    throw new Error("Stripe publishable key is not configured");
  }
  if (!stripePromise) {
    const { loadStripe } = await import("@stripe/stripe-js");
    stripePromise = loadStripe(publishableKey);
  }
  const stripe = await stripePromise;
  if (!stripe) {
    throw new Error("Stripe.js failed to load");
  }
  return stripe;
}

// startStripeConnection launches Stripe Financial Connections and completes the session.
export async function startStripeConnection(clientSecret: string, publishableKey: string): Promise<void> {
  const stripe = await loadStripeClient(publishableKey);
  const result = await stripe.collectFinancialConnectionsAccounts({ clientSecret });
  if (result.error) {
    throw new Error(result.error.message || "Stripe Financial Connections failed");
  }

  const sessionID = result.financialConnectionsSession?.id;
  if (!sessionID) {
    throw new Error("No Financial Connections session returned");
  }

  await completeConnection({ session_id: sessionID });
}

// completeConnection finalizes a provider link flow on the server.
export async function completeConnection(body: CompleteConnectionRequest): Promise<void> {
  await clientApiRequest("/api/v1/connections/complete", {
    method: "POST",
    body,
  });
}
