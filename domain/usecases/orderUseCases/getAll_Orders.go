package orderUseCases

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

func (uc *GetAllOrdersUseCase) ExecutePage(filter models.OrderFilter) (*models.OrderPage, error) {
	return uc.orderRepo.List(filter)
}
