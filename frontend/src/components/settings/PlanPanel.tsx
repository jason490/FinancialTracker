import { Show, createSignal } from "solid-js";
import PlanPicker, { formatPlanPrice } from "~/components/settings/PlanPicker";
import {
  changeSubscription,
  createBillingPortal,
  createCheckoutSession,
} from "~/lib/subscription";
import { formatDate } from "~/lib/format";
import type { Resource } from "solid-js";
import type { ConnectionsPayload, SubscriptionPayload, TierPlan } from "~/lib/types";
import styles from "~/styles/settings.module.css";

type PlanPanelProps = {
  subscription: Resource<SubscriptionPayload>;
  refetchSubscription: () => void;
  connections: Resource<ConnectionsPayload>;
  onMessage: (message: string, type: "ok" | "error" | "info") => void;
};

// usagePercent returns a capped percentage for meter fills.
function usagePercent(used: number, max: number) {
  if (max <= 0) {
    return 0;
  }
  return Math.min(100, Math.round((used / max) * 100));
}

// usageTone picks a meter color based on how close usage is to the limit.
function usageTone(used: number, max: number) {
  const percent = usagePercent(used, max);
  if (percent >= 100) {
    return styles.usageMeterFillDanger;
  }
  if (percent >= 80) {
    return styles.usageMeterFillWarn;
  }
  return styles.usageMeterFillOk;
}

// formatLimitLabel renders quota labels, showing Unlimited for operator overrides.
function formatLimitLabel(value: number) {
  return value >= 1_000_000 ? "Unlimited" : String(value);
}

