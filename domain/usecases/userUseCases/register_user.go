package userUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"errors"
	"os"
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
		CreatedAt: time.Now(),
	}

	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	if role == models.RoleCustomer && uc.customerRepo != nil {
		customer := &models.Customer{
			ID:   user.ID,
			Name: name,
		}
		if phone != "" {
			customer.Phone = &phone
		}
		if email != "" {
			customer.Email = &email
		}

		if err := uc.customerRepo.Create(customer); err != nil {
			return nil, err
		}
	}

	return user, nil
}
