import type { DashboardWidget } from "./types";

export type WidgetMeta = {
  label: string;
  span: 1 | 2;
  minRows: number;
  maxRows: number;
};

export const WIDGET_IDS = {
  netWorth: "net_worth",
  cashAccounts: "cash_accounts",
  savingsAccounts: "savings_accounts",
  creditAccounts: "credit_accounts",
  loanAccounts: "loan_accounts",
  investmentAccounts: "investment_accounts",
  recentTransactions: "recent_transactions",
  quickActions: "quick_actions",
  spendingTrend: "spending_trend",
  monthCashflow: "month_cashflow",
  spendingByTag: "spending_by_tag",
  incomeByTag: "income_by_tag",
} as const;

export const WIDGET_META: Record<string, WidgetMeta> = {
  [WIDGET_IDS.netWorth]: { label: "Net Worth", span: 2, minRows: 2, maxRows: 3 },
  [WIDGET_IDS.monthCashflow]: { label: "This Month", span: 2, minRows: 2, maxRows: 2 },
  [WIDGET_IDS.spendingByTag]: { label: "Spending by Tag", span: 1, minRows: 2, maxRows: 3 },
  [WIDGET_IDS.incomeByTag]: { label: "Income by Tag", span: 1, minRows: 2, maxRows: 3 },
  [WIDGET_IDS.cashAccounts]: { label: "Cash & Checking", span: 1, minRows: 2, maxRows: 4 },
  [WIDGET_IDS.savingsAccounts]: { label: "Savings", span: 1, minRows: 2, maxRows: 4 },
  [WIDGET_IDS.creditAccounts]: { label: "Credit Cards", span: 1, minRows: 2, maxRows: 4 },
  [WIDGET_IDS.loanAccounts]: { label: "Loans", span: 1, minRows: 2, maxRows: 4 },
  [WIDGET_IDS.investmentAccounts]: { label: "Investments", span: 1, minRows: 2, maxRows: 4 },
  [WIDGET_IDS.quickActions]: { label: "Quick Actions", span: 1, minRows: 2, maxRows: 2 },
  [WIDGET_IDS.spendingTrend]: { label: "Spending Trend", span: 2, minRows: 2, maxRows: 3 },
  [WIDGET_IDS.recentTransactions]: { label: "Recent Transactions", span: 2, minRows: 2, maxRows: 4 },
};

export const ACCOUNT_BUCKET_COPY: Record<
  string,
  { title: string; subtitle: string; liability?: boolean; showMonthly?: boolean }
> = {
  cash: {
    title: "Cash & Checking",
    subtitle: "Checking, prepaid, and everyday accounts",
  },
  savings: {
    title: "Savings",
    subtitle: "Savings, CDs, HSA, and money market",
  },
  credit: {
    title: "Credit Cards",
    subtitle: "Outstanding credit card balances",
    liability: true,
  },
  loans: {
    title: "Loans",
    subtitle: "Mortgage, auto, student, and other loans",
    liability: true,
    showMonthly: true,
  },
  investments: {
    title: "Investments",
    subtitle: "Brokerage, retirement, and investment accounts",
  },
  other: {
    title: "Other Accounts",
    subtitle: "Accounts outside standard categories",
  },
};

export function widgetLabel(id: string): string {
  return WIDGET_META[id]?.label ?? id;
}

export function widgetMeta(id: string): WidgetMeta | undefined {
  return WIDGET_META[id];
}

export function widgetsForRender(
  widgets: DashboardWidget[],
  editMode: boolean
) {
  const sorted = [...widgets].sort((a, b) => a.order - b.order);
  return editMode ? sorted : sorted.filter((widget) => widget.visible);
}
