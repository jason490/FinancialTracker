import { ParentProps, Show } from "solid-js";
import { XIcon } from "~/components/icons";
import styles from "~/styles/tags.module.css";

type ModalProps = ParentProps<{
  open: boolean;
  title: string;
  subtitle?: string;
  narrow?: boolean;
  onClose: () => void;
}>;

// Modal renders a centered overlay panel for tag management forms.
export default function Modal(props: ModalProps) {
  return (
    <Show when={props.open}>
      <div
        class={styles.modalBackdrop}
        role="presentation"
        onClick={(event) => {
          if (event.target === event.currentTarget) props.onClose();
        }}
      >
        <div
          class={`${styles.modalPanel} ${props.narrow ? styles.modalPanelNarrow : ""}`}
          role="dialog"
          aria-modal="true"
          aria-labelledby="tags-modal-title"
        >
          <div class={styles.modalHeader}>
            <div>
              <h2 id="tags-modal-title" class={styles.modalTitle}>
                {props.title}
              </h2>
              <Show when={props.subtitle}>
                <p class={styles.modalSubtitle}>{props.subtitle}</p>
              </Show>
            </div>
            <button type="button" class={styles.iconBtn} aria-label="Close" onClick={props.onClose}>
              <XIcon size={18} />
            </button>
          </div>
          {props.children}
        </div>
      </div>
    </Show>
  );
}
