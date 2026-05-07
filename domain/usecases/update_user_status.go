package usecases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"errors"
	"time"
)

type UpdateOrderStatusUseCase struct {
	orderRepo ports.OrderRepository
	notifier  ports.Notifier
}

func NewUpdateOrderStatusUseCase(
	orderRepo ports.OrderRepository,
	notifier ports.Notifier,
) *UpdateOrderStatusUseCase {
	return &UpdateOrderStatusUseCase{
		orderRepo: orderRepo,
		notifier:  notifier,
	}
}

func (uc *UpdateOrderStatusUseCase) Execute(orderID string, newStatus models.OrderStatus) error {
	// 1. Verificar que la orden exista
	order, err := uc.orderRepo.FindByID(orderID)
	if err != nil {
		return errors.New("orden no encontrada")
	}

	// 2. Validar transición de estado
	if !isValidTransition(order.Status, newStatus) {
		return errors.New("transición de estado inválida")
	}

	// 3. Actualizar estado
	if err := uc.orderRepo.UpdateStatus(orderID, newStatus); err != nil {
		return err
	}

	// 4. Si la orden está lista, notificar al cliente
	if newStatus == models.StatusReady {
		_ = uc.notifier.NotifyOrderReady(order.CustomerID, orderID)
		now := time.Now()
		order.ReadyAt = &now
	}

	return nil
}

// Solo se permiten estas transiciones en orden
func isValidTransition(current, next models.OrderStatus) bool {
	transitions := map[models.OrderStatus]models.OrderStatus{
		models.StatusReceived:   models.StatusProcessing,
		models.StatusProcessing: models.StatusReady,
		models.StatusReady:      models.StatusDelivered,
	}
	return transitions[current] == next
}
