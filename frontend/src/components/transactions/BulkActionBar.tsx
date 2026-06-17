import { For, Show, Accessor } from "solid-js";
import type { TransactionListPayload } from "~/lib/types";
import styles from "~/styles/transactions.module.css";

interface BulkActionBarProps {
  selectedCount: number;
  bulkAction: Accessor<"add" | "remove">;
  setBulkAction: (val: "add" | "remove") => void;
  bulkTagId: Accessor<number | undefined>;
  setBulkTagId: (val: number | undefined) => void;
  bulkLoading: Accessor<boolean>;
  metadata: TransactionListPayload | undefined;
  onApply: () => void;
  onClear: () => void;
}

export function BulkActionBar(props: BulkActionBarProps) {
  return (
    <Show when={props.selectedCount > 0}>
      <div class={styles.bulkBar}>
        <span class={styles.bulkCount}>
          {props.selectedCount} selected
        </span>
        <select
          class={styles.bulkSelect}
          value={props.bulkAction()}
          onChange={(e) => props.setBulkAction(e.currentTarget.value as "add" | "remove")}
        >
          <option value="add">Add tag</option>
          <option value="remove">Remove tag</option>
        </select>
        <select
          class={styles.bulkSelect}
          value={props.bulkTagId() ?? ""}
          onChange={(e) => {
            const v = e.currentTarget.value;
            props.setBulkTagId(v ? Number(v) : undefined);
          }}
        >
          <option value="">Choose tag…</option>
          <For each={props.metadata?.tags ?? []}>
            {(tag) => <option value={tag.id}>{tag.name}</option>}
          </For>
        </select>
        <button
          type="button"
          class={`${styles.bulkBtn} ${props.bulkAction() === "add" ? styles.bulkBtnApply : styles.bulkBtnRemove}`}
          disabled={props.bulkTagId() == null || props.bulkLoading()}
          onClick={props.onApply}
        >
          {props.bulkLoading() ? "Applying…" : "Apply"}
        </button>
        <button
          type="button"
          class={`${styles.bulkBtn} ${styles.bulkBtnClear}`}
          onClick={props.onClear}
        >
          Clear
        </button>
      </div>
    </Show>
  );
}
