const currency = new Intl.NumberFormat("en-US", {
  style: "currency",
  currency: "USD",
  minimumFractionDigits: 2,
});

// formatCurrency formats a positive amount as USD.
export function formatCurrency(amount: number): string {
  return currency.format(amount);
}

// formatCurrencyCompact abbreviates large USD amounts for tight chart labels.
export function formatCurrencyCompact(amount: number): string {
  const sign = amount < 0 ? "-" : "";
  const abs = Math.abs(amount);

  if (abs >= 1_000_000) {
    const scaled = abs / 1_000_000;
    const digits = scaled >= 100 ? 0 : scaled >= 10 ? 1 : 2;
    return `${sign}$${scaled.toFixed(digits)}M`;
  }
  if (abs >= 10_000) {
    return `${sign}$${(abs / 1_000).toFixed(1)}K`;
  }
  if (abs >= 1_000) {
    return `${sign}$${(abs / 1_000).toFixed(2)}K`;
  }

  return formatCurrency(amount);
}

// donutCenterFontSize scales donut center labels so long totals stay inside the ring.
export function donutCenterFontSize(text: string): string {
  const len = text.length;
  if (len <= 8) return "1.35rem";
  if (len <= 10) return "1.15rem";
  if (len <= 12) return "1rem";
  return "0.88rem";
}

// formatNetWorth formats net worth with a leading minus when negative.
export function formatNetWorth(amount: number): string {
  if (amount < 0) {
    return `-${currency.format(Math.abs(amount))}`;
  }
  return currency.format(amount);
}

// formatAmount formats a Plaid transaction amount for display.
export function formatAmount(amount: number): string {
  const displayAmount = -amount;
  if (displayAmount >= 0) {
    return `+${currency.format(displayAmount)}`;
  }
  return `-${currency.format(Math.abs(displayAmount))}`;
}

// formatDate converts a unix timestamp to a readable date.
export function formatDate(unix: number): string {
  return new Date(unix * 1000).toLocaleDateString("en-US", {
    month: "short",
    day: "2-digit",
    year: "numeric",
    timeZone: "UTC",
  });
}

// formatMonthLabel turns YYYY-MM into a short month label.
export function formatMonthLabel(month: string): string {
  const [year, monthNum] = month.split("-").map(Number);
  return new Date(year, monthNum - 1, 1).toLocaleDateString("en-US", {
    month: "short",
  });
}

// formatMonths formats runway months for display.
export function formatMonths(value: number): string {
  if (value >= 99) return "99+";
  if (value < 10) return value.toFixed(1);
  return String(Math.round(value));
}
