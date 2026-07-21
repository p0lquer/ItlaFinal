package orderUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"database/sql"
	"errors"
	"time"
)

type CreateServiceTypeUseCase struct {
	repo ports.ServiceTypeRepository
}

func NewCreateServiceTypeUseCase(repo ports.ServiceTypeRepository) *CreateServiceTypeUseCase {
	return &CreateServiceTypeUseCase{repo: repo}
}

func (uc *CreateServiceTypeUseCase) Execute(name, description string) (*models.ServiceType, error) {
	existing, err := uc.repo.FindByName(name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("service type with the same name already exists")
	}
	serviceType := &models.ServiceType{
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
	}

	if err := uc.repo.Create(serviceType); err != nil {
		return nil, err
	}

	return serviceType, nil
}
