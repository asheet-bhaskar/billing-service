package customer

import "time"

type Customer struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
