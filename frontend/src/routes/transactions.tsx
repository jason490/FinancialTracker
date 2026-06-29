import { Title } from "@solidjs/meta";
import { useNavigate } from "@solidjs/router";
import {
  Show,
  createSignal,
  createResource,
  Suspense,
  createEffect,
  on,
  onCleanup,
  batch,
  useTransition,
} from "solid-js";
import { Portal } from "solid-js/web";
import AppLayout from "~/layouts/AppLayout";
import { SyncIcon } from "~/components/icons";
import PageStatusBanner, { type PageStatus } from "~/components/PageStatusBanner";
import { useAuth } from "~/lib/auth-context";
import { getTransactions, bulkAddTag, bulkRemoveTag } from "~/lib/transactions";
import { syncAllConnections } from "~/lib/connections";
import { reportSyncError } from "~/lib/api-error";
import { useMinLoadingHold } from "~/lib/loading-transition";
import type { TransactionQueryParams } from "~/lib/transactions";
import { TransactionHeader } from "~/components/transactions/TransactionHeader";
import { TransactionToolbar } from "~/components/transactions/TransactionToolbar";
import { FilterPanel } from "~/components/transactions/FilterPanel";
import { TransactionList } from "~/components/transactions/TransactionList";
import { Pagination } from "~/components/transactions/Pagination";
import { BulkActionBar } from "~/components/transactions/BulkActionBar";
import { TransactionSkeletonList } from "~/components/transactions/TransactionSkeletons";
import styles from "~/styles/transactions.module.css";

const PAGE_SIZE = 25;
const DEBOUNCE_MS = 300;

