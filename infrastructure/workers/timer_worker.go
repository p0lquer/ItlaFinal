package workers

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"ITLAFINAL/domain/usecases/orderUseCases"
	"context"
	"log"
	"time"
)

// TimerWorker escanea órdenes en "en_proceso" y, cuando su tiempo estimado
// ya transcurrió (medido desde StartedAt), las marca como "lista" llamando al
// MISMO UpdateOrderStatusUseCase que el endpoint manual. La transición
// guardada de ese use case impide duplicar registros de entrenamiento.
type TimerWorker struct {
	orderRepo         ports.OrderRepository
	notifier          ports.Notifier
	updateOrderStatus *orderUseCases.UpdateOrderStatusUseCase
	interval          time.Duration
}

func NewTimerWorker(
	orderRepo ports.OrderRepository,
	updateOrderStatus *orderUseCases.UpdateOrderStatusUseCase,
	notifier ports.Notifier,
) *TimerWorker {
	return &TimerWorker{
		orderRepo:         orderRepo,
		updateOrderStatus: updateOrderStatus,
		notifier:          notifier,
		interval:          30 * time.Second, // cada 30s
	}
}

// Start lanza el worker en background — se llama con go worker.Start()
func (w *TimerWorker) Start(ctx context.Context) {
	log.Println("⏱️  TimerWorker iniciado")
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.checkOrders()
		}
	}
}

func (w *TimerWorker) checkOrders() {
	orders, err := w.orderRepo.FindAll()
	if err != nil {
		log.Printf("TimerWorker error: %v", err)
		return
	}

	for _, order := range orders {
		if order.Status != models.StatusProcessing {
			continue
		}
		actualDuration := time.Since(order.CreatedAt)
		if time.Since(order.StartedBase()) < order.EstimatedTime {
			continue
		}
		if order.Status != models.StatusProcessing {
			continue
		}

		if actualDuration < order.EstimatedTime {
			continue
		}

		// Pasa por el mismo use case que usa el endpoint manual de status:
		// así también queda guardado el dato de entrenamiento (peso + tiempo
		// real) y se marca ready_at, sin duplicar esa lógica aquí. Además
		// evita que sigamos re-notificando cada 30s: una vez que el status
		// cambia a "lista", este mismo filtro de arriba ya la ignora.
		if err := w.updateOrderStatus.Execute(order.ID, models.StatusReady); err != nil {
			log.Printf("TimerWorker: error marcando %s como lista: %v", order.ID, err)
			continue
		}

		log.Printf("🔔 Orden %s lista (tiempo estimado alcanzado)", order.ID)
		// _ = w.notifier.NotifyOrderReady(order.CustomerID, order.ID)
	}
}
