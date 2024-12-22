package models

import (
	"time"
)

type Currency struct {
	ID        int64
	Code      string
	Name      string
	Symbol    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *Currency) IsValid() bool {
	if r.Code == "" || r.Name == "" || r.Symbol == "" {
		return false
	}
	return true
}
