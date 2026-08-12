package repository

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"database/sql"
	"errors"
	"strconv"
	"strings"
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
		`INSERT INTO users (id, name, email, password, role, is_active, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		user.ID, user.Name, user.Email, user.Password, user.Role, user.IsActive, user.CreatedAt,
	)
	return err
}

func (r *userRepositoryPG) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(
		`SELECT id, name, email, password, role, is_active, created_at FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.IsActive, &user.CreatedAt)
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
		`SELECT id, name, email, password, role, is_active, created_at FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.IsActive, &user.CreatedAt)
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

func (r *userRepositoryPG) List(filter models.UserFilter) (*models.UserPage, error) {
	where := make([]string, 0, 3)
	args := make([]any, 0, 5)
	if search := strings.TrimSpace(filter.Search); search != "" {
		args = append(args, "%"+search+"%")
		placeholder := "$" + strconv.Itoa(len(args))
		where = append(where, "(name ILIKE "+placeholder+" OR email ILIKE "+placeholder+")")
	}
	if filter.Role != "" {
		args = append(args, filter.Role)
		where = append(where, "role = $"+strconv.Itoa(len(args)))
	}
	if filter.IsActive != nil {
		args = append(args, *filter.IsActive)
		where = append(where, "is_active = $"+strconv.Itoa(len(args)))
	}
	clause := ""
	if len(where) > 0 {
		clause = " WHERE " + strings.Join(where, " AND ")
	}
	var total int
	if err := r.db.QueryRow("SELECT COUNT(*) FROM users"+clause, args...).Scan(&total); err != nil {
		return nil, err
	}
	args = append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)
	query := "SELECT id, name, email, password, role, is_active, created_at FROM users" + clause +
		" ORDER BY created_at DESC LIMIT $" + strconv.Itoa(len(args)-1) + " OFFSET $" + strconv.Itoa(len(args))
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]*models.User, 0)
	for rows.Next() {
		user := new(models.User)
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.IsActive, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &models.UserPage{Users: users, Total: total, Page: filter.Page, PageSize: filter.PageSize}, nil
}

func (r *userRepositoryPG) SetActive(id string, active bool) error {
	result, err := r.db.Exec(`UPDATE users SET is_active = $1 WHERE id = $2`, active, id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("usuario no encontrado")
	}
	return nil
}

// DeleteManaged removes an account only when its customer profile has no
// orders. Order history is never cascade-deleted by an administrative action.
func (r *userRepositoryPG) DeleteManaged(id string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var orders int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM orders o JOIN customers c ON c.id = o.customer_id WHERE c.user_id = $1`, id).Scan(&orders); err != nil {
		return err
	}
	if orders > 0 {
		return errors.New("no se puede eliminar un usuario con Ã³rdenes; bloquÃ©elo para conservar su historial")
	}
	if _, err := tx.Exec(`DELETE FROM customers WHERE user_id = $1`, id); err != nil {
		return err
	}
	result, err := tx.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("usuario no encontrado")
	}
	return tx.Commit()
}

func (r *userRepositoryPG) EnsureBootstrapAdmin(user *models.User) error {
	existing, err := r.FindByEmail(user.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		if existing.Role != models.RoleAdmin {
			return errors.New("ADMIN_EMAIL ya pertenece a un usuario que no es administrador")
		}
		return nil
	}
	return r.Create(user)
}
