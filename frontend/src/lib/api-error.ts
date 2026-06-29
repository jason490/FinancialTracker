export class ClientApiError extends Error {
  readonly status: number;
  readonly code: string;

  constructor(message: string, status: number, code: string) {
    super(message);
    this.name = "ClientApiError";
    this.status = status;
    this.code = code;
  }
}

export const PLAID_SYNC_LAG_HINT =
  "Bank data is usually a few hours to a day behind. Manual sync runs at most once per minute.";

export const SYNC_RATE_LIMIT_MESSAGE =
  "You already synced within the last minute. Bank data may be a few hours to a day behind—wait and try again.";

export function isSyncRateLimited(err: unknown): boolean {
  return (
    err instanceof ClientApiError &&
    (err.code === "connection_sync_rate_limit" || err.code === "plaid_sync_rate_limit")
  );
}

export function syncErrorMessage(err: unknown, fallback = "Sync failed."): string {
  if (err instanceof ClientApiError) {
    return err.message;
  }
  if (err instanceof Error) {
    return err.message;
  }
  return fallback;
}

export type StatusMessageType = "ok" | "error" | "info";

// reportSyncError maps sync failures to user-facing status messages.
export function reportSyncError(
  err: unknown,
  onMessage: (text: string, type: StatusMessageType) => void,
  fallback = "Sync failed."
): void {
  if (isSyncRateLimited(err)) {
    onMessage(SYNC_RATE_LIMIT_MESSAGE, "info");
    return;
  }
  onMessage(syncErrorMessage(err, fallback), "error");
}
