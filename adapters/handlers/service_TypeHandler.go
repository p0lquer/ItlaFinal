package handlers

import (
	"ITLAFINAL/adapters/dto"
	"ITLAFINAL/domain/usecases/orderUseCases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServiceTypeHandler struct {
	getAll *orderUseCases.GetServiceTypeUseCase
	create *orderUseCases.CreateServiceTypeUseCase
}

func NewServiceTypeHandler(
	getAll *orderUseCases.GetServiceTypeUseCase,
	create *orderUseCases.CreateServiceTypeUseCase,
) *ServiceTypeHandler {
	return &ServiceTypeHandler{
		getAll: getAll,
		create: create,
	}
}

func (h *ServiceTypeHandler) GetAll(c *gin.Context) {
	types, err := h.getAll.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, types)
}

func (h *ServiceTypeHandler) Create(c *gin.Context) {
	var req dto.CreateServiceTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	serviceType, err := h.create.Execute(req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, serviceType)
}
