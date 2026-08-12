package repository

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"database/sql"
)

type ServiceTypeRepository struct {
	db *sql.DB
}

var _ ports.ServiceTypeRepository = (*ServiceTypeRepository)(nil)

func NewServiceTypeRepository(db *sql.DB) *ServiceTypeRepository {
	return &ServiceTypeRepository{db: db}
}

func (r *ServiceTypeRepository) Create(serviceType *models.ServiceType) error {
	err := r.db.QueryRow("INSERT INTO service_types (name, description, base_price, price_per_weight, price_per_piece, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		serviceType.Name, serviceType.Description, serviceType.BasePrice, serviceType.PricePerWeight, serviceType.PricePerPiece, serviceType.CreatedAt,
	).Scan(&serviceType.ID)
	return err
}

func (r *ServiceTypeRepository) FindAll() ([]*models.ServiceType, error) {
	rows, err := r.db.Query("SELECT id, name, description, base_price, price_per_weight, price_per_piece, created_at FROM service_types ORDER BY name ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.ServiceType
	for rows.Next() {
		var serviceType models.ServiceType
		if err := rows.Scan(&serviceType.ID, &serviceType.Name, &serviceType.Description, &serviceType.BasePrice, &serviceType.PricePerWeight, &serviceType.PricePerPiece, &serviceType.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, &serviceType)
	}
	return result, nil
}

func (r *ServiceTypeRepository) FindByName(name string) (*models.ServiceType, error) {
	var serviceType models.ServiceType
	err := r.db.QueryRow("SELECT id, name, description, base_price, price_per_weight, price_per_piece, created_at FROM service_types WHERE lower(name) = lower($1)", name).Scan(&serviceType.ID, &serviceType.Name, &serviceType.Description, &serviceType.BasePrice, &serviceType.PricePerWeight, &serviceType.PricePerPiece, &serviceType.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &serviceType, nil

}

func (r *ServiceTypeRepository) Update(serviceType *models.ServiceType) error {
	_, err := r.db.Exec("UPDATE service_types SET name = $1, description = $2, base_price = $3, price_per_weight = $4, price_per_piece = $5 WHERE id = $6",
		serviceType.Name, serviceType.Description, serviceType.BasePrice, serviceType.PricePerWeight, serviceType.PricePerPiece, serviceType.ID,
	)
	return err
}

func (r *ServiceTypeRepository) Delete(name string) error {
	_, err := r.db.Exec("DELETE FROM service_types WHERE name = $1", name)
	return err
}
