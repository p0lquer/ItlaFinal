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
	"context"
	"log"
	"os"

	_ "ITLAFINAL/docs"

	"github.com/gin-contrib/cors"
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
	stRepo := repository.NewServiceTypeRepository(db)

	// 5. Use Cases
	createCustomer := customerUseCases.NewCreateCustomerUseCase(customerRepo)
	getAllCustomers := customerUseCases.NewGetAllCustomersUseCase(customerRepo)
	getAllServiceTypes := orderUseCases.NewGetServiceTypeUseCase(stRepo)
	createServiceType := orderUseCases.NewCreateServiceTypeUseCase(stRepo)

	createOrder := orderUseCases.NewCreateOrderUseCase(orderRepo, predRepo, stRepo)
	updateOrderStatus := orderUseCases.NewUpdateOrderStatusUseCase(orderRepo, predRepo)
	// 6. Handlers
	customerHandler := handlers.NewCustomerHandler(createCustomer, getAllCustomers, customerUseCases.NewDeleteCustomerUseCase(customerRepo))
	getAllOrders := orderUseCases.NewGetAllOrdersUseCase(orderRepo)
	getMyOrders := orderUseCases.NewGetMyOrdersUseCase(orderRepo)
	deleteOrder := orderUseCases.NewDeleteOrderUseCase(orderRepo, predRepo)
	serviceTypeHandler := handlers.NewServiceTypeHandler(getAllServiceTypes, createServiceType)

	orderHandler := handlers.NewOrderHandler(createOrder, updateOrderStatus, getAllOrders, getMyOrders, deleteOrder)

	// 7. Timer Worker en background
	worker := workers.NewTimerWorker(orderRepo, updateOrderStatus, hub)
	go worker.Start(context.Background())

	//8. Auth
	userRepo := repository.NewUserRepository(db)
	registerUC := userUseCases.NewRegisterUserUseCase(userRepo, customerRepo)
	loginUC := userUseCases.NewLoginUserUseCase(userRepo)
	deleteUserUC := userUseCases.NewDeleteUserUseCase(userRepo, customerRepo)
	authHandler := handlers.NewAuthHandler(registerUC, loginUC, deleteUserUC)

	// 9. Router
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000", "http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: false,
	}))

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
		api.GET("/service-types", serviceTypeHandler.GetAll)
		api.GET("/orders/mine", orderHandler.GetMy)
		api.POST("/orders", orderHandler.Create)

		// Órdenes — cualquier usuario autenticado puede ver
		api.GET("/orders", orderHandler.GetAll)

		// Órdenes — solo operadores pueden crear/modificar/eliminar
		operator := api.Group("/", middleware.OperatorOnly())
		{
			operator.DELETE("/users/:id", authHandler.DeleteUser)
			operator.PATCH("/orders/:id/status", orderHandler.UpdateStatus)
			operator.DELETE("/orders/:id", orderHandler.Delete)
			operator.POST("/service-types", serviceTypeHandler.Create)
		}
		// Clientes — solo operadores pueden crear/modificar/eliminar
		operator = api.Group("/", middleware.OperatorOnly())
		{
			operator.POST("/customers", customerHandler.Create)
			operator.GET("/customers", customerHandler.GetAll)
			operator.GET("/customers/:id", customerHandler.GetByID)
			operator.DELETE("/customers/:id", customerHandler.Delete)

		}

	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
