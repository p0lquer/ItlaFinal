package ports

import (
	"ITLAFINAL/domain/models"

	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(order *models.Order) error
	FindByID(id string) (*models.Order, error)
	FindAll() ([]*models.Order, error)
	List(filter models.OrderFilter) (*models.OrderPage, error)
	FindByUserID(userID uuid.UUID) ([]*models.Order, error)
	FindByCustomerID(customerID string) ([]*models.Order, error)
	FindStatusHistory(orderID string) ([]*models.OrderStatusChange, error)
	UpdateStatus(id string, status models.OrderStatus) error
	Delete(id string) error
	Transition(orderID string, target models.OrderStatus, change models.OrderStatusChange) (bool, error)
}
