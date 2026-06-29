import { Title } from "@solidjs/meta";
import { useNavigate } from "@solidjs/router";
import { Show, createEffect, createResource, createSignal, onMount } from "solid-js";
import OnboardingConnectStep from "~/components/onboarding/OnboardingConnectStep";
import OnboardingPlanStep from "~/components/onboarding/OnboardingPlanStep";
import OnboardingWelcomeStep from "~/components/onboarding/OnboardingWelcomeStep";
import type { OnboardingStep } from "~/components/onboarding/OnboardingProgress";
import OnboardingLayout from "~/layouts/OnboardingLayout";
import { completeOnboarding } from "~/lib/auth";
import { useAuth } from "~/lib/auth-context";
import {
  beginAuthTransition,
  endAuthTransition,
  prefetchDashboardForAuth,
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

  onMount(() => {
    endAuthTransition();
  });

  createEffect(() => {
    if (loading()) {
      return;
    }

    const profile = user();
    if (!profile) {
      navigate("/login", { replace: true });
      return;
    }

    if (profile.onboarding_completed) {
      navigate("/dashboard", { replace: true });
    }
  });

  const handleFinish = async () => {
    setError(null);
    setFinishing(true);
    try {
      await completeOnboarding();
      await refetch();
      beginAuthTransition({
        title: "You're all set",
        hint: "Opening your dashboard",
      });
      await prefetchDashboardForAuth();
      navigate("/dashboard", { replace: true });
    } catch (err) {
      endAuthTransition();
      setError(err instanceof Error ? err.message : "Failed to finish onboarding");
      setFinishing(false);
    }
  };

  const firstName = () => user()?.first_name || "";

  return (
    <Show when={!loading() && user() && !user()!.onboarding_completed}>
      <OnboardingLayout step={step()}>
        <Title>Get Started | Financial Tracker</Title>

        <Show when={error()}>
          <p class={styles.errorBanner} role="alert">
            {error()}
          </p>
        </Show>

        <Show when={step() === "welcome"}>
          <OnboardingWelcomeStep firstName={firstName()} onContinue={() => setStep("plan")} />
        </Show>

        <Show when={step() === "plan"}>
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
            refetchConnections={refetchConnections}
            onBack={() => setStep("plan")}
            onFinish={handleFinish}
            onError={(message) => setError(message)}
            finishing={finishing()}
          />
        </Show>
      </OnboardingLayout>
    </Show>
  );
}
