package userUseCases

import (
	"ITLAFINAL/domain/models"
	"testing"
)

type adminUserRepoStub struct {
	users     map[string]*models.User
	lastID    string
	lastState bool
}

func (s *adminUserRepoStub) Create(*models.User) error                { return nil }
func (s *adminUserRepoStub) FindByEmail(string) (*models.User, error) { return nil, nil }
func (s *adminUserRepoStub) Delete(string) error                      { return nil }
func (s *adminUserRepoStub) List(models.UserFilter) (*models.UserPage, error) {
	return &models.UserPage{}, nil
}
func (s *adminUserRepoStub) EnsureBootstrapAdmin(*models.User) error  { return nil }
func (s *adminUserRepoStub) FindByID(id string) (*models.User, error) { return s.users[id], nil }
func (s *adminUserRepoStub) SetActive(id string, active bool) error {
	s.lastID, s.lastState = id, active
	return nil
}
func (s *adminUserRepoStub) DeleteManaged(id string) error {
	s.lastID = id
	return nil
}

func TestAdminUsersCannotBlockThemselvesOrAdmins(t *testing.T) {
	repo := &adminUserRepoStub{users: map[string]*models.User{
		"admin": {ID: "admin", Role: models.RoleAdmin},
		"user":  {ID: "user", Role: models.RoleCustomer},
	}}
	uc := NewAdminUsersUseCase(repo)
	if err := uc.SetActive("admin", "admin", false); err == nil {
		t.Fatal("expected self-block to be rejected")
	}
	if err := uc.SetActive("other", "admin", false); err == nil {
		t.Fatal("expected admin target to be protected")
	}
	if err := uc.SetActive("admin", "user", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.lastID != "user" || repo.lastState {
		t.Fatal("expected target user to be blocked")
	}
}

func TestAdminUsersCannotDeleteThemselves(t *testing.T) {
	repo := &adminUserRepoStub{users: map[string]*models.User{
		"admin": {ID: "admin", Role: models.RoleAdmin},
		"user":  {ID: "user", Role: models.RoleOperator},
	}}
	uc := NewAdminUsersUseCase(repo)
	if err := uc.Delete("admin", "admin"); err == nil {
		t.Fatal("expected self-delete to be rejected")
	}
	if err := uc.Delete("admin", "user"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.lastID != "user" {
		t.Fatal("expected managed deletion for target user")
	}
}
