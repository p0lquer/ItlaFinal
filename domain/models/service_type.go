package models

import "time"

type ServiceType struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
}
