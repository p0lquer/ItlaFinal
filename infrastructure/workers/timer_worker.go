package workers

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"log"
	"time"
)

// TimerWorker corre en background con goroutines
// Monitorea órdenes y notifica cuando el tiempo estimado llega
type TimerWorker struct {
	orderRepo ports.OrderRepository
	notifier  ports.Notifier
	interval  time.Duration
}

func NewTimerWorker(orderRepo ports.OrderRepository, notifier ports.Notifier) *TimerWorker {
	return &TimerWorker{
		orderRepo: orderRepo,
		notifier:  notifier,
		interval:  30 * time.Second, // revisa cada 30 segundos
	}
}

// Start lanza el worker en background — se llama con go worker.Start()
func (w *TimerWorker) Start() {
	log.Println("⏱️  TimerWorker iniciado")
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for range ticker.C {
		w.checkOrders()
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

		deadline := order.CreatedAt.Add(order.EstimatedTime)

		// Si ya pasó el tiempo estimado, notificar
		if time.Now().After(deadline) {
			log.Printf("🔔 Orden %s lista (tiempo estimado alcanzado)", order.ID)
			_ = w.notifier.NotifyOrderReady(order.CustomerID, order.ID)
		}
	}
}
