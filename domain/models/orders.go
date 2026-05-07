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
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ReadyAt       *time.Time // nil hasta que esté lista
}
