package repository

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"database/sql"
)

type userRepositoryPG struct {
	db *sql.DB
}

var _ ports.UserRepository = (*userRepositoryPG)(nil)

func NewUserRepository(db *sql.DB) ports.UserRepository {
	return &userRepositoryPG{db: db}
}

func (r *userRepositoryPG) Create(user *models.User) error {
	_, err := r.db.Exec(
		`INSERT INTO users (id, name, email, password, role, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		user.ID, user.Name, user.Email, user.Password, user.Role, user.CreatedAt,
	)
	return err
}

func (r *userRepositoryPG) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(
		`SELECT id, name, email, password, role, created_at FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryPG) FindByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(
		`SELECT id, name, email, password, role, created_at FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryPG) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	return err
}
