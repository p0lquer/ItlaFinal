package models

import "time"

// Payment is deliberately provider-agnostic: this demo records a completed
// local payment without ever receiving or storing card details.
type Payment struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"order_id"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	Method    string    `json:"method"`
	Status    string    `json:"status"`
	ReceiptNo string    `json:"receipt_number"`
	PaidAt    time.Time `json:"paid_at"`
}

type PaymentSummary struct {
	Order   *Order   `json:"order"`
	Payment *Payment `json:"payment,omitempty"`
	CanPay  bool     `json:"can_pay"`
}
