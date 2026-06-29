import { useNavigate } from "@solidjs/router";
import { Show, createEffect, createSignal } from "solid-js";
import { getPublicApiUrl, getFrontendUrl } from "~/lib/env";
import { updatePassword, updateProfile, unlinkSSO } from "~/lib/settings";
import { logout } from "~/lib/auth";
import { useAuth } from "~/lib/auth-context";
import type { SettingsProfile } from "~/lib/types";
import DeleteAccountSection from "~/components/settings/DeleteAccountSection";
import { LogOutIcon } from "~/components/icons";
import styles from "~/styles/settings.module.css";

type AccountPanelProps = {
  profile: SettingsProfile;
  onUpdated: (profile: SettingsProfile) => void;
  onMessage: (message: string, type: "ok" | "error" | "info") => void;
  reauthSuccess?: boolean;
  onReauthHandled?: () => void;
};

// AccountPanel manages profile details, password, and Google SSO linking.
export default function AccountPanel(props: AccountPanelProps) {
  const navigate = useNavigate();
  const { refetch: refetchAuth } = useAuth();
  const [firstName, setFirstName] = createSignal(props.profile.first_name);
  const [lastName, setLastName] = createSignal(props.profile.last_name);

  createEffect(() => {
    setFirstName(props.profile.first_name);
    setLastName(props.profile.last_name);
  });
  const [showPasswordForm, setShowPasswordForm] = createSignal(false);
  const [currentPassword, setCurrentPassword] = createSignal("");
  const [newPassword, setNewPassword] = createSignal("");
  const [confirmPassword, setConfirmPassword] = createSignal("");
  const [pendingProfile, setPendingProfile] = createSignal(false);
  const [pendingPassword, setPendingPassword] = createSignal(false);
  const [pendingUnlink, setPendingUnlink] = createSignal(false);
  const [pendingLogout, setPendingLogout] = createSignal(false);

  const hasGoogle = () => props.profile.sso_providers.includes("google");

  const handleLogout = async () => {
    setPendingLogout(true);
    try {
      await logout();
      await refetchAuth();
      navigate("/login", { replace: true });
    } catch (err) {
      props.onMessage(err instanceof Error ? err.message : "Failed to log out", "error");
    } finally {
      setPendingLogout(false);
    }
  };

  const handleProfileSubmit = async (event: SubmitEvent) => {
    event.preventDefault();
    setPendingProfile(true);
    try {
      const profile = await updateProfile(firstName(), lastName());
      props.onUpdated(profile);
      props.onMessage("Profile updated.", "ok");
    } catch (err) {
      props.onMessage(err instanceof Error ? err.message : "Failed to update profile", "error");
    } finally {
      setPendingProfile(false);
    }
  };

  const handlePasswordSubmit = async (event: SubmitEvent) => {
    event.preventDefault();
    if (
      !window.confirm(
        "Changing your password will sign you out on all devices. Continue?"
      )
    ) {
      return;
    }

    setPendingPassword(true);
    try {
      await updatePassword({
        current_password: props.profile.has_password ? currentPassword() : undefined,
        new_password: newPassword(),
        confirm_password: confirmPassword(),
      });
      navigate("/login", { replace: true });
    } catch (err) {
      props.onMessage(err instanceof Error ? err.message : "Failed to update password", "error");
    } finally {
      setPendingPassword(false);
    }
  };

  const handleUnlinkGoogle = async () => {
    if (!window.confirm("Unlink Google from this account?")) {
      return;
    }

    setPendingUnlink(true);
    try {
      const profile = await unlinkSSO("google");
      props.onUpdated(profile);
      props.onMessage("Google account unlinked.", "ok");
    } catch (err) {
      props.onMessage(err instanceof Error ? err.message : "Failed to unlink Google", "error");
    } finally {
      setPendingUnlink(false);
    }
  };

  const googleLinkUrl = () => {
    const returnTo = `${getFrontendUrl()}/settings?tab=account&success=linked`;
    return `${getPublicApiUrl()}/api/v1/auth/google?return_to=${encodeURIComponent(returnTo)}&action=link`;
  };

  return (
    <div class={styles.panelInner}>
      <section>
        <h2 class={styles.sectionTitle}>Profile</h2>
        <p class={styles.sectionHint}>Update how your name appears across FinancialTracker.</p>

        <form class={`${styles.card} ${styles.fieldGrid}`} onSubmit={handleProfileSubmit}>
          <div class={`${styles.fieldGrid} ${styles.twoCol}`}>
            <div class={styles.field}>
              <label class={styles.label} for="first-name">
                First name
              </label>
              <input
                id="first-name"
                class={styles.input}
                value={firstName()}
                onInput={(event) => setFirstName(event.currentTarget.value)}
                required
              />
            </div>
            <div class={styles.field}>
              <label class={styles.label} for="last-name">
                Last name
              </label>
              <input
                id="last-name"
                class={styles.input}
                value={lastName()}
                onInput={(event) => setLastName(event.currentTarget.value)}
                required
              />
            </div>
          </div>

          <div class={styles.field}>
            <label class={styles.label} for="email">
              Email
            </label>
            <input
              id="email"
              class={`${styles.input} ${styles.inputReadonly}`}
              value={props.profile.email}
              readOnly
            />
          </div>

          <div class={styles.actions}>
            <button class={styles.buttonPrimary} type="submit" disabled={pendingProfile()}>
              {pendingProfile() ? "Saving..." : "Save profile"}
            </button>
          </div>
        </form>
      </section>

      <section>
        <h2 class={styles.sectionTitle}>Security</h2>
        <p class={styles.sectionHint}>Manage password login and connected sign-in providers.</p>

        <div class={styles.card}>
          <div class={styles.ssoRow}>
            <div>
              <p class={styles.connectionName}>Password login</p>
              <Show
                when={props.profile.has_password}
                fallback={<span class={`${styles.badge} ${styles.badgeMuted}`}>Not set</span>}
              >
                <span class={`${styles.badge} ${styles.badgeOk}`}>Enabled</span>
              </Show>
            </div>
            <button
              type="button"
              class={styles.buttonGhost}
              onClick={() => setShowPasswordForm((value) => !value)}
            >
              {props.profile.has_password ? "Change password" : "Add password"}
            </button>
          </div>

          <form
            class={showPasswordForm() ? styles.collapsedForm : styles.hidden}
            onSubmit={handlePasswordSubmit}
          >
            <Show when={props.profile.has_password}>
              <div class={styles.field}>
                <label class={styles.label} for="current-password">
                  Current password
                </label>
                <input
                  id="current-password"
                  class={styles.input}
                  type="password"
                  value={currentPassword()}
                  onInput={(event) => setCurrentPassword(event.currentTarget.value)}
                  required
                />
              </div>
            </Show>
            <div class={`${styles.fieldGrid} ${styles.twoCol}`}>
              <div class={styles.field}>
                <label class={styles.label} for="new-password">
                  New password
                </label>
                <input
                  id="new-password"
                  class={styles.input}
                  type="password"
                  value={newPassword()}
                  onInput={(event) => setNewPassword(event.currentTarget.value)}
                  required
                />
              </div>
              <div class={styles.field}>
                <label class={styles.label} for="confirm-password">
                  Confirm password
                </label>
                <input
                  id="confirm-password"
                  class={styles.input}
                  type="password"
                  value={confirmPassword()}
                  onInput={(event) => setConfirmPassword(event.currentTarget.value)}
                  required
                />
              </div>
            </div>
            <div class={styles.actions}>
              <button class={styles.buttonPrimary} type="submit" disabled={pendingPassword()}>
                {pendingPassword() ? "Saving..." : "Save password"}
              </button>
            </div>
          </form>
        </div>

        <div class={styles.card}>
          <div class={styles.ssoRow}>
            <div class={styles.ssoIdentity}>
              <div class={styles.ssoIcon}>
                <img
                  src="https://www.gstatic.com/firebasejs/ui/2.0.0/images/auth/google.svg"
                  alt=""
                  width="18"
                  height="18"
                />
              </div>
              <div>
                <p class={styles.connectionName}>Google</p>
                <Show
                  when={hasGoogle()}
                  fallback={<span class={`${styles.badge} ${styles.badgeMuted}`}>Not connected</span>}
                >
                  <span class={`${styles.badge} ${styles.badgeOk}`}>Connected</span>
                </Show>
              </div>
            </div>

            <Show
              when={hasGoogle()}
              fallback={
                <a class={styles.buttonSecondary} href={googleLinkUrl()} rel="external">
                  Link Google
                </a>
              }
            >
              <button
                type="button"
                class={styles.buttonDanger}
                onClick={handleUnlinkGoogle}
                disabled={pendingUnlink()}
              >
                {pendingUnlink() ? "Unlinking..." : "Unlink"}
              </button>
            </Show>
          </div>
        </div>
      </section>

      <section>
        <h2 class={styles.sectionTitle}>Sign out</h2>
        <p class={styles.sectionHint}>Log out of your current session on this device.</p>
        <div class={styles.card}>
          <div class={styles.actions}>
            <button
              class={styles.buttonSecondary}
              type="button"
              onClick={handleLogout}
              disabled={pendingLogout()}
            >
              <LogOutIcon size={18} style={{ "margin-right": "0.5rem", "vertical-align": "middle" }} />
              {pendingLogout() ? "Logging out..." : "Log out"}
            </button>
          </div>
        </div>
      </section>

      <DeleteAccountSection
        profile={props.profile}
        reauthSuccess={props.reauthSuccess}
        onMessage={props.onMessage}
        onReauthHandled={props.onReauthHandled}
      />
    </div>
  );
}
