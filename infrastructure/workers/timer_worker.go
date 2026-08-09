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
	testMode          bool
}

func NewTimerWorker(
	orderRepo ports.OrderRepository,
	updateOrderStatus *orderUseCases.UpdateOrderStatusUseCase,
	notifier ports.Notifier,
	testMode bool,
) *TimerWorker {
	interval := 10 * time.Second
	if testMode {
		interval = 3 * time.Second // 🧪 modo prueba: avanza status cada 3s
	}
	return &TimerWorker{
		orderRepo:         orderRepo,
		updateOrderStatus: updateOrderStatus,
		notifier:          notifier,
		interval:          interval,
		testMode:          testMode,
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
			if w.testMode {
				w.advanceAllForTest()
			} else {
				w.checkOrders()
			}
		}
	}
}

// nextTestStatus devuelve el siguiente estado para el modo de prueba.
var nextTestStatus = map[models.OrderStatus]models.OrderStatus{
	models.StatusReceived:   models.StatusProcessing,
	models.StatusProcessing: models.StatusReady,
	models.StatusReady:      models.StatusDelivered,
}

// advanceAllForTest (modo prueba) avanza TODAS las órdenes al siguiente
// estado cada tick, para poder ver en consola los logs del notifier y la
// creación de predicciones sin esperar el tiempo estimado real.
func (w *TimerWorker) advanceAllForTest() {
	orders, err := w.orderRepo.FindAll()
	if err != nil {
		log.Printf("TimerWorker test: %v", err)
		return
	}

	for _, order := range orders {
		target, ok := nextTestStatus[order.Status]
		if !ok {
			continue
		}
		log.Printf("🧪 [TEST] Orden %s (%s): %s → %s", order.ID, order.ServiceType, order.Status, target)
		if err := w.updateOrderStatus.Execute(order.ID, target); err != nil {
			log.Printf("TimerWorker test: error orden %s: %v", order.ID, err)
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
