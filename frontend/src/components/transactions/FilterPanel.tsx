import { For, Show, Accessor } from "solid-js";
import { XIcon } from "~/components/icons";
import { tagColorHex } from "~/lib/tag-colors";
import type { TransactionListPayload } from "~/lib/types";
import styles from "~/styles/transactions.module.css";

interface FilterPanelProps {
  metadata: TransactionListPayload | undefined;
  categoryId: Accessor<number | undefined>;
  setCategoryId: (id: number | undefined) => void;
  sortBy: Accessor<string>;
  setSortBy: (val: string) => void;
  sortDir: Accessor<string>;
  setSortDir: (val: string) => void;
  startDate: Accessor<string>;
  setStartDate: (val: string) => void;
  endDate: Accessor<string>;
  setEndDate: (val: string) => void;
  minAmount: Accessor<string>;
  setMinAmount: (val: string) => void;
  maxAmount: Accessor<string>;
  setMaxAmount: (val: string) => void;
  selectedTags: Accessor<number[]>;
  toggleTag: (id: number) => void;
  clearFilters: () => void;
  closeFilters: () => void;
}

export function FilterPanel(props: FilterPanelProps) {
  return (
    <div class={styles.filterRoot}>
      <button
        type="button"
        class={styles.filterBackdrop}
        aria-label="Close filters"
        onClick={props.closeFilters}
      />
      <div
        id="transaction-filters"
        class={styles.filterPanel}
        role="dialog"
        aria-modal="true"
        aria-label="Transaction filters"
      >
        <div class={styles.filterSheetHeader}>
          <div class={styles.filterSheetHandle} aria-hidden="true" />
          <div class={styles.filterSheetTitleRow}>
            <h2 class={styles.filterSheetTitle}>Filters</h2>
            <button
              type="button"
              class={styles.filterSheetClose}
              aria-label="Close filters"
              onClick={props.closeFilters}
            >
              <XIcon />
            </button>
          </div>
        </div>

        <div class={styles.filterPanelBody}>
          {/* Category */}
          <div class={styles.filterGroup}>
            <label class={styles.filterLabel}>Category</label>
            <select
              class={styles.filterSelect}
              value={props.categoryId() ?? ""}
              onChange={(e) => {
                const v = e.currentTarget.value;
                props.setCategoryId(v ? Number(v) : undefined);
              }}
            >
              <option value="">All categories</option>
              <For each={props.metadata?.categories ?? []}>
                {(cat) => <option value={cat.id}>{cat.name}</option>}
              </For>
            </select>
          </div>

          {/* Sort */}
          <div class={styles.filterGroup}>
            <label class={styles.filterLabel}>Sort By</label>
            <select
              class={styles.filterSelect}
              value={`${props.sortBy()}_${props.sortDir()}`}
              onChange={(e) => {
                const [by, dir] = e.currentTarget.value.split("_");
                props.setSortBy(by);
                props.setSortDir(dir);
              }}
            >
              <option value="date_desc">Newest first</option>
              <option value="date_asc">Oldest first</option>
              <option value="amount_desc">Highest amount</option>
              <option value="amount_asc">Lowest amount</option>
              <option value="name_asc">Name A-Z</option>
              <option value="name_desc">Name Z-A</option>
            </select>
          </div>

          {/* Date range */}
          <div class={`${styles.filterGroup} ${styles.filterGroupFull}`}>
            <label class={styles.filterLabel}>Date Range</label>
            <div class={styles.filterRow}>
              <div class={styles.filterGroup}>
                <span class={styles.filterSubLabel}>From</span>
                <input
                  type="date"
                  class={styles.filterInput}
                  value={props.startDate()}
                  onInput={(e) => props.setStartDate(e.currentTarget.value)}
                />
              </div>
              <div class={styles.filterGroup}>
                <span class={styles.filterSubLabel}>To</span>
                <input
                  type="date"
                  class={styles.filterInput}
                  value={props.endDate()}
                  onInput={(e) => props.setEndDate(e.currentTarget.value)}
                />
              </div>
            </div>
          </div>

          {/* Amount range — filters by magnitude (absolute value) so it works for both expenses and income */}
          <div class={`${styles.filterGroup} ${styles.filterGroupFull}`}>
            <label class={styles.filterLabel}>Amount Range</label>
            <div class={styles.filterRow}>
              <div class={styles.filterGroup}>
                <span class={styles.filterSubLabel}>Min</span>
                <input
                  type="number"
                  class={styles.filterInput}
                  placeholder="0.00"
                  value={props.minAmount()}
                  onInput={(e) => {
                    const raw = e.currentTarget.value;
                    if (raw === "") {
                      props.setMinAmount("");
                      return;
                    }
                    const n = parseFloat(raw);
                    if (!Number.isFinite(n)) return;
                    props.setMinAmount(String(Math.abs(n)));
                  }}
                  min="0"
                  step="0.01"
                  inputmode="decimal"
                />
              </div>
              <div class={styles.filterGroup}>
                <span class={styles.filterSubLabel}>Max</span>
                <input
                  type="number"
                  class={styles.filterInput}
                  placeholder="0.00"
                  value={props.maxAmount()}
                  onInput={(e) => {
                    const raw = e.currentTarget.value;
                    if (raw === "") {
                      props.setMaxAmount("");
                      return;
                    }
                    const n = parseFloat(raw);
                    if (!Number.isFinite(n)) return;
                    props.setMaxAmount(String(Math.abs(n)));
                  }}
                  min="0"
                  step="0.01"
                  inputmode="decimal"
                />
              </div>
            </div>
          </div>

          {/* Tags */}
          <Show when={(props.metadata?.tags ?? []).length > 0}>
            <div class={`${styles.filterGroup} ${styles.filterGroupFull}`}>
              <label class={styles.filterLabel}>Tags</label>
              <div class={styles.tagOptions}>
                <For each={props.metadata?.tags ?? []}>
                  {(tag) => (
                    <button
                      type="button"
                      class={`${styles.tagOption} ${props.selectedTags().includes(tag.id) ? styles.tagOptionSelected : ""}`}
                      style={{ "--tag-color": tagColorHex(tag.color) }}
                      onClick={() => props.toggleTag(tag.id)}
                    >
                      {tag.name}
                    </button>
                  )}
                </For>
              </div>
            </div>
          </Show>
        </div>

        <div class={styles.filterSheetFooter}>
          <button type="button" class={styles.filterClear} onClick={props.clearFilters}>
            Clear all
          </button>
          <button type="button" class={styles.filterDone} onClick={props.closeFilters}>
            Show results
          </button>
        </div>

        <div class={styles.filterActions}>
          <button type="button" class={styles.filterClear} onClick={props.clearFilters}>
            Clear all
          </button>
        </div>
      </div>
    </div>
  );
}
