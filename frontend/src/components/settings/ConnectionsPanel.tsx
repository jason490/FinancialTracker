import { Index, Show, Suspense, createSignal } from "solid-js";
import { EyeIcon, EyeOffIcon, SyncIcon, TrashIcon } from "~/components/icons";
import { formatCurrency, formatDate } from "~/lib/format";
import {
  disconnectConnection,
  manageConnection,
  removeDisconnectedAccount,
  startNewConnection,
  syncAllConnections,
  toggleAccountVisibility,
} from "~/lib/connections";
import {
  PLAID_SYNC_LAG_HINT,
  reportSyncError,
} from "~/lib/api-error";
import type { Resource } from "solid-js";
import type { ConnectionView, ConnectionsPayload } from "~/lib/types";
import styles from "~/styles/settings.module.css";

type ConnectionsPanelProps = {
  connections: Resource<ConnectionsPayload>;
  refetchConnections: () => void;
  onMessage: (message: string, type: "ok" | "error" | "info") => void;
};

// statusBadge returns a badge class for a connection status.
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

// statusLabel humanizes a connection status.
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

// ConnectionsPanel manages linked institutions and accounts.
export default function ConnectionsPanel(props: ConnectionsPanelProps) {
  const [pendingAction, setPendingAction] = createSignal<string | null>(null);
  const connections = props.connections;
  const refetch = props.refetchConnections;

  const runAction = async (key: string, action: () => Promise<void>, message: string) => {
    setPendingAction(key);
    try {
      await action();
      refetch();
      props.onMessage(message, "ok");
      setPendingAction(null);
    } catch (err) {
      reportSyncError(err, props.onMessage);
      setPendingAction(null);
    }
  };

  const handleConnect = () =>
    runAction("connect", startNewConnection, "Bank connection updated.");

  const handleSyncAll = () =>
    runAction("sync-all", syncAllConnections, "All connections synced.");

  const handleManage = (connection: ConnectionView) =>
    runAction(
      `manage-${connection.row_id}`,
      () => manageConnection(connection.row_id),
      `${connection.institution_name} updated.`
    );

  const handleDisconnect = (connection: ConnectionView) => {
    if (
      !window.confirm(
        `Disconnect ${connection.institution_name}? This removes all associated accounts and transactions.`
      )
    ) {
      return;
    }

    runAction(
      `disconnect-${connection.row_id}`,
      () => disconnectConnection(connection.row_id),
      `${connection.institution_name} disconnected.`
    );
  };

  const handleToggleVisibility = (accountId: string) =>
    runAction(
      `toggle-${accountId}`,
      async () => {
        await toggleAccountVisibility(accountId);
      },
      "Account visibility updated."
    );

  const handleRemoveAccount = (accountId: string, accountName: string) => {
    if (!window.confirm(`Remove ${accountName}? This permanently deletes its transactions.`)) {
      return;
    }

    runAction(
      `remove-${accountId}`,
      () => removeDisconnectedAccount(accountId),
      `${accountName} removed.`
    );
  };

  return (
    <div class={styles.panelInner}>
      <Suspense fallback={<div class={styles.statusInfo}>Loading connections...</div>}>
        <section class={styles.connectionsIntro}>
          <div class={styles.toolbar}>
            <div>
              <h2 class={styles.sectionTitle}>Bank connections</h2>
              <p class={styles.sectionHint}>
                Link institutions, manage accounts, and keep balances in sync.
              </p>
              <Show when={connections()?.provider === "plaid"}>
                <p class={styles.syncLagHint}>{PLAID_SYNC_LAG_HINT}</p>
              </Show>
            </div>
            <div class={styles.actions}>
              <button
                type="button"
                class={styles.buttonSecondary}
                disabled={pendingAction() === "sync-all"}
                onClick={handleSyncAll}
              >
                <SyncIcon size={16} />
                {pendingAction() === "sync-all" ? "Syncing..." : "Sync all"}
              </button>
              <button
                type="button"
                class={styles.buttonPrimary}
                disabled={
                  pendingAction() === "connect" ||
                  (!!connections()?.usage &&
                    connections()!.usage.active_items >= connections()!.usage.limits.max_items)
                }
                onClick={handleConnect}
              >
                {pendingAction() === "connect" ? "Opening link flow..." : "Connect bank"}
              </button>
            </div>
          </div>

          <div class={styles.callout}>
            <p class={styles.calloutTitle}>Managing linked accounts</p>
            <ul class={styles.calloutList}>
              <li>
                Use <strong>Manage accounts</strong> to add accounts, re-authenticate, or unlink
                accounts you no longer want synced.
              </li>
              <li>
                Disconnected accounts stay in your history until you remove them from
                FinancialTracker.
              </li>
              <li>
                To permanently remove a disconnected account and its transactions, use the remove
                action after disconnecting it at your bank provider.
              </li>
            </ul>
          </div>
        </section>

        <Show when={connections() !== undefined}>
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
                  disabled={
                    pendingAction() !== null ||
                    (connections()?.usage &&
                      connections()!.usage.active_items >= connections()!.usage.limits.max_items)
                  }
                  onClick={handleConnect}
                >
                  Connect bank
                </button>
              </div>
            }
          >
            <div class={styles.connectionList}>
              <Index each={connections()?.connections || []}>
                {(connection) => (
                  <article class={styles.connectionCard}>
                    <div class={styles.connectionHeader}>
                      <div class={styles.connectionMeta}>
                        <div class={styles.institutionMark}>
                          {(connection().institution_name || "B").charAt(0).toUpperCase()}
                        </div>
                        <div>
                          <div class={styles.connectionTitleRow}>
                            <h3 class={styles.connectionName}>{connection().institution_name}</h3>
                            <span class={`${styles.badge} ${statusBadge(connection().status)}`}>
                              {statusLabel(connection().status)}
                            </span>
                          </div>
                          <p class={styles.connectionSub}>
                            Linked {formatDate(connection().created_at)}
                            {connection().last_synced
                              ? ` · Last synced ${formatDate(connection().last_synced)}`
                              : ""}
                          </p>
                        </div>
                      </div>

                      <div class={styles.connectionActions}>
                        <Show when={connection().status !== "disconnected"}>
                          <button
                            type="button"
                            class={styles.buttonSecondary}
                            disabled={pendingAction() === `manage-${connection().row_id}`}
                            onClick={() => handleManage(connection())}
                          >
                            {pendingAction() === `manage-${connection().row_id}`
                              ? "Opening..."
                              : "Manage accounts"}
                          </button>
                        </Show>
                        <button
                          type="button"
                          class={styles.buttonDanger}
                          disabled={pendingAction() === `disconnect-${connection().row_id}`}
                          onClick={() => handleDisconnect(connection())}
                        >
                          {pendingAction() === `disconnect-${connection().row_id}`
                            ? "Disconnecting..."
                            : "Disconnect"}
                        </button>
                      </div>
                    </div>

                    <div class={styles.accountList}>
                      <Show
                        when={connection().accounts.length > 0}
                        fallback={
                          <p class={styles.connectionSub}>No accounts found for this connection.</p>
                        }
                      >
                        <Index each={connection().accounts}>
                          {(account) => (
                            <div
                              class={`${styles.accountRow} ${
                                account().status === "disconnected" || account().is_hidden
                                  ? styles.accountRowMuted
                                  : ""
                              }`}
                            >
                              <div>
                                <p class={styles.accountName}>{account().name}</p>
                                <p class={styles.accountMeta}>
                                  ****{account().mask} · {account().subtype}
                                </p>
                                <Show when={account().status === "disconnected"}>
                                  <span class={`${styles.badge} ${styles.badgeDanger}`}>
                                    Disconnected
                                  </span>
                                </Show>
                                <Show when={account().is_hidden}>
                                  <span class={`${styles.badge} ${styles.badgeMuted}`}>Hidden</span>
                                </Show>
                              </div>

                              <div class={styles.accountActions}>
                                <span class={styles.accountBalance}>
                                  {formatCurrency(account().balance)}
                                </span>
                                <button
                                  type="button"
                                  class={styles.iconButton}
                                  title={account().is_hidden ? "Unhide account" : "Hide account"}
                                  disabled={pendingAction() === `toggle-${account().account_id}`}
                                  onClick={() => handleToggleVisibility(account().account_id)}
                                >
                                  {account().is_hidden ? (
                                    <EyeIcon size={16} />
                                  ) : (
                                    <EyeOffIcon size={16} />
                                  )}
                                </button>
                                <button
                                  type="button"
                                  class={styles.iconButton}
                                  title={
                                    account().status === "disconnected"
                                      ? "Remove account"
                                      : "Disconnect to delete"
                                  }
                                  disabled={
                                    pendingAction() === `remove-${account().account_id}` ||
                                    account().status !== "disconnected"
                                  }
                                  onClick={() =>
                                    handleRemoveAccount(account().account_id, account().name)
                                  }
                                >
                                  <TrashIcon size={16} />
                                </button>
                              </div>
                            </div>
                          )}
                        </Index>
                      </Show>
                    </div>
                  </article>
                )}
              </Index>
            </div>
          </Show>
        </Show>
      </Suspense>
    </div>
  );
}
