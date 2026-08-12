package userUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// BootstrapAdmin is deliberately server-only. Call it with environment values;
// no HTTP endpoint or registration field can assign the admin role.
func BootstrapAdmin(userRepo ports.UserRepository, name, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return userRepo.EnsureBootstrapAdmin(&models.User{
		ID:        uuid.NewString(),
		Name:      strings.TrimSpace(name),
		Email:     strings.ToLower(strings.TrimSpace(email)),
		Password:  string(hash),
		Role:      models.RoleAdmin,
		IsActive:  true,
		CreatedAt: time.Now(),
	})
}
