package models

import "time"

type Customer struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateCustomerRequest struct {
	FirstName string
	LastName  string
	Email     string
}

func (r *CreateCustomerRequest) IsValid() bool {
	if r.FirstName == "" || r.LastName == "" || r.Email == "" {
		return false
	}
	return true
}

func (r *CreateCustomerRequest) ToCustomer() *Customer {
	return &Customer{
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Email:     r.Email,
	}
}
