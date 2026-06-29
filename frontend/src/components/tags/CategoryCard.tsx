import { For, Show, createSignal } from "solid-js";
import { CheckIcon, PlusIcon, TrashIcon, XIcon } from "~/components/icons";
import { tagColorHex } from "~/lib/tag-colors";
import { deleteTag, updateCategory } from "~/lib/tags";
import type { CategoryWithTagsView, TagView } from "~/lib/types";
import styles from "~/styles/tags.module.css";

type CategoryCardProps = {
  category: CategoryWithTagsView;
  index: number;
  draggingTag: TagView | null;
  hoverCategoryId: number | null;
  onAddTag: (category: CategoryWithTagsView) => void;
  onEditTag: (tag: TagView) => void;
  onDeleteCategory: (category: CategoryWithTagsView) => void;
  onRefresh: () => void;
  onError: (message: string) => void;
  onTagPointerDown: (tag: TagView, event: PointerEvent) => void;
};

// CategoryCard renders a category column with inline rename, tag pills, and
// drag-and-drop support for moving tags between categories.
export default function CategoryCard(props: CategoryCardProps) {
  // ── Rename state ───────────────────────────────────────────
  const [editingName, setEditingName] = createSignal(false);
  const [name, setName] = createSignal(props.category.name);
  const [renamePending, setRenamePending] = createSignal(false);

  const cancelEdit = () => {
    setName(props.category.name);
    setEditingName(false);
  };

  const startEdit = () => {
    setName(props.category.name);
    setEditingName(true);
  };

  const saveName = async () => {
    const trimmed = name().trim();
    if (!trimmed || trimmed === props.category.name) {
      cancelEdit();
      return;
    }

    setRenamePending(true);
    try {
      await updateCategory(props.category.id, { name: trimmed });
      props.onRefresh();
      setEditingName(false);
    } catch (err) {
      props.onError(err instanceof Error ? err.message : "Failed to update category");
      setName(props.category.name);
    } finally {
      setRenamePending(false);
    }
  };

  // ── Tag deletion ───────────────────────────────────────────
  const handleDeleteTag = async (tag: TagView) => {
    if (!window.confirm(`Delete tag "${tag.name}"?`)) return;
    try {
      await deleteTag(tag.id);
      props.onRefresh();
    } catch (err) {
      props.onError(err instanceof Error ? err.message : "Failed to delete tag");
    }
  };

  // ── Drag & drop ────────────────────────────────────────────
  // Drag state is owned by the page-level `useTagDrag` controller; this card
  // only derives its visual state from the props it receives.
  const isDropTarget = () =>
    props.draggingTag != null && props.draggingTag.category_id !== props.category.id;

  const isDropActive = () => isDropTarget() && props.hoverCategoryId === props.category.id;

  return (
    <article
      classList={{
        [styles.categoryCard]: true,
        [styles.categoryCardDropTarget]: isDropTarget(),
        [styles.categoryCardDropActive]: isDropActive(),
      }}
      style={{ "animation-delay": `${props.index * 60}ms` }}
      data-category-id={props.category.id}
    >
      <div class={styles.categoryHeader}>
        <div class={styles.categoryTitleWrap}>
          <Show
            when={editingName()}
            fallback={
              <h2 class={styles.categoryName} onDblClick={startEdit}>
                {props.category.name}
              </h2>
            }
          >
            <input
              class={styles.categoryNameInput}
              value={name()}
              disabled={renamePending()}
              onInput={(event) => setName(event.currentTarget.value)}
              onKeyDown={(event) => {
                if (event.key === "Enter") void saveName();
                else if (event.key === "Escape") cancelEdit();
              }}
            />
            <button
              type="button"
              class={styles.iconBtn}
              aria-label="Save category name"
              disabled={renamePending()}
              onClick={() => void saveName()}
            >
              <CheckIcon size={16} />
            </button>
            <button
              type="button"
              class={styles.iconBtn}
              aria-label="Cancel rename"
              onClick={cancelEdit}
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
                classList={{
                  [styles.tagChip]: true,
                  [styles.tagChipDragging]: props.draggingTag?.id === tag.id,
                }}
                style={{ "--tag-color": tagColorHex(tag.color) }}
                title="Drag to another category to move this tag"
                onPointerDown={(event) => props.onTagPointerDown(tag, event)}
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
