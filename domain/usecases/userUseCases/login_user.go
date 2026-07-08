package userUseCases

import (
	"ITLAFINAL/domain/ports"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginUserUseCase struct {
	userRepo ports.UserRepository
}

func NewLoginUserUseCase(userRepo ports.UserRepository) *LoginUserUseCase {
	return &LoginUserUseCase{userRepo: userRepo}
}

type LoginResult struct {
	Token string
	Name  string
	Email string
	Role  string
}

func (uc *LoginUserUseCase) Execute(email, password string) (*LoginResult, error) {
	// Buscar el usuario
	user, err := uc.userRepo.FindByEmail(email)
	if err != nil || user == nil {
		return nil, errors.New("credenciales inválidas")
	}

	// Comparar la contraseña con el hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	// Generar JWT
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // expira en 24 horas
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, errors.New("error al generar el token")
	}

	return &LoginResult{
		Token: tokenString,
		Name:  user.Name,
		Email: user.Email,
		Role:  string(user.Role),
	}, nil
}
