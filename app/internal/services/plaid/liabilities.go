package plaid

import (
	"context"

	"github.com/labstack/gommon/log"
	"github.com/plaid/plaid-go/v42/plaid"
)

// SyncLiabilities fetches loan payment amounts from Plaid and stores them on matching accounts.
func (p *PlaidService) SyncLiabilities(ctx *context.Context, accessToken string) {
	req := plaid.NewLiabilitiesGetRequest(accessToken)
	resp, _, err := p.client.PlaidApi.LiabilitiesGet(*ctx).LiabilitiesGetRequest(*req).Execute()
	if err != nil {
		log.Warnf("Liabilities sync skipped (product may be unavailable): %v", err)
		return
	}

	liabilities := resp.GetLiabilities()
	for _, mortgage := range liabilities.GetMortgage() {
		if amount, ok := mortgage.GetNextMonthlyPaymentOk(); ok && amount != nil {
			if err := p.store.UpdateAccountMonthlyPayment(mortgage.AccountId, *amount); err != nil {
				log.Errorf("Failed to update monthly payment for %s: %v", mortgage.AccountId, err)
			}
		}
	}
	for _, student := range liabilities.GetStudent() {
		accountID, ok := student.GetAccountIdOk()
		if !ok || accountID == nil {
			continue
		}
		if amount, ok := student.GetMinimumPaymentAmountOk(); ok && amount != nil {
			if err := p.store.UpdateAccountMonthlyPayment(*accountID, *amount); err != nil {
				log.Errorf("Failed to update monthly payment for %s: %v", *accountID, err)
			}
		}
	}
}
