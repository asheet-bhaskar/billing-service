package models

import (
	"time"
)

type Bill struct {
	ID          int64
	Description string
	CustomerID  int64
	CurrencyID  int64
	Status      string
	TotalAmount float64
	PeriodStart time.Time
	PeriodEnd   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type LineItem struct {
	ID          string
	BillID      int64
	Description string
	Amount      float64
	CreatedAt   time.Time
	Removed     bool
}

type BillRequest struct {
	Description  string
	CustomerID   int64
	CurrencyCode string
	PeriodStart  time.Time
	PeriodEnd    time.Time
}

func (r *BillRequest) IsValid() bool {
	if (r.Description == "" || r.CustomerID <= 0 || r.CurrencyCode == "" || r.PeriodStart == time.Time{} || r.PeriodEnd == time.Time{}) {
		return false
	}
	return true
}
