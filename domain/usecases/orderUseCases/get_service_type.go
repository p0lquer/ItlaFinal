package orderUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
)

type GetServiceTypeUseCase struct {
	repo ports.ServiceTypeRepository
}

func NewGetServiceTypeUseCase(repo ports.ServiceTypeRepository) *GetServiceTypeUseCase {
	return &GetServiceTypeUseCase{repo: repo}
}

func (uc *GetServiceTypeUseCase) Execute() ([]*models.ServiceType, error) {
	return uc.repo.FindAll()
}
