import type { ApexOptions } from "apexcharts";
import { For, Show } from "solid-js";
import ApexChart from "~/components/dashboard/ApexChart";
import {
  ACCOUNT_BUCKET_COPY,
  WIDGET_IDS,
  widgetLabel,
} from "~/lib/dashboard-widgets";
import {
  formatAmount,
  formatCurrency,
  formatCurrencyCompact,
  formatDate,
  donutCenterFontSize,
  formatMonthLabel,
  formatMonths,
  formatNetWorth,
} from "~/lib/format";
import { tagColorHex } from "~/lib/tag-colors";
import { readCssVar } from "~/lib/themes";
import type { DashboardPayload } from "~/lib/types";
import styles from "~/styles/dashboard.module.css";

type WidgetBodyProps = {
  data: DashboardPayload;
  widgetId: string;
};

function chartBaseOptions(): ApexOptions {
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
      theme: document.documentElement.dataset.theme === "dark" ? "dark" : "light",
    },
  };
}

function AccountBucketWidget(props: { data: DashboardPayload; bucket: string }) {
  const copy = () => ACCOUNT_BUCKET_COPY[props.bucket] ?? ACCOUNT_BUCKET_COPY.other;
  const accounts = () => props.data.groups?.[props.bucket] ?? [];
  const total = () => {
    switch (props.bucket) {
      case "cash":
        return props.data.summary.cash;
      case "savings":
        return props.data.summary.savings;
      case "credit":
        return props.data.summary.credit_debt;
      case "loans":
        return props.data.summary.loan_debt;
      case "investments":
        return props.data.summary.investments;
      default:
        return accounts().reduce((sum, account) => sum + account.balance, 0);
    }
  };

  return (
    <article class={styles.widget}>
      <header class={styles.widgetHeader}>
        <div>
          <h2 class={styles.widgetTitle}>{copy().title}</h2>
          <p class={styles.widgetSubtitle}>{copy().subtitle}</p>
        </div>
        <div class={styles.widgetAside}>
          <span class={styles.widgetAsideLabel}>
            {copy().liability ? "Amount owed" : "Subtotal"}
          </span>
          <strong class={copy().liability ? styles.debtValue : styles.valueStrong}>
            {formatCurrency(total())}
          </strong>
          <Show when={copy().showMonthly && props.data.summary.loan_monthly_payments > 0}>
            <span class={styles.widgetAsideMeta}>
              {formatCurrency(props.data.summary.loan_monthly_payments)}/mo total
            </span>
          </Show>
        </div>
      </header>

      <div class={styles.scrollList}>
        <Show
          when={accounts().length > 0}
          fallback={<p class={styles.emptyState}>No accounts in this category.</p>}
        >
          <For each={accounts()}>
            {(account) => (
              <div
                class={styles.accountRow}
                classList={{
                  [styles.dimmed]: account.is_hidden || account.status === "disconnected",
                }}
              >
                <div>
                  <div class={styles.accountTitleRow}>
                    <p
                      class={styles.accountName}
                      classList={{ [styles.strike]: account.status === "disconnected" }}
                    >
                      {account.name}
                    </p>
                    <Show when={account.status === "disconnected"}>
                      <span class={styles.badgeDanger}>Disconnected</span>
                    </Show>
                    <Show when={account.is_hidden}>
                      <span class={styles.badgeMuted}>Hidden</span>
                    </Show>
                  </div>
                  <p class={styles.accountMeta}>
                    {account.subtype} • ****{account.mask}
                  </p>
                </div>
                <div class={styles.accountAmounts}>
                  <strong class={copy().liability ? styles.debtValue : styles.valueStrong}>
                    {formatCurrency(account.balance)}
                  </strong>
                  <Show when={copy().showMonthly && account.monthly_payment > 0}>
                    <span class={styles.widgetAsideMeta}>
                      {formatCurrency(account.monthly_payment)}/mo
                    </span>
                  </Show>
                </div>
              </div>
            )}
          </For>
        </Show>
      </div>
    </article>
  );
}

