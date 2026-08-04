package models

import "time"

type ServiceType struct {
	ID             int
	Name           string
	Description    string
	BasePrice      float64
	PricePerWeight float64
	PricePerPiece  float64
	CreatedAt      time.Time
}
