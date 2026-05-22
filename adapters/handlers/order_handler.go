package handlers

import (
	"ITLAFINAL/adapters/dto"
	"ITLAFINAL/domain/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	createOrder       *usecases.CreateOrderUseCase
	updateOrderStatus *usecases.UpdateOrderStatusUseCase
	getAllOrders      *usecases.GetAllOrdersUseCase
	deleteOrder       *usecases.DeleteOrderUseCase
}

func NewOrderHandler(
	create *usecases.CreateOrderUseCase,
	update *usecases.UpdateOrderStatusUseCase,
	getAll *usecases.GetAllOrdersUseCase,
	delete *usecases.DeleteOrderUseCase,
) *OrderHandler {
	return &OrderHandler{
		createOrder:       create,
		updateOrderStatus: update,
		getAllOrders:      getAll,
		deleteOrder:       delete,
	}
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.createOrder.Execute(req.CustomerID, req.ServiceType, req.PiecesCount, req.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.OrderResponse{
		ID:            order.ID,
		Status:        string(order.Status),
		EstimatedTime: order.EstimatedTime.Minutes(),
		CreatedAt:     order.CreatedAt,
	})
}

func (h *OrderHandler) GetAll(c *gin.Context) {
	orders, err := h.getAllOrders.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

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

func (h *OrderHandler) Delete(c *gin.Context) {
	orderID := c.Param("id")

	if err := h.deleteOrder.Execute(orderID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "orden eliminada correctamente"})

}
