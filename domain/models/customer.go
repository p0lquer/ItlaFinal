package models

import "time"

type Customer struct {
	ID        string
	Name      string
	Phone     string
	Email     string
	CreatedAt time.Time
}
