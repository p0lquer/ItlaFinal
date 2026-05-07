package main

import (
	"ITLAFINAL/adapters/handlers"
	"ITLAFINAL/domain/usecases"
	"ITLAFINAL/infrastructure/database"
	"ITLAFINAL/infrastructure/repository"
	"ITLAFINAL/infrastructure/workers"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Cargar variables de entorno
	godotenv.Load()

	// 2. Conectar base de datos
	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatalf("❌ Error conectando a PostgreSQL: %v", err)
	}

	// 3. Repositories (infrastructure implementa los ports)
	orderRepo := repository.NewOrderRepository(db)
	predRepo := repository.NewPredictionRepository(db)

	// 4. Notifier (implementación simple por ahora)
	notifier := &LogNotifier{} // luego reemplazas con WebSocket real

	// 5. Use Cases (dominio puro)
	createOrder := usecases.NewCreateOrderUseCase(orderRepo, predRepo)
	updateOrderStatus := usecases.NewUpdateOrderStatusUseCase(orderRepo, notifier)

	// 6. Handlers (adapters)
	orderHandler := handlers.NewOrderHandler(createOrder, updateOrderStatus)

	// 7. Timer Worker en background
	worker := workers.NewTimerWorker(orderRepo, notifier)
	go worker.Start()

	// 8. Router
	r := gin.Default()
	api := r.Group("/api")
	{
		api.POST("/orders", orderHandler.Create)
		api.PATCH("/orders/:id/status", orderHandler.UpdateStatus)
	}

	log.Println("🚀 Servidor corriendo en :8080")
	r.Run(":8080")
}

// LogNotifier — implementación temporal del port Notifier
type LogNotifier struct{}

func (n *LogNotifier) NotifyOrderReady(customerID, orderID string) error {
	log.Printf("🔔 Orden %s lista para cliente %s", orderID, customerID)
	return nil
}
func (n *LogNotifier) NotifyStatusChange(customerID, orderID, status string) error {
	log.Printf("📦 Orden %s cambió a '%s' para cliente %s", orderID, status, customerID)
	return nil
}
