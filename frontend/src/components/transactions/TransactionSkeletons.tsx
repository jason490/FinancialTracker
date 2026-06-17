import { For } from "solid-js";
import styles from "~/styles/transactions.module.css";

export function TransactionSkeletonCard() {
  return (
    <div class={styles.skeletonCard}>
      <div class={styles.skeletonCheck} />
      <div class={styles.skeletonBody}>
        <div class={`${styles.skeletonLine} ${styles.skeletonLineMed}`} />
        <div class={`${styles.skeletonLine} ${styles.skeletonLineShort}`} />
      </div>
      <div class={styles.skeletonAmount} />
    </div>
  );
}

export function TransactionSkeletonList(props: { count?: number }) {
  return (
    <div class={styles.skeletonList}>
      <For each={Array(props.count || 5)}>{() => <TransactionSkeletonCard />}</For>
    </div>
  );
}
