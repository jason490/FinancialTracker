import { Title } from "@solidjs/meta";
import { useNavigate, useSearchParams } from "@solidjs/router";
import { createSignal, Show } from "solid-js";
import AuthLayout, { type AuthTransitionPhase } from "~/layouts/AuthLayout";
import FormError from "~/components/auth/FormError";
import SSOButtons from "~/components/auth/SSOButtons";
import { login } from "~/lib/auth";
import {
  beginAuthTransition,
  authTransitionActive,
  prefetchDashboardForAuth,
} from "~/lib/auth-transition";
import { RedirectIfAuth, useAuth } from "~/lib/auth-context";
import styles from "~/styles/auth.module.css";

const SUCCESS_MS = 520;
const EXIT_MS = 420;

export default function LoginPage() {
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();
  const { refetch } = useAuth();
  const [error, setError] = createSignal<string | null>(null);
  const [pending, setPending] = createSignal(false);
  const [transitionPhase, setTransitionPhase] = createSignal<AuthTransitionPhase>("idle");
  const passwordResetSuccess = () => searchParams.reset === "success";

  const dismissResetSuccess = () => {
    if (searchParams.reset) {
      setSearchParams({ reset: undefined }, { replace: true });
    }
  };

  const wait = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

  const handleSubmit = async (e: SubmitEvent) => {
    e.preventDefault();
    setError(null);
    setPending(true);

    const formData = new FormData(e.currentTarget as HTMLFormElement);
    const data = {
      email: formData.get("email"),
      password: formData.get("password"),
      remember: formData.get("remember") === "on",
    };

    try {
      await login(data);
      setTransitionPhase("success");
      beginAuthTransition();
      await refetch();
      await wait(SUCCESS_MS);
      setTransitionPhase("exiting");
      await Promise.all([wait(EXIT_MS), prefetchDashboardForAuth()]);
      navigate("/dashboard");
    } catch (err: any) {
      setTransitionPhase("idle");
      setError(err.message || "Login failed");
      setPending(false);
    }
  };

  const isBusy = () => pending() || transitionPhase() !== "idle";

  return (
    <RedirectIfAuth>
      <AuthLayout
        eyebrow="Welcome back"
        title="Sign in"
        subtitle="Access your accounts, tags, and dashboard."
        transitionPhase={transitionPhase()}
      >
        <Title>Sign In | Financial Tracker</Title>

        <div class={styles.form}>
          <SSOButtons label="Continue with Google" />

          <div class={styles.divider}>or email</div>

          <form onSubmit={handleSubmit} class={styles.form}>
            <Show when={passwordResetSuccess()}>
              <div class={styles.success} role="status">
                Your password was reset. Sign in with your new password.
                <button class={styles.inlineDismiss} type="button" onClick={dismissResetSuccess}>
                  Dismiss
                </button>
              </div>
            </Show>

            <div class={styles.field}>
              <label class={styles.label} for="email">
                Email address
              </label>
              <input
                class={styles.input}
                id="email"
                name="email"
                type="email"
                autocomplete="email"
                required
              />
            </div>

            <div class={styles.field}>
              <label class={styles.label} for="password">
                Password
              </label>
              <input
                class={styles.input}
                id="password"
                name="password"
                type="password"
                autocomplete="current-password"
                required
              />
            </div>

            <FormError message={error() || undefined} />

            <div class={styles.linkRow}>
              <label class={styles.checkboxRow}>
                <input type="checkbox" id="remember" name="remember" />
                <span>Remember me</span>
              </label>
              <a href="/forgot-password">Forgot password?</a>
            </div>

            <button class={styles.button} type="submit" disabled={isBusy()}>
              <Show when={pending()} fallback="Sign in">
                Signing in...
              </Show>
            </button>
          </form>

          <p class={styles.footer}>
            Don't have an account? <a href="/register">Create one</a>
          </p>
        </div>
      </AuthLayout>
    </RedirectIfAuth>
  );
}
