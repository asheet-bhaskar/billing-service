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
