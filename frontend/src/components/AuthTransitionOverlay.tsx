import { Show } from "solid-js";
import { authTransitionActive, authTransitionCopy } from "~/lib/auth-transition";
import styles from "~/styles/auth.module.css";

// AuthTransitionOverlay keeps a loading screen visible while auth routes hand off to the app.
export default function AuthTransitionOverlay() {
  return (
    <Show when={authTransitionActive()}>
      <div
        class={styles.globalTransitionOverlay}
        role="status"
        aria-live="polite"
        aria-busy="true"
      >
        <div class={styles.transitionRing}>
          <span class={styles.transitionCheck} />
        </div>
        <p class={styles.transitionLabel}>{authTransitionCopy().title}</p>
        <p class={styles.transitionHint}>{authTransitionCopy().hint}</p>
      </div>
    </Show>
  );
}
