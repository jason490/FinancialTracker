import { clientApiRequest } from "./api";
import { getPublicApiUrl, isCapacitorClient } from "./env";
import { ClientApiError } from "./api-error";
import type { APIError, TransactionListPayload } from "./types";

export type TransactionQueryParams = {
  search?: string;
  min_amount?: number;
  max_amount?: number;
  start_date?: number;
  end_date?: number;
  category_id?: number;
  tags?: number[];
  sort_by?: string;
  sort_dir?: string;
  page?: number;
  page_size?: number;
};

// getTransactions fetches a filtered, paginated list of transactions.
export async function getTransactions(params: TransactionQueryParams = {}): Promise<TransactionListPayload> {
  const query = new URLSearchParams();
  if (params.search) query.set("search", params.search);
  if (params.min_amount != null) query.set("min_amount", String(params.min_amount));
  if (params.max_amount != null) query.set("max_amount", String(params.max_amount));
  if (params.start_date != null) query.set("start_date", String(params.start_date));
  if (params.end_date != null) query.set("end_date", String(params.end_date));
  if (params.category_id != null) query.set("category_id", String(params.category_id));
  if (params.tags?.length) query.set("tags", params.tags.join(","));
  if (params.sort_by) query.set("sort_by", params.sort_by);
  if (params.sort_dir) query.set("sort_dir", params.sort_dir);
  if (params.page != null) query.set("page", String(params.page));
  if (params.page_size != null) query.set("page_size", String(params.page_size));

  const qs = query.toString();
  const path = qs ? `/api/v1/transactions?${qs}` : "/api/v1/transactions";
  const payload = await clientApiRequest<TransactionListPayload>(path);
  return {
    ...payload,
    transactions: payload.transactions ?? [],
    tags: payload.tags ?? [],
    categories: payload.categories ?? [],
  };
}

// bulkAddTag adds a tag to multiple transactions.
export async function bulkAddTag(transactionIds: number[], tagId: number): Promise<void> {
  await clientApiRequest("/api/v1/transactions/bulk-add-tag", {
    method: "POST",
    body: { transaction_ids: transactionIds.map(String), tag_id: tagId },
  });
}

// bulkRemoveTag removes a tag from multiple transactions.
export async function bulkRemoveTag(transactionIds: number[], tagId: number): Promise<void> {
  await clientApiRequest("/api/v1/transactions/bulk-remove-tag", {
    method: "POST",
    body: { transaction_ids: transactionIds.map(String), tag_id: tagId },
  });
}

export type TransactionExportResult = {
  blob: Blob;
  filename: string;
};

// buildExportUrl resolves the full URL for the CSV export endpoint, honoring
// the same browser vs. Capacitor routing rules as the JSON client.
function buildExportUrl(): string {
  const path = "/api/v1/transactions/export";
  if (isCapacitorClient()) {
    const base = getPublicApiUrl().replace(/\/$/, "");
    return `${base}${path}`;
  }
  if (typeof window !== "undefined") {
    return path;
  }
  const base = getPublicApiUrl().replace(/\/$/, "");
  return `${base}${path}`;
}

// filenameFromDisposition extracts the suggested filename from a
// Content-Disposition header, falling back to a date-stamped default.
function filenameFromDisposition(header: string | null): string {
  if (!header) return defaultExportFilename();
  const match = /filename="?([^";]+)"?/i.exec(header);
  return match?.[1] ?? defaultExportFilename();
}

function defaultExportFilename(): string {
  const now = new Date();
  const stamp = now.toISOString().slice(0, 10);
  return `financial-tracker-transactions-${stamp}.csv`;
}

// fetchTransactionsExport downloads the user's full transaction history as a
// CSV blob, including assigned tags and the categories those tags belong to.
export async function fetchTransactionsExport(): Promise<TransactionExportResult> {
  const url = buildExportUrl();
  const response = await fetch(url, {
    method: "GET",
    credentials: "include",
    headers: { Accept: "text/csv" },
  });

  if (!response.ok) {
    let message = `Export failed with status ${response.status}`;
    let code = "export_failed";
    try {
      const body = (await response.json()) as APIError;
      message = body.message || message;
      code = body.code || code;
    } catch {
      // Non-JSON error body; keep defaults.
    }
    throw new ClientApiError(message, response.status, code);
  }

  const blob = await response.blob();
  const filename = filenameFromDisposition(response.headers.get("Content-Disposition"));
  return { blob, filename };
}
