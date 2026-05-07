package ports

import "ITLAFINAL/domain/models"

type OrderRepository interface {
	Create(order *models.Order) error
	FindByID(id string) (*models.Order, error)
	FindAll() ([]*models.Order, error)
	FindByCustomerID(customerID string) ([]*models.Order, error)
	UpdateStatus(id string, status models.OrderStatus) error
	Delete(id string) error
}
