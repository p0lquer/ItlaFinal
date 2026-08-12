package ports

import "ITLAFINAL/domain/models"

type PaymentRepository interface {
	FindByOrderID(orderID string) (*models.Payment, error)
	Create(payment *models.Payment) error
}
