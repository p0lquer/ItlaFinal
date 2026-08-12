package handlers

import (
	"ITLAFINAL/adapters/dto"
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/usecases/orderUseCases"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	createOrder       *orderUseCases.CreateOrderUseCase
	updateOrderStatus *orderUseCases.UpdateOrderStatusUseCase
	getAllOrders      *orderUseCases.GetAllOrdersUseCase
	getMyOrders       *orderUseCases.GetMyOrdersUseCase
	deleteOrder       *orderUseCases.DeleteOrderUseCase
	getOrder          *orderUseCases.GetOrderUseCase
	getOrderHistory   *orderUseCases.GetOrderHistoryUseCase
	payment           *orderUseCases.PaymentUseCase
}

func NewOrderHandler(
	create *orderUseCases.CreateOrderUseCase,
	update *orderUseCases.UpdateOrderStatusUseCase,
	getAll *orderUseCases.GetAllOrdersUseCase,
	getMy *orderUseCases.GetMyOrdersUseCase,
	delete *orderUseCases.DeleteOrderUseCase,
	getOrder *orderUseCases.GetOrderUseCase,
	getHistory *orderUseCases.GetOrderHistoryUseCase,
	payment *orderUseCases.PaymentUseCase,
) *OrderHandler {
	return &OrderHandler{
		createOrder:       create,
		updateOrderStatus: update,
		getAllOrders:      getAll,
		getMyOrders:       getMy,
		deleteOrder:       delete,
		getOrder:          getOrder,
		getOrderHistory:   getHistory,
		payment:           payment,
	}
}

// PaymentSummary returns the invoice-ready price and any completed payment.
func (h *OrderHandler) PaymentSummary(c *gin.Context) {
	order, err := h.loadAuthorizedOrder(c)
	if err != nil {
		writeOrderError(c, err)
		return
	}
	summary, err := h.payment.Summary(order.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, paymentPayload(summary))
}

// Invoice returns the immutable payment receipt plus the order line items.
func (h *OrderHandler) Invoice(c *gin.Context) {
	order, err := h.loadAuthorizedOrder(c)
	if err != nil {
		writeOrderError(c, err)
		return
	}
	summary, err := h.payment.Summary(order.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if summary.Payment == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "la factura estara disponible despues del pago"})
		return
	}
	c.JSON(http.StatusOK, paymentPayload(summary))
}

type payOrderRequest struct {
	Method string `json:"method" binding:"required"`
}

// Pay creates one simulated, final payment. No sensitive payment data is accepted.
func (h *OrderHandler) Pay(c *gin.Context) {
	order, err := h.loadAuthorizedOrder(c)
	if err != nil {
		writeOrderError(c, err)
		return
	}
	var req payOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	payment, err := h.payment.Pay(order.ID, req.Method)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, paymentPayload(&models.PaymentSummary{Order: order, Payment: payment}))
}

func paymentPayload(summary *models.PaymentSummary) gin.H {
	status := "pending"
	payload := gin.H{
		"order":          dto.NewOrderResponse(summary.Order),
		"subtotal":       summary.Order.EstimatedCost,
		"total":          summary.Order.EstimatedCost,
		"currency":       "DOP",
		"payment_status": status,
		"can_pay":        summary.CanPay,
	}
	if summary.Payment != nil {
		status = summary.Payment.Status
		payload["payment_status"] = status
		payload["invoice_number"] = summary.Payment.ReceiptNo
		payload["paid_at"] = summary.Payment.PaidAt
		payload["payment_method"] = paymentMethodPayload(summary.Payment.Method)
		payload["receipt"] = summary.Payment
	}
	return payload
}

