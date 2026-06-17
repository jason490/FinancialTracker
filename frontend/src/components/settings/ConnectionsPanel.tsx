import { For, Show, createResource, createSignal } from "solid-js";
import { EyeIcon, EyeOffIcon, SyncIcon, TrashIcon } from "~/components/icons";
import { formatCurrency, formatDate } from "~/lib/format";
import {
  disconnectPlaidConnection,
  getPlaidConnections,
  managePlaidConnection,
  removeDisconnectedPlaidAccount,
  startNewPlaidConnection,
  syncAllPlaidConnections,
  togglePlaidAccountVisibility,
} from "~/lib/plaid";
import type { PlaidConnectionView } from "~/lib/types";
import styles from "~/styles/settings.module.css";

type ConnectionsPanelProps = {
  onMessage: (message: string, type: "ok" | "error" | "info") => void;
};

// statusBadge returns a badge class for a Plaid connection status.
function statusBadge(status: string) {
  switch (status) {
    case "needs_reauth":
      return styles.badgeWarn;
    case "disconnected":
    case "error":
      return styles.badgeDanger;
    default:
      return styles.badgeOk;
  }
}

// statusLabel humanizes a Plaid connection status.
function statusLabel(status: string) {
  switch (status) {
    case "needs_reauth":
      return "Re-authentication required";
    case "disconnected":
      return "Disconnected";
    case "error":
      return "Connection error";
    default:
      return "Connected";
  }
}

