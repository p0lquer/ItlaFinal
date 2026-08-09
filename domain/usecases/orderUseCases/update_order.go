package orderUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
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

func NewUpdateOrderStatusUseCase(orderRepo ports.OrderRepository, predRepo ports.PredictionRepository, notifier ports.Notifier) *UpdateOrderStatusUseCase {
	return &UpdateOrderStatusUseCase{orderRepo: orderRepo, predRepo: predRepo, notifier: notifier}
}

func (uc *UpdateOrderStatusUseCase) Execute(orderID string, target models.OrderStatus) error {
	updated, err := uc.orderRepo.Transition(orderID, target)
	if err != nil {
		return err
	}
	if !updated {
		return nil
	}
	order, err := uc.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}

	if uc.notifier != nil {
		_ = uc.notifier.NotifyStatusChange(order.CustomerID, order.ID, string(target))
	}

	if target != models.StatusReady {
		return nil
	}

	actual := order.ElapsedTime()
	if actual < minRecordedMinutes {
		actual = minRecordedMinutes
	}

	prediction := &models.Prediction{
		ID:          uuid.NewString(),
		OrderID:     order.ID,
		ServiceType: normalizeServiceType(order.ServiceType), // misma clave que GetHistoricalData
		PiecesCount: order.PiecesCount,
		Weight:      order.Weight,
		Estimated:   order.EstimatedTime,
		Actual:      &actual,
		CreatedAt:   time.Now(),
	}

	if err := uc.predRepo.UpsertForOrder(order.ID, prediction); err != nil {
		return err
	}
	log.Printf("📊 Predicción creada/actualizada para orden %s: %s, %d piezas, %.1f kg → actual=%.0fm estimado=%.0fm",
		order.ID, prediction.ServiceType, prediction.PiecesCount, prediction.Weight,
		prediction.Actual.Minutes(), prediction.Estimated.Minutes())
	return nil
}
