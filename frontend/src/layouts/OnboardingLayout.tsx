import type { JSX } from "solid-js";
import OnboardingProgress, { type OnboardingStep } from "~/components/onboarding/OnboardingProgress";
import styles from "~/styles/onboarding.module.css";

type OnboardingLayoutProps = {
  step: OnboardingStep;
  subscriptionsEnabled?: boolean;
  children: JSX.Element;
};

// OnboardingLayout renders a focused shell for the new-user wizard without app navigation.
export default function OnboardingLayout(props: OnboardingLayoutProps) {
  return (
    <div class={styles.shell}>
      <header class={styles.header}>
        <p class={styles.brand}>FinancialTracker</p>
        <p class={styles.headerHint}>Account setup</p>
      </header>

      <OnboardingProgress
        current={props.step}
        subscriptionsEnabled={props.subscriptionsEnabled}
      />

      <main class={styles.main}>{props.children}</main>
    </div>
  );
}
