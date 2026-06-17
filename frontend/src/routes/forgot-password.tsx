import { Title } from "@solidjs/meta";
import { createSignal, Show } from "solid-js";
import AuthLayout from "~/layouts/AuthLayout";
import FormError from "~/components/auth/FormError";
import { forgotPassword } from "~/lib/auth";
import { RedirectIfAuth } from "~/lib/auth-context";
import styles from "~/styles/auth.module.css";

export default function ForgotPasswordPage() {
  const [error, setError] = createSignal<string | null>(null);
  const [pending, setPending] = createSignal(false);
  const [submitted, setSubmitted] = createSignal(false);

  const handleSubmit = async (e: SubmitEvent) => {
    e.preventDefault();
    setError(null);
    setPending(true);
    setSubmitted(false);

    const formData = new FormData(e.currentTarget as HTMLFormElement);
    const email = formData.get("email") as string;

    try {
      await forgotPassword(email);
      setSubmitted(true);
    } catch (err: any) {
      setError(err.message || "Request failed");
    } finally {
      setPending(false);
    }
  };

  return (
    <RedirectIfAuth>
      <AuthLayout
        eyebrow="Account recovery"
        title="Reset password"
        subtitle="We'll send reset instructions if an account exists for that email."
      >
        <Title>Reset Password | Financial Tracker</Title>

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

          <FormError message={error() || undefined} />

          <Show when={submitted()}>
            <div class={styles.success} role="status">
              If an account exists for that email, a reset link will be sent shortly.
            </div>
          </Show>

          <button class={styles.button} type="submit" disabled={pending()}>
            <Show when={pending()} fallback="Send reset link">
              Sending...
            </Show>
          </button>
        </form>

        <p class={styles.footer}>
          Remembered your password? <a href="/login">Back to sign in</a>
        </p>
      </AuthLayout>
    </RedirectIfAuth>
  );
}
