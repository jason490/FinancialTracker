import { For, Index } from "solid-js";
import { PlusIcon, XIcon } from "~/components/icons";
import type { TagFilterView } from "~/lib/types";
import styles from "~/styles/tags.module.css";

const FILTER_TYPES = [
  { value: "string", label: "String Match" },
  { value: "regex", label: "Regex" },
  { value: "amount_greater", label: "Amount >" },
  { value: "amount_less", label: "Amount <" },
  { value: "amount_equal", label: "Amount =" },
];

type FilterEditorProps = {
  filters: TagFilterView[];
  onChange: (filters: TagFilterView[]) => void;
};

// FilterEditor manages auto-tagging rules for a tag.
export default function FilterEditor(props: FilterEditorProps) {
  const addFilter = () => {
    props.onChange([...props.filters, { pattern: "", filter_type: "string" }]);
  };

  const updateFilter = (index: number, patch: Partial<TagFilterView>) => {
    props.onChange(
      props.filters.map((filter, i) => (i === index ? { ...filter, ...patch } : filter))
    );
  };

  const removeFilter = (index: number) => {
    props.onChange(props.filters.filter((_, i) => i !== index));
  };

  return (
    <div class={styles.field}>
      <div class={styles.sectionHeader}>
        <span class={styles.label}>Auto-tagging Filters</span>
        <button type="button" class={styles.addFilterBtn} onClick={addFilter}>
          <PlusIcon size={14} />
          Add Filter
        </button>
      </div>

      <div class={styles.filterList}>
        <Index each={props.filters}>
          {(filter, index) => (
            <div class={styles.filterRow}>
              <select
                class={styles.select}
                value={filter().filter_type}
                onChange={(event) =>
                  updateFilter(index, { filter_type: event.currentTarget.value })
                }
              >
                <For each={FILTER_TYPES}>
                  {(option) => <option value={option.value}>{option.label}</option>}
                </For>
              </select>
              <input
                class={styles.input}
                type="text"
                placeholder="Pattern or value"
                value={filter().pattern}
                onInput={(event) =>
                  updateFilter(index, { pattern: event.currentTarget.value })
                }
              />
              <button
                type="button"
                class={styles.iconBtn}
                aria-label="Remove filter"
                onClick={() => removeFilter(index)}
              >
                <XIcon size={16} />
              </button>
            </div>
          )}
        </Index>
        {props.filters.length === 0 && (
          <div class={styles.filterEmpty}>No filters defined for this tag.</div>
        )}
      </div>
    </div>
  );
}
