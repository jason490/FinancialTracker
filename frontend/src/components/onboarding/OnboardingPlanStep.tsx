import { Show, createSignal } from "solid-js";
import PlanPicker from "~/components/settings/PlanPicker";
import { onboardingStepNumber } from "~/components/onboarding/OnboardingProgress";
import { changeSubscription, createCheckoutSession } from "~/lib/subscription";
import type { ConnectionsPayload, SubscriptionPayload, TierPlan } from "~/lib/types";
import type { Resource } from "solid-js";
import styles from "~/styles/onboarding.module.css";

type OnboardingPlanStepProps = {
  subscription: Resource<SubscriptionPayload>;
  connections: Resource<ConnectionsPayload>;
  refetchSubscription: () => void;
  onContinue: () => void;
  onBack: () => void;
  onError: (message: string) => void;
};

// OnboardingPlanStep lets new users pick a subscription tier before connecting banks.
export default function OnboardingPlanStep(props: OnboardingPlanStepProps) {
  const [pendingTier, setPendingTier] = createSignal<string | null>(null);

  const isPlaidProvider = () => props.connections()?.provider === "plaid";

  const handleSelectTier = async (plan: TierPlan) => {
    const sub = props.subscription();
    if (!sub || plan.id === sub.tier) {
      return;
    }

    if (sub.billing_enabled && plan.id !== "free") {
      setPendingTier(plan.id);
      try {
        const { url } = await createCheckoutSession(plan.id);
        window.location.assign(url);
      } catch (err) {
        props.onError(err instanceof Error ? err.message : "Failed to start checkout");
        setPendingTier(null);
      }
      return;
    }

    setPendingTier(plan.id);
    try {
      await changeSubscription(plan.id);
      props.refetchSubscription();
    } catch (err) {
      props.onError(err instanceof Error ? err.message : "Failed to select plan");
    } finally {
      setPendingTier(null);
    }
  };

  const stepLabel = () => {
    const { current, total } = onboardingStepNumber("plan", true);
    return `Step ${current} of ${total}`;
  };

  return (
    <section class={styles.stepPanel}>
      <p class={styles.stepEyebrow}>{stepLabel()}</p>
      <h2 class={styles.stepTitle}>Choose your plan</h2>
      <p class={styles.stepLead}>
        Start on Free and upgrade anytime. Limits reset on your signup date each month.
      </p>

      <Show
        when={!props.subscription.loading && props.subscription()}
        fallback={<div class={styles.statusInfo}>Loading plans...</div>}
      >
        {(sub) => (
          <>
            <Show when={sub().billing_enabled}>
              <p class={styles.planNotice}>
                Paid plans checkout through Stripe. You can stay on Free and upgrade later.
              </p>
            </Show>
            <Show when={!sub().billing_enabled && sub().can_change_plan}>
              <p class={styles.planNotice}>
                Development mode: you can switch tiers instantly without payment.
              </p>
            </Show>
            <PlanPicker
              plans={sub().plans}
              currentTier={sub().tier}
              canChangePlan={sub().can_change_plan}
              billingEnabled={sub().billing_enabled}
              hasActiveSubscription={sub().has_active_subscription}
              isPlaidProvider={isPlaidProvider()}
              mode="onboarding"
              pendingTier={pendingTier()}
              onSelectTier={handleSelectTier}
            />
          </>
        )}
      </Show>

      <div class={styles.stepActions}>
        <button type="button" class={styles.secondaryButton} onClick={() => props.onBack()}>
          Back
        </button>
        <button
          type="button"
          class={styles.primaryButton}
          disabled={props.subscription.loading || pendingTier() !== null}
          onClick={() => props.onContinue()}
        >
          Continue
        </button>
      </div>
    </section>
  );
}
