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

func (r *Customer) IsValid() bool {
	if r.FirstName == "" || r.LastName == "" || r.Email == "" {
		return false
	}
	return true
}
