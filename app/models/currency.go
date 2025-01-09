package models

import (
	"time"
)

type Currency struct {
	ID        string
	Code      string
	Name      string
	Symbol    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateCurrencyRequest struct {
	Code   string
	Name   string
	Symbol string
}

func (r *CreateCurrencyRequest) IsValid() bool {
	if r.Code == "" || r.Name == "" || r.Symbol == "" {
		return false
	}
	return true
}

func (r *CreateCurrencyRequest) ToCurrency() *Currency {
	return &Currency{
		Code:   r.Code,
		Name:   r.Name,
		Symbol: r.Symbol,
	}
}
