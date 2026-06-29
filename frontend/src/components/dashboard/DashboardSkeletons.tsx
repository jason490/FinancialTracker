import { For } from "solid-js";
import { WIDGET_IDS, widgetMeta } from "~/lib/dashboard-widgets";
import styles from "~/styles/dashboard.module.css";

type SkeletonVariant = "hero" | "chart" | "list" | "actions";

type SkeletonItem = {
  id: string;
  variant: SkeletonVariant;
};

const DEFAULT_SKELETON_LAYOUT: SkeletonItem[] = [
  { id: WIDGET_IDS.monthCashflow, variant: "hero" },
  { id: WIDGET_IDS.quickActions, variant: "actions" },
  { id: WIDGET_IDS.netWorth, variant: "hero" },
  { id: WIDGET_IDS.spendingByTag, variant: "chart" },
  { id: WIDGET_IDS.spendingTrend, variant: "chart" },
  { id: WIDGET_IDS.incomeByTag, variant: "chart" },
  { id: WIDGET_IDS.cashAccounts, variant: "list" },
  { id: WIDGET_IDS.recentTransactions, variant: "list" },
];

function SkeletonBlock(props: {
  class?: string;
  style?: Record<string, string>;
}) {
  return (
    <div
      class={`${styles.skeletonBlock} ${props.class ?? ""}`}
      style={props.style}
      aria-hidden="true"
    />
  );
}

function DashboardSkeletonWidget(props: { variant: SkeletonVariant }) {
  return (
    <div class={styles.skeletonWidget}>
      <div class={styles.skeletonWidgetHeader}>
        <SkeletonBlock class={styles.skeletonTitle} />
        <SkeletonBlock class={styles.skeletonAside} />
      </div>

      {props.variant === "hero" && (
        <div class={styles.skeletonHeroBody}>
          <SkeletonBlock class={styles.skeletonHeroValue} />
          <div class={styles.skeletonMetricGrid}>
            <SkeletonBlock class={styles.skeletonMetricTile} />
            <SkeletonBlock class={styles.skeletonMetricTile} />
            <SkeletonBlock class={styles.skeletonMetricTile} />
          </div>
        </div>
      )}

      {props.variant === "chart" && (
        <div class={styles.skeletonChartBody}>
          <SkeletonBlock class={styles.skeletonChartArea} />
        </div>
      )}

      {props.variant === "list" && (
        <div class={styles.skeletonListBody}>
          <SkeletonBlock class={styles.skeletonListRow} />
          <SkeletonBlock class={styles.skeletonListRow} />
          <SkeletonBlock class={styles.skeletonListRow} />
          <SkeletonBlock class={styles.skeletonListRowShort} />
        </div>
      )}

      {props.variant === "actions" && (
        <div class={styles.skeletonActionBody}>
          <SkeletonBlock class={styles.skeletonActionButton} />
          <SkeletonBlock class={styles.skeletonActionButton} />
          <SkeletonBlock class={styles.skeletonActionButtonShort} />
        </div>
      )}
    </div>
  );
}

// DashboardSkeletonGrid mirrors the default widget grid while dashboard data loads.
export function DashboardSkeletonGrid() {
  return (
    <section class={styles.gridSection} aria-busy="true" aria-label="Loading dashboard">
      <div class={styles.skeletonGrid}>
        <For each={DEFAULT_SKELETON_LAYOUT}>
          {(item, index) => {
            const meta = widgetMeta(item.id);
            return (
              <div
                class={styles.gridItem}
                classList={{ [styles.span2]: meta?.span === 2 }}
                style={{
                  "--skeleton-delay": `${index() * 70}ms`,
                  "--widget-min-rows": String(meta?.minRows ?? 2),
                  "--widget-max-rows": String(meta?.maxRows ?? 3),
                }}
              >
                <DashboardSkeletonWidget variant={item.variant} />
              </div>
            );
          }}
        </For>
      </div>
    </section>
  );
}
