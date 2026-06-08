package repository

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"database/sql"
)

type customerRepositoryPG struct {
	db *sql.DB
}

var _ ports.CustomerRepository = (*customerRepositoryPG)(nil)

func NewCustomerRepository(db *sql.DB) ports.CustomerRepository {
	return &customerRepositoryPG{db: db}
}

func (r *customerRepositoryPG) Create(customer *models.Customer) error {
	query := `
		INSERT INTO customers (id, name, email, phone, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query,
		customer.ID,
		customer.Name,
		customer.Email,
		customer.Phone,
		customer.CreatedAt,
		customer.UpdatedAt)
	return err
}

func (r *customerRepositoryPG) FindAll() ([]*models.Customer, error) {
	query := `SELECT id, name, email, phone, created_at, updated_at FROM customers`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []*models.Customer
	for rows.Next() {
		var customer models.Customer
		err := rows.Scan(
			&customer.ID,
			&customer.Name,
			&customer.Email,
			&customer.Phone,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		customers = append(customers, &customer)
	}
	return customers, nil
}

func (r *customerRepositoryPG) UpdateStatus(id string, status models.OrderStatus) error {
	// This method is not applicable for customers, so we can return an error or leave it unimplemented.
	return nil
}

func (r *customerRepositoryPG) Delete(id string) error {
	// This method is not applicable for customers, so we can return an error or leave it unimplemented.
	return nil
}