// ConnectionsPanel manages Plaid institutions and linked accounts.
export default function ConnectionsPanel(props: ConnectionsPanelProps) {
  const [pendingAction, setPendingAction] = createSignal<string | null>(null);
  const [connections, { refetch }] = createResource(getPlaidConnections);

  const runAction = async (key: string, action: () => Promise<void>, message: string) => {
    setPendingAction(key);
    try {
      await action();
      await refetch();
      props.onMessage(message, "ok");
    } catch (err) {
      props.onMessage(err instanceof Error ? err.message : "Action failed", "error");
    } finally {
      setPendingAction(null);
    }
  };

  const handleConnect = () =>
    runAction("connect", startNewPlaidConnection, "Bank connection updated.");

  const handleSyncAll = () =>
    runAction("sync-all", syncAllPlaidConnections, "All connections synced.");

  const handleManage = (connection: PlaidConnectionView) =>
    runAction(
      `manage-${connection.row_id}`,
      () => managePlaidConnection(connection.row_id),
      `${connection.institution_name} updated.`
    );

  const handleDisconnect = (connection: PlaidConnectionView) => {
    if (
      !window.confirm(
        `Disconnect ${connection.institution_name}? This removes all associated accounts and transactions.`
      )
    ) {
      return;
    }

    runAction(
      `disconnect-${connection.row_id}`,
      () => disconnectPlaidConnection(connection.row_id),
      `${connection.institution_name} disconnected.`
    );
  };

  const handleToggleVisibility = (accountId: string) =>
    runAction(
      `toggle-${accountId}`,
      async () => {
        await togglePlaidAccountVisibility(accountId);
      },
      "Account visibility updated."
    );

  const handleRemoveAccount = (accountId: string, accountName: string) => {
    if (
      !window.confirm(
        `Remove ${accountName}? This permanently deletes its transactions.`
      )
    ) {
      return;
    }

    runAction(
      `remove-${accountId}`,
      () => removeDisconnectedPlaidAccount(accountId),
      `${accountName} removed.`
    );
  };

  return (
    <div class={styles.panelInner}>
      <section>
        <div class={styles.toolbar}>
          <div>
            <h2 class={styles.sectionTitle}>Bank connections</h2>
            <p class={styles.sectionHint}>
              Link institutions, manage accounts in Plaid, and keep balances in sync.
            </p>
          </div>
          <div class={styles.actions}>
            <button
              type="button"
              class={styles.buttonSecondary}
              disabled={pendingAction() !== null}
              onClick={handleSyncAll}
            >
              <SyncIcon size={16} />
              {pendingAction() === "sync-all" ? "Syncing..." : "Sync all"}
            </button>
            <button
              type="button"
              class={styles.buttonPrimary}
              disabled={pendingAction() !== null}
              onClick={handleConnect}
            >
              {pendingAction() === "connect" ? "Opening Plaid..." : "Connect bank"}
            </button>
          </div>
        </div>

        <div class={styles.callout}>
          <p class={styles.calloutTitle}>Managing accounts through Plaid</p>
          <ul class={styles.calloutList}>
            <li>
              Use <strong>Manage accounts</strong> to add accounts, re-authenticate, or uncheck
              accounts you no longer want linked.
            </li>
            <li>
              Unchecking an account in Plaid marks it disconnected here. Active accounts cannot be
              deleted from FinancialTracker.
            </li>
            <li>
              To permanently remove a disconnected account and its transactions, use the remove
              action after disconnecting it in Plaid.
            </li>
          </ul>
        </div>
      </section>

      <Show
        when={!connections.loading}
        fallback={<div class={styles.statusInfo}>Loading connections...</div>}
      >
        <Show
          when={(connections()?.connections.length || 0) > 0}
          fallback={
            <div class={styles.emptyState}>
              <h3 class={styles.emptyTitle}>No banks connected</h3>
              <p class={styles.emptyCopy}>
                Connect your first institution to start syncing transactions automatically.
              </p>
              <button
                type="button"
                class={styles.buttonPrimary}
                disabled={pendingAction() !== null}
                onClick={handleConnect}
              >
                Connect bank
              </button>
            </div>
          }
        >
          <div class={styles.connectionList}>
            <For each={connections()?.connections || []}>
              {(connection) => (
                <article class={styles.connectionCard}>
                  <div class={styles.connectionHeader}>
                    <div class={styles.connectionMeta}>
                      <div class={styles.institutionMark}>
                        {(connection.institution_name || "B").charAt(0).toUpperCase()}
                      </div>
                      <div>
                        <div class={styles.connectionTitleRow}>
                          <h3 class={styles.connectionName}>{connection.institution_name}</h3>
                          <span class={`${styles.badge} ${statusBadge(connection.status)}`}>
                            {statusLabel(connection.status)}
                          </span>
                        </div>
                        <p class={styles.connectionSub}>
                          Linked {formatDate(connection.created_at)}
                          {connection.last_synced
                            ? ` · Last synced ${formatDate(connection.last_synced)}`
                            : ""}
                        </p>
                      </div>
                    </div>

                    <div class={styles.connectionActions}>
                      <Show when={connection.status !== "disconnected"}>
                        <button
                          type="button"
                          class={styles.buttonSecondary}
                          disabled={pendingAction() !== null}
                          onClick={() => handleManage(connection)}
                        >
                          {pendingAction() === `manage-${connection.row_id}`
                            ? "Opening..."
                            : "Manage accounts"}
                        </button>
                      </Show>
                      <button
                        type="button"
                        class={styles.buttonDanger}
                        disabled={pendingAction() !== null}
                        onClick={() => handleDisconnect(connection)}
                      >
                        {pendingAction() === `disconnect-${connection.row_id}`
                          ? "Disconnecting..."
                          : "Disconnect"}
                      </button>
                    </div>
                  </div>

                  <div class={styles.accountList}>
                    <Show
                      when={connection.accounts.length > 0}
                      fallback={<p class={styles.connectionSub}>No accounts found for this connection.</p>}
                    >
                      <For each={connection.accounts}>
                        {(account) => (
                          <div
                            class={`${styles.accountRow} ${
                              account.status === "disconnected" || account.is_hidden
                                ? styles.accountRowMuted
                                : ""
                            }`}
                          >
                            <div>
                              <p class={styles.accountName}>{account.name}</p>
                              <p class={styles.accountMeta}>
                                ****{account.mask} · {account.subtype}
                              </p>
                              <Show when={account.status === "disconnected"}>
                                <span class={`${styles.badge} ${styles.badgeDanger}`}>
                                  Disconnected in Plaid
                                </span>
                              </Show>
                              <Show when={account.is_hidden}>
                                <span class={`${styles.badge} ${styles.badgeMuted}`}>Hidden</span>
                              </Show>
                            </div>

                            <div class={styles.accountActions}>
                              <span class={styles.accountBalance}>
                                {formatCurrency(account.balance)}
                              </span>
                              <button
                                type="button"
                                class={styles.iconButton}
                                title={account.is_hidden ? "Unhide account" : "Hide account"}
                                disabled={pendingAction() !== null}
                                onClick={() => handleToggleVisibility(account.plaid_account_id)}
                              >
                                {account.is_hidden ? <EyeIcon size={16} /> : <EyeOffIcon size={16} />}
                              </button>
                              <button
                                type="button"
                                class={styles.iconButton}
                                title={
                                  account.status === "disconnected"
                                    ? "Remove account"
                                    : "Disconnect in Plaid to delete"
                                }
                                disabled={
                                  pendingAction() !== null || account.status !== "disconnected"
                                }
                                onClick={() =>
                                  handleRemoveAccount(account.plaid_account_id, account.name)
                                }
                              >
                                <TrashIcon size={16} />
                              </button>
                            </div>
                          </div>
                        )}
                      </For>
                    </Show>
                  </div>
                </article>
              )}
            </For>
          </div>
        </Show>
      </Show>
    </div>
  );
}
