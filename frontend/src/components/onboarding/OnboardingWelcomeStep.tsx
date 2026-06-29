import styles from "~/styles/onboarding.module.css";

type OnboardingWelcomeStepProps = {
  firstName: string;
  onContinue: () => void;
};

// OnboardingWelcomeStep greets new users and explains the product value.
export default function OnboardingWelcomeStep(props: OnboardingWelcomeStepProps) {
  const greeting = () => {
    const name = props.firstName.trim();
    return name ? `Welcome, ${name}` : "Welcome aboard";
  };

  return (
    <section class={styles.stepPanel}>
      <p class={styles.stepEyebrow}>Step 1 of 3</p>
      <h1 class={styles.stepTitle}>{greeting()}</h1>
      <p class={styles.stepLead}>
        FinancialTracker syncs your accounts, tags every transaction, and builds a dashboard you
        can shape around how you actually spend.
      </p>

      <ul class={styles.featureGrid}>
        <li class={styles.featureCard}>
          <span class={styles.featureIcon} aria-hidden="true">
            ◈
          </span>
          <div>
            <p class={styles.featureTitle}>Automatic sync</p>
            <p class={styles.featureCopy}>
              Link banks through Plaid and keep balances and transactions up to date.
            </p>
          </div>
        </li>
        <li class={styles.featureCard}>
          <span class={styles.featureIcon} aria-hidden="true">
            ◎
          </span>
          <div>
            <p class={styles.featureTitle}>Smart tagging</p>
            <p class={styles.featureCopy}>
              Rules categorize spending so you spend less time sorting receipts.
            </p>
          </div>
        </li>
        <li class={styles.featureCard}>
          <span class={styles.featureIcon} aria-hidden="true">
            ◐
          </span>
          <div>
            <p class={styles.featureTitle}>Your dashboard</p>
            <p class={styles.featureCopy}>
              Widgets snap into place so net worth, cash flow, and trends stay in view.
            </p>
          </div>
        </li>
      </ul>

      <div class={styles.stepActions}>
        <button type="button" class={styles.primaryButton} onClick={() => props.onContinue()}>
          Get started
        </button>
      </div>
    </section>
  );
}
