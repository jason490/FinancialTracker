import { clientApiRequest } from "./api";
import type { PlaidConnectionsPayload } from "./types";

type LinkMode = "exchange" | "sync";

type PlaidLinkDetail = {
  token: string;
  mode: LinkMode;
  rowId?: string;
};

declare global {
  interface Window {
    Plaid?: {
      create: (config: {
        token: string;
        onSuccess: (publicToken: string) => void;
        onExit: (err: { display_message?: string; error_message?: string } | null) => void;
      }) => { open: () => void };
    };
  }
}

// getPlaidConnections returns all linked institutions and accounts.
export async function getPlaidConnections(): Promise<PlaidConnectionsPayload> {
  return clientApiRequest<PlaidConnectionsPayload>("/api/v1/plaid/connections");
}

// createPlaidLinkToken requests a token for linking a new bank.
export async function createPlaidLinkToken(): Promise<string> {
  const result = await clientApiRequest<{ link_token: string }>("/api/v1/plaid/create-link-token", {
    method: "POST",
  });
  return result.link_token;
}

// createPlaidUpdateToken requests a token for managing an existing connection.
export async function createPlaidUpdateToken(rowId: string): Promise<string> {
  const result = await clientApiRequest<{ link_token: string }>(
    `/api/v1/plaid/create-update-token/${rowId}`,
    { method: "POST" }
  );
  return result.link_token;
}

// exchangePlaidToken exchanges a public token from Plaid Link.
export async function exchangePlaidToken(publicToken: string): Promise<void> {
  await clientApiRequest("/api/v1/plaid/exchange", {
    method: "POST",
    body: { public_token: publicToken },
  });
}

// syncPlaidItem syncs a specific institution after update mode.
export async function syncPlaidItem(rowId: string): Promise<void> {
  await clientApiRequest(`/api/v1/plaid/sync-item/${rowId}`, { method: "POST" });
}

// syncAllPlaidConnections syncs every linked institution.
export async function syncAllPlaidConnections(): Promise<void> {
  await clientApiRequest("/api/v1/plaid/sync", { method: "POST" });
}

// disconnectPlaidConnection removes an entire institution connection.
export async function disconnectPlaidConnection(rowId: string): Promise<void> {
  await clientApiRequest(`/api/v1/plaid/disconnect/${rowId}`, { method: "POST" });
}

// togglePlaidAccountVisibility flips whether an account is hidden.
export async function togglePlaidAccountVisibility(accountId: string): Promise<boolean> {
  const result = await clientApiRequest<{ is_hidden: boolean }>(
    `/api/v1/plaid/toggle-visibility/${accountId}`,
    { method: "POST" }
  );
  return result.is_hidden;
}

// removeDisconnectedPlaidAccount permanently deletes a disconnected account.
export async function removeDisconnectedPlaidAccount(accountId: string): Promise<void> {
  await clientApiRequest(`/api/v1/plaid/remove-account/${accountId}`, { method: "POST" });
}

let plaidScriptPromise: Promise<void> | undefined;

// loadPlaidScript injects the Plaid Link SDK once per page load.
function loadPlaidScript(): Promise<void> {
  if (typeof window === "undefined") {
    return Promise.resolve();
  }
  if (window.Plaid) {
    return Promise.resolve();
  }
  if (!plaidScriptPromise) {
    plaidScriptPromise = new Promise((resolve, reject) => {
      const script = document.createElement("script");
      script.src = "https://cdn.plaid.com/link/v2/stable/link-initialize.js";
      script.async = true;
      script.onload = () => resolve();
      script.onerror = () => reject(new Error("Failed to load Plaid Link"));
      document.head.appendChild(script);
    });
  }
  return plaidScriptPromise;
}

// openPlaidLink launches Plaid Link for exchange or update flows.
export async function openPlaidLink(detail: PlaidLinkDetail): Promise<void> {
  await loadPlaidScript();
  if (!window.Plaid) {
    throw new Error("Plaid Link is unavailable");
  }

  return new Promise((resolve, reject) => {
    const handler = window.Plaid!.create({
      token: detail.token,
      onSuccess: async (publicToken) => {
        try {
          if (detail.mode === "exchange") {
            await exchangePlaidToken(publicToken);
          } else if (detail.rowId) {
            await syncPlaidItem(detail.rowId);
          }
          resolve();
        } catch (err) {
          reject(err);
        }
      },
      onExit: (err) => {
        if (err) {
          reject(new Error(err.display_message || err.error_message || "Plaid Link closed"));
          return;
        }
        resolve();
      },
    });
    handler.open();
  });
}

// startNewPlaidConnection opens Plaid Link for a new institution.
export async function startNewPlaidConnection(): Promise<void> {
  const token = await createPlaidLinkToken();
  await openPlaidLink({ token, mode: "exchange" });
}

// managePlaidConnection opens Plaid Link in update mode for an institution.
export async function managePlaidConnection(rowId: string): Promise<void> {
  const token = await createPlaidUpdateToken(rowId);
  await openPlaidLink({ token, mode: "sync", rowId });
}
