import { Show, createSignal } from "solid-js";
import { TicketIcon } from "~/components/icons";
import { createRegistrationCode } from "~/lib/admin";
import { ClientApiError } from "~/lib/api-error";
import styles from "~/styles/settings.module.css";

type RegistrationInvitesPanelProps = {
  onMessage: (message: string, type: "ok" | "error" | "info") => void;
};

// formatExpiry renders a unix timestamp as a readable local datetime.
function formatExpiry(unixSeconds: number): string {
  return new Date(unixSeconds * 1000).toLocaleString(undefined, {
    dateStyle: "medium",
    timeStyle: "short",
  });
}

// RegistrationInvitesPanel lets admins issue single-use registration invite codes.
export default function RegistrationInvitesPanel(props: RegistrationInvitesPanelProps) {
  const [pending, setPending] = createSignal(false);
  const [issuedCode, setIssuedCode] = createSignal<string | null>(null);
  const [expiresAt, setExpiresAt] = createSignal<number | null>(null);
  const [copied, setCopied] = createSignal(false);

  const handleGenerate = async () => {
    setPending(true);
    setCopied(false);
    try {
      const result = await createRegistrationCode();
      setIssuedCode(result.code);
      setExpiresAt(result.expires_at);
      props.onMessage("Invite code created. Share it once — it expires in 48 hours.", "ok");
    } catch (err) {
      const message =
        err instanceof ClientApiError
          ? err.message
          : err instanceof Error
            ? err.message
            : "Failed to create invite code";
      props.onMessage(message, "error");
    } finally {
      setPending(false);
    }
  };

  const handleCopy = async () => {
    const code = issuedCode();
    if (!code) {
      return;
    }
    
    const fallbackCopy = () => {
      const textArea = document.createElement("textarea");
      textArea.value = code;
      
      // Avoid scrolling to bottom
      textArea.style.top = "0";
      textArea.style.left = "0";
      textArea.style.position = "fixed";
      
      document.body.appendChild(textArea);
      textArea.focus();
      textArea.select();
      
      try {
        const successful = document.execCommand("copy");
        if (!successful) {
          throw new Error("Fallback copy failed");
        }
      } finally {
        document.body.removeChild(textArea);
      }
    };

    try {
      if (navigator.clipboard && window.isSecureContext) {
        await navigator.clipboard.writeText(code);
      } else {
        fallbackCopy();
      }
      setCopied(true);
      props.onMessage("Code copied to clipboard.", "info");
      window.setTimeout(() => setCopied(false), 2200);
    } catch {
      try {
        // Try fallback one last time if clipboard API failed unexpectedly
        fallbackCopy();
        setCopied(true);
        props.onMessage("Code copied to clipboard.", "info");
        window.setTimeout(() => setCopied(false), 2200);
      } catch {
        props.onMessage("Could not copy automatically. Select and copy the code manually.", "error");
      }
    }
  };

  return (
    <section class={styles.inviteSection}>
      <div class={styles.inviteHeader}>
        <span class={styles.inviteIcon} aria-hidden="true">
          <TicketIcon size={20} />
        </span>
        <div>
          <h2 class={styles.sectionTitle}>Invite codes</h2>
          <p class={styles.sectionHint}>
            Registration is invite-only. Generate a temporary code for someone who needs a new
            account — each code works once and expires after 48 hours.
          </p>
        </div>
      </div>

      <div class={styles.inviteCard}>
        <Show
          when={issuedCode()}
          fallback={
            <div class={styles.inviteEmpty}>
              <p class={styles.inviteEmptyCopy}>
                No active code on screen. Generate a fresh code when you are ready to share it with
                a new user.
              </p>
              <button
                type="button"
                class={styles.buttonPrimary}
                onClick={handleGenerate}
                disabled={pending()}
              >
                {pending() ? "Generating..." : "Generate invite code"}
              </button>
            </div>
          }
        >
          {(code) => (
            <div class={styles.inviteReveal}>
              <p class={styles.inviteLabel}>Share this code once</p>
              <div class={styles.inviteCodeRow}>
                <code class={styles.inviteCode}>{code()}</code>
                <button type="button" class={styles.buttonGhost} onClick={handleCopy}>
                  {copied() ? "Copied" : "Copy"}
                </button>
              </div>
              <Show when={expiresAt()}>
                {(expiry) => (
                  <p class={styles.inviteExpiry}>Expires {formatExpiry(expiry())}</p>
                )}
              </Show>
              <div class={styles.inviteActions}>
                <button
                  type="button"
                  class={styles.buttonSecondary}
                  onClick={handleGenerate}
                  disabled={pending()}
                >
                  {pending() ? "Generating..." : "Generate another"}
                </button>
              </div>
            </div>
          )}
        </Show>
      </div>
    </section>
  );
}
