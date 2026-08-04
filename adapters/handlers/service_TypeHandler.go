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

// @GetAllServiceTypes godoc
// @Summary Obtener todos los tipos de servicio
// @Description Recupera la lista de todos los tipos de servicio registrados
// @Tags service-types
// @Security BearerAuth
// @Accept  json
// @Success 200 {array} dto.ServiceTypeResponse
// @Router /service-types [get]
// @Produce  json
func (h *ServiceTypeHandler) GetAll(c *gin.Context) {
	types, err := h.getAll.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, types)
}

// @CreateServiceType godoc
// @Summary Crear un tipo de servicio
// @Description Registra un nuevo tipo de servicio en la base de datos
// @Tags service-types
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Success 201 "Tipo de servicio creado con éxito"
// @Router /service-types [post]
// @Param serviceType body dto.CreateServiceTypeRequest true "Service type data"
func (h *ServiceTypeHandler) Create(c *gin.Context) {
	var req dto.CreateServiceTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	serviceType, err := h.create.Execute(req.Name, req.Description, req.BasePrice, req.PricePerWeight, req.PricePerPiece)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, serviceType)
}
