package dashboard

import (
	"FinancialTracker/internal/models"
	"testing"
)

func TestClassifyAccount(t *testing.T) {
	tests := []struct {
		acc    models.Account
		bucket string
	}{
		{models.Account{Type: "depository", Subtype: "checking"}, BucketCash},
		{models.Account{Type: "depository", Subtype: "savings"}, BucketSavings},
		{models.Account{Type: "credit", Subtype: "credit card"}, BucketCredit},
		{models.Account{Type: "loan", Subtype: "mortgage"}, BucketLoans},
		{models.Account{Type: "investment", Subtype: "401k"}, BucketInvestments},
	}
	for _, tc := range tests {
		if got := ClassifyAccount(tc.acc); got != tc.bucket {
			t.Errorf("ClassifyAccount(%s/%s) = %q, want %q", tc.acc.Type, tc.acc.Subtype, got, tc.bucket)
		}
	}
}

func TestNetWorthExcludesLiabilitiesFromAssets(t *testing.T) {
	accounts := []models.Account{
		{Type: "depository", Subtype: "checking", Balance: 1000, Status: "active"},
		{Type: "credit", Subtype: "credit card", Balance: 200, Status: "active"},
	}
	summary := BuildSummary(accounts)
	if summary.NetWorth != 800 {
		t.Errorf("NetWorth = %v, want 800", summary.NetWorth)
	}
	if summary.CreditDebt != 200 {
		t.Errorf("CreditDebt = %v, want 200", summary.CreditDebt)
	}
}

func TestCreditDebtIgnoresAvailableCredit(t *testing.T) {
	accounts := []models.Account{
		{Type: "credit", Subtype: "credit card", Balance: 0, AvailableBalance: 10000, Status: "active"},
		{Type: "credit", Subtype: "credit card", Balance: 500, AvailableBalance: 9500, Status: "active"},
	}
	summary := BuildSummary(accounts)
	if summary.CreditDebt != 500 {
		t.Errorf("CreditDebt = %v, want 500 (available credit must not count as debt)", summary.CreditDebt)
	}
}

func TestHiddenExcludedFromTotals(t *testing.T) {
	accounts := []models.Account{
		{Type: "depository", Subtype: "checking", Balance: 1000, Status: "active"},
		{Type: "depository", Subtype: "savings", Balance: 500, Status: "active", IsHidden: true},
	}
	summary := BuildSummary(accounts)
	if summary.Cash != 1000 {
		t.Errorf("Cash = %v, want 1000 (hidden excluded)", summary.Cash)
	}
	if summary.AccountCount != 1 {
		t.Errorf("AccountCount = %v, want 1", summary.AccountCount)
	}
}
