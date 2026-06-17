import { For, Show, Accessor } from "solid-js";
import { TransactionCard } from "./TransactionCard";
import { FilePlusIcon } from "~/components/icons";
import type { TransactionListPayload } from "~/lib/types";
import styles from "~/styles/transactions.module.css";

interface TransactionListProps {
  payload: TransactionListPayload;
  loading: boolean;
  visualLoading: Accessor<boolean>;
  selected: Accessor<Set<number>>;
  onToggleSelect: (id: number) => void;
  onToggleSelectAll: () => void;
}

export function TransactionList(props: TransactionListProps) {
  const allSelected = () => 
    props.payload.transactions.length > 0 &&
    props.payload.transactions.every((t) => props.selected().has(t.id));

  return (
    <div class={props.loading ? styles.loadingOverlay : undefined}>
      <div class={`${styles.loadingBarContainer} ${props.visualLoading() ? styles.loadingBarContainerActive : ""}`}>
        <div class={`${styles.loadingBar} ${props.visualLoading() ? styles.loadingBarActive : ""}`} />
      </div>

      {/* Select all row */}
      <Show when={props.payload.transactions.length > 0}>
        <div class={styles.selectAllRow}>
          <label>
            <input
              type="checkbox"
              class={styles.checkbox}
              checked={allSelected()}
              onChange={props.onToggleSelectAll}
            />
            Select all on page
          </label>
        </div>
      </Show>

      {/* Empty state */}
      <Show when={props.payload.transactions.length === 0}>
        <div class={styles.emptyState}>
          <FilePlusIcon class={styles.emptyIcon} />
          <h3 class={styles.emptyTitle}>No transactions found</h3>
          <p class={styles.emptyText}>
            Try adjusting your search or filters, or link a bank account to start syncing transactions.
          </p>
        </div>
      </Show>

      {/* Transaction cards */}
      <div class={styles.transactionList}>
        <For each={props.payload.transactions}>
          {(txn) => (
            <TransactionCard
              transaction={txn}
              selected={props.selected().has(txn.id)}
              onToggleSelect={props.onToggleSelect}
            />
          )}
        </For>
      </div>
    </div>
  );
}
