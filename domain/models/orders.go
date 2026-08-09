package models

import "time"

type OrderStatus string

const (
	StatusReceived   OrderStatus = "recibida"
	StatusProcessing OrderStatus = "en_proceso"
	StatusReady      OrderStatus = "lista"
	StatusDelivered  OrderStatus = "entregada"
)

type Order struct {
	ID            string
	CustomerID    string
	ServiceType   string
	Status        OrderStatus
	PiecesCount   int
	Notes         string
	EstimatedTime time.Duration
	EstimatedCost float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ReadyAt       *time.Time // nil hasta que esté lista
	Weight        float64
	StartedAt     *time.Time // nil hasta que esté en proceso
}

func (o *Order) ElapsedTime() time.Duration {
	return time.Since(o.StartedBase())
}

func (o *Order) StartedBase() time.Time {
	if o.StartedAt != nil {
		return *o.StartedAt
	}
	return o.CreatedAt
}
