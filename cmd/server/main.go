package main

import (
	"ITLAFINAL/adapters/handlers"
	"ITLAFINAL/adapters/middleware"
	"ITLAFINAL/adapters/websocket"
	"ITLAFINAL/domain/usecases/customerUseCases"
	"ITLAFINAL/domain/usecases/orderUseCases"
	"ITLAFINAL/domain/usecases/userUseCases"
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
	// @securityDefinitions.apikey BearerAuth
	// @in header
	// @name Authorization
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

	//8. Auth
	userRepo := repository.NewUserRepository(db)
	registerUC := userUseCases.NewRegisterUserUseCase(userRepo, customerRepo)
	loginUC := userUseCases.NewLoginUserUseCase(userRepo)
	deleteUserUC := userUseCases.NewDeleteUserUseCase(userRepo, customerRepo)
	authHandler := handlers.NewAuthHandler(registerUC, loginUC, deleteUserUC)

	// 9. Router
	r := gin.Default()

	// Públicas — sin middleware
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	//protected
	api := r.Group("/api", middleware.AuthRequired())
	{
		api.GET("/auth/me", authHandler.Me)

		// Órdenes — cualquier usuario autenticado puede ver
		api.GET("/orders", orderHandler.GetAll)

		// Órdenes — solo operadores pueden crear/modificar/eliminar
		operator := api.Group("/", middleware.OperatorOnly())
		{
			operator.DELETE("/users/:id", authHandler.DeleteUser)
			operator.POST("/orders", orderHandler.Create)
			operator.PATCH("/orders/:id/status", orderHandler.UpdateStatus)
			operator.DELETE("/orders/:id", orderHandler.Delete)
		}
	}

	{
		api.POST("/customers", customerHandler.Create)
		api.GET("/customers", customerHandler.GetAll)
		api.GET("/customers/:id", customerHandler.GetByID)
	}

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
