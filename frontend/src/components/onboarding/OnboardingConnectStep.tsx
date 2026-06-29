import { Show, createSignal } from "solid-js";
import { BankIcon } from "~/components/icons";
import { PLAID_SYNC_LAG_HINT, reportSyncError } from "~/lib/api-error";
import { startNewConnection } from "~/lib/connections";
import type { ConnectionsPayload } from "~/lib/types";
import type { Resource } from "solid-js";
import styles from "~/styles/onboarding.module.css";

type OnboardingConnectStepProps = {
  connections: Resource<ConnectionsPayload>;
  refetchConnections: () => void;
  onBack: () => void;
  onFinish: () => void;
  onError: (message: string) => void;
  finishing: boolean;
};

// OnboardingConnectStep guides users to link their first bank or skip for later.
export default function OnboardingConnectStep(props: OnboardingConnectStepProps) {
  const [pending, setPending] = createSignal(false);

  const hasConnections = () => (props.connections()?.connections.length || 0) > 0;
  const isPlaidProvider = () => props.connections()?.provider === "plaid";

  const handleConnect = async () => {
    setPending(true);
    try {
      await startNewConnection();
      props.refetchConnections();
    } catch (err) {
      reportSyncError(err, (message, type) => {
        if (type === "error") {
          props.onError(message);
        }
      });
    } finally {
      setPending(false);
    }
  };

  return (
    <section class={styles.stepPanel}>
      <p class={styles.stepEyebrow}>Step 3 of 3</p>
      <h2 class={styles.stepTitle}>Connect your bank</h2>
      <p class={styles.stepLead}>
        Link an institution to start syncing transactions automatically. You can always add more
        from Settings.
      </p>

      <Show when={isPlaidProvider()}>
        <p class={styles.syncHint}>{PLAID_SYNC_LAG_HINT}</p>
      </Show>

      <div class={styles.connectCard}>
        <div class={styles.connectIcon} aria-hidden="true">
          <BankIcon size={28} />
        </div>
        <Show
          when={hasConnections()}
          fallback={
            <>
              <h3 class={styles.connectTitle}>No banks linked yet</h3>
              <p class={styles.connectCopy}>
                Securely connect through our financial data partner. Credentials are never stored
                on FinancialTracker servers.
              </p>
            </>
          }
        >
          <h3 class={styles.connectTitle}>Bank connected</h3>
          <p class={styles.connectCopy}>
            Your institution is linked. Transactions will appear on your dashboard after the first
            sync.
          </p>
        </Show>

        <button
          type="button"
          class={styles.primaryButton}
          disabled={pending() || props.finishing}
          onClick={() => handleConnect()}
        >
          {pending()
            ? "Opening link flow..."
            : hasConnections()
              ? "Connect another bank"
              : "Connect bank"}
        </button>
      </div>

      <div class={styles.stepActions}>
        <button
          type="button"
          class={styles.secondaryButton}
          disabled={props.finishing}
          onClick={() => props.onBack()}
        >
          Back
        </button>
        <button
          type="button"
          class={styles.primaryButton}
          disabled={props.finishing || pending()}
          onClick={() => props.onFinish()}
        >
          {props.finishing ? "Finishing..." : hasConnections() ? "Go to dashboard" : "Skip for now"}
        </button>
      </div>
    </section>
  );
}
