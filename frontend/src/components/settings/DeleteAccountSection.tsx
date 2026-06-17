import { useNavigate } from "@solidjs/router";
import { Show, createEffect, createSignal } from "solid-js";
import { getPublicApiUrl, getFrontendUrl } from "~/lib/env";
import {
  confirmDeleteAccount,
  getDeleteAccountReauthStatus,
  verifyDeleteAccount,
} from "~/lib/settings";
import type { SettingsProfile } from "~/lib/types";
import styles from "~/styles/settings.module.css";

type DeleteStep = "idle" | "reauth" | "confirm";

type DeleteAccountSectionProps = {
  profile: SettingsProfile;
  reauthSuccess?: boolean;
  onMessage: (message: string, type: "ok" | "error") => void;
  onReauthHandled?: () => void;
};

// DeleteAccountSection guides users through re-auth and permanent account deletion.
export default function DeleteAccountSection(props: DeleteAccountSectionProps) {
  const navigate = useNavigate();
  const [step, setStep] = createSignal<DeleteStep>("idle");
  const [email, setEmail] = createSignal(props.profile.email);
  const [password, setPassword] = createSignal("");
  const [pendingVerify, setPendingVerify] = createSignal(false);
  const [pendingDelete, setPendingDelete] = createSignal(false);

  createEffect(() => {
    setEmail(props.profile.email);
  });

  const hasGoogle = () => props.profile.sso_providers.includes("google");
  const hasPassword = () => props.profile.has_password;

  const googleReauthUrl = () => {
    const returnTo = `${getFrontendUrl()}/settings?tab=account&reauth_success=true`;
    return `${getPublicApiUrl()}/api/v1/auth/google?return_to=${encodeURIComponent(returnTo)}&action=reauth-delete`;
  };

  const advanceIfVerified = async () => {
    try {
      const status = await getDeleteAccountReauthStatus();
      if (status.reauth_verified) {
        setStep("confirm");
        return true;
      }
    } catch (err) {
      props.onMessage(
        err instanceof Error ? err.message : "Failed to verify re-authentication",
        "error"
      );
    }
    return false;
  };

  createEffect(() => {
    if (!props.reauthSuccess) {
      return;
    }

    void (async () => {
      const verified = await advanceIfVerified();
      if (verified) {
        props.onMessage("Identity verified. Confirm deletion to continue.", "ok");
      }
      props.onReauthHandled?.();
    })();
  });

  const handleStart = () => {
    setPassword("");
    setStep("reauth");
  };

  const handleCancel = () => {
    setPassword("");
    setStep("idle");
  };

  const handleVerifySubmit = async (event: SubmitEvent) => {
    event.preventDefault();
    setPendingVerify(true);
    try {
      await verifyDeleteAccount({ email: email(), password: password() });
      setPassword("");
      setStep("confirm");
      props.onMessage("Identity verified. Confirm deletion to continue.", "ok");
    } catch (err) {
      props.onMessage(err instanceof Error ? err.message : "Verification failed", "error");
    } finally {
      setPendingVerify(false);
    }
  };

  const handleConfirmDelete = async () => {
    setPendingDelete(true);
    try {
      await confirmDeleteAccount();
      navigate("/", { replace: true });
    } catch (err) {
      const message = err instanceof Error ? err.message : "Failed to delete account";
      props.onMessage(message, "error");
      if (message.toLowerCase().includes("verification expired")) {
        setStep("reauth");
      }
    } finally {
      setPendingDelete(false);
    }
  };

  return (
    <section class={styles.dangerZone}>
      <div class={styles.dangerIntro}>
        <h2 class={styles.sectionTitle}>Delete account</h2>
        <p class={styles.sectionHint}>
          Permanently remove your profile, bank connections, transactions, tags, and all saved
          preferences. This action cannot be undone.
        </p>
      </div>

      <Show when={step() === "idle"}>
        <div class={styles.dangerCard}>
          <p class={styles.dangerCopy}>
            Once deleted, your data is erased from FinancialTracker immediately. Export anything
            you need before continuing.
          </p>
          <button type="button" class={styles.buttonDanger} onClick={handleStart}>
            Delete my account
          </button>
        </div>
      </Show>

      <Show when={step() === "reauth"}>
        <div class={styles.dangerCard}>
          <div class={styles.dangerBanner}>
            <p class={styles.dangerBannerTitle}>Re-authentication required</p>
            <p class={styles.dangerBannerCopy}>
              Verify your identity before we proceed with account deletion.
            </p>
          </div>

          <Show when={hasPassword()}>
            <form class={styles.fieldGrid} onSubmit={handleVerifySubmit}>
              <div class={styles.field}>
                <label class={styles.label} for="delete-email">
                  Confirm email
                </label>
                <input
                  id="delete-email"
                  class={styles.input}
                  type="email"
                  value={email()}
                  onInput={(event) => setEmail(event.currentTarget.value)}
                  required
                />
              </div>
              <div class={styles.field}>
                <label class={styles.label} for="delete-password">
                  Password
                </label>
                <input
                  id="delete-password"
                  class={styles.input}
                  type="password"
                  value={password()}
                  onInput={(event) => setPassword(event.currentTarget.value)}
                  required
                />
              </div>
              <div class={styles.actions}>
                <button class={styles.buttonDangerSolid} type="submit" disabled={pendingVerify()}>
                  {pendingVerify() ? "Verifying..." : "Verify with password"}
                </button>
              </div>
            </form>
          </Show>

          <Show when={hasGoogle()}>
            <Show when={hasPassword()}>
              <div class={styles.divider}>
                <span>Or verify with SSO</span>
              </div>
            </Show>
            <a class={styles.ssoVerifyButton} href={googleReauthUrl()} rel="external">
              <img
                src="https://www.gstatic.com/firebasejs/ui/2.0.0/images/auth/google.svg"
                alt=""
                width="18"
                height="18"
              />
              Verify with Google
            </a>
          </Show>

          <Show when={!hasPassword() && !hasGoogle()}>
            <p class={styles.dangerCopy}>
              Add a password or connect Google before you can delete this account.
            </p>
          </Show>

          <button type="button" class={styles.cancelLink} onClick={handleCancel}>
            Cancel
          </button>
        </div>
      </Show>

      <Show when={step() === "confirm"}>
        <div class={`${styles.dangerCard} ${styles.dangerCardFinal}`}>
          <div class={styles.dangerFinalHeader}>
            <p class={styles.dangerFinalEyebrow}>Permanent action</p>
            <p class={styles.dangerFinalTitle}>This will delete everything.</p>
          </div>
          <p class={styles.dangerCopy}>
            All financial data, linked accounts, transaction history, tags, and dashboard layouts
            tied to <strong>{props.profile.email}</strong> will be removed permanently.
          </p>
          <div class={styles.actions}>
            <button
              type="button"
              class={styles.buttonDangerSolid}
              onClick={handleConfirmDelete}
              disabled={pendingDelete()}
            >
              {pendingDelete() ? "Deleting..." : "Yes, delete everything"}
            </button>
            <button type="button" class={styles.buttonSecondary} onClick={handleCancel}>
              Keep my account
            </button>
          </div>
        </div>
      </Show>
    </section>
  );
}
