package models

import "time"

// OrderStatusChange is the immutable audit entry produced whenever an order
// advances through the production workflow.
type OrderStatusChange struct {
	OrderID       string      `json:"order_id"`
	FromStatus    OrderStatus `json:"from_status"`
	ToStatus      OrderStatus `json:"to_status"`
	ChangedBy     string      `json:"changed_by"`
	ChangedByRole string      `json:"changed_by_role"`
	Description   string      `json:"description,omitempty"`
	ChangedAt     time.Time   `json:"changed_at"`
}

// OrderFilter describes the operator's searchable, paginated order list.
// From and To filter by order creation time, inclusive.
type OrderFilter struct {
	Search      string
	CustomerID  string
	Status      OrderStatus
	ServiceType string
	From        *time.Time
	To          *time.Time
	Page        int
	PageSize    int
}

type OrderPage struct {
	Orders   []*Order
	Total    int
	Page     int
	PageSize int
}
