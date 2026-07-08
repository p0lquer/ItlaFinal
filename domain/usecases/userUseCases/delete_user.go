package userUseCases

import (
	"ITLAFINAL/domain/ports"
	"errors"
)

type DeleteUserUseCase struct {
	userRepo     ports.UserRepository
	customerRepo ports.CustomerRepository
}

func NewDeleteUserUseCase(userRepo ports.UserRepository, customerRepo ports.CustomerRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{userRepo: userRepo, customerRepo: customerRepo}
}

func (uc *DeleteUserUseCase) Execute(userID string) error {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("usuario no encontrado")
	}

	if uc.customerRepo != nil {
		if err := uc.customerRepo.Delete(userID); err != nil {
			return err
		}
	}

	return uc.userRepo.Delete(userID)
}
