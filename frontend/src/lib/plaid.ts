import { clientApiRequest } from "./api";

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
            await clientApiRequest("/api/v1/connections/complete", {
              method: "POST",
              body: { public_token: publicToken },
            });
          } else if (detail.rowId) {
            await clientApiRequest(`/api/v1/connections/sync-item/${detail.rowId}`, {
              method: "POST",
            });
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
