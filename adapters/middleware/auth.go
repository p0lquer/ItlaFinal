package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token requerido"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "formato inválido: Bearer <token>"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido o expirado"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_id", claims["user_id"].(string))
		c.Set("email", claims["email"].(string))
		c.Set("role", claims["role"].(string))
		c.Next()
	}
}

// OperatorOnly — solo operadores
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

// CustomerOnly — solo customers
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
