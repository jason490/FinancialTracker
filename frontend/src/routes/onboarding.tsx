import { Title } from "@solidjs/meta";
import { useNavigate } from "@solidjs/router";
import { Show, createEffect, createResource, createSignal } from "solid-js";
import OnboardingConnectStep from "~/components/onboarding/OnboardingConnectStep";
import OnboardingPlanStep from "~/components/onboarding/OnboardingPlanStep";
import OnboardingWelcomeStep from "~/components/onboarding/OnboardingWelcomeStep";
import type { OnboardingStep } from "~/components/onboarding/OnboardingProgress";
import OnboardingLayout from "~/layouts/OnboardingLayout";
import { completeOnboarding } from "~/lib/auth";
import { useAuth } from "~/lib/auth-context";
import {
  endAuthTransition,
  preloadDashboardRoute,
} from "~/lib/auth-transition";
import { getConnections } from "~/lib/connections";
import { getSubscription } from "~/lib/subscription";
import styles from "~/styles/onboarding.module.css";

// OnboardingPage guides new users through plan selection and optional bank connection.
export default function OnboardingPage() {
  const navigate = useNavigate();
  const { user, loading, refetch } = useAuth();
  const [step, setStep] = createSignal<OnboardingStep>("welcome");
  const [error, setError] = createSignal<string | null>(null);
  const [finishing, setFinishing] = createSignal(false);

  const [subscription, { refetch: refetchSubscription }] = createResource(getSubscription);
  const [connections, { refetch: refetchConnections }] = createResource(getConnections);

  const subscriptionsEnabled = () => subscription()?.subscriptions_enabled !== false;

  createEffect(() => {
    if (loading()) {
      return;
    }

    const profile = user();
    if (!profile) {
      endAuthTransition();
      navigate("/login", { replace: true });
      return;
    }

    if (profile.onboarding_completed) {
      // Let the dashboard take over the overlay handoff.
      navigate("/dashboard", { replace: true });
      return;
    }

    // Hold the post-auth overlay until the first wizard step can render with
    // its data resolved, so the handoff from login/SSO has no blank flash.
    if (subscription.loading) {
      return;
    }
    endAuthTransition();
  });

  createEffect(() => {
    if (subscription.loading || subscriptionsEnabled()) {
      return;
    }
    if (step() === "plan") {
      setStep("connect");
    }
  });

  const handleFinish = async () => {
    setError(null);
    setFinishing(true);
    try {
      await completeOnboarding();
      await refetch();
      // Preload the dashboard route so it mounts instantly with its own
      // skeleton loader — no "You're all set" overlay text.
      await preloadDashboardRoute();
      navigate("/dashboard", { replace: true });
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to finish onboarding");
      setFinishing(false);
    }
  };

  const firstName = () => user()?.first_name || "";

  return (
    <Show when={!loading() && user() && !user()!.onboarding_completed}>
      <OnboardingLayout step={step()} subscriptionsEnabled={subscriptionsEnabled()}>
        <Title>Get Started | Financial Tracker</Title>

        <Show when={error()}>
          <p class={styles.errorBanner} role="alert">
            {error()}
          </p>
        </Show>

        <Show when={step() === "welcome"}>
          <OnboardingWelcomeStep
            firstName={firstName()}
            subscriptionsEnabled={subscriptionsEnabled()}
            onContinue={() => setStep(subscriptionsEnabled() ? "plan" : "connect")}
          />
        </Show>

        <Show when={step() === "plan" && subscriptionsEnabled()}>
          <OnboardingPlanStep
            subscription={subscription}
            connections={connections}
            refetchSubscription={refetchSubscription}
            onContinue={() => setStep("connect")}
            onBack={() => setStep("welcome")}
            onError={(message) => setError(message)}
          />
        </Show>

        <Show when={step() === "connect"}>
          <OnboardingConnectStep
            connections={connections}
            subscriptionsEnabled={subscriptionsEnabled()}
            refetchConnections={refetchConnections}
            onBack={() => setStep(subscriptionsEnabled() ? "plan" : "welcome")}
            onFinish={handleFinish}
            onError={(message) => setError(message)}
            finishing={finishing()}
          />
        </Show>
      </OnboardingLayout>
    </Show>
  );
}
