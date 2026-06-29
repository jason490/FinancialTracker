import { For, Show } from "solid-js";
import type { TierPlan } from "~/lib/types";
import styles from "~/styles/settings.module.css";

type PlanPickerProps = {
  plans: TierPlan[];
  currentTier: string;
  canChangePlan: boolean;
  billingEnabled: boolean;
  hasActiveSubscription: boolean;
  isPlaidProvider: boolean;
  mode: "settings" | "onboarding";
  pendingTier: string | null;
  onSelectTier: (plan: TierPlan) => void;
  onManageBilling?: () => void;
  managingBilling?: boolean;
};

// formatPrice renders monthly pricing in dollars.
export function formatPlanPrice(cents: number) {
  if (cents <= 0) {
    return "Free";
  }
  return `$${(cents / 100).toFixed(2)}/mo`;
}

// planActionLabel returns the CTA label for a plan card.
function planActionLabel(
  plan: TierPlan,
  currentTier: string,
  billingEnabled: boolean,
  canChangePlan: boolean,
  pendingTier: string | null
) {
  if (pendingTier === plan.id) {
    return billingEnabled && plan.id !== "free" ? "Redirecting..." : "Updating...";
  }
  if (plan.id === currentTier) {
    return "Current plan";
  }
  if (plan.id === "free") {
    return billingEnabled ? "Manage in billing portal" : "Coming soon";
  }
  if (billingEnabled) {
    return `Upgrade to ${plan.name}`;
  }
  if (canChangePlan) {
    return `Switch to ${plan.name}`;
  }
  return "Coming soon";
}

// PlanPicker renders selectable subscription tier cards.
export default function PlanPicker(props: PlanPickerProps) {
  const isCurrent = (plan: TierPlan) => plan.id === props.currentTier;

  return (
    <div class={styles.planGrid}>
      <For each={props.plans}>
        {(plan) => (
          <article
            class={`${styles.planCard} ${isCurrent(plan) ? styles.planCardActive : ""}`}
          >
            <div class={styles.planHeader}>
              <h3 class={styles.planName}>{plan.name}</h3>
              <p class={styles.planPrice}>{formatPlanPrice(plan.price_monthly_cents)}</p>
            </div>
            <ul class={styles.calloutList}>
              <li>
                Up to <strong>{plan.limits.max_items}</strong> bank connections
              </li>
              <Show
                when={props.isPlaidProvider}
                fallback={
                  <li>
                    <strong>{plan.limits.max_api_calls_month}</strong> sync API calls per billing
                    period
                  </li>
                }
              >
                <li>Unlimited account &amp; transaction sync (once per minute)</li>
                <li>
                  <strong>{plan.limits.max_api_calls_month}</strong> link API calls per billing
                  period
                </li>
              </Show>
            </ul>
            <Show
              when={props.mode === "settings"}
              fallback={
                <button
                  type="button"
                  class={isCurrent(plan) ? styles.buttonSecondary : styles.buttonPrimary}
                  disabled={
                    props.pendingTier !== null ||
                    (plan.id === "free" && isCurrent(plan)) ||
                    (!props.canChangePlan && !props.billingEnabled && !isCurrent(plan))
                  }
                  onClick={() => props.onSelectTier(plan)}
                >
                  {planActionLabel(
                    plan,
                    props.currentTier,
                    props.billingEnabled,
                    props.canChangePlan,
                    props.pendingTier
                  )}
                </button>
              }
            >
              <Show
                when={plan.id === "free" && props.billingEnabled && props.hasActiveSubscription}
                fallback={
                  <Show
                    when={props.canChangePlan || props.billingEnabled}
                    fallback={
                      <button type="button" class={styles.buttonSecondary} disabled>
                        {isCurrent(plan) ? "Current plan" : "Coming soon"}
                      </button>
                    }
                  >
                    <button
                      type="button"
                      class={isCurrent(plan) ? styles.buttonSecondary : styles.buttonPrimary}
                      disabled={
                        props.pendingTier !== null ||
                        (isCurrent(plan) && plan.id !== "free") ||
                        (plan.id === "free" && !props.hasActiveSubscription)
                      }
                      onClick={() => props.onSelectTier(plan)}
                    >
                      {planActionLabel(
                        plan,
                        props.currentTier,
                        props.billingEnabled,
                        props.canChangePlan,
                        props.pendingTier
                      )}
                    </button>
                  </Show>
                }
              >
                <button
                  type="button"
                  class={styles.buttonSecondary}
                  disabled={props.managingBilling}
                  onClick={() => props.onManageBilling?.()}
                >
                  {props.managingBilling ? "Opening portal..." : "Manage in billing portal"}
                </button>
              </Show>
            </Show>
          </article>
        )}
      </For>
    </div>
  );
}
