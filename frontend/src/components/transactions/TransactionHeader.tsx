import { Show, JSX } from "solid-js";
import styles from "~/styles/transactions.module.css";

interface TransactionHeaderProps {
  totalCount?: number;
  children?: JSX.Element;
}

export function TransactionHeader(props: TransactionHeaderProps) {
  return (
    <header class={styles.header}>
      <div class={styles.headerTop}>
        <div>
          <p class={styles.eyebrow}>Activity</p>
          <h1 class={styles.title}>Transactions</h1>
          <p class={styles.subtitle}>
            Search, filter, and manage all your bank transactions in one place.
          </p>
          <Show when={props.totalCount !== undefined}>
            <p class={styles.resultCount}>
              <strong>{props.totalCount}</strong> transaction{props.totalCount !== 1 ? "s" : ""} found
            </p>
          </Show>
        </div>
        <div class={styles.headerActions}>
          {props.children}
        </div>
      </div>
    </header>
  );
}
