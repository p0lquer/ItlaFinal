package dto

import (
	"ITLAFINAL/domain/models"
	"time"
)

type CreateOrderRequest struct {
	CustomerID  string `json:"customer_id"  binding:"required"`
	ServiceType string `json:"service_type" binding:"required"`
	PiecesCount int    `json:"pieces_count" binding:"required,min=1"`
	Notes       string `json:"notes"`
}

type UpdateStatusRequest struct {
	Status models.OrderStatus `json:"status" binding:"required"`
}

type OrderResponse struct {
	ID            string    `json:"id"`
	Status        string    `json:"status"`
	EstimatedTime float64   `json:"estimated_time_minutes"`
	CreatedAt     time.Time `json:"created_at"`
}
