package models

import "time"

type Prediction struct {
	ID          string
	ServiceType string
	PiecesCount int
	Estimated   time.Duration  // tiempo estimado calculado
	Actual      *time.Duration // tiempo real (se llena al finalizar)
	CreatedAt   time.Time
}
