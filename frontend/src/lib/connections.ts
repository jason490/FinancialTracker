import { clientApiRequest } from "./api";
import type {
  ConnectionsPayload,
  CreateSessionResponse,
  ProviderInfoResponse,
} from "./types";

let cachedProvider: ProviderInfoResponse | undefined;

// getConnectionProvider returns the active financial provider from the backend.
export async function getConnectionProvider(): Promise<ProviderInfoResponse> {
  if (cachedProvider) {
    return cachedProvider;
  }
  cachedProvider = await clientApiRequest<ProviderInfoResponse>("/api/v1/connections/provider");
  return cachedProvider;
}

// getConnections returns linked institutions and accounts for the active provider.
export async function getConnections(): Promise<ConnectionsPayload> {
  return clientApiRequest<ConnectionsPayload>("/api/v1/connections");
}

// createConnectionSession requests a provider link session from the backend.
async function createConnectionSession(): Promise<CreateSessionResponse> {
  return clientApiRequest<CreateSessionResponse>("/api/v1/connections/create-session", {
    method: "POST",
  });
}

// createConnectionUpdateSession requests an update/relink session for a connection.
async function createConnectionUpdateSession(rowId: string): Promise<CreateSessionResponse> {
  return clientApiRequest<CreateSessionResponse>(`/api/v1/connections/create-update-session/${rowId}`, {
    method: "POST",
  });
}

// syncAllConnections syncs every linked institution.
export async function syncAllConnections(): Promise<void> {
  await clientApiRequest("/api/v1/connections/sync", { method: "POST" });
}

// syncConnection syncs a specific institution after update mode.
export async function syncConnection(rowId: string): Promise<void> {
  await clientApiRequest(`/api/v1/connections/sync-item/${rowId}`, { method: "POST" });
}

// disconnectConnection removes an entire institution connection.
export async function disconnectConnection(rowId: string): Promise<void> {
  await clientApiRequest(`/api/v1/connections/disconnect/${rowId}`, { method: "POST" });
}

// toggleAccountVisibility flips whether an account is hidden.
export async function toggleAccountVisibility(accountId: string): Promise<boolean> {
  const result = await clientApiRequest<{ is_hidden: boolean }>(
    `/api/v1/connections/toggle-visibility/${accountId}`,
    { method: "POST" }
  );
  return result.is_hidden;
}

// removeDisconnectedAccount permanently deletes a disconnected account.
export async function removeDisconnectedAccount(accountId: string): Promise<void> {
  await clientApiRequest(`/api/v1/connections/remove-account/${accountId}`, { method: "POST" });
}

async function openProviderLink(session: CreateSessionResponse, mode: "exchange" | "sync", rowId?: string) {
  const provider = await getConnectionProvider();
  if (provider.provider === "plaid") {
    if (!session.link_token) {
      throw new Error("Plaid link token missing");
    }
    const { openPlaidLink } = await import("./plaid");
    await openPlaidLink({ token: session.link_token, mode, rowId });
    return;
  }

  if (!session.client_secret) {
    throw new Error("Stripe client secret missing");
  }
  const { startStripeConnection } = await import("./stripe");
  await startStripeConnection(session.client_secret, provider.publishable_key || "");
  if (mode === "sync" && rowId) {
    await syncConnection(rowId);
  }
}

// startNewConnection opens the active provider flow for a new institution.
export async function startNewConnection(): Promise<void> {
  const session = await createConnectionSession();
  await openProviderLink(session, "exchange");
}

// manageConnection opens update/relink mode for an institution.
export async function manageConnection(rowId: string): Promise<void> {
  const session = await createConnectionUpdateSession(rowId);
  await openProviderLink(session, "sync", rowId);
}