// PlanPanel shows the current plan, billing reset date, and upgrade actions.
export default function PlanPanel(props: PlanPanelProps) {
  const [pendingTier, setPendingTier] = createSignal<string | null>(null);
  const [managingBilling, setManagingBilling] = createSignal(false);
  const subscription = props.subscription;
  const refetch = props.refetchSubscription;
  const connections = props.connections;

  const subscriptionsEnabled = () => subscription()?.subscriptions_enabled !== false;
  const isPlaidProvider = () => connections()?.provider === "plaid";

  const currentPlanName = () => {
    const sub = subscription();
    if (!sub) {
      return "";
    }
    return sub.plans.find((plan) => plan.id === sub.tier)?.name ?? sub.tier;
  };

  const openBillingPortal = async () => {
    setManagingBilling(true);
    try {
      const { url } = await createBillingPortal();
      window.location.assign(url);
    } catch (err) {
      props.onMessage(err instanceof Error ? err.message : "Failed to open billing portal", "error");
      setManagingBilling(false);
    }
  };

  const handleChangeTier = async (plan: TierPlan) => {
    const sub = subscription();
    if (!sub || plan.id === sub.tier) {
      return;
    }

    if (sub.billing_enabled) {
      if (plan.id === "free") {
        if (!sub.has_active_subscription) {
          return;
        }
        await openBillingPortal();
        return;
      }

      setPendingTier(plan.id);
      try {
        const { url } = await createCheckoutSession(plan.id);
        window.location.assign(url);
      } catch (err) {
        props.onMessage(err instanceof Error ? err.message : "Failed to start checkout", "error");
        setPendingTier(null);
      }
      return;
    }

    if (
      !window.confirm(
        `Switch to the ${plan.name} plan? Usage limits reset on your billing anniversary.`
      )
    ) {
      return;
    }

    setPendingTier(plan.id);
    try {
      await changeSubscription(plan.id);
      await refetch();
      props.onMessage(`Plan updated to ${plan.name}.`, "ok");
    } catch (err) {
      props.onMessage(err instanceof Error ? err.message : "Failed to change plan", "error");
    } finally {
      setPendingTier(null);
    }
  };

  return (
    <div class={styles.panelInner}>
      <section>
        <h2 class={styles.sectionTitle}>Plan & billing</h2>

        <Show
          when={subscriptionsEnabled()}
          fallback={
            <p class={styles.sectionHint}>
              Subscriptions and usage limits are disabled on this server. Bank connections and sync
              are available without plan restrictions.
            </p>
          }
        >
          <p class={styles.sectionHint}>
            Free usage resets on your signup date each month. Paid plans reset on your subscription
            start date.
          </p>

          <Show
            when={!subscription.loading}
            fallback={<div class={styles.statusInfo}>Loading plan...</div>}
          >
            <article class={styles.billingSummary}>
              <div class={styles.billingSummaryHeader}>
                <div class={styles.billingSummaryIdentity}>
                  <p class={styles.label}>Current plan</p>
                  <h3 class={styles.planTierName}>{currentPlanName()}</h3>
                  <p class={styles.billingSummaryHint}>
                    {formatPlanPrice(
                      subscription()?.plans.find((plan) => plan.id === subscription()?.tier)
                        ?.price_monthly_cents || 0
                    )}
                  </p>
                </div>
                <div class={styles.billingCycle}>
                  <p class={styles.label}>Billing period</p>
                  <p class={styles.billingCycleDate}>
                    Resets {formatDate(subscription()?.billing.period_end || 0)}
                  </p>
                </div>
              </div>

              <Show when={subscription()?.privileges.unlimited_limits}>
                <p class={styles.billingNotice}>
                  Your account has unlimited bank connections and API usage.
                </p>
              </Show>

              <Show when={subscription()?.privileges.has_discount}>
                <p class={styles.billingNotice}>
                  A billing discount will be applied automatically at checkout.
                </p>
              </Show>

              <Show when={!subscription()?.billing_enabled && !subscription()?.can_change_plan}>
                <p class={styles.billingNotice}>
                  Paid upgrades will be available once Stripe billing is fully configured.
                </p>
              </Show>

              <Show
                when={subscription()?.billing_enabled && subscription()?.has_active_subscription}
              >
                <div class={styles.planSection}>
                  <button
                    type="button"
                    class={styles.buttonSecondary}
                    disabled={managingBilling()}
                    onClick={() => void openBillingPortal()}
                  >
                    {managingBilling() ? "Opening portal..." : "Manage billing"}
                  </button>
                </div>
              </Show>

              <Show
                when={connections()?.usage}
                fallback={<div class={styles.usageLoading}>Loading usage...</div>}
              >
                {(usage) => (
                  <div class={styles.usageGrid}>
                    <div class={styles.usageMeter}>
                      <div class={styles.usageMeterHeader}>
                        <span class={styles.usageMeterLabel}>Bank connections</span>
                        <span class={styles.usageMeterValue}>
                          {usage().active_items} / {formatLimitLabel(usage().limits.max_items)}
                        </span>
                      </div>
                      <Show when={!subscription()?.privileges.unlimited_limits}>
                        <div
                          class={styles.usageMeterTrack}
                          role="progressbar"
                          aria-valuemin={0}
                          aria-valuemax={usage().limits.max_items}
                          aria-valuenow={usage().active_items}
                          aria-label="Bank connections used"
                        >
                          <div
                            class={`${styles.usageMeterFill} ${usageTone(
                              usage().active_items,
                              usage().limits.max_items
                            )}`}
                            style={{
                              width: `${usagePercent(usage().active_items, usage().limits.max_items)}%`,
                            }}
                          />
                        </div>
                      </Show>
                    </div>

                    <div class={styles.usageMeter}>
                      <div class={styles.usageMeterHeader}>
                        <span class={styles.usageMeterLabel}>
                          {isPlaidProvider() ? "Link API calls" : "Sync API calls"}
                        </span>
                        <span class={styles.usageMeterValue}>
                          {usage().api_calls_used} /{" "}
                          {formatLimitLabel(usage().limits.max_api_calls_month)}
                        </span>
                      </div>
                      <Show when={!subscription()?.privileges.unlimited_limits}>
                        <div
                          class={styles.usageMeterTrack}
                          role="progressbar"
                          aria-valuemin={0}
                          aria-valuemax={usage().limits.max_api_calls_month}
                          aria-valuenow={usage().api_calls_used}
                          aria-label={
                            isPlaidProvider()
                              ? "Link API calls used this billing period"
                              : "Sync API calls used this billing period"
                          }
                        >
                          <div
                            class={`${styles.usageMeterFill} ${usageTone(
                              usage().api_calls_used,
                              usage().limits.max_api_calls_month
                            )}`}
                            style={{
                              width: `${usagePercent(
                                usage().api_calls_used,
                                usage().limits.max_api_calls_month
                              )}%`,
                            }}
                          />
                        </div>
                      </Show>
                      <Show when={isPlaidProvider()}>
                        <p class={styles.usageFootnote}>
                          Account and transaction sync is unlimited on every plan (manual sync once
                          per minute). Bank data is often a few hours to a day behind.
                        </p>
                      </Show>
                    </div>
                  </div>
                )}
              </Show>
            </article>

            <div class={styles.planSection}>
              <p class={styles.label}>Available plans</p>
              <PlanPicker
                plans={subscription()?.plans || []}
                currentTier={subscription()?.tier || "free"}
                canChangePlan={subscription()?.can_change_plan || false}
                billingEnabled={subscription()?.billing_enabled || false}
                hasActiveSubscription={subscription()?.has_active_subscription || false}
                isPlaidProvider={isPlaidProvider()}
                mode="settings"
                pendingTier={pendingTier()}
                managingBilling={managingBilling()}
                onSelectTier={handleChangeTier}
                onManageBilling={() => void openBillingPortal()}
              />
            </div>
          </Show>
        </Show>
      </section>
    </div>
  );
}
