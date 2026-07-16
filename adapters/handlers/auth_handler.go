package handlers

import (
	"ITLAFINAL/adapters/dto"
	"ITLAFINAL/domain/usecases/userUseCases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	register *userUseCases.RegisterUserUseCase
	login    *userUseCases.LoginUserUseCase
	delete   *userUseCases.DeleteUserUseCase
}

func NewAuthHandler(
	register *userUseCases.RegisterUserUseCase,
	login *userUseCases.LoginUserUseCase,
	delete *userUseCases.DeleteUserUseCase,
) *AuthHandler {
	return &AuthHandler{register: register, login: login, delete: delete}
}

// Register godoc
// @Summary Registrar un usuario
// @Description Registra un nuevo usuario en la base de datos
// @Tags auth
// @Accept  json
// @Produce  json
// @Success 201 "Usuario registrado con éxito"
// @Router /auth/register [post]
// @Param user body dto.RegisterRequest true "User data"
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.register.Execute(req.Name, req.Email, req.Password, req.Phone, req.OperatorKey)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "usuario registrado correctamente",
		"id":      user.ID,
		"email":   user.Email,
		"role":    user.Role,
	})
}

// Login godoc
// @Summary Iniciar sesión
// @Description Autentica a un usuario y devuelve un token JWT
// @Tags auth
// @Accept  json
// @Produce  json
// @Success 200 "Inicio de sesión exitoso"
// @Router /auth/login [post]
// @Param credentials body dto.LoginRequest true "User credentials"
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.login.Execute(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{
		Token: result.Token,
		Name:  result.Name,
		Email: result.Email,
		Role:  result.Role,
	})
}

// Me — retorna el perfil del usuario autenticado usando el JWT
// @Summary Obtener perfil autenticado
// @Description Retorna los datos del usuario autenticado usando el JWT
// @Tags auth
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Success 200 "Perfil obtenido con éxito"
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"user_id": c.GetString("user_id"),
		"email":   c.GetString("email"),
		"role":    c.GetString("role"),
	})
}

// DeleteUser godoc
// @Summary Eliminar un usuario
// @Description Elimina un usuario y su cliente asociado si existe
// @Tags auth
// @Accept  json
// @Produce  json
// @Param id path string true "ID del usuario"
// @Success 200 "Usuario eliminado con éxito"
// @Router /users/{id} [delete]
func (h *AuthHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if err := h.delete.Execute(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "usuario eliminado correctamente"})
}
