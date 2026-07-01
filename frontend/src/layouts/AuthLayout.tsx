import type { JSX } from "solid-js";
import { authTransitionActive } from "~/lib/auth-transition";
import styles from "~/styles/auth.module.css";

export type AuthTransitionPhase = "idle" | "success" | "exiting";

type AuthLayoutProps = {
  eyebrow: string;
  title: string;
  subtitle: string;
  children: JSX.Element;
  transitionPhase?: AuthTransitionPhase;
};

// AuthLayout renders the shared split layout for authentication pages.
export default function AuthLayout(props: AuthLayoutProps) {
  const phase = () => props.transitionPhase ?? "idle";

  return (
    <div
      class={styles.shell}
      classList={{
        [styles.shellSuccess]: phase() === "success",
        [styles.shellExiting]: phase() === "exiting",
      }}
    >
      <section class={styles.hero}>
        <p class={styles.eyebrow}>Financial Tracker</p>
        <h1 class={styles.brand}>Clarity for every dollar you move.</h1>
        <p class={styles.lede}>
          Sync accounts through Plaid, automate tagging, and shape a dashboard that
          matches how you actually manage money.
        </p>
        <ul class={styles.featureList}>
          <li>Real-time transaction sync</li>
          <li>Smart tagging rules</li>
          <li>Customizable dashboard widgets</li>
        </ul>
      </section>

      <section class={styles.panel}>
        <div
          class={styles.card}
          classList={{
            [styles.cardSuccess]: phase() === "success",
            [styles.cardExiting]: phase() === "exiting",
          }}
        >
          <p class={styles.eyebrow}>{props.eyebrow}</p>
          <h2 class={styles.title}>{props.title}</h2>
          <p class={styles.subtitle}>{props.subtitle}</p>
          {props.children}
        </div>

        <div
          class={styles.transitionOverlay}
          classList={{
            [styles.transitionOverlayVisible]:
              phase() !== "idle" && !authTransitionActive(),
          }}
          aria-hidden={phase() === "idle" || authTransitionActive()}
        >
          <div class={styles.transitionRing}>
            <span class={styles.transitionCheck} />
          </div>

        </div>
      </section>
    </div>
  );
}
