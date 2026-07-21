package models

import "time"

type ServiceType struct {
	ID          int
	Name        string
	Description string
	CreatedAt   time.Time
}
