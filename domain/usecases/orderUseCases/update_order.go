package orderUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
)

type UpdateOrderStatusUseCase struct {
	orderRepo ports.OrderRepository
}

func NewUpdateOrderStatusUseCase(orderRepo ports.OrderRepository) *UpdateOrderStatusUseCase {
	return &UpdateOrderStatusUseCase{
		orderRepo: orderRepo,
	}
}

func (uc *UpdateOrderStatusUseCase) Execute(orderID string, status models.OrderStatus) error {
	// 1. Obtener la orden para conocer su tipo de servicio
	err := uc.orderRepo.UpdateStatus(orderID, status)
	if err != nil {
		return err
	}
	return nil
}
