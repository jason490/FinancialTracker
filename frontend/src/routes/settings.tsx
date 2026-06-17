import { Title } from "@solidjs/meta";
import { useSearchParams } from "@solidjs/router";
import { Match, Show, Switch, createEffect, createResource, createSignal } from "solid-js";
import AccountPanel from "~/components/settings/AccountPanel";
import AppearancePanel from "~/components/settings/AppearancePanel";
import ConnectionsPanel from "~/components/settings/ConnectionsPanel";
import { BankIcon, PaletteIcon, UserIcon } from "~/components/icons";
import AppLayout from "~/layouts/AppLayout";
import { useAuth } from "~/lib/auth-context";
import { getSettings } from "~/lib/settings";
import type { SettingsProfile } from "~/lib/types";
import styles from "~/styles/settings.module.css";

type SettingsTab = "account" | "connections" | "appearance";

const TABS: Array<{ id: SettingsTab; label: string; icon: typeof UserIcon }> = [
  { id: "account", label: "Account", icon: UserIcon },
  { id: "connections", label: "Connections", icon: BankIcon },
  { id: "appearance", label: "Appearance", icon: PaletteIcon },
];

// SettingsPage centralizes profile, security, Plaid, and theme preferences.
export default function SettingsPage() {
  const { user: sessionProfile, refetch: refetchAuth } = useAuth();
  const [searchParams, setSearchParams] = useSearchParams();
  const [activeTab, setActiveTab] = createSignal<SettingsTab>("account");
  const [message, setMessage] = createSignal<{ text: string; type: "ok" | "error" | "info" } | null>(
    null
  );
  const [settings, { refetch: refetchSettings }] = createResource(
    sessionProfile,
    async () => getSettings()
  );

  createEffect(() => {
    const tab = searchParams.tab;
    if (tab === "account" || tab === "connections" || tab === "appearance") {
      setActiveTab(tab);
    }
  });

  createEffect(() => {
    if (searchParams.success === "linked") {
      setMessage({ text: "Google account linked successfully.", type: "ok" });
      void refetchSettings();
      void refetchAuth();
      setSearchParams({ tab: "account" }, { replace: true });
    }
  });

  const reauthSuccess = () => searchParams.reauth_success === "true";

  const clearReauthParam = () => {
    if (searchParams.reauth_success) {
      setSearchParams({ tab: "account" }, { replace: true });
    }
  };

  createEffect(() => {
    const error = searchParams.error;
    if (!error) {
      return;
    }

    const messages: Record<string, string> = {
      identity_mismatch: "Google account does not match your profile.",
      reauth_failed: "Google verification failed. Try again.",
      authentication_failed: "Authentication failed. Sign in and try again.",
    };

    setMessage({
      text: messages[error] || "Verification failed. Try again.",
      type: "error",
    });
    setSearchParams({ tab: "account" }, { replace: true });
  });

  const handleUpdated = (profile: SettingsProfile) => {
    void refetchSettings();
    void refetchAuth();
    setMessage({
      text: `${profile.first_name} ${profile.last_name}`.trim()
        ? "Settings updated."
        : "Settings updated.",
      type: "ok",
    });
  };

  const handleMessage = (text: string, type: "ok" | "error" | "info") => {
    setMessage({ text, type });
  };

  const selectTab = (tab: SettingsTab) => {
    setActiveTab(tab);
    setSearchParams({ tab }, { replace: true });
    setMessage(null);
  };

  return (
    <AppLayout>
      <Title>Settings | Financial Tracker</Title>

      <div class={styles.page}>
        <header class={styles.header}>
          <p class={styles.eyebrow}>Preferences</p>
          <h1 class={styles.title}>Settings</h1>
          <p class={styles.subtitle}>
            Tune your profile, secure sign-in, bank connections, and the look of your workspace.
          </p>
        </header>

        <Show when={message()}>
          {(current) => (
            <div
              class={
                current().type === "error"
                  ? styles.statusError
                  : current().type === "info"
                    ? styles.statusInfo
                    : styles.statusOk
              }
              role="status"
            >
              {current().text}
            </div>
          )}
        </Show>

        <div class={styles.layout}>
          <div class={styles.tabRail} role="tablist" aria-label="Settings sections">
            {TABS.map((tab) => {
              const Icon = tab.icon;
              return (
                <button
                  type="button"
                  role="tab"
                  aria-selected={activeTab() === tab.id}
                  class={`${styles.tabButton} ${activeTab() === tab.id ? styles.tabButtonActive : ""}`}
                  onClick={() => selectTab(tab.id)}
                >
                  <Icon size={18} />
                  <span>{tab.label}</span>
                </button>
              );
            })}
          </div>

          <section class={styles.panel} role="tabpanel">
            <Show
              when={!settings.loading}
              fallback={<div class={`${styles.panelInner} ${styles.tabPanel}`}>Loading settings...</div>}
            >
              <Show when={settings()}>
                {(profile) => (
                  <Show keyed when={activeTab()}>
                    {(tab) => (
                      <div class={styles.tabPanel}>
                        <Switch>
                          <Match when={tab === "account"}>
                            <AccountPanel
                              profile={profile()}
                              onUpdated={handleUpdated}
                              onMessage={handleMessage}
                              reauthSuccess={reauthSuccess()}
                              onReauthHandled={clearReauthParam}
                            />
                          </Match>
                          <Match when={tab === "connections"}>
                            <ConnectionsPanel onMessage={handleMessage} />
                          </Match>
                          <Match when={tab === "appearance"}>
                            <AppearancePanel
                              profile={profile()}
                              onUpdated={handleUpdated}
                              onMessage={handleMessage}
                            />
                          </Match>
                        </Switch>
                      </div>
                    )}
                  </Show>
                )}
              </Show>
            </Show>
          </section>
        </div>
      </div>
    </AppLayout>
  );
}
