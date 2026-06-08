package handlers

import (
	"ITLAFINAL/domain/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	createCustomer  *usecases.CreateCustomerUseCase
	getAllCustomers *usecases.GetAllCustomersUseCase
}

func NewCustomerHandler(
	create *usecases.CreateCustomerUseCase,
	getAll *usecases.GetAllCustomersUseCase,
) *CustomerHandler {
	return &CustomerHandler{
		createCustomer:  create,
		getAllCustomers: getAll,
	}
}

func (h *CustomerHandler) Create(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Phone string `json:"phone" binding:"required"`
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	customer, err := h.createCustomer.Execute(req.Name, req.Phone, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, customer)

}

func (h *CustomerHandler) GetAll(c *gin.Context) {
	customers, err := h.getAllCustomers.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customers)

}