func paymentMethodPayload(method string) string {
	switch method {
	case "tarjeta":
		return "card"
	case "transferencia":
		return "transfer"
	case "efectivo":
		return "cash"
	default:
		return method
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

	customerID := strings.TrimSpace(req.CustomerID)
	if c.GetString("role") == "customer" {
		// A customer may only create an order for their own profile; any ID sent
		// by the client is intentionally ignored.
		customerID = c.GetString("user_id")
	}
	if customerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer_id requerido"})
		return
	}

	order, err := h.createOrder.Execute(customerID, req.ServiceType, req.PiecesCount, req.Notes, req.Weight)
	if err != nil {
		// Most failures at this boundary are business validation failures
		// (service availability, pieces, weight and notes). Returning 400 lets
		// clients show a useful correction instead of a generic outage.
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.NewOrderResponse(order))
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
	filter, err := orderFilterFromQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	page, err := h.getAllOrders.ExecutePage(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	orders := make([]dto.OrderResponse, 0, len(page.Orders))
	for _, order := range page.Orders {
		orders = append(orders, dto.NewOrderResponse(order))
	}
	totalPages := 0
	if page.Total > 0 {
		totalPages = (page.Total + page.PageSize - 1) / page.PageSize
	}
	c.JSON(http.StatusOK, dto.OrderPageResponse{
		Data: orders,
		Pagination: dto.PaginationResponse{
			Page: page.Page, PageSize: page.PageSize, Total: page.Total, TotalPages: totalPages,
		},
	})
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
	response := make([]dto.OrderResponse, 0, len(orders))
	for _, order := range orders {
		response = append(response, dto.NewOrderResponse(order))
	}
	c.JSON(http.StatusOK, response)
}

// GetDetail returns the complete business data for one order plus its state
// history. Customers can only access their own orders.
func (h *OrderHandler) GetDetail(c *gin.Context) {
	order, err := h.loadAuthorizedOrder(c)
	if err != nil {
		writeOrderError(c, err)
		return
	}
	history, err := h.getOrderHistory.Execute(order.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.OrderDetailResponse{Order: dto.NewOrderResponse(order), History: history})
}

// GetHistory exposes the audit log on its own for clients that only need to
// refresh production activity.
func (h *OrderHandler) GetHistory(c *gin.Context) {
	order, err := h.loadAuthorizedOrder(c)
	if err != nil {
		writeOrderError(c, err)
		return
	}
	history, err := h.getOrderHistory.Execute(order.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
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

	change := models.OrderStatusChange{
		ChangedBy:     c.GetString("user_id"),
		ChangedByRole: c.GetString("role"),
		Description:   strings.TrimSpace(req.Description),
	}
	if err := h.updateOrderStatus.Execute(orderID, req.Status, change); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "estado actualizado"})
}

func (h *OrderHandler) loadAuthorizedOrder(c *gin.Context) (*models.Order, error) {
	if h.getOrder == nil {
		return nil, errors.New("detalle de orden no configurado")
	}
	order, err := h.getOrder.Execute(c.Param("id"))
	if err != nil {
		return nil, err
	}
	if c.GetString("role") == "customer" && order.CustomerID != c.GetString("user_id") {
		return nil, errOrderForbidden
	}
	return order, nil
}

var errOrderForbidden = errors.New("no tiene permiso para acceder a esta orden")

func writeOrderError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, errOrderForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, sql.ErrNoRows):
		c.JSON(http.StatusNotFound, gin.H{"error": "orden no encontrada"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func orderFilterFromQuery(c *gin.Context) (models.OrderFilter, error) {
	filter := models.OrderFilter{
		Search:      firstQuery(c, "q", "search"),
		CustomerID:  strings.TrimSpace(c.Query("customer_id")),
		Status:      models.OrderStatus(strings.TrimSpace(c.Query("status"))),
		ServiceType: strings.TrimSpace(c.Query("service_type")),
		Page:        1,
		PageSize:    20,
	}
	if filter.Status != "" && !isKnownOrderStatus(filter.Status) {
		return models.OrderFilter{}, errors.New("estado de orden invÃ¡lido")
	}
	var err error
	if filter.Page, err = positiveQueryInt(c, "page", 1, 1, 1000000); err != nil {
		return models.OrderFilter{}, err
	}
	if filter.PageSize, err = positiveQueryInt(c, "page_size", 20, 1, 100); err != nil {
		return models.OrderFilter{}, err
	}
	if filter.From, err = orderQueryDate(c.Query("from"), false); err != nil {
		return models.OrderFilter{}, errors.New("from debe usar YYYY-MM-DD o RFC3339")
	}
	if filter.To, err = orderQueryDate(c.Query("to"), true); err != nil {
		return models.OrderFilter{}, errors.New("to debe usar YYYY-MM-DD o RFC3339")
	}
	if filter.From != nil && filter.To != nil && filter.From.After(*filter.To) {
		return models.OrderFilter{}, errors.New("from no puede ser posterior a to")
	}
	return filter, nil
}

func firstQuery(c *gin.Context, names ...string) string {
	for _, name := range names {
		if value := strings.TrimSpace(c.Query(name)); value != "" {
			return value
		}
	}
	return ""
}

func positiveQueryInt(c *gin.Context, name string, fallback, min, max int) (int, error) {
	raw := strings.TrimSpace(c.Query(name))
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < min || value > max {
		return 0, errors.New(name + " debe estar entre " + strconv.Itoa(min) + " y " + strconv.Itoa(max))
	}
	return value, nil
}

func orderQueryDate(raw string, endOfDay bool) (*time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	if date, err := time.Parse(time.DateOnly, raw); err == nil {
		if endOfDay {
			date = date.AddDate(0, 0, 1).Add(-time.Nanosecond)
		}
		return &date, nil
	}
	date, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return nil, err
	}
	return &date, nil
}

func isKnownOrderStatus(status models.OrderStatus) bool {
	return status == models.StatusReceived || status == models.StatusProcessing || status == models.StatusReady || status == models.StatusDelivered
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
