import { For, Show } from "solid-js";
import { EyeOffIcon } from "~/components/icons";
import type { TransactionView } from "~/lib/types";
import { formatAmount, formatDate } from "~/lib/format";
import { tagColorHex } from "~/lib/tag-colors";
import styles from "~/styles/transactions.module.css";

interface TransactionCardProps {
  transaction: TransactionView;
  selected: boolean;
  onToggleSelect: (id: number) => void;
}

export function TransactionCard(props: TransactionCardProps) {
  const amountClass = (amount: number) => {
    return amount > 0 ? styles.amountExpense : styles.amountIncome;
  };

  // A transaction is excluded from monthly totals when any of its tags is hidden.
  const isExcluded = () => props.transaction.tags?.some((tag) => tag.is_hidden) ?? false;

  return (
    <div
      classList={{
        [styles.transactionCard]: true,
        [styles.transactionCardSelected]: props.selected,
        [styles.transactionCardExcluded]: isExcluded(),
      }}
    >
      <input
        type="checkbox"
        class={styles.checkbox}
        checked={props.selected}
        onChange={() => props.onToggleSelect(props.transaction.id)}
      />
      <div class={styles.transactionInfo}>
        <p class={styles.transactionName}>
          {props.transaction.merchant_name || props.transaction.name}
        </p>
        <div class={styles.transactionMeta}>
          <span>{formatDate(props.transaction.date)}</span>
          <Show when={props.transaction.merchant_name && props.transaction.merchant_name !== props.transaction.name}>
            <span class={styles.metaDivider} />
            <span>{props.transaction.name}</span>
          </Show>
        </div>
        <Show when={props.transaction.tags && props.transaction.tags.length > 0}>
          <div class={styles.transactionTags}>
            <For each={props.transaction.tags}>
              {(tag) => (
                <span
                  classList={{ [styles.tagChip]: true, [styles.tagChipHidden]: tag.is_hidden }}
                  style={{ "--tag-color": tagColorHex(tag.color) }}
                >
                  <Show when={tag.is_hidden}>
                    <EyeOffIcon size={11} />
                  </Show>
                  {tag.name}
                </span>
              )}
            </For>
          </div>
        </Show>
      </div>
      <div class={styles.transactionRight}>
        <div class={`${styles.transactionAmount} ${amountClass(props.transaction.amount)}`}>
          {formatAmount(props.transaction.amount)}
        </div>
        <Show when={isExcluded()}>
          <span class={styles.excludedBadge} title="Excluded from income & spending totals">
            <EyeOffIcon size={11} />
            Not counted
          </span>
        </Show>
        <Show when={props.transaction.pending}>
          <span class={styles.pendingBadge}>Pending</span>
        </Show>
      </div>
    </div>
  );
}
