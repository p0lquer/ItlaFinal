package dto

import (
	"ITLAFINAL/domain/models"
	"time"
)

type CreateOrderRequest struct {
	CustomerID  string  `json:"customer_id"  binding:"omitempty,uuid4"`
	ServiceType string  `json:"service_type" binding:"required"`
	PiecesCount int     `json:"pieces_count" binding:"required,min=1,max=200"`
	Weight      float64 `json:"weight" binding:"required,gt=0,lte=100"`
	Notes       string  `json:"notes" binding:"max=1000"`
}

type UpdateStatusRequest struct {
	Status      models.OrderStatus `json:"status" binding:"required"`
	Description string             `json:"description" binding:"max=500"`
}

type OrderResponse struct {
	ID            string     `json:"id"`
	CustomerID    string     `json:"customer_id"`
	ServiceType   string     `json:"service_type"`
	PiecesCount   int        `json:"pieces_count"`
	Weight        float64    `json:"weight"`
	Notes         string     `json:"notes"`
	Status        string     `json:"status"`
	EstimatedTime float64    `json:"estimated_time_minutes"`
	EstimatedCost float64    `json:"estimated_cost"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	ReadyAt       *time.Time `json:"ready_at"`
	StartedAt     *time.Time `json:"started_at"`
}

type OrderDetailResponse struct {
	Order   OrderResponse               `json:"order"`
	History []*models.OrderStatusChange `json:"history"`
}

type PaginationResponse struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type OrderPageResponse struct {
	Data       []OrderResponse    `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

func NewOrderResponse(order *models.Order) OrderResponse {
	return OrderResponse{
		ID:            order.ID,
		CustomerID:    order.CustomerID,
		ServiceType:   order.ServiceType,
		PiecesCount:   order.PiecesCount,
		Weight:        order.Weight,
		Notes:         order.Notes,
		Status:        string(order.Status),
		EstimatedTime: order.EstimatedTime.Minutes(),
		EstimatedCost: order.EstimatedCost,
		CreatedAt:     order.CreatedAt,
		UpdatedAt:     order.UpdatedAt,
		ReadyAt:       order.ReadyAt,
		StartedAt:     order.StartedAt,
	}
}
