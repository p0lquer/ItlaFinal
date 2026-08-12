package orderUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

const minRecordedMinutes = 15 * time.Minute

type UpdateOrderStatusUseCase struct {
	orderRepo ports.OrderRepository
	predRepo  ports.PredictionRepository
	notifier  ports.Notifier
}

// The notifier is optional to preserve the existing composition. main.go should
// inject the WebSocket hub as the third argument so HTTP and worker changes use
// the same notification path.
func NewUpdateOrderStatusUseCase(orderRepo ports.OrderRepository, predRepo ports.PredictionRepository, notifiers ...ports.Notifier) *UpdateOrderStatusUseCase {
	uc := &UpdateOrderStatusUseCase{orderRepo: orderRepo, predRepo: predRepo}
	if len(notifiers) > 0 {
		uc.notifier = notifiers[0]
	}
	return uc
}

func (uc *UpdateOrderStatusUseCase) Execute(orderID string, target models.OrderStatus) error {
	order, err := uc.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}
	if order.Status == target {
		return nil // retry-safe: no duplicate prediction or notification.
	}
	if !IsValidStatusTransition(order.Status, target) {
		return errors.New("transición de estado inválida")
	}

	updated, err := uc.orderRepo.Transition(orderID, target)
	if err != nil {
		return err
	}
	if !updated {
		return errors.New("la orden cambió de estado; vuelva a intentarlo")
	}

	if target == models.StatusReady {
		actual := order.ElapsedTime()
		if actual < minRecordedMinutes {
			actual = minRecordedMinutes
		}
		prediction := &models.Prediction{
			ID:          uuid.NewString(),
			OrderID:     order.ID,
			ServiceType: normalizeServiceType(order.ServiceType),
			PiecesCount: order.PiecesCount,
			Weight:      order.Weight,
			Estimated:   order.EstimatedTime,
			Actual:      &actual,
			CreatedAt:   time.Now(),
		}
		if err := uc.predRepo.UpsertForOrder(order.ID, prediction); err != nil {
			return err
		}
	}

	uc.notify(order, target)
	return nil
}

func (uc *UpdateOrderStatusUseCase) notify(order *models.Order, target models.OrderStatus) {
	if uc.notifier == nil {
		return
	}
	if err := uc.notifier.NotifyStatusChange(order.CustomerID, order.ID, string(target)); err != nil {
		log.Printf("order %s: status notification: %v", order.ID, err)
	}
	if target == models.StatusReady {
		if err := uc.notifier.NotifyOrderReady(order.CustomerID, order.ID); err != nil {
			log.Printf("order %s: ready notification: %v", order.ID, err)
		}
	}
}

// IsValidStatusTransition contains the only permitted business workflow.
func IsValidStatusTransition(current, target models.OrderStatus) bool {
	return (current == models.StatusReceived && target == models.StatusProcessing) ||
		(current == models.StatusProcessing && target == models.StatusReady) ||
		(current == models.StatusReady && target == models.StatusDelivered)
}
