package models

import (
	"time"
)

type Bill struct {
	ID          string
	Description string
	CustomerID  string
	CurrencyID  string
	Status      string
	TotalAmount float64
	PeriodStart time.Time
	PeriodEnd   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type LineItem struct {
	ID          string
	BillID      string
	Description string
	Amount      float64
	CreatedAt   time.Time
	Removed     bool
}

type BillRequest struct {
	Description  string
	CustomerID   string
	CurrencyCode string
	PeriodStart  time.Time
	PeriodEnd    time.Time
}

func (r *BillRequest) IsValid() bool {
	if (r.Description == "" || r.CustomerID == "" || r.CurrencyCode == "" || r.PeriodStart == time.Time{} ||
		r.PeriodEnd == time.Time{} || r.PeriodStart.After(r.PeriodEnd)) {
		return false
	}
	return true
}
