import { Show, createEffect, onCleanup, type Accessor } from "solid-js";
import { XIcon } from "~/components/icons";
import styles from "~/styles/page-status.module.css";

export type PageStatus = {
  text: string;
  type: "ok" | "error" | "info";
};

const DISMISS_MS = 30_000;

type PageStatusBannerProps = {
  message: Accessor<PageStatus | null>;
  onDismiss: () => void;
};

// PageStatusBanner renders a dismissible status line that auto-hides after 30 seconds.
export default function PageStatusBanner(props: PageStatusBannerProps) {
  createEffect(() => {
    if (!props.message()) {
      return;
    }

    const timer = window.setTimeout(() => props.onDismiss(), DISMISS_MS);
    onCleanup(() => window.clearTimeout(timer));
  });

  return (
    <Show when={props.message()}>
      {(current) => (
        <div
          class={`${styles.banner} ${
            current().type === "error"
              ? styles.statusError
              : current().type === "info"
                ? styles.statusInfo
                : styles.statusOk
          }`}
          role="status"
        >
          <p class={styles.statusText}>{current().text}</p>
          <button
            type="button"
            class={styles.dismissButton}
            aria-label="Dismiss message"
            onClick={() => props.onDismiss()}
          >
            <XIcon size={16} />
          </button>
        </div>
      )}
    </Show>
  );
}
