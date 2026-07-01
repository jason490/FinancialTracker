import { Title } from "@solidjs/meta";
import { useNavigate, useSearchParams } from "@solidjs/router";
import { Show, onMount } from "solid-js";
import AuthLayout from "~/layouts/AuthLayout";
import { useAuth } from "~/lib/auth-context";
import { postAuthPath } from "~/lib/auth";
import { beginAuthTransition, endAuthTransition, prefetchDashboardForAuth } from "~/lib/auth-transition";
import styles from "~/styles/auth.module.css";

const errorMessages: Record<string, string> = {
  authentication_failed: "Google sign-in failed. Please try again.",
  sso_account_exists:
    "An account with this email already exists. Sign in with your password, then link Google from Settings.",
  invalid_registration_code:
    "Invalid or expired registration code. Contact an administrator for a new code.",
  registration_code_required:
    "A registration code from an administrator is required. Use the register page to sign up.",
};

// SSOCompletePage verifies that the SSO session was established directly with the API.
export default function SSOCompletePage() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const { refetch } = useAuth();

  onMount(async () => {
    const error = searchParams.error;
    if (error) {
      return;
    }

    // Use the centralized refetch to verify the session and update global state.
    // The browser automatically includes the 'Session' cookie set during the redirect.
    try {
      beginAuthTransition({
        title: "Welcome aboard",
        hint: "Setting up your workspace",
      });
      const user = await refetch();
      if (user) {
        const destination = postAuthPath(user);
        if (destination === "/dashboard") {
          await prefetchDashboardForAuth();
        }
        // Keep the overlay up across the route change; the destination page
        // (dashboard or onboarding) dismisses it once its content is ready.
        navigate(destination, { replace: true });
      } else {
        endAuthTransition();
        navigate("/login", { replace: true });
      }
    } catch {
      endAuthTransition();
      navigate("/login", { replace: true });
    }
  });

  const errorMessage = () => {
    const code = searchParams.error;
    if (!code || typeof code !== "string") {
      return undefined;
    }
    return errorMessages[code] || "Sign-in could not be completed.";
  };

  return (
    <AuthLayout
      eyebrow="Single sign-on"
      title="Finishing sign-in"
      subtitle="Hold on while we secure your session."
    >
      <Title>Signing In | Financial Tracker</Title>

      <Show
        when={errorMessage()}
        fallback={<div class={styles.success}>Completing Google sign-in...</div>}
      >
        <div class={styles.error} role="alert">
          {errorMessage()}
        </div>
        <p class={styles.footer}>
          <a href="/login">Return to sign in</a>
        </p>
      </Show>
    </AuthLayout>
  );
}
