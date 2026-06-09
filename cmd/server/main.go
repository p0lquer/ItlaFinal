package main

import (
	"ITLAFINAL/adapters/handlers"
	"ITLAFINAL/adapters/websocket"
	"ITLAFINAL/domain/usecases/customerUseCases"
	"ITLAFINAL/domain/usecases/orderUseCases"
	"ITLAFINAL/infrastructure/database"
	"ITLAFINAL/infrastructure/repository"
	"ITLAFINAL/infrastructure/workers"
	"log"
	"net/http"
	"os"

	_ "ITLAFINAL/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"

	ginSwagger "github.com/swaggo/gin-swagger"
	// @title ITLAFINAL API
	// @version 1.0
	// @description API para la gestión de ordenes y clientes.
	// @host localhost:8080
	// @BasePath /api
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
	customerRepo := repository.NewCustomerRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	predRepo := repository.NewPredictionRepository(db)

	// 5. Use Cases
	createCustomer := customerUseCases.NewCreateCustomerUseCase(customerRepo)
	getAllCustomers := customerUseCases.NewGetAllCustomersUseCase(customerRepo)

	createOrder := orderUseCases.NewCreateOrderUseCase(orderRepo, predRepo)
	updateOrderStatus := orderUseCases.NewUpdateOrderStatusUseCase(orderRepo)

	// 6. Handlers
	customerHandler := handlers.NewCustomerHandler(createCustomer, getAllCustomers)
	getAllOrders := orderUseCases.NewGetAllOrdersUseCase(orderRepo)
	deleteOrder := orderUseCases.NewDeleteOrderUseCase(orderRepo, predRepo)

	orderHandler := handlers.NewOrderHandler(createOrder, updateOrderStatus, getAllOrders, deleteOrder)

	// 7. Timer Worker en background
	worker := workers.NewTimerWorker(orderRepo, hub)
	go worker.Start()

	// 8. Router
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
		api.POST("/customers", customerHandler.Create)
		api.GET("/customers", customerHandler.GetAll)

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

	log.Print(" 🌐 http://localhost:8080/swagger/index.html#/")
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
