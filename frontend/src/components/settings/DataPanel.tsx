import { createSignal } from "solid-js";
import { DatabaseIcon, DownloadIcon } from "~/components/icons";
import { fetchTransactionsExport } from "~/lib/transactions";
import styles from "~/styles/settings.module.css";

type DataPanelProps = {
  onMessage: (message: string, type: "ok" | "error" | "info") => void;
};

// triggerBrowserDownload saves the blob to disk via a temporary anchor link.
// Works in modern browsers and inside the Capacitor Android WebView.
function triggerBrowserDownload(blob: Blob, filename: string) {
  const url = URL.createObjectURL(blob);
  const anchor = document.createElement("a");
  anchor.href = url;
  anchor.download = filename;
  anchor.rel = "noopener";
  document.body.appendChild(anchor);
  anchor.click();
  document.body.removeChild(anchor);
  setTimeout(() => URL.revokeObjectURL(url), 1_000);
}

// DataPanel exposes data-portability tools for the authenticated user. The
// CSV export streams every visible transaction together with each transaction's
// assigned tags and their parent categories.
export default function DataPanel(props: DataPanelProps) {
  const [pending, setPending] = createSignal(false);

  const handleExport = async () => {
    setPending(true);
    try {
      const { blob, filename } = await fetchTransactionsExport();
      triggerBrowserDownload(blob, filename);
      props.onMessage("Export started. Check your downloads folder.", "ok");
    } catch (err) {
      props.onMessage(
        err instanceof Error ? err.message : "Failed to export transactions",
        "error"
      );
    } finally {
      setPending(false);
    }
  };

  return (
    <div class={styles.panelInner}>
      <section>
        <h2 class={styles.sectionTitle}>Data</h2>
        <p class={styles.sectionHint}>
          Take your numbers with you. Exports are generated on demand and never stored on our servers.
        </p>

        <div class={styles.exportCard}>
          <div class={styles.exportCardAura} aria-hidden="true" />

          <header class={styles.exportHeader}>
            <div class={styles.exportMark}>
              <DatabaseIcon size={22} />
            </div>
            <div class={styles.exportHeading}>
              <p class={styles.exportEyebrow}>Export · ZIP (CSV inside)</p>
              <h3 class={styles.exportTitle}>Transactions &amp; tags archive</h3>
              <p class={styles.exportSubtitle}>
                Every visible transaction together with every tag you have applied — packaged as a ZIP
                with a single CSV ready for Excel, Numbers, Google Sheets, or your own scripts.
              </p>
            </div>
          </header>

          <ul class={styles.exportSchema}>
            <li>
              <span class={styles.exportSchemaKey}>date</span>
              <span class={styles.exportSchemaHint}>YYYY-MM-DD (UTC)</span>
            </li>
            <li>
              <span class={styles.exportSchemaKey}>name · merchant_name</span>
              <span class={styles.exportSchemaHint}>raw descriptors from your bank</span>
            </li>
            <li>
              <span class={styles.exportSchemaKey}>amount · currency_sign</span>
              <span class={styles.exportSchemaHint}>negative for debits, positive for credits</span>
            </li>
            <li>
              <span class={styles.exportSchemaKey}>pending · provider · provider_category</span>
              <span class={styles.exportSchemaHint}>sync status and source classification</span>
            </li>
            <li>
              <span class={styles.exportSchemaKey}>tags · tag_categories</span>
              <span class={styles.exportSchemaHint}>semicolon-separated, aligned by index</span>
            </li>
          </ul>

          <div class={styles.exportFooter}>
            <p class={styles.exportFootnote}>
              Hidden accounts are excluded. Each download spends <strong>1 API call</strong> from your
              monthly quota (shared with bank syncs); see the Plan tab for current usage.
            </p>
            <button
              type="button"
              class={styles.buttonPrimary}
              onClick={handleExport}
              disabled={pending()}
            >
              <DownloadIcon size={18} />
              <span>{pending() ? "Preparing archive…" : "Download ZIP"}</span>
            </button>
          </div>
        </div>
      </section>
    </div>
  );
}
