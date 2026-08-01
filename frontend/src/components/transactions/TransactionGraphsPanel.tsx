import type { ApexOptions } from "apexcharts";
import { Show } from "solid-js";
import ApexChart from "~/components/dashboard/ApexChart";
import {
  donutCenterFontSize,
  formatCurrency,
  formatCurrencyCompact,
  formatMonthLabel,
} from "~/lib/format";
import { tagColorHex } from "~/lib/tag-colors";
import { readCssVar } from "~/lib/themes";
import type { TransactionAnalyticsPayload } from "~/lib/types";
import styles from "~/styles/transactions.module.css";

type TransactionGraphsPanelProps = {
  data: TransactionAnalyticsPayload;
};

// chartBaseOptions returns theme-aware ApexCharts defaults for transactions graphs.
function chartBaseOptions(): ApexOptions {
  const theme = typeof document !== "undefined" ? document.documentElement.dataset.theme : "light";
  const isDark = ["dark", "tokyo-night", "forest", "rose", "midnight"].includes(theme || "light");

  return {
    chart: {
      background: "transparent",
      foreColor: readCssVar("--text"),
      toolbar: { show: false },
      fontFamily: "var(--font-body)",
    },
    grid: {
      borderColor: readCssVar("--border"),
      strokeDashArray: 4,
    },
    dataLabels: { enabled: false },
    legend: {
      labels: { colors: readCssVar("--text-muted") },
    },
    tooltip: {
      theme: isDark ? "dark" : "light",
    },
  };
}

// TransactionGraphsPanel renders filter-scoped totals, a category donut, and monthly cashflow.
export function TransactionGraphsPanel(props: TransactionGraphsPanelProps) {
  const slices = () => props.data.spending_by_category ?? [];
  const trend = () => props.data.spending_trend ?? [];
  const donutLabels = () => slices().map((slice) => slice.category_name);
  const donutValues = () => slices().map((slice) => slice.total);
  const donutColors = () =>
    slices().map((slice) => tagColorHex(slice.color) || readCssVar("--accent"));
  const donutTotal = () => donutValues().reduce((sum, value) => sum + value, 0);
  const centerTotal = () => formatCurrencyCompact(donutTotal());

  return (
    <div class={styles.graphsPanel}>
      <div class={styles.graphsMetrics}>
        <div class={`${styles.graphsMetricTile} ${styles.graphsMetricSpend}`}>
          <span>Total spend</span>
          <strong>{formatCurrency(props.data.cashflow.spend)}</strong>
        </div>
        <div class={`${styles.graphsMetricTile} ${styles.graphsMetricIncome}`}>
          <span>Total income</span>
          <strong>{formatCurrency(props.data.cashflow.income)}</strong>
        </div>
      </div>

      <div class={styles.graphsCharts}>
        <section class={styles.graphsCard}>
          <header class={styles.graphsCardHeader}>
            <div>
              <h2 class={styles.graphsCardTitle}>Spending by Category</h2>
              <p class={styles.graphsCardSubtitle}>Outflows across your tag categories</p>
            </div>
            <Show when={slices().length > 0}>
              <div class={styles.graphsCardAside}>
                <span>Total</span>
                <strong>{formatCurrency(donutTotal())}</strong>
              </div>
            </Show>
          </header>

          <Show
            when={slices().length > 0}
            fallback={<p class={styles.graphsEmpty}>No spending in this filter set.</p>}
          >
            <ApexChart
              class={styles.graphsChartWrap}
              height={260}
              series={donutValues()}
              options={{
                ...chartBaseOptions(),
                chart: { ...chartBaseOptions().chart, type: "donut" },
                labels: donutLabels(),
                colors: donutColors(),
                stroke: { width: 0 },
                tooltip: {
                  ...chartBaseOptions().tooltip,
                  y: {
                    formatter: (value: number) => formatCurrency(value),
                  },
                },
                plotOptions: {
                  pie: {
                    donut: {
                      size: "62%",
                      background: readCssVar("--surface"),
                      labels: {
                        show: true,
                        name: {
                          show: true,
                          color: readCssVar("--text-muted"),
                          fontSize: "0.72rem",
                          fontWeight: 600,
                          offsetY: -6,
                        },
                        value: {
                          show: true,
                          color: readCssVar("--text"),
                          fontSize: donutCenterFontSize(centerTotal()),
                          fontWeight: 700,
                          fontFamily: "var(--font-display)",
                          offsetY: 8,
                        },
                        total: {
                          show: true,
                          showAlways: true,
                          label: "Total",
                          color: readCssVar("--text-muted"),
                          fontWeight: 600,
                          formatter: () => centerTotal(),
                        },
                      },
                    },
                  },
                },
              }}
            />
          </Show>
        </section>

        <section class={styles.graphsCard}>
          <header class={styles.graphsCardHeader}>
            <div>
              <h2 class={styles.graphsCardTitle}>Monthly Cashflow</h2>
              <p class={styles.graphsCardSubtitle}>Income and spending over time</p>
            </div>
          </header>

          <Show
            when={trend().length > 0}
            fallback={<p class={styles.graphsEmpty}>Not enough transaction history yet.</p>}
          >
            <ApexChart
              class={styles.graphsChartWrap}
              height={280}
              series={[
                {
                  name: "Income",
                  data: trend().map((point) => point.income),
                },
                {
                  name: "Spend",
                  data: trend().map((point) => point.total),
                },
              ]}
              options={{
                ...chartBaseOptions(),
                chart: { ...chartBaseOptions().chart, type: "area" },
                colors: ["#10b981", "#f59e0b"],
                stroke: { curve: "smooth", width: 2.5 },
                fill: {
                  type: "gradient",
                  gradient: {
                    shadeIntensity: 0.85,
                    opacityFrom: 0.42,
                    opacityTo: 0.04,
                    stops: [0, 88, 100],
                  },
                },
                legend: {
                  show: true,
                  position: "top",
                  horizontalAlign: "right",
                  labels: { colors: readCssVar("--text-muted") },
                  markers: { size: 5 },
                },
                xaxis: {
                  categories: trend().map((point) => formatMonthLabel(point.month)),
                  labels: { style: { colors: readCssVar("--text-muted") } },
                  axisBorder: { show: false },
                  axisTicks: { show: false },
                },
                yaxis: {
                  labels: {
                    style: { colors: readCssVar("--text-muted") },
                    formatter: (value) => formatCurrency(value),
                  },
                },
              }}
            />
          </Show>
        </section>
      </div>
    </div>
  );
}