export default function TransactionsPage() {
  const navigate = useNavigate();
  const { user: profile } = useAuth();

  // ── Metadata (Categories & Tags) ────────────────────────
  const [metadata] = createResource(
    () => (profile() ? "metadata" : undefined),
    async () => {
      return getTransactions({ page_size: 1 });
    }
  );

  // ── Filter state ────────────────────────────────────────
  const [search, setSearch] = createSignal("");
  const [debouncedSearch, setDebouncedSearch] = createSignal("");
  const [categoryId, setCategoryId] = createSignal<number | undefined>();
  const [selectedTags, setSelectedTags] = createSignal<number[]>([]);
  const [startDate, setStartDate] = createSignal("");
  const [debouncedStartDate, setDebouncedStartDate] = createSignal("");
  const [endDate, setEndDate] = createSignal("");
  const [debouncedEndDate, setDebouncedEndDate] = createSignal("");
  const [minAmount, setMinAmount] = createSignal("");
  const [debouncedMinAmount, setDebouncedMinAmount] = createSignal("");
  const [maxAmount, setMaxAmount] = createSignal("");
  const [debouncedMaxAmount, setDebouncedMaxAmount] = createSignal("");
  const [sortBy, setSortBy] = createSignal("date");
  const [sortDir, setSortDir] = createSignal("desc");
  const [page, setPage] = createSignal(1);
  const [filtersOpen, setFiltersOpen] = createSignal(false);

  // ── Responsive State ─────────────────────────────────────
  const [isMobile, setIsMobile] = createSignal(false);
  createEffect(() => {
    const mql = window.matchMedia("(max-width: 767px)");
    const handler = (e: MediaQueryListEvent | MediaQueryList) => setIsMobile(e.matches);
    handler(mql);
    mql.addEventListener("change", handler);
    onCleanup(() => mql.removeEventListener("change", handler));
  });

  const [pending, startTransition] = useTransition();
  const visualLoading = useMinLoadingHold(pending, 250);

  // ── Selection state ─────────────────────────────────────
  const [selected, setSelected] = createSignal<Set<number>>(new Set());
  const [bulkTagId, setBulkTagId] = createSignal<number | undefined>();
  const [bulkAction, setBulkAction] = createSignal<"add" | "remove">("add");
  const [bulkLoading, setBulkLoading] = createSignal(false);
  const [syncing, setSyncing] = createSignal(false);
  const [message, setMessage] = createSignal<PageStatus | null>(null);

  const handleSync = async () => {
    setSyncing(true);
    try {
      await syncAllConnections();
      refetch();
      setMessage({ text: "All connections synced.", type: "ok" });
    } catch (err) {
      reportSyncError(err, (text, type) => setMessage({ text, type }));
    } finally {
      setSyncing(false);
    }
  };

  // ── Filter debouncing ───────────────────────────────────
  let debounceTimer: number | undefined;
  
  createEffect(() => {
    const s = search();
    const sd = startDate();
    const ed = endDate();
    const min = minAmount();
    const max = maxAmount();

    clearTimeout(debounceTimer);
    debounceTimer = window.setTimeout(() => {
      startTransition(() => {
        batch(() => {
          setDebouncedSearch(s);
          setDebouncedStartDate(sd);
          setDebouncedEndDate(ed);
          setDebouncedMinAmount(min);
          setDebouncedMaxAmount(max);
          setPage(1);
        });
      });
    }, DEBOUNCE_MS);
  });

  onCleanup(() => clearTimeout(debounceTimer));

  // ── Build query params reactively ───────────────────────
  const queryParams = (): TransactionQueryParams => {
    const params: TransactionQueryParams = {
      page: page(),
      page_size: PAGE_SIZE,
      sort_by: sortBy(),
      sort_dir: sortDir(),
    };
    const s = debouncedSearch().trim();
    if (s) params.search = s;
    const cat = categoryId();
    if (cat != null) params.category_id = cat;
    const tags = selectedTags();
    if (tags.length) params.tags = tags;
    
    if (debouncedStartDate()) {
      const d = new Date(debouncedStartDate() + "T00:00:00Z");
      if (!isNaN(d.getTime())) params.start_date = Math.floor(d.getTime() / 1000);
    }
    if (debouncedEndDate()) {
      const d = new Date(debouncedEndDate() + "T23:59:59Z");
      if (!isNaN(d.getTime())) params.end_date = Math.floor(d.getTime() / 1000);
    }

    const minA = debouncedMinAmount();
    if (minA !== "") {
      const n = parseFloat(minA);
      if (Number.isFinite(n)) params.min_amount = Math.abs(n);
    }
    const maxA = debouncedMaxAmount();
    if (maxA !== "") {
      const n = parseFloat(maxA);
      if (Number.isFinite(n)) params.max_amount = Math.abs(n);
    }
    return params;
  };

  // ── Fetch data ──────────────────────────────────────────
  const [data, { refetch }] = createResource(
    () => (profile() ? queryParams() : undefined),
    (params) => getTransactions(params)
  );

  const updateFilter = (fn: () => void) => {
    startTransition(() => {
      batch(() => {
        fn();
        setPage(1);
      });
    });
  };

  // ── Selection helpers ───────────────────────────────────
  const toggleSelect = (id: number) => {
    const next = new Set(selected());
    if (next.has(id)) next.delete(id);
    else next.add(id);
    setSelected(next);
  };

  const toggleSelectAll = () => {
    const txns = data()?.transactions ?? [];
    const allOnPage = txns.map((t) => t.id);
    const current = selected();
    const allSelected = allOnPage.length > 0 && allOnPage.every((id) => current.has(id));
    if (allSelected) {
      setSelected(new Set());
    } else {
      setSelected(new Set(allOnPage));
    }
  };

  const clearSelection = () => setSelected(new Set());

  // ── Tag toggle for filters ─────────────────────────────
  const toggleTag = (id: number) => {
    updateFilter(() => {
      const current = selectedTags();
      if (current.includes(id)) {
        setSelectedTags(current.filter((t) => t !== id));
      } else {
        setSelectedTags([...current, id]);
      }
    });
  };

  // ── Clear all filters ──────────────────────────────────
  const clearFilters = () => {
    updateFilter(() => {
      setSearch("");
      setDebouncedSearch("");
      setCategoryId(undefined);
      setSelectedTags([]);
      setStartDate("");
      setDebouncedStartDate("");
      setEndDate("");
      setDebouncedEndDate("");
      setMinAmount("");
      setDebouncedMinAmount("");
      setMaxAmount("");
      setDebouncedMaxAmount("");
      setSortBy("date");
      setSortDir("desc");
    });
  };

  const activeFilterCount = () => {
    let count = 0;
    if (categoryId() != null) count++;
    if (selectedTags().length > 0) count++;
    if (startDate()) count++;
    if (endDate()) count++;
    if (minAmount() !== "") count++;
    if (maxAmount() !== "") count++;
    if (sortBy() !== "date" || sortDir() !== "desc") count++;
    return count;
  };

  const closeFilters = () => setFiltersOpen(false);

  createEffect(() => {
    if (!filtersOpen()) return;
    const mq = window.matchMedia("(max-width: 767px)");
    if (!mq.matches) return;

    const prev = document.body.style.overflow;
    document.body.style.overflow = "hidden";
    onCleanup(() => {
      document.body.style.overflow = prev;
    });
  });

  // ── Bulk operations ─────────────────────────────────────
  const handleBulkApply = async () => {
    const tagId = bulkTagId();
    if (tagId == null) return;
    const ids = Array.from(selected());
    if (!ids.length) return;

    setBulkLoading(true);
    try {
      if (bulkAction() === "add") {
        await bulkAddTag(ids, tagId);
      } else {
        await bulkRemoveTag(ids, tagId);
      }
      clearSelection();
      refetch();
    } catch (err) {
      console.error("Bulk operation failed:", err);
    } finally {
      setBulkLoading(false);
    }
  };

  return (
    <AppLayout>
      <Title>Transactions | Financial Tracker</Title>

      <div class={styles.page}>
        <PageStatusBanner message={message} onDismiss={() => setMessage(null)} />
        <Suspense fallback={<TransactionSkeletonList count={8} />}>
          <Show when={metadata()}>
            <TransactionHeader totalCount={data.latest?.total_count}>
              <button
                type="button"
                class={styles.primaryButton}
                disabled={syncing()}
                onClick={() => void handleSync()}
              >
                <Show when={syncing()} fallback={<SyncIcon size={18} style={{ "margin-right": "0.5rem" }} />}>
                  <span class={styles.loadingSpinner} style={{ "margin-right": "0.5rem" }} />
                </Show>
                {syncing() ? "Syncing..." : "Sync transactions"}
              </button>
            </TransactionHeader>

            <TransactionToolbar
              search={search}
              setSearch={setSearch}
              filtersOpen={filtersOpen}
              setFiltersOpen={setFiltersOpen}
              activeFilterCount={activeFilterCount}
            />

            <Show when={filtersOpen()}>
              <Show
                when={isMobile()}
                fallback={
                  <FilterPanel
                    metadata={metadata()}
                    categoryId={categoryId}
                    setCategoryId={(id) => updateFilter(() => setCategoryId(id))}
                    sortBy={sortBy}
                    setSortBy={(val) => updateFilter(() => setSortBy(val))}
                    sortDir={sortDir}
                    setSortDir={(val) => updateFilter(() => setSortDir(val))}
                    startDate={startDate}
                    setStartDate={setStartDate}
                    endDate={endDate}
                    setEndDate={setEndDate}
                    minAmount={minAmount}
                    setMinAmount={setMinAmount}
                    maxAmount={maxAmount}
                    setMaxAmount={setMaxAmount}
                    selectedTags={selectedTags}
                    toggleTag={toggleTag}
                    clearFilters={clearFilters}
                    closeFilters={closeFilters}
                  />
                }
              >
                <Portal>
                  <FilterPanel
                    metadata={metadata()}
                    categoryId={categoryId}
                    setCategoryId={(id) => updateFilter(() => setCategoryId(id))}
                    sortBy={sortBy}
                    setSortBy={(val) => updateFilter(() => setSortBy(val))}
                    sortDir={sortDir}
                    setSortDir={(val) => updateFilter(() => setSortDir(val))}
                    startDate={startDate}
                    setStartDate={setStartDate}
                    endDate={endDate}
                    setEndDate={setEndDate}
                    minAmount={minAmount}
                    setMinAmount={setMinAmount}
                    maxAmount={maxAmount}
                    setMaxAmount={setMaxAmount}
                    selectedTags={selectedTags}
                    toggleTag={toggleTag}
                    clearFilters={clearFilters}
                    closeFilters={closeFilters}
                  />
                </Portal>
              </Show>
            </Show>

            <Suspense fallback={<TransactionSkeletonList count={5} />}>
              <Show when={data()}>
                {(payload) => (
                  <>
                    <TransactionList
                      payload={payload()}
                      loading={data.loading}
                      visualLoading={visualLoading}
                      selected={selected}
                      onToggleSelect={toggleSelect}
                      onToggleSelectAll={toggleSelectAll}
                    />

                    <Pagination
                      currentPage={page}
                      totalPages={payload().total_pages}
                      onPageChange={(p) => startTransition(() => setPage(p))}
                    />
                  </>
                )}
              </Show>
            </Suspense>

            <BulkActionBar
              selectedCount={selected().size}
              bulkAction={bulkAction}
              setBulkAction={setBulkAction}
              bulkTagId={bulkTagId}
              setBulkTagId={setBulkTagId}
              bulkLoading={bulkLoading}
              metadata={metadata()}
              onApply={() => void handleBulkApply()}
              onClear={clearSelection}
            />
          </Show>
        </Suspense>
      </div>
    </AppLayout>
  );
}

