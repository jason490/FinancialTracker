import { createEffect, createSignal, Show } from "solid-js";
import Modal from "~/components/tags/Modal";
import { deleteCategory } from "~/lib/tags";
import type { CategoryWithTagsView } from "~/lib/types";
import styles from "~/styles/tags.module.css";

type DeleteCategoryModalProps = {
  open: boolean;
  category?: CategoryWithTagsView;
  categories: CategoryWithTagsView[];
  onClose: () => void;
  onSaved: () => void;
  onError: (message: string) => void;
};

// DeleteCategoryModal confirms category deletion and tag disposition.
export default function DeleteCategoryModal(props: DeleteCategoryModalProps) {
  const [action, setAction] = createSignal<"move_to_misc" | "move_to" | "delete_all">("move_to_misc");
  const [targetCategoryId, setTargetCategoryId] = createSignal(0);
  const [pending, setPending] = createSignal(false);
  const [error, setError] = createSignal<string | null>(null);

  const otherCategories = () =>
    props.categories.filter((category) => category.id !== props.category?.id);

  createEffect(() => {
    if (!props.open) return;
    setAction("move_to_misc");
    setTargetCategoryId(otherCategories()[0]?.id ?? 0);
    setError(null);
  });

  const submit = async () => {
    if (!props.category) return;
    setPending(true);
    setError(null);
    try {
      await deleteCategory(props.category.id, {
        action: action(),
        target_category_id: action() === "move_to" ? targetCategoryId() : undefined,
      });
      props.onSaved();
      props.onClose();
    } catch (err) {
      const message = err instanceof Error ? err.message : "Failed to delete category";
      setError(message);
      props.onError(message);
    } finally {
      setPending(false);
    }
  };

  return (
    <Modal
      open={props.open}
      title="Delete Category"
      subtitle="What would you like to do with the tags in this category?"
      narrow
      onClose={props.onClose}
    >
      <div class={styles.formGrid}>
        <Show when={error()}>
          <div class={styles.statusError} role="alert">
            {error()}
          </div>
        </Show>

        <div class={styles.radioGroup}>
          <label class={styles.radioOption}>
            <input
              type="radio"
              name="delete-action"
              checked={action() === "move_to_misc"}
              onChange={() => setAction("move_to_misc")}
            />
            <span class={styles.radioLabel}>Move tags to "Misc" category</span>
          </label>

          <label class={styles.radioOption}>
            <input
              type="radio"
              name="delete-action"
              checked={action() === "move_to"}
              onChange={() => setAction("move_to")}
            />
            <span class={styles.radioLabel}>
              Move tags to:
              <Show when={action() === "move_to"}>
                <select
                  class={styles.select}
                  value={targetCategoryId()}
                  onChange={(event) => setTargetCategoryId(Number(event.currentTarget.value))}
                >
                  {otherCategories().map((category) => (
                    <option value={category.id}>{category.name}</option>
                  ))}
                </select>
              </Show>
            </span>
          </label>

          <label class={`${styles.radioOption} ${styles.radioOptionDanger}`}>
            <input
              type="radio"
              name="delete-action"
              checked={action() === "delete_all"}
              onChange={() => setAction("delete_all")}
            />
            <span class={styles.radioLabel}>Delete all tags in this category</span>
          </label>
        </div>
      </div>

      <div class={styles.modalFooter}>
        <button type="button" class={styles.ghostBtn} onClick={props.onClose} disabled={pending()}>
          Cancel
        </button>
        <button type="button" class={styles.dangerBtn} disabled={pending()} onClick={() => void submit()}>
          Confirm Delete
        </button>
      </div>
    </Modal>
  );
}
