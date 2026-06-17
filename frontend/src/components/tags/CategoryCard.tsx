import { For, Show, createSignal } from "solid-js";
import { CheckIcon, PlusIcon, TrashIcon, XIcon } from "~/components/icons";
import { tagColorHex } from "~/lib/tag-colors";
import { deleteTag, updateCategory } from "~/lib/tags";
import type { CategoryWithTagsView, TagView } from "~/lib/types";
import styles from "~/styles/tags.module.css";

type CategoryCardProps = {
  category: CategoryWithTagsView;
  index: number;
  onAddTag: (category: CategoryWithTagsView) => void;
  onEditTag: (tag: TagView) => void;
  onDeleteCategory: (category: CategoryWithTagsView) => void;
  onRefresh: () => void;
  onError: (message: string) => void;
};

// CategoryCard renders a category column with inline rename and tag pills.
export default function CategoryCard(props: CategoryCardProps) {
  const [editingName, setEditingName] = createSignal(false);
  const [name, setName] = createSignal(props.category.name);
  const [pending, setPending] = createSignal(false);

  const saveName = async () => {
    const trimmed = name().trim();
    if (!trimmed || trimmed === props.category.name) {
      setEditingName(false);
      setName(props.category.name);
      return;
    }

    setPending(true);
    try {
      await updateCategory(props.category.id, { name: trimmed });
      props.onRefresh();
      setEditingName(false);
    } catch (err) {
      props.onError(err instanceof Error ? err.message : "Failed to update category");
      setName(props.category.name);
    } finally {
      setPending(false);
    }
  };

  const handleDeleteTag = async (tag: TagView) => {
    if (!window.confirm(`Delete tag "${tag.name}"?`)) return;
    try {
      await deleteTag(tag.id);
      props.onRefresh();
    } catch (err) {
      props.onError(err instanceof Error ? err.message : "Failed to delete tag");
    }
  };

  return (
    <article
      class={styles.categoryCard}
      style={{ "animation-delay": `${props.index * 60}ms` }}
    >
      <div class={styles.categoryHeader}>
        <div class={styles.categoryTitleWrap}>
          <Show
            when={editingName()}
            fallback={
              <h2
                class={styles.categoryName}
                onDblClick={() => {
                  setName(props.category.name);
                  setEditingName(true);
                }}
              >
                {props.category.name}
              </h2>
            }
          >
            <input
              class={styles.categoryNameInput}
              value={name()}
              disabled={pending()}
              onInput={(event) => setName(event.currentTarget.value)}
              onKeyDown={(event) => {
                if (event.key === "Enter") void saveName();
                if (event.key === "Escape") {
                  setName(props.category.name);
                  setEditingName(false);
                }
              }}
            />
            <button
              type="button"
              class={styles.iconBtn}
              aria-label="Save category name"
              disabled={pending()}
              onClick={() => void saveName()}
            >
              <CheckIcon size={16} />
            </button>
            <button
              type="button"
              class={styles.iconBtn}
              aria-label="Cancel rename"
              onClick={() => {
                setName(props.category.name);
                setEditingName(false);
              }}
            >
              <XIcon size={16} />
            </button>
          </Show>
        </div>

        <div class={styles.categoryActions}>
          <button
            type="button"
            class={styles.iconBtn}
            aria-label={`Delete category ${props.category.name}`}
            onClick={() => props.onDeleteCategory(props.category)}
          >
            <TrashIcon size={16} />
          </button>
          <button
            type="button"
            class={styles.iconBtn}
            aria-label={`Add tag to ${props.category.name}`}
            onClick={() => props.onAddTag(props.category)}
          >
            <PlusIcon size={18} />
          </button>
        </div>
      </div>

      <div class={styles.tagList}>
        <Show
          when={props.category.tags.length > 0}
          fallback={<div class={styles.emptyTags}>No tags yet — add one</div>}
        >
          <For each={props.category.tags}>
            {(tag) => (
              <span
                class={styles.tagChip}
                style={{ "--tag-color": tagColorHex(tag.color) }}
              >
                <span class={styles.tagDot} />
                <button
                  type="button"
                  class={styles.tagNameBtn}
                  onClick={() => props.onEditTag(tag)}
                >
                  {tag.name}
                </button>
                <button
                  type="button"
                  class={styles.tagDeleteBtn}
                  aria-label={`Delete tag ${tag.name}`}
                  onClick={() => void handleDeleteTag(tag)}
                >
                  <XIcon size={14} />
                </button>
              </span>
            )}
          </For>
        </Show>
      </div>
    </article>
  );
}
