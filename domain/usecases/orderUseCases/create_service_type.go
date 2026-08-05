package orderUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"database/sql"
	"errors"
	"time"
)

var ErrServiceTypeAlreadyExists = errors.New("service type with the same name already exists")

type CreateServiceTypeUseCase struct {
	repo ports.ServiceTypeRepository
}

func NewCreateServiceTypeUseCase(repo ports.ServiceTypeRepository) *CreateServiceTypeUseCase {
	return &CreateServiceTypeUseCase{repo: repo}
}

func (uc *CreateServiceTypeUseCase) Execute(name, description string, basePrice, pricePerWeight, pricePerPiece float64) (*models.ServiceType, error) {
	existing, err := uc.repo.FindByName(name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrServiceTypeAlreadyExists
	}
	serviceType := &models.ServiceType{
		Name:           name,
		Description:    description,
		BasePrice:      basePrice,
		PricePerWeight: pricePerWeight,
		PricePerPiece:  pricePerPiece,
		CreatedAt:      time.Now(),
	}

	if err := uc.repo.Create(serviceType); err != nil {
		return nil, err
	}

	return serviceType, nil
}
