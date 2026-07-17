package ports

import "ITLAFINAL/domain/models"

type ServiceTypeRepository interface {
	Create(serviceType *models.ServiceType) error
	FindAll() ([]*models.ServiceType, error)
	FindByName(name string) (*models.ServiceType, error)
	Update(serviceType *models.ServiceType) error
	Delete(name string) error
}
