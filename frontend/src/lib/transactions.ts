import { clientApiRequest } from "./api";
import type { TransactionListPayload } from "./types";

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
