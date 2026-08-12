package userUseCases

import (
	"ITLAFINAL/domain/ports"
	"ITLAFINAL/pkg/authjwt"
	"ITLAFINAL/pkg/inputvalidation"
	"errors"
	"time"

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
	normalizedEmail, err := inputvalidation.Email(email)
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}
	if err := inputvalidation.LoginPassword(password); err != nil {
		return nil, err
	}
	// Buscar el usuario
	user, err := uc.userRepo.FindByEmail(normalizedEmail)
	if err != nil || user == nil {
		return nil, errors.New("credenciales inválidas")
	}

	// Comparar la contraseña con el hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	if !user.IsActive {
		return nil, errors.New("cuenta bloqueada")
	}

	tokenString, err := authjwt.NewToken(user.ID, user.Email, string(user.Role), time.Now())
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
