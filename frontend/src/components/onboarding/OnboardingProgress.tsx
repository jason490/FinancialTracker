import styles from "~/styles/onboarding.module.css";

export type OnboardingStep = "welcome" | "plan" | "connect";

const STEPS: Array<{ id: OnboardingStep; label: string }> = [
  { id: "welcome", label: "Welcome" },
  { id: "plan", label: "Plan" },
  { id: "connect", label: "Connect" },
];

type OnboardingProgressProps = {
  current: OnboardingStep;
};

// OnboardingProgress renders the horizontal step indicator for the wizard.
export default function OnboardingProgress(props: OnboardingProgressProps) {
  const currentIndex = () => STEPS.findIndex((step) => step.id === props.current);

  return (
    <ol class={styles.progress} aria-label="Onboarding progress">
      {STEPS.map((step, index) => {
        const state =
          index < currentIndex() ? "done" : index === currentIndex() ? "active" : "upcoming";
        return (
          <li
            class={styles.progressItem}
            classList={{
              [styles.progressItemDone]: state === "done",
              [styles.progressItemActive]: state === "active",
            }}
            aria-current={state === "active" ? "step" : undefined}
          >
            <span class={styles.progressDot} aria-hidden="true">
              {state === "done" ? "✓" : index + 1}
            </span>
            <span class={styles.progressLabel}>{step.label}</span>
          </li>
        );
      })}
    </ol>
  );
}
