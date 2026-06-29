import styles from "~/styles/onboarding.module.css";

export type OnboardingStep = "welcome" | "plan" | "connect";

const ALL_STEPS: Array<{ id: OnboardingStep; label: string }> = [
  { id: "welcome", label: "Welcome" },
  { id: "plan", label: "Plan" },
  { id: "connect", label: "Connect" },
];

type OnboardingProgressProps = {
  current: OnboardingStep;
  subscriptionsEnabled?: boolean;
};

// OnboardingProgress renders the horizontal step indicator for the wizard.
export default function OnboardingProgress(props: OnboardingProgressProps) {
  const steps = () =>
    props.subscriptionsEnabled === false
      ? ALL_STEPS.filter((step) => step.id !== "plan")
      : ALL_STEPS;

  const currentIndex = () => steps().findIndex((step) => step.id === props.current);

  return (
    <ol class={styles.progress} aria-label="Onboarding progress">
      {steps().map((step, index) => {
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

export function onboardingStepNumber(
  step: OnboardingStep,
  subscriptionsEnabled: boolean
): { current: number; total: number } {
  const steps =
    subscriptionsEnabled === false
      ? ALL_STEPS.filter((item) => item.id !== "plan")
      : ALL_STEPS;
  const current = steps.findIndex((item) => item.id === step) + 1;
  return { current, total: steps.length };
}