function TagDonutWidget(props: {
  title: string;
  subtitle: string;
  slices: DashboardPayload["spending_by_tag"];
}) {
  const slices = () => props.slices ?? [];
  const labels = () => slices().map((slice) => slice.tag_name);
  const values = () => slices().map((slice) => slice.total);
  const colors = () => slices().map((slice) => tagColorHex(slice.color) || readCssVar("--accent"));
  const total = () => values().reduce((sum, value) => sum + value, 0);
  const centerTotal = () => formatCurrencyCompact(total());

  return (
    <article class={styles.widget}>
      <header class={styles.widgetHeader}>
        <div>
          <h2 class={styles.widgetTitle}>{props.title}</h2>
          <p class={styles.widgetSubtitle}>{props.subtitle}</p>
        </div>
        <div class={styles.widgetAside}>
          <span class={styles.widgetAsideLabel}>Total</span>
          <strong class={styles.valueStrong}>{formatCurrency(total())}</strong>
        </div>
      </header>

      <Show
        when={slices().length > 0}
        fallback={<p class={styles.emptyState}>No tagged activity this month.</p>}
      >
        <ApexChart
          class={styles.chartWrap}
          height={240}
          series={values()}
          options={{
            ...chartBaseOptions(),
            chart: { ...chartBaseOptions().chart, type: "donut" },
            labels: labels(),
            colors: colors(),
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
    </article>
  );
}

export default function WidgetBody(props: WidgetBodyProps) {
  const data = () => props.data;
  const spendingTrend = () => data().spending_trend ?? [];
  const transactions = () => data().transactions ?? [];

  switch (props.widgetId) {
    case WIDGET_IDS.netWorth:
      return (
        <article class={styles.widget}>
          <header class={styles.widgetHeader}>
            <div>
              <h2 class={styles.widgetTitle}>Net Worth</h2>
              <p class={styles.widgetSubtitle}>Assets minus liabilities across linked accounts</p>
            </div>
          </header>
          <p
            class={styles.netWorthValue}
            classList={{ [styles.debtValue]: data().summary.net_worth < 0 }}
          >
            {formatNetWorth(data().summary.net_worth)}
          </p>
          <div class={styles.metricGrid}>
            <div class={`${styles.metricTile} ${styles.metricTileCash}`}>
              <span>Cash</span>
              <strong>{formatCurrency(data().summary.cash)}</strong>
            </div>
            <div class={`${styles.metricTile} ${styles.metricTileSavings}`}>
              <span>Savings</span>
              <strong>{formatCurrency(data().summary.savings)}</strong>
            </div>
            <div class={`${styles.metricTile} ${styles.metricTileInvestments}`}>
              <span>Investments</span>
              <strong>{formatCurrency(data().summary.investments)}</strong>
            </div>
            <div class={`${styles.metricTile} ${styles.metricTileCredit}`}>
              <span>Credit debt</span>
              <strong>{formatCurrency(data().summary.credit_debt)}</strong>
            </div>
            <div class={`${styles.metricTile} ${styles.metricTileLoans}`}>
              <span>Loans</span>
              <strong>{formatCurrency(data().summary.loan_debt)}</strong>
              <Show when={data().summary.loan_monthly_payments > 0}>
                <small>{formatCurrency(data().summary.loan_monthly_payments)}/mo</small>
              </Show>
            </div>
            <div class={`${styles.metricTile} ${styles.metricTileAccounts}`}>
              <span>Accounts</span>
              <strong>{data().summary.account_count}</strong>
            </div>
          </div>
        </article>
      );

    case WIDGET_IDS.monthCashflow:
      return (
        <article class={styles.widget}>
          <header class={styles.widgetHeader}>
            <div>
              <h2 class={styles.widgetTitle}>This Month</h2>
              <p class={styles.widgetSubtitle}>
                {new Date().toLocaleDateString("en-US", { month: "long", year: "numeric" })}
              </p>
            </div>
          </header>
          <div class={styles.cashflowGrid}>
            <div class={styles.cashflowTile}>
              <span>Spend</span>
              <strong>{formatCurrency(data().month_cashflow.spend)}</strong>
            </div>
            <div class={`${styles.cashflowTile} ${styles.cashflowIncome}`}>
              <span>Income</span>
              <strong>{formatCurrency(data().month_cashflow.income)}</strong>
            </div>
          </div>
        </article>
      );

    case WIDGET_IDS.spendingByTag:
      return (
        <TagDonutWidget
          title="Spending by Tag"
          subtitle="This month outflows by tag"
          slices={data().spending_by_tag ?? []}
        />
      );

    case WIDGET_IDS.incomeByTag:
      return (
        <TagDonutWidget
          title="Income by Tag"
          subtitle="This month inflows by tag"
          slices={data().income_by_tag ?? []}
        />
      );

    case WIDGET_IDS.cashAccounts:
      return <AccountBucketWidget data={data()} bucket="cash" />;
    case WIDGET_IDS.savingsAccounts:
      return <AccountBucketWidget data={data()} bucket="savings" />;
    case WIDGET_IDS.creditAccounts:
      return <AccountBucketWidget data={data()} bucket="credit" />;
    case WIDGET_IDS.loanAccounts:
      return <AccountBucketWidget data={data()} bucket="loans" />;
    case WIDGET_IDS.investmentAccounts:
      return <AccountBucketWidget data={data()} bucket="investments" />;

    case WIDGET_IDS.spendingTrend:
      return (
        <article class={styles.widget}>
          <header class={styles.widgetHeader}>
            <div>
              <h2 class={styles.widgetTitle}>Spending Trend</h2>
              <p class={styles.widgetSubtitle}>Monthly outflows over the last six months</p>
            </div>
          </header>

          <Show
            when={spendingTrend().length > 0}
            fallback={<p class={styles.emptyState}>Not enough transaction history yet.</p>}
          >
            <div class={styles.metricGrid}>
              <div class={`${styles.metricTile} ${styles.metricTileSpend}`}>
                <span>Avg monthly spend</span>
                <strong>{formatCurrency(data().summary.avg_monthly_spend)}</strong>
              </div>
              <Show when={data().summary.months_to_zero > 0}>
                <div class={`${styles.metricTile} ${styles.metricTileRunway}`}>
                  <span>Cash runway (est.)</span>
                  <strong>{formatMonths(data().summary.months_to_zero)} mo</strong>
                </div>
              </Show>
            </div>

            <ApexChart
              class={styles.chartWrap}
              height={260}
              series={[
                {
                  name: "Spend",
                  data: spendingTrend().map((point) => point.total),
                },
              ]}
              options={{
                ...chartBaseOptions(),
                chart: { ...chartBaseOptions().chart, type: "bar" },
                colors: [readCssVar("--accent")],
                xaxis: {
                  categories: spendingTrend().map((point) => formatMonthLabel(point.month)),
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
                plotOptions: {
                  bar: {
                    borderRadius: 4,
                    columnWidth: "60%",
                  },
                },
              }}
            />
          </Show>
        </article>
      );

    case WIDGET_IDS.recentTransactions:
      return (
        <article class={styles.widget}>
          <header class={styles.widgetHeader}>
            <div>
              <h2 class={styles.widgetTitle}>Recent Transactions</h2>
              <p class={styles.widgetSubtitle}>Latest synced activity from your accounts</p>
            </div>
            <a class={styles.widgetLink} href="/transactions">
              View all
            </a>
          </header>

          <div class={styles.scrollList}>
            <Show
              when={transactions().length > 0}
              fallback={<p class={styles.emptyState}>No transactions found.</p>}
            >
              <For each={transactions()}>
                {(transaction) => (
                  <div class={styles.transactionRow}>
                    <div>
                      <p class={styles.transactionName}>{transaction.name}</p>
                      <p class={styles.transactionMeta}>
                        {formatDate(transaction.date)}
                        {transaction.pending ? " • Pending" : ""}
                      </p>
                      <Show when={transaction.tags && transaction.tags.length > 0}>
                        <div class={styles.tagRow}>
                          <For each={transaction.tags!}>
                            {(tag) => (
                              <span
                                class={styles.tagChip}
                                style={{ "--tag-color": tagColorHex(tag.color) || readCssVar("--accent") }}
                              >
                                {tag.name}
                              </span>
                            )}
                          </For>
                        </div>
                      </Show>
                    </div>
                    <strong
                      class={styles.transactionAmount}
                      classList={{ [styles.incomeAmount]: -transaction.amount >= 0 }}
                    >
                      {formatAmount(transaction.amount)}
                    </strong>
                  </div>
                )}
              </For>
            </Show>
          </div>
        </article>
      );

    case WIDGET_IDS.quickActions:
      return (
        <article class={styles.widget}>
          <header class={styles.widgetHeader}>
            <div>
              <h2 class={styles.widgetTitle}>Quick Actions</h2>
              <p class={styles.widgetSubtitle}>Common tasks while you review your dashboard</p>
            </div>
          </header>
          <div class={styles.actionStack}>
            <a class={styles.primaryAction} href="/manage">
              Manage connections
            </a>
            <a class={styles.secondaryAction} href="/settings">
              Settings
            </a>
          </div>
        </article>
      );

    default:
      return (
        <article class={styles.widget}>
          <p class={styles.emptyState}>Unknown widget: {props.widgetId}</p>
        </article>
      );
  }
}

export { widgetLabel };
