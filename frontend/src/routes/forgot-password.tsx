import { Title } from "@solidjs/meta";
import { useNavigate } from "@solidjs/router";
import { createSignal, onCleanup, onMount, Show } from "solid-js";
import AuthLayout from "~/layouts/AuthLayout";
import FormError from "~/components/auth/FormError";
import { forgotPassword, resetPassword, verifyResetCode } from "~/lib/auth";
import { RedirectIfAuth } from "~/lib/auth-context";
import styles from "~/styles/auth.module.css";

const RESEND_COOLDOWN_SECONDS = 60;

function formatCountdown(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins}:${secs.toString().padStart(2, "0")}`;
}

export default function ForgotPasswordPage() {
  const navigate = useNavigate();
  const [step, setStep] = createSignal<1 | 2 | 3>(1);
  const [email, setEmail] = createSignal("");
  const [verifiedCode, setVerifiedCode] = createSignal("");
  const [expiresAt, setExpiresAt] = createSignal<number | null>(null);
  const [secondsLeft, setSecondsLeft] = createSignal(0);
  const [resendAvailableAt, setResendAvailableAt] = createSignal<number | null>(null);
  const [resendSecondsLeft, setResendSecondsLeft] = createSignal(0);
  const [resendNotice, setResendNotice] = createSignal<string | null>(null);
  const [code, setCode] = createSignal("");
  const [error, setError] = createSignal<string | null>(null);
  const [pending, setPending] = createSignal(false);

  const syncTimers = () => {
    const expiry = expiresAt();
    if (!expiry) {
      setSecondsLeft(0);
    } else {
      setSecondsLeft(Math.max(0, expiry - Math.floor(Date.now() / 1000)));
    }

    const resendAt = resendAvailableAt();
    if (!resendAt) {
      setResendSecondsLeft(0);
    } else {
      setResendSecondsLeft(Math.max(0, resendAt - Math.floor(Date.now() / 1000)));
    }
  };

  onMount(() => {
    syncTimers();
    const timer = window.setInterval(syncTimers, 1000);
    onCleanup(() => window.clearInterval(timer));
  });

  const codeExpired = () => step() === 2 && secondsLeft() <= 0 && expiresAt() !== null;
  const canResend = () => resendSecondsLeft() <= 0 && !pending();

  const armResendCooldown = () => {
    setResendAvailableAt(Math.floor(Date.now() / 1000) + RESEND_COOLDOWN_SECONDS);
    setResendSecondsLeft(RESEND_COOLDOWN_SECONDS);
  };

  const requestResetCode = async (targetEmail: string, options?: { resend?: boolean }) => {
    const result = await forgotPassword(targetEmail);
    const expiry = Math.floor(Date.now() / 1000) + result.code_expires_in_seconds;
    setExpiresAt(expiry);
    setSecondsLeft(result.code_expires_in_seconds);
    setVerifiedCode("");
    setStep(2);

    if (options?.resend) {
      setResendNotice("A new code has been sent. Check your email.");
      setCode("");
      armResendCooldown();
    } else {
      setResendNotice(null);
      setResendAvailableAt(null);
      setResendSecondsLeft(0);
    }
  };

  const handleEmailSubmit = async (e: SubmitEvent) => {
    e.preventDefault();
    setError(null);
    setPending(true);

    const formData = new FormData(e.currentTarget as HTMLFormElement);
    const nextEmail = (formData.get("email") as string).trim();
    setEmail(nextEmail);

    try {
      await requestResetCode(nextEmail);
    } catch (err: any) {
      setError(err.message || "Request failed");
    } finally {
      setPending(false);
    }
  };

  const handleResendCode = async () => {
    if (!canResend()) {
      return;
    }

    setError(null);
    setResendNotice(null);
    setPending(true);

    try {
      await requestResetCode(email(), { resend: true });
    } catch (err: any) {
      setError(err.message || "Failed to resend code");
    } finally {
      setPending(false);
    }
  };

  const handleVerifySubmit = async (e: SubmitEvent) => {
    e.preventDefault();
    if (codeExpired()) {
      setError("This reset code has expired. Resend a new code to continue.");
      return;
    }

    setError(null);
    setResendNotice(null);
    setPending(true);

    const formData = new FormData(e.currentTarget as HTMLFormElement);
    const submittedCode = (formData.get("code") as string).trim();

    try {
      await verifyResetCode({ email: email(), code: submittedCode });
      setVerifiedCode(submittedCode);
      setExpiresAt(null);
      setSecondsLeft(0);
      setResendAvailableAt(null);
      setResendSecondsLeft(0);
      setStep(3);
    } catch (err: any) {
      setError(err.message || "Verification failed");
    } finally {
      setPending(false);
    }
  };

  const handlePasswordSubmit = async (e: SubmitEvent) => {
    e.preventDefault();
    setError(null);
    setPending(true);

    const formData = new FormData(e.currentTarget as HTMLFormElement);
    const newPassword = formData.get("new_password") as string;
    const confirmPassword = formData.get("confirm_password") as string;

    try {
      await resetPassword({
        email: email(),
        code: verifiedCode(),
        new_password: newPassword,
        confirm_password: confirmPassword,
      });
      navigate("/login?reset=success");
    } catch (err: any) {
      setError(err.message || "Reset failed");
    } finally {
      setPending(false);
    }
  };

  const resetToEmailStep = () => {
    setStep(1);
    setError(null);
    setResendNotice(null);
    setExpiresAt(null);
    setSecondsLeft(0);
    setResendAvailableAt(null);
    setResendSecondsLeft(0);
    setVerifiedCode("");
    setCode("");
  };

  const subtitle = () => {
    if (step() === 1) {
      return "Enter the email for your account and we'll send you a reset code.";
    }
    if (step() === 2) {
      return `We sent a 6-digit code to ${email()}.`;
    }
    return "Create a new password for your account.";
  };

  return (
    <RedirectIfAuth>
      <AuthLayout eyebrow="Account recovery" title="Forgot password" subtitle={subtitle()}>
        <Title>Forgot Password | Financial Tracker</Title>

        <Show when={step() === 1}>
          <form onSubmit={handleEmailSubmit} class={styles.form}>
            <div class={styles.field}>
              <label class={styles.label} for="email">
                Email address
              </label>
              <input
                class={styles.input}
                id="email"
                name="email"
                type="email"
                autoComplete="email"
                value={email()}
                required
              />
            </div>

            <p class={styles.hint}>
              Google-only accounts should sign in with Google instead of resetting a password.
            </p>

            <FormError message={error() || undefined} />

            <button class={styles.button} type="submit" disabled={pending()}>
              <Show when={pending()} fallback="Send reset code">
                Sending...
              </Show>
            </button>
          </form>
        </Show>

        <Show when={step() === 2}>
          <form onSubmit={handleVerifySubmit} class={styles.form}>
            <Show when={resendNotice()}>
              <p class={styles.success} role="status">
                {resendNotice()}
              </p>
            </Show>

            <Show when={expiresAt()}>
              <p
                class={codeExpired() ? styles.error : styles.hint}
                role="status"
                aria-live="polite"
              >
                {codeExpired()
                  ? "This code has expired. Resend a new code to continue."
                  : `Code expires in ${formatCountdown(secondsLeft())}`}
              </p>
            </Show>

            <div class={styles.field}>
              <label class={styles.label} for="code">
                Reset code
              </label>
              <input
                class={styles.input}
                id="code"
                name="code"
                type="text"
                inputMode="numeric"
                autoComplete="one-time-code"
                pattern="[0-9]{6}"
                maxLength={6}
                minLength={6}
                required
                disabled={codeExpired()}
                value={code()}
                onInput={(e) => setCode(e.currentTarget.value)}
              />
            </div>

            <FormError message={error() || undefined} />

            <button class={styles.button} type="submit" disabled={pending() || codeExpired()}>
              <Show when={pending()} fallback="Continue">
                Verifying...
              </Show>
            </button>

            <div class={styles.assistLinks}>
              <p>
                Didn't receive a code?{" "}
                <button
                  class={styles.textLink}
                  type="button"
                  disabled={!canResend()}
                  onClick={handleResendCode}
                >
                  {resendSecondsLeft() > 0
                    ? `Resend in ${formatCountdown(resendSecondsLeft())}`
                    : "Resend code"}
                </button>
              </p>
              <p>
                <button class={styles.textLink} type="button" disabled={pending()} onClick={resetToEmailStep}>
                  Use a different email
                </button>
              </p>
            </div>
          </form>
        </Show>

        <Show when={step() === 3}>
          <form onSubmit={handlePasswordSubmit} class={styles.form}>
            <div class={styles.field}>
              <label class={styles.label} for="new_password">
                New password
              </label>
              <input
                class={styles.input}
                id="new_password"
                name="new_password"
                type="password"
                autoComplete="new-password"
                required
              />
            </div>

            <div class={styles.field}>
              <label class={styles.label} for="confirm_password">
                Confirm new password
              </label>
              <input
                class={styles.input}
                id="confirm_password"
                name="confirm_password"
                type="password"
                autoComplete="new-password"
                required
              />
            </div>

            <FormError message={error() || undefined} />

            <button class={styles.button} type="submit" disabled={pending()}>
              <Show when={pending()} fallback="Reset password">
                Resetting...
              </Show>
            </button>
          </form>
        </Show>

        <p class={styles.footer}>
          Remembered your password? <a href="/login">Back to sign in</a>
        </p>
      </AuthLayout>
    </RedirectIfAuth>
  );
}
