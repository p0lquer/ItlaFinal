package orderUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"errors"
	"time"

	"github.com/google/uuid"
)

type CreateServiceTypeUseCase struct {
	repo ports.ServiceTypeRepository
}

func NewCreateServiceTypeUseCase(repo ports.ServiceTypeRepository) *CreateServiceTypeUseCase {
	return &CreateServiceTypeUseCase{repo: repo}
}

func (uc *CreateServiceTypeUseCase) Execute(name, description string) (*models.ServiceType, error) {
	//verificar si el nombre ya existe
	existing, _ := uc.repo.FindByName(name)
	if existing != nil {
		return nil, errors.New("service type already exists")
	}
	serviceType := &models.ServiceType{
		ID:          uuid.NewString(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
	}

	if err := uc.repo.Create(serviceType); err != nil {
		return nil, err
	}

	return serviceType, nil
}
