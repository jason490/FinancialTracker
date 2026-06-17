import { Title } from "@solidjs/meta";
import { useNavigate } from "@solidjs/router";
import { createSignal, Show } from "solid-js";
import AuthLayout from "~/layouts/AuthLayout";
import FormError from "~/components/auth/FormError";
import SSOButtons from "~/components/auth/SSOButtons";
import { login } from "~/lib/auth";
import { RedirectIfAuth, useAuth } from "~/lib/auth-context";
import styles from "~/styles/auth.module.css";

export default function LoginPage() {
  const navigate = useNavigate();
  const { refetch } = useAuth();
  const [error, setError] = createSignal<string | null>(null);
  const [pending, setPending] = createSignal(false);

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
      await refetch();
      navigate("/dashboard");
    } catch (err: any) {
      setError(err.message || "Login failed");
    } finally {
      setPending(false);
    }
  };

  return (
    <RedirectIfAuth>
      <AuthLayout
        eyebrow="Welcome back"
        title="Sign in"
        subtitle="Access your accounts, tags, and dashboard."
      >
        <Title>Sign In | Financial Tracker</Title>

        <div class={styles.form}>
          <SSOButtons label="Continue with Google" />

          <div class={styles.divider}>or email</div>

          <form onSubmit={handleSubmit} class={styles.form}>
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

            <button class={styles.button} type="submit" disabled={pending()}>
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
