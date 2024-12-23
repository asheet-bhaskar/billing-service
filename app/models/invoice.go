package models

import "time"

type Invoice struct {
	BillID      string
	Description string
	CustomerID  string
	CurrencyID  string
	Status      string
	TotalAmount float64
	PeriodStart time.Time
	PeriodEnd   time.Time
	LineItems   []LineItem
}

func CreateInvoice(bill *Bill, lineItems []*LineItem, currencyCode string) *Invoice {
	dereferencedLineItems := []LineItem{}
	for _, item := range lineItems {
		if !item.Removed {
			dereferencedLineItems = append(dereferencedLineItems, *item)
		}
	}

	return &Invoice{
		BillID:      bill.ID,
		Description: bill.Description,
		CustomerID:  bill.CustomerID,
		CurrencyID:  bill.CurrencyID,
		Status:      bill.Status,
		TotalAmount: bill.TotalAmount,
		PeriodStart: bill.PeriodStart,
		PeriodEnd:   bill.PeriodEnd,
		LineItems:   dereferencedLineItems,
	}
}
