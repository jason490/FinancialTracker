import { Title } from "@solidjs/meta";
import { useNavigate } from "@solidjs/router";
import { createResource, createSignal, Show } from "solid-js";
import AuthLayout from "~/layouts/AuthLayout";
import FormError from "~/components/auth/FormError";
import SSOButtons from "~/components/auth/SSOButtons";
import { ClientApiError } from "~/lib/api-error";
import { getRegistrationConfig, register } from "~/lib/auth";
import { beginAuthTransition } from "~/lib/auth-transition";
import { RedirectIfAuth, useAuth } from "~/lib/auth-context";
import styles from "~/styles/auth.module.css";

export default function RegisterPage() {
  const navigate = useNavigate();
  const { refetch } = useAuth();
  const [error, setError] = createSignal<string | null>(null);
  const [pending, setPending] = createSignal(false);
  const [registrationCode, setRegistrationCode] = createSignal("");
  const [registrationConfig] = createResource(getRegistrationConfig);

  const inviteRequired = () => registrationConfig()?.registration_code_required === true;
  const inviteReady = () => {
    if (registrationConfig.loading) {
      return false;
    }
    return !inviteRequired() || registrationCode().trim().length > 0;
  };

  const handleSubmit = async (e: SubmitEvent) => {
    e.preventDefault();
    setError(null);
    setPending(true);

    const formData = new FormData(e.currentTarget as HTMLFormElement);
    const data = Object.fromEntries(formData.entries());

    try {
      await register(data);
      beginAuthTransition({
        title: "Welcome aboard",
        hint: "Setting up your workspace",
      });
      await refetch();
      navigate("/onboarding", { replace: true });
    } catch (err: unknown) {
      if (err instanceof ClientApiError) {
        if (err.code === "invalid_registration_code") {
          setError("Invalid or expired registration code. Contact an administrator for a new code.");
        } else if (err.code === "registration_code_required") {
          setError("A registration code from an administrator is required.");
        } else {
          setError(err.message || "Registration failed");
        }
      } else {
        setError(err instanceof Error ? err.message : "Registration failed");
      }
    } finally {
      setPending(false);
    }
  };

  return (
    <RedirectIfAuth>
      <AuthLayout
        eyebrow="Get started"
        title="Create account"
        subtitle="Start tracking spending with automated sync and smart tags."
      >
        <Title>Create Account | Financial Tracker</Title>

        <div class={styles.form}>
          <Show when={inviteRequired()}>
            <div class={styles.callout} role="note">
              Registration is invite-only. Contact an administrator to receive a temporary code
              before creating an account.
            </div>
          </Show>

          <SSOButtons
            label="Sign up with Google"
            mode="register"
            registrationCode={registrationCode()}
            disabled={!inviteReady()}
          />

          <div class={styles.divider}>or email</div>

          <form onSubmit={handleSubmit} class={styles.form}>
            <Show when={inviteRequired()}>
              <div class={styles.field}>
                <label class={styles.label} for="registration_code">
                  Registration code
                </label>
                <input
                  class={styles.input}
                  id="registration_code"
                  name="registration_code"
                  type="text"
                  autocomplete="off"
                  autocapitalize="characters"
                  spellcheck={false}
                  required
                  value={registrationCode()}
                  onInput={(event) => setRegistrationCode(event.currentTarget.value.toUpperCase())}
                  placeholder="Enter code from admin"
                />
              </div>
            </Show>

            <div class={`${styles.row} ${styles.form}`}>
              <div class={styles.field}>
                <label class={styles.label} for="first_name">
                  First name
                </label>
                <input
                  class={styles.input}
                  id="first_name"
                  name="first_name"
                  type="text"
                  autocomplete="given-name"
                  required
                />
              </div>

              <div class={styles.field}>
                <label class={styles.label} for="last_name">
                  Last name
                </label>
                <input
                  class={styles.input}
                  id="last_name"
                  name="last_name"
                  type="text"
                  autocomplete="family-name"
                  required
                />
              </div>
            </div>

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
                autocomplete="new-password"
                required
              />
            </div>

            <div class={styles.field}>
              <label class={styles.label} for="confirm_password">
                Confirm password
              </label>
              <input
                class={styles.input}
                id="confirm_password"
                name="confirm_password"
                type="password"
                autocomplete="new-password"
                required
              />
            </div>

            <div class={styles.hint}>
              <strong>Password requirements:</strong> 8-30 characters with at least one
              special character from{" "}
              <code>!@#$%^&*(),.?":{}|&lt;&gt;</code>
            </div>

            <FormError message={error() || undefined} />

            <button class={styles.button} type="submit" disabled={pending() || !inviteReady()}>
              <Show when={pending()} fallback="Create account">
                Creating account...
              </Show>
            </button>
          </form>

          <p class={styles.footer}>
            Already have an account? <a href="/login">Sign in</a>
          </p>
        </div>
      </AuthLayout>
    </RedirectIfAuth>
  );
}
