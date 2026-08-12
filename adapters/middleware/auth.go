package middleware

import (
	"ITLAFINAL/domain/ports"
	"ITLAFINAL/pkg/authjwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "formato inválido: Bearer <token>"})
			c.Abort()
			return
		}

		claims, err := authjwt.Parse(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido o expirado"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// ActiveUserRequired revokes access immediately after an administrator blocks
// an account, including tokens issued before the block.
func ActiveUserRequired(userRepo ports.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := userRepo.FindByID(c.GetString("user_id"))
		if err != nil || user == nil || !user.IsActive {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "cuenta bloqueada o no disponible"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func OperatorOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("role") != "operator" {
			c.JSON(http.StatusForbidden, gin.H{"error": "acceso restringido a operadores"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func CustomerOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("role") != "customer" {
			c.JSON(http.StatusForbidden, gin.H{"error": "acceso restringido a clientes"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("role") != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "acceso restringido a administradores"})
			c.Abort()
			return
		}
		c.Next()
	}
}
