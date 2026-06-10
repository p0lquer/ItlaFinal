package dto

import "ITLAFINAL/domain/models"

type CreateCustomerRequest struct {
	ID    string `json:"id"`
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type UpdateCustomerRequest struct {
	models.Customer
}

type CustomerResponse struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Phone *string `json:"phone,omitempty"`
	Email *string `json:"email,omitempty"`
}
