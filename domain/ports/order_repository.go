package ports

import (
	"ITLAFINAL/domain/models"

	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(order *models.Order) error
	FindByID(id string) (*models.Order, error)
	FindAll() ([]*models.Order, error)
	FindByUserID(userID uuid.UUID) ([]*models.Order, error)
	UpdateStatus(id string, status models.OrderStatus) error
	Delete(id string) error
}
