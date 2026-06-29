import styles from "~/styles/settings.module.css";

function SkeletonBlock(props: { class?: string; style?: Record<string, string> }) {
  return (
    <div
      class={`${styles.skeletonBlock} ${props.class ?? ""}`}
      style={props.style}
      aria-hidden="true"
    />
  );
}

function SkeletonField() {
  return (
    <div class={styles.skeletonField}>
      <SkeletonBlock class={styles.skeletonLabel} />
      <SkeletonBlock class={styles.skeletonInput} />
    </div>
  );
}

function SkeletonSsoRow() {
  return (
    <div class={styles.skeletonSsoRow}>
      <div class={styles.skeletonSsoIdentity}>
        <SkeletonBlock class={styles.skeletonSsoIcon} />
        <div class={styles.skeletonSsoText}>
          <SkeletonBlock class={styles.skeletonSsoName} />
          <SkeletonBlock class={styles.skeletonSsoBadge} />
        </div>
      </div>
      <SkeletonBlock class={styles.skeletonSsoAction} />
    </div>
  );
}

// AccountPanelSkeleton mirrors the account tab layout while settings data loads.
export default function AccountPanelSkeleton() {
  return (
    <div
      class={`${styles.panelInner} ${styles.tabPanel}`}
      aria-busy="true"
      aria-label="Loading account settings"
    >
      <section class={styles.skeletonSection} style={{ "--skeleton-delay": "0ms" }}>
        <SkeletonBlock class={styles.skeletonSectionTitle} />
        <SkeletonBlock class={styles.skeletonSectionHint} />

        <div class={styles.card}>
          <div class={`${styles.fieldGrid} ${styles.twoCol}`}>
            <SkeletonField />
            <SkeletonField />
          </div>
          <SkeletonField />
          <SkeletonBlock class={styles.skeletonButtonPrimary} />
        </div>
      </section>

      <section class={styles.skeletonSection} style={{ "--skeleton-delay": "70ms" }}>
        <SkeletonBlock class={styles.skeletonSectionTitle} />
        <SkeletonBlock class={styles.skeletonSectionHintShort} />

        <div class={styles.card}>
          <SkeletonSsoRow />
        </div>

        <div class={styles.card}>
          <SkeletonSsoRow />
        </div>
      </section>

      <section class={styles.skeletonSection} style={{ "--skeleton-delay": "140ms" }}>
        <SkeletonBlock class={styles.skeletonSectionTitle} />
        <SkeletonBlock class={styles.skeletonSectionHintShort} />

        <div class={styles.card}>
          <SkeletonBlock class={styles.skeletonButtonSecondary} />
        </div>
      </section>

      <section
        class={`${styles.dangerZone} ${styles.skeletonSection}`}
        style={{ "--skeleton-delay": "210ms" }}
      >
        <SkeletonBlock class={styles.skeletonSectionTitle} />
        <SkeletonBlock class={styles.skeletonSectionHint} />

        <div class={styles.dangerCard}>
          <SkeletonBlock class={styles.skeletonDangerCopy} />
          <SkeletonBlock class={styles.skeletonButtonDanger} />
        </div>
      </section>
    </div>
  );
}
