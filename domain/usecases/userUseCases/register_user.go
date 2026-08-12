package userUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUserUseCase struct {
	userRepo     ports.UserRepository
	customerRepo ports.CustomerRepository
}

func NewRegisterUserUseCase(
	userRepo ports.UserRepository,
	customerRepo ports.CustomerRepository,
) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepo:     userRepo,
		customerRepo: customerRepo,
	}
}

func (uc *RegisterUserUseCase) Execute(
	name, email, password, phone, operatorKey string,
) (*models.User, error) {
	email = strings.ToLower(email)

	existing, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("el email ya está registrado")
	}

	role := models.RoleCustomer
	if operatorKey != "" && operatorKey == os.Getenv("OPERATOR_KEY") {
		role = models.RoleOperator
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error al procesar la contraseña")
	}
	user := &models.User{
		ID:        uuid.NewString(),
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword),
		Role:      role,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	if role == models.RoleCustomer && uc.customerRepo != nil {
		customer := &models.Customer{
			// Mantener la misma clave para el perfil autocreado permite que el
			// token del cliente identifique de forma directa sus órdenes. Los
			// clientes creados por un operador siguen recibiendo UUID del servidor.
			ID:     user.ID,
			UserID: &user.ID,
			Name:   name,
		}
		if phone != "" {
			customer.Phone = &phone
		}
		if email != "" {
			customer.Email = &email
		}

		if err := uc.customerRepo.Create(customer); err != nil {
			// Las interfaces actuales no exponen una transacción compartida. Se
			// compensa la creación de usuario para no dejar una cuenta sin perfil.
			_ = uc.userRepo.Delete(user.ID)
			return nil, err
		}
	}

	return user, nil
}
