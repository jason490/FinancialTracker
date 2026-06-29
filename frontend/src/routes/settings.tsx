import { Title } from "@solidjs/meta";
import { useSearchParams } from "@solidjs/router";
import { Show, createEffect, createResource, createSignal } from "solid-js";
import LoadingCrossfade from "~/components/LoadingCrossfade";
import PageStatusBanner, { type PageStatus } from "~/components/PageStatusBanner";
import AccountPanel from "~/components/settings/AccountPanel";
import AccountPanelSkeleton from "~/components/settings/AccountPanelSkeleton";
import AppearancePanel from "~/components/settings/AppearancePanel";
import ConnectionsPanel from "~/components/settings/ConnectionsPanel";
import DataPanel from "~/components/settings/DataPanel";
import PlanPanel from "~/components/settings/PlanPanel";
import { BankIcon, DatabaseIcon, PaletteIcon, UserIcon, CreditCardIcon } from "~/components/icons";
import AppLayout from "~/layouts/AppLayout";
import { useAuth } from "~/lib/auth-context";
import { getConnections } from "~/lib/connections";
import { getSettings } from "~/lib/settings";
import { getSubscription } from "~/lib/subscription";
import type { SettingsProfile } from "~/lib/types";
import styles from "~/styles/settings.module.css";

type SettingsTab = "account" | "connections" | "appearance" | "plan" | "data";

const TABS: Array<{ id: SettingsTab; label: string; icon: typeof UserIcon }> = [
  { id: "account", label: "Account", icon: UserIcon },
  { id: "connections", label: "Connections", icon: BankIcon },
  { id: "appearance", label: "Appearance", icon: PaletteIcon },
  { id: "plan", label: "Plan", icon: CreditCardIcon },
  { id: "data", label: "Data", icon: DatabaseIcon },
];

// SettingsPage centralizes profile, security, Plaid, and theme preferences.
export default function SettingsPage() {
  const { user: sessionProfile, refetch: refetchAuth } = useAuth();
  const [searchParams, setSearchParams] = useSearchParams();
  const [activeTab, setActiveTab] = createSignal<SettingsTab>("account");
  const [message, setMessage] = createSignal<PageStatus | null>(null);
  const [settings, { refetch: refetchSettings }] = createResource(
    sessionProfile,
    async () => getSettings()
  );
  const [connections, { refetch: refetchConnections }] = createResource(getConnections);
  const [subscription, { refetch: refetchSubscription }] = createResource(getSubscription);

  createEffect(() => {
    const tab = searchParams.tab;
    if (
      tab === "account" ||
      tab === "connections" ||
      tab === "appearance" ||
      tab === "plan" ||
      tab === "data"
    ) {
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

  createEffect(() => {
    const checkout = searchParams.checkout;
    if (checkout === "success") {
      setMessage({ text: "Subscription updated. It may take a moment to reflect.", type: "ok" });
      void refetchSubscription();
      setSearchParams({ tab: "plan" }, { replace: true });
    } else if (checkout === "cancelled") {
      setMessage({ text: "Checkout cancelled.", type: "info" });
      setSearchParams({ tab: "plan" }, { replace: true });
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
    setSearchParams({ tab }, { replace: true, scroll: false });
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

        <PageStatusBanner message={message} onDismiss={() => setMessage(null)} />

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
            <LoadingCrossfade
              loading={() => settings.loading}
              ready={() => settings() !== undefined}
              skeleton={<AccountPanelSkeleton />}
            >
              <Show when={settings()}>
                {(profile) => (
                  <div class={styles.tabPanel}>
                    <Show when={activeTab() === "account"}>
                      <AccountPanel
                        profile={profile()}
                        onUpdated={handleUpdated}
                        onMessage={handleMessage}
                        reauthSuccess={reauthSuccess()}
                        onReauthHandled={clearReauthParam}
                      />
                    </Show>
                    <Show when={activeTab() === "connections"}>
                      <ConnectionsPanel
                        connections={connections}
                        refetchConnections={refetchConnections}
                        onMessage={handleMessage}
                      />
                    </Show>
                    <Show when={activeTab() === "appearance"}>
                      <AppearancePanel
                        profile={profile()}
                        onUpdated={handleUpdated}
                        onMessage={handleMessage}
                      />
                    </Show>
                    <Show when={activeTab() === "plan"}>
                      <PlanPanel
                        subscription={subscription}
                        refetchSubscription={refetchSubscription}
                        connections={connections}
                        onMessage={handleMessage}
                      />
                    </Show>
                    <Show when={activeTab() === "data"}>
                      <DataPanel onMessage={handleMessage} />
                    </Show>
                  </div>
                )}
              </Show>
            </LoadingCrossfade>
          </section>
        </div>
      </div>
    </AppLayout>
  );
}
