package handlers

import (
	"ITLAFINAL/adapters/dto"
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/usecases/userUseCases"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	users *userUseCases.AdminUsersUseCase
}

func NewAdminHandler(users *userUseCases.AdminUsersUseCase) *AdminHandler {
	return &AdminHandler{users: users}
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, pageSize := positiveQuery(c, "page", 1), positiveQuery(c, "page_size", 5)
	if pageSize > 50 {
		pageSize = 50
	}
	filter := models.UserFilter{Search: c.Query("search"), Page: page, PageSize: pageSize}
	if role := strings.ToLower(strings.TrimSpace(c.Query("role"))); role != "" {
		if role != string(models.RoleCustomer) && role != string(models.RoleOperator) && role != string(models.RoleAdmin) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "role debe ser customer, operator o admin"})
			return
		}
		filter.Role = models.UserRole(role)
	}
	if status := strings.ToLower(strings.TrimSpace(c.Query("status"))); status != "" {
		switch status {
		case "active":
			value := true
			filter.IsActive = &value
		case "blocked":
			value := false
			filter.IsActive = &value
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "status debe ser active o blocked"})
			return
		}
	}
	result, err := h.users.List(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron listar los usuarios"})
		return
	}
	users := make([]dto.AdminUserResponse, 0, len(result.Users))
	for _, user := range result.Users {
		users = append(users, dto.AdminUserResponse{ID: user.ID, Name: user.Name, Email: user.Email, Role: string(user.Role), IsActive: user.IsActive, CreatedAt: user.CreatedAt})
	}
	totalPages := (result.Total + result.PageSize - 1) / result.PageSize
	c.JSON(http.StatusOK, gin.H{"users": users, "pagination": gin.H{"page": result.Page, "page_size": result.PageSize, "total": result.Total, "total_pages": totalPages}})
}

func (h *AdminHandler) BlockUser(c *gin.Context)   { h.setActive(c, false) }
func (h *AdminHandler) UnblockUser(c *gin.Context) { h.setActive(c, true) }

func (h *AdminHandler) setActive(c *gin.Context, active bool) {
	if err := h.users.SetActive(c.GetString("user_id"), c.Param("id"), active); err != nil {
		adminError(c, err)
		return
	}
	message := "usuario bloqueado correctamente"
	if active {
		message = "usuario desbloqueado correctamente"
	}
	c.JSON(http.StatusOK, gin.H{"message": message})
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	if err := h.users.Delete(c.GetString("user_id"), c.Param("id")); err != nil {
		adminError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "usuario eliminado correctamente"})
}

func positiveQuery(c *gin.Context, key string, fallback int) int {
	value, err := strconv.Atoi(c.DefaultQuery(key, strconv.Itoa(fallback)))
	if err != nil || value < 1 {
		return fallback
	}
	return value
}

func adminError(c *gin.Context, err error) {
	if strings.Contains(err.Error(), "no encontrado") {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if strings.Contains(err.Error(), "Ã³rdenes") || strings.Contains(err.Error(), "ordenes") {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}
