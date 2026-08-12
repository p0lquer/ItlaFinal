package userUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"errors"
)

type AdminUsersUseCase struct {
	userRepo ports.UserRepository
}

func NewAdminUsersUseCase(userRepo ports.UserRepository) *AdminUsersUseCase {
	return &AdminUsersUseCase{userRepo: userRepo}
}

func (uc *AdminUsersUseCase) List(filter models.UserFilter) (*models.UserPage, error) {
	return uc.userRepo.List(filter)
}

func (uc *AdminUsersUseCase) SetActive(actorID, targetID string, active bool) error {
	target, err := uc.userRepo.FindByID(targetID)
	if err != nil {
		return err
	}
	if target == nil {
		return errors.New("usuario no encontrado")
	}
	if actorID == targetID {
		return errors.New("no puedes bloquear ni desbloquear tu propia cuenta")
	}
	if target.Role == models.RoleAdmin {
		return errors.New("las cuentas administradoras no pueden bloquearse desde este panel")
	}
	return uc.userRepo.SetActive(targetID, active)
}

func (uc *AdminUsersUseCase) Delete(actorID, targetID string) error {
	target, err := uc.userRepo.FindByID(targetID)
	if err != nil {
		return err
	}
	if target == nil {
		return errors.New("usuario no encontrado")
	}
	if actorID == targetID {
		return errors.New("no puedes eliminar tu propia cuenta")
	}
	if target.Role == models.RoleAdmin {
		return errors.New("las cuentas administradoras no pueden eliminarse desde este panel")
	}
	return uc.userRepo.DeleteManaged(targetID)
}
