package handlers

import (
	"ITLAFINAL/adapters/dto"
	"ITLAFINAL/domain/usecases/orderUseCases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	createOrder       *orderUseCases.CreateOrderUseCase
	updateOrderStatus *orderUseCases.UpdateOrderStatusUseCase
	getAllOrders      *orderUseCases.GetAllOrdersUseCase
	getMyOrders       *orderUseCases.GetMyOrdersUseCase
	deleteOrder       *orderUseCases.DeleteOrderUseCase
}

func NewOrderHandler(
	create *orderUseCases.CreateOrderUseCase,
	update *orderUseCases.UpdateOrderStatusUseCase,
	getAll *orderUseCases.GetAllOrdersUseCase,
	getMy *orderUseCases.GetMyOrdersUseCase,
	delete *orderUseCases.DeleteOrderUseCase,
) *OrderHandler {
	return &OrderHandler{
		createOrder:       create,
		updateOrderStatus: update,
		getAllOrders:      getAll,
		getMyOrders:       getMy,
		deleteOrder:       delete,
	}
}

// @CreateOrder godoc
// @Summary Crear una orden
// @Description Registra una nueva orden en la base de datos
// @Tags orders
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Success 201 "Orden creada con éxito"
// @Router /orders [post]
// @Param order body dto.CreateOrderRequest true "Order data"
func (h *OrderHandler) Create(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customerID := req.CustomerID
	if customerID == "" || (c.GetString("role") == "customer" && customerID != c.GetString("user_id")) {
		customerID = c.GetString("user_id")
	}
	if customerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer_id requerido"})
		return
	}

	order, err := h.createOrder.Execute(customerID, req.ServiceType, req.PiecesCount, req.Notes, req.Weight)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.OrderResponse{
		ID:            order.ID,
		Status:        string(order.Status),
		EstimatedTime: order.EstimatedTime.Minutes(),
		EstimatedCost: order.EstimatedCost,
		CreatedAt:     order.CreatedAt,
	})
}

// GetAllOrders godoc
// @Summary Obtener todas las órdenes
// @Description Recupera la lista de todas las órdenes registradas
// @Tags orders
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Success 200 "Lista de órdenes obtenida con éxito"
// @Router /orders [get]
func (h *OrderHandler) GetAll(c *gin.Context) {
	orders, err := h.getAllOrders.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// GetMy godoc
// @Summary Obtener mis órdenes
// @Description Devuelve únicamente las órdenes del cliente autenticado
// @Tags orders
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.OrderResponse
// @Router /orders/mine [get]
func (h *OrderHandler) GetMy(c *gin.Context) {
	customerID := c.GetString("user_id")
	orders, err := h.getMyOrders.Execute(customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// @UpdateOrderStatus godoc
// @Summary Actualizar el estado de una orden
// @Description Actualiza el estado de una orden existente
// @Tags orders
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path string true "ID de la orden"
// @Param status body string true "Nuevo estado de la orden"
// @Success 200 "Estado actualizado con éxito"
// @Router /orders/{id}/status [patch]
func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	orderID := c.Param("id")

	var req dto.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.updateOrderStatus.Execute(orderID, req.Status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "estado actualizado"})
}

// @DeleteOrder godoc
// @Summary Eliminar una orden
// @Description Elimina una orden existente de la base de datos
// @Tags orders
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path string true "ID de la orden"
// @Success 200 "Orden eliminada con éxito"
// @Router /orders/{id} [delete]
func (h *OrderHandler) Delete(c *gin.Context) {
	orderID := c.Param("id")

	if err := h.deleteOrder.Execute(orderID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "orden eliminada correctamente"})

}
