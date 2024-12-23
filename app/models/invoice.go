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
