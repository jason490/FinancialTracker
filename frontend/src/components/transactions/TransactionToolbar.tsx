import { Show, Accessor } from "solid-js";
import { SearchIcon, FilterIcon } from "~/components/icons";
import styles from "~/styles/transactions.module.css";

interface TransactionToolbarProps {
  search: Accessor<string>;
  setSearch: (val: string) => void;
  filtersOpen: Accessor<boolean>;
  setFiltersOpen: (val: boolean) => void;
  activeFilterCount: Accessor<number>;
}

export function TransactionToolbar(props: TransactionToolbarProps) {
  return (
    <div class={styles.toolbar}>
      <div class={styles.searchWrap}>
        <SearchIcon class={styles.searchIcon} />
        <input
          type="text"
          class={styles.searchInput}
          placeholder="Search transactions..."
          value={props.search()}
          onInput={(e) => props.setSearch(e.currentTarget.value)}
        />
      </div>
      <button
        type="button"
        class={`${styles.filterToggle} ${props.filtersOpen() ? styles.filterToggleActive : ""}`}
        onClick={() => props.setFiltersOpen(!props.filtersOpen())}
        aria-expanded={props.filtersOpen()}
        aria-controls="transaction-filters"
      >
        <FilterIcon class={styles.filterToggleIcon} />
        <span class={styles.filterToggleLabel}>Filters</span>
        <Show when={props.activeFilterCount() > 0}>
          <span class={styles.filterBadge}>{props.activeFilterCount()}</span>
        </Show>
      </button>
    </div>
  );
}
