package usecases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
)

type GetAllOrdersUseCase struct {
	orderRepo ports.OrderRepository
}

func NewGetAllOrdersUseCase(orderRepo ports.OrderRepository) *GetAllOrdersUseCase {
	return &GetAllOrdersUseCase{orderRepo: orderRepo}
}

func (uc *GetAllOrdersUseCase) Execute() ([]*models.Order, error) {
	return uc.orderRepo.FindAll()
}
