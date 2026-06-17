import { createEffect, createSignal } from "solid-js";
import Modal from "~/components/tags/Modal";
import { createCategory } from "~/lib/tags";
import styles from "~/styles/tags.module.css";

type NewCategoryModalProps = {
  open: boolean;
  onClose: () => void;
  onSaved: () => void;
  onError: (message: string) => void;
};

// NewCategoryModal creates a new tag category.
export default function NewCategoryModal(props: NewCategoryModalProps) {
  const [name, setName] = createSignal("");
  const [pending, setPending] = createSignal(false);
  const [error, setError] = createSignal<string | null>(null);

  createEffect(() => {
    if (props.open) {
      setName("");
      setError(null);
    }
  });

  const submit = async () => {
    setPending(true);
    setError(null);
    try {
      await createCategory({ name: name().trim() });
      props.onSaved();
      props.onClose();
    } catch (err) {
      const message = err instanceof Error ? err.message : "Failed to create category";
      setError(message);
      props.onError(message);
    } finally {
      setPending(false);
    }
  };

  return (
    <Modal open={props.open} title="New Category" narrow onClose={props.onClose}>
      <div class={styles.formGrid}>
        {error() && (
          <div class={styles.statusError} role="alert">
            {error()}
          </div>
        )}
        <div class={styles.field}>
          <label class={styles.label} for="category-name">
            Category Name
          </label>
          <input
            id="category-name"
            class={styles.input}
            type="text"
            required
            value={name()}
            onInput={(event) => setName(event.currentTarget.value)}
          />
        </div>
      </div>
      <div class={styles.modalFooter}>
        <button type="button" class={styles.ghostBtn} onClick={props.onClose} disabled={pending()}>
          Cancel
        </button>
        <button
          type="button"
          class={styles.primaryBtn}
          disabled={pending() || !name().trim()}
          onClick={() => void submit()}
        >
          Create Category
        </button>
      </div>
    </Modal>
  );
}
