package handlers

import (
	"ITLAFINAL/adapters/dto"
	"ITLAFINAL/domain/usecases/customerUseCases"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	createCustomer  *customerUseCases.CreateCustomerUseCase
	getAllCustomers *customerUseCases.GetAllCustomersUseCase
	deleteCustomer  *customerUseCases.DeleteCustomerUseCase
}

func NewCustomerHandler(
	create *customerUseCases.CreateCustomerUseCase,
	getAll *customerUseCases.GetAllCustomersUseCase,
	delete *customerUseCases.DeleteCustomerUseCase,
) *CustomerHandler {
	return &CustomerHandler{
		createCustomer:  create,
		getAllCustomers: getAll,
		deleteCustomer:  delete,
	}
}

// CreateCustomer godoc
// @Summary Crear un cliente
// @Description Registra un nuevo cliente en la base de datos
// @Tags customers
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Success 201 "Cliente creado con éxito"
// @Router /customers [post]
// @Param customer body dto.CreateCustomerRequest true "Customer data"
func (h *CustomerHandler) Create(c *gin.Context) {
	var req dto.CreateCustomerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	customer, err := h.createCustomer.Execute(req.ID, req.Name, req.Phone, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, customer)

}

// GetAllCustomers godoc
// @Summary Obtener todos los clientes
// @Description Recupera la lista de todos los clientes registrados
// @Tags customers
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Success 200 "Lista de clientes obtenida con éxito"
// @Router /customers [get]
func (h *CustomerHandler) GetAll(c *gin.Context) {
	customers, err := h.getAllCustomers.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customers)

}

// GetCustomerByID godoc
// @Summary Obtener un cliente por ID
// @Description Recupera los detalles de un cliente específico utilizando su ID
// @Tags customers
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path string true "ID del cliente"
// @Success 200 "Cliente obtenido con éxito"
// @Router /customers/{id} [get]
func (h *CustomerHandler) GetByID(c *gin.Context) {
	customerID := c.Param("id")
	customer, err := h.getAllCustomers.ExecuteByID(customerID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cliente no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customer)
}

// DeleteCustomer godoc
// @Summary Eliminar un cliente
// @Description Elimina un cliente específico utilizando su ID
// @Tags customers
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path string true "ID del cliente"
// @Success 200 "Cliente eliminado con éxito"
// @Router /customers/{id} [delete]
func (h *CustomerHandler) Delete(c *gin.Context) {
	customerID := c.Param("id")
	err := h.deleteCustomer.Execute(customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cliente eliminado con éxito"})
}
