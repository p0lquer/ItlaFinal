package main

import (
	"ITLAFINAL/adapters/handlers"
	"ITLAFINAL/adapters/websocket"
	"ITLAFINAL/domain/usecases"
	"ITLAFINAL/infrastructure/database"
	"ITLAFINAL/infrastructure/repository"
	"ITLAFINAL/infrastructure/workers"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Variables de entorno
	godotenv.Load()

	// 2. Base de datos
	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatalf("❌ PostgreSQL: %v", err)
	}
	defer db.Close()

	// 3. WebSocket Hub (implementa el port Notifier)
	hub := websocket.NewHub()

	// 4. Repositories
	orderRepo := repository.NewOrderRepository(db)
	predRepo := repository.NewPredictionRepository(db)

	// 5. Use Cases
	createOrder := usecases.NewCreateOrderUseCase(orderRepo, predRepo)
	updateOrderStatus := usecases.NewUpdateOrderStatusUseCase(orderRepo, hub)

	// 6. Handlers
	getAllOrders := usecases.NewGetAllOrdersUseCase(orderRepo)
	deleteOrder := usecases.NewDeleteOrderUseCase(orderRepo, predRepo)

	orderHandler := handlers.NewOrderHandler(createOrder, updateOrderStatus, getAllOrders, deleteOrder)

	// 7. Timer Worker en background
	worker := workers.NewTimerWorker(orderRepo, hub)
	go worker.Start()

	// 8. Router
	r := gin.Default()

	// CORS para el frontend React
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// WebSocket endpoint
	r.GET("/ws", func(c *gin.Context) {
		hub.HandleConnection(c.Writer, c.Request)
	})

	// API REST
	api := r.Group("/api")
	{
		api.POST("/orders", orderHandler.Create)
		api.GET("/orders", orderHandler.GetAll)
		api.PATCH("/orders/:id/status", orderHandler.UpdateStatus)
		api.DELETE("/orders/:id", orderHandler.Delete)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Servidor corriendo en :%s", port)
	r.Run(":" + port)
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
