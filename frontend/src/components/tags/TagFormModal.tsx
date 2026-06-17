import { createEffect, createSignal, Show } from "solid-js";
import Modal from "~/components/tags/Modal";
import FilterEditor from "~/components/tags/FilterEditor";
import TagColorPicker from "~/components/tags/TagColorPicker";
import { defaultTagColor } from "~/lib/tag-colors";
import { createTag, getTagFilters, updateTag } from "~/lib/tags";
import type { CategoryWithTagsView, TagFilterView, TagView } from "~/lib/types";
import styles from "~/styles/tags.module.css";

type TagFormModalProps = {
  open: boolean;
  mode: "create" | "edit";
  category?: CategoryWithTagsView;
  tag?: TagView;
  categories: CategoryWithTagsView[];
  onClose: () => void;
  onSaved: () => void;
  onError: (message: string) => void;
};

// TagFormModal handles creating and editing tags with optional filters.
export default function TagFormModal(props: TagFormModalProps) {
  const [name, setName] = createSignal("");
  const [color, setColor] = createSignal(defaultTagColor());
  const [categoryId, setCategoryId] = createSignal(0);
  const [filters, setFilters] = createSignal<TagFilterView[]>([]);
  const [loadingFilters, setLoadingFilters] = createSignal(false);
  const [pending, setPending] = createSignal(false);
  const [error, setError] = createSignal<string | null>(null);

  createEffect(() => {
    if (!props.open) return;

    setError(null);
    if (props.mode === "create") {
      setName("");
      setColor(defaultTagColor());
      setCategoryId(props.category?.id ?? props.categories[0]?.id ?? 0);
      setFilters([]);
      return;
    }

    const tag = props.tag;
    if (!tag) return;

    setName(tag.name);
    setColor(tag.color || defaultTagColor());
    setCategoryId(tag.category_id);
    setFilters([]);
    setLoadingFilters(true);
    void getTagFilters(tag.id)
      .then((loaded) => setFilters(loaded))
      .catch((err) => setError(err instanceof Error ? err.message : "Failed to load filters"))
      .finally(() => setLoadingFilters(false));
  });

  const submit = async (apply: boolean) => {
    setPending(true);
    setError(null);
    try {
      const payload = {
        name: name().trim(),
        color: color(),
        category_id: categoryId(),
        filters: filters().filter((filter) => filter.pattern.trim() !== ""),
        apply,
      };

      if (props.mode === "create") {
        await createTag(payload);
      } else if (props.tag) {
        await updateTag(props.tag.id, payload);
      }

      props.onSaved();
      props.onClose();
    } catch (err) {
      const message = err instanceof Error ? err.message : "Failed to save tag";
      setError(message);
      props.onError(message);
    } finally {
      setPending(false);
    }
  };

  const title =
    props.mode === "create"
      ? `Add Tag${props.category ? ` to ${props.category.name}` : ""}`
      : `Edit Tag: ${props.tag?.name ?? ""}`;

  return (
    <Modal open={props.open} title={title} onClose={props.onClose}>
      <div class={styles.formGrid}>
        <Show when={error()}>
          <div class={styles.statusError} role="alert">
            {error()}
          </div>
        </Show>

        <div class={styles.field}>
          <label class={styles.label} for="tag-name">
            Tag Name
          </label>
          <input
            id="tag-name"
            class={styles.input}
            type="text"
            required
            value={name()}
            onInput={(event) => setName(event.currentTarget.value)}
          />
        </div>

        <div class={styles.field}>
          <label class={styles.label} for="tag-category">
            Category
          </label>
          <select
            id="tag-category"
            class={styles.select}
            value={categoryId()}
            onChange={(event) => setCategoryId(Number(event.currentTarget.value))}
          >
            {props.categories.map((category) => (
              <option value={category.id}>{category.name}</option>
            ))}
          </select>
        </div>

        <TagColorPicker value={color()} onChange={setColor} />

        <Show
          when={!loadingFilters()}
          fallback={<div class={styles.filterEmpty}>Loading filters...</div>}
        >
          <FilterEditor filters={filters()} onChange={setFilters} />
        </Show>
      </div>

      <div class={styles.modalFooter}>
        <button type="button" class={styles.ghostBtn} onClick={props.onClose} disabled={pending()}>
          Cancel
        </button>
        <button
          type="button"
          class={styles.primaryBtn}
          disabled={pending() || !name().trim()}
          onClick={() => void submit(false)}
        >
          Save Tag
        </button>
        <button
          type="button"
          class={`${styles.primaryBtn} ${styles.successBtn}`}
          disabled={pending() || !name().trim()}
          onClick={() => void submit(true)}
        >
          Save & Apply
        </button>
      </div>
    </Modal>
  );
}
