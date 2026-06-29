import styles from "~/styles/home.module.css";

// DashboardPreview renders a decorative dashboard mock for the landing hero.
export default function DashboardPreview() {
  return (
    <div class={styles.previewShell} aria-hidden="true">
      <div class={styles.previewGlow} />

      <div class={styles.previewCard}>
        <header class={styles.previewHeader}>
          <span class={styles.previewDot} />
          <span class={styles.previewDot} />
          <span class={styles.previewDot} />
          <span class={styles.previewTitle}>Your dashboard</span>
        </header>

        <div class={styles.previewGrid}>
          <article class={`${styles.previewWidget} ${styles.previewWidgetWide}`}>
            <p class={styles.previewLabel}>Net worth</p>
            <p class={styles.previewValue}>$48,290</p>
            <div class={styles.previewSparkline}>
              <span style={{ "--h": "42%" }} />
              <span style={{ "--h": "58%" }} />
              <span style={{ "--h": "48%" }} />
              <span style={{ "--h": "72%" }} />
              <span style={{ "--h": "66%" }} />
              <span style={{ "--h": "84%" }} />
              <span style={{ "--h": "78%" }} />
            </div>
          </article>

          <article class={styles.previewWidget}>
            <p class={styles.previewLabel}>This month</p>
            <p class={styles.previewValueSmall}>−$2,140</p>
            <div class={styles.previewBars}>
              <span style={{ "--w": "72%" }} />
              <span style={{ "--w": "48%" }} />
              <span style={{ "--w": "36%" }} />
            </div>
          </article>

          <article class={styles.previewWidget}>
            <p class={styles.previewLabel}>Top tag</p>
            <p class={styles.previewTag}>Groceries</p>
            <div class={styles.previewDonut} />
          </article>

          <article class={`${styles.previewWidget} ${styles.previewWidgetList}`}>
            <p class={styles.previewLabel}>Recent</p>
            <ul class={styles.previewList}>
              <li>
                <span>Whole Foods</span>
                <span>−$84.20</span>
              </li>
              <li>
                <span>Payroll</span>
                <span>+$3,200</span>
              </li>
              <li>
                <span>Metro Transit</span>
                <span>−$12.50</span>
              </li>
            </ul>
          </article>
        </div>
      </div>

      <div class={styles.previewFloatTag}>Auto-tagged</div>
      <div class={styles.previewFloatSync}>Synced just now</div>
    </div>
  );
}
