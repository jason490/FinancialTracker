import { Show } from "solid-js";
import styles from "~/styles/auth.module.css";

// FormError displays a submission error when present.
export default function FormError(props: { message?: string }) {
  return (
    <Show when={props.message}>
      <div class={styles.error} role="alert">
        {props.message}
      </div>
    </Show>
  );
}
