package repository

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"database/sql"
)

type paymentRepositoryPG struct{ db *sql.DB }

var _ ports.PaymentRepository = (*paymentRepositoryPG)(nil)

func NewPaymentRepository(db *sql.DB) ports.PaymentRepository { return &paymentRepositoryPG{db: db} }

func (r *paymentRepositoryPG) FindByOrderID(orderID string) (*models.Payment, error) {
	p := &models.Payment{}
	err := r.db.QueryRow(`SELECT id, order_id, amount, currency, method, status, receipt_number, paid_at FROM payments WHERE order_id=$1`, orderID).
		Scan(&p.ID, &p.OrderID, &p.Amount, &p.Currency, &p.Method, &p.Status, &p.ReceiptNo, &p.PaidAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

// FindByCustomerID provides the customer's durable payment history without
// exposing payments that belong to another account.
func (r *paymentRepositoryPG) FindByCustomerID(customerID string) ([]*models.Payment, error) {
	rows, err := r.db.Query(`SELECT p.id, p.order_id, p.amount, p.currency, p.method, p.status, p.receipt_number, p.paid_at
		FROM payments p JOIN orders o ON o.id = p.order_id
		WHERE o.customer_id=$1 ORDER BY p.paid_at DESC`, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	payments := make([]*models.Payment, 0)
	for rows.Next() {
		payment := &models.Payment{}
		if err := rows.Scan(&payment.ID, &payment.OrderID, &payment.Amount, &payment.Currency, &payment.Method, &payment.Status, &payment.ReceiptNo, &payment.PaidAt); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, rows.Err()
}

func (r *paymentRepositoryPG) Create(p *models.Payment) error {
	_, err := r.db.Exec(`INSERT INTO payments (id, order_id, amount, currency, method, status, receipt_number, paid_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		p.ID, p.OrderID, p.Amount, p.Currency, p.Method, p.Status, p.ReceiptNo, p.PaidAt)
	return err
}
