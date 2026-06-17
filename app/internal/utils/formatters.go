package utils

import (
    "fmt"
    "time"
    "golang.org/x/text/language"
    "golang.org/x/text/message"
)

// printer is a localized printer for number formatting
var printer = message.NewPrinter(language.English)

// FormatDate converts a unix timestamp to a formatted date string
func FormatDate(unix int64) string {
    return time.Unix(unix, 0).Format("Jan 02, 2006")
}

// FormatAmount formats a Plaid transaction amount for display
func FormatAmount(amount float64) string {
    // Plaid returns positive for debits and negative for credits
    // We want to show credits (income) as positive/green and debits as negative/black
    displayAmount := -amount
    if displayAmount >= 0 {
        return printer.Sprintf("+$%.2f", displayAmount)
    }
    return printer.Sprintf("-$%.2f", -displayAmount)
}

// FormatCurrency formats a positive amount as currency
func FormatCurrency(amount float64) string {
    return printer.Sprintf("$%.2f", amount)
}

// FormatCurrencyAbs formats the absolute value as currency (no minus sign).
func FormatCurrencyAbs(amount float64) string {
    if amount < 0 {
        amount = -amount
    }
    return FormatCurrency(amount)
}

// FormatNetWorth formats net worth with a leading minus when negative.
func FormatNetWorth(amount float64) string {
    if amount < 0 {
        return printer.Sprintf("-$%.2f", -amount)
    }
    return FormatCurrency(amount)
}

// FormatNumber formats a float with thousands separators and no currency symbol
func FormatNumber(num float64, precision int) string {
    format := fmt.Sprintf("%%.%df", precision)
    return printer.Sprintf(format, num)
}

// FormatInt formats an integer with thousands separators
func FormatInt(num int) string {
    return printer.Sprintf("%d", num)
}
