import { For, Show, Accessor } from "solid-js";
import styles from "~/styles/transactions.module.css";

interface PaginationProps {
  currentPage: Accessor<number>;
  totalPages: number;
  onPageChange: (page: number) => void;
}

export function Pagination(props: PaginationProps) {
  const pageNumbers = () => {
    const total = props.totalPages;
    const current = props.currentPage();
    const pages: (number | "ellipsis")[] = [];

    if (total <= 7) {
      for (let i = 1; i <= total; i++) pages.push(i);
      return pages;
    }

    pages.push(1);
    if (current > 3) pages.push("ellipsis");
    const start = Math.max(2, current - 1);
    const end = Math.min(total - 1, current + 1);
    for (let i = start; i <= end; i++) pages.push(i);
    if (current < total - 2) pages.push("ellipsis");
    pages.push(total);
    return pages;
  };

  return (
    <Show when={props.totalPages > 1}>
      <nav class={styles.pagination}>
        <button
          type="button"
          class={styles.pageBtn}
          disabled={props.currentPage() <= 1}
          onClick={() => props.onPageChange(props.currentPage() - 1)}
        >
          ‹
        </button>
        <For each={pageNumbers()}>
          {(p) => (
            <Show
              when={typeof p === "number"}
              fallback={<span class={styles.pageEllipsis}>…</span>}
            >
              <button
                type="button"
                class={`${styles.pageBtn} ${props.currentPage() === p ? styles.pageBtnActive : ""}`}
                onClick={() => props.onPageChange(p as number)}
              >
                {p}
              </button>
            </Show>
          )}
        </For>
        <button
          type="button"
          class={styles.pageBtn}
          disabled={props.currentPage() >= props.totalPages}
          onClick={() => props.onPageChange(props.currentPage() + 1)}
        >
          ›
        </button>
      </nav>
    </Show>
  );
}
