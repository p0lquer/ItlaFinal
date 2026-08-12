package orderUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
)

type GetOrderUseCase struct {
	orderRepo ports.OrderRepository
}

func NewGetOrderUseCase(orderRepo ports.OrderRepository) *GetOrderUseCase {
	return &GetOrderUseCase{orderRepo: orderRepo}
}

func (uc *GetOrderUseCase) Execute(orderID string) (*models.Order, error) {
	return uc.orderRepo.FindByID(orderID)
}
