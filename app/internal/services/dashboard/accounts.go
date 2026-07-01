package dashboard

import (
	"FinancialTracker/internal/models"
	"strings"
)

// Account bucket identifiers used for grouping and widgets.
const (
	BucketCash        = "cash"
	BucketSavings     = "savings"
	BucketCredit      = "credit"
	BucketLoans       = "loans"
	BucketInvestments = "investments"
	BucketOther       = "other"
)

// AllBuckets is the canonical order for grouping accounts.
var AllBuckets = []string{
	BucketCash,
	BucketSavings,
	BucketCredit,
	BucketLoans,
	BucketInvestments,
	BucketOther,
}

// ClassifyAccount maps a Plaid account into a display bucket using type and subtype.
func ClassifyAccount(acc models.Account) string {
	t := strings.ToLower(strings.TrimSpace(acc.Type))
	st := strings.ToLower(strings.TrimSpace(acc.Subtype))

	switch t {
	case "credit":
		return BucketCredit
	case "loan":
		return BucketLoans
	case "investment", "brokerage":
		return BucketInvestments
	case "depository":
		switch st {
		case "savings", "money market", "cd", "hsa", "cash management":
			return BucketSavings
		default:
			return BucketCash
		}
	default:
		if t == "brokerage" {
			return BucketInvestments
		}
		return BucketOther
	}
}

// ExcludedFromTotals reports whether an account should be omitted from summary math.
func ExcludedFromTotals(acc models.Account) bool {
	return acc.IsHidden || acc.Status == "disconnected"
}

// IsLiabilityBucket reports whether the bucket represents debt.
func IsLiabilityBucket(bucket string) bool {
	return bucket == BucketCredit || bucket == BucketLoans
}

// DisplayBalance returns the balance shown in the UI for an account.
// Assets use current balance with available as fallback; liabilities show amount owed as a positive value.
func DisplayBalance(acc models.Account, bucket string) float64 {
	if IsLiabilityBucket(bucket) {
		bal := acc.Balance
		if bal < 0 {
			return -bal
		}
		return bal
	}
	bal := acc.Balance
	if bal == 0 {
		bal = acc.AvailableBalance
	}
	return bal
}

// GroupAccounts organizes accounts by bucket while preserving input order within each bucket.
func GroupAccounts(accounts []models.Account) map[string][]models.Account {
	groups := make(map[string][]models.Account, len(AllBuckets))
	for _, b := range AllBuckets {
		groups[b] = []models.Account{}
	}
	for _, acc := range accounts {
		bucket := ClassifyAccount(acc)
		groups[bucket] = append(groups[bucket], acc)
	}
	return groups
}

// BuildSummary computes typed totals from accounts. Hidden and disconnected accounts are excluded from totals.
func BuildSummary(accounts []models.Account) models.DashboardSummary {
	groups := GroupAccounts(accounts)
	summary := models.DashboardSummary{}

	for _, acc := range accounts {
		if !ExcludedFromTotals(acc) {
			summary.AccountCount++
		}
	}

	for bucket, list := range groups {
		for _, acc := range list {
			if ExcludedFromTotals(acc) {
				continue
			}
			amount := DisplayBalance(acc, bucket)
			switch bucket {
			case BucketCash:
				summary.Cash += amount
			case BucketSavings:
				summary.Savings += amount
			case BucketCredit:
				summary.CreditDebt += amount
			case BucketLoans:
				summary.LoanDebt += amount
				if acc.MonthlyPayment > 0 {
					summary.LoanMonthlyPayments += acc.MonthlyPayment
				}
			case BucketInvestments:
				summary.Investments += amount
			}
		}
	}

	summary.NetWorth = summary.Cash + summary.Savings + summary.Investments - summary.CreditDebt - summary.LoanDebt
	return summary
}

// BuildDashboardData assembles the full dashboard view model.
func BuildDashboardData(accounts []models.Account, transactions []models.Transaction, layout models.DashboardLayout) *models.DashboardData {
	return &models.DashboardData{
		Summary:      BuildSummary(accounts),
		Groups:       GroupAccounts(accounts),
		Transactions: transactions,
		Layout:       NormalizeLayout(layout),
	}
}
