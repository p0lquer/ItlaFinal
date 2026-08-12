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

func (r *paymentRepositoryPG) Create(p *models.Payment) error {
	_, err := r.db.Exec(`INSERT INTO payments (id, order_id, amount, currency, method, status, receipt_number, paid_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		p.ID, p.OrderID, p.Amount, p.Currency, p.Method, p.Status, p.ReceiptNo, p.PaidAt)
	return err
}
