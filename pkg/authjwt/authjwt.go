// Package authjwt centralizes JWT creation and validation for HTTP and WebSocket entry points.
package authjwt

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	secretEnv        = "JWT_SECRET"
	minimumSecretLen = 32
)

var (
	ErrSecretNotConfigured = errors.New("JWT_SECRET no está configurado")
	ErrInvalidClaims       = errors.New("claims JWT inválidos")
)

// Claims is strict: every application token must identify a user and role.
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func SecretFromEnv() ([]byte, error) {
	secret := strings.TrimSpace(os.Getenv(secretEnv))
	if secret == "" {
		return nil, ErrSecretNotConfigured
	}
	if len(secret) < minimumSecretLen {
		return nil, errors.New("JWT_SECRET debe tener al menos 32 caracteres")
	}
	return []byte(secret), nil
}

func NewToken(userID, email, role string, now time.Time) (string, error) {
	secret, err := SecretFromEnv()
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(userID) == "" || strings.TrimSpace(email) == "" || !validRole(role) {
		return "", ErrInvalidClaims
	}
	claims := Claims{UserID: userID, Email: email, Role: role, RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(now),
	}}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

func Parse(raw string) (*Claims, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, errors.New("token vacío")
	}
	secret, err := SecretFromEnv()
	if err != nil {
		return nil, err
	}
	claims := new(Claims)
	token, err := jwt.ParseWithClaims(raw, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("algoritmo de firma no permitido: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil || token == nil || !token.Valid {
		if err == nil {
			err = errors.New("token inválido")
		}
		return nil, err
	}
	if strings.TrimSpace(claims.UserID) == "" || strings.TrimSpace(claims.Email) == "" || !validRole(claims.Role) {
		return nil, ErrInvalidClaims
	}
	return claims, nil
}

func validRole(role string) bool { return role == "customer" || role == "operator" }
