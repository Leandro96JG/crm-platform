package rest

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type CreateCustomerDTO struct {
	Name      string `json:"name" validate:"required"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Document  string `json:"document"`
	Address   string `json:"address"`
	Notes     string `json:"notes"`
	CreatedBy string `json:"created_by" validate:"required"`
}

type UpdateCustomerDTO struct {
	Name      *string `json:"name"`
	Phone     *string `json:"phone"`
	Email     *string `json:"email"`
	Document  *string `json:"document"`
	Address   *string `json:"address"`
	Notes     *string `json:"notes"`
	IsActive  *bool   `json:"is_active"`
	UpdatedBy string  `json:"updated_by" validate:"required"`
}

type CustomerDTO struct {
	CustomerID string    `json:"customer_id"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	Document   string    `json:"document"`
	Address    string    `json:"address"`
	Notes      string    `json:"notes"`
	IsActive   bool      `json:"is_active"`
	CreatedBy  string    `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedBy  string    `json:"updated_by"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func mapCustomerToCustomerDTO(c domain.Customer) CustomerDTO {
	return CustomerDTO{
		CustomerID: c.CustomerID,
		Name:       c.Name,
		Phone:      c.Phone,
		Email:      c.Email,
		Document:   c.Document,
		Address:    c.Address,
		Notes:      c.Notes,
		IsActive:   c.IsActive,
		CreatedBy:  c.CreatedBy,
		CreatedAt:  c.CreatedAt,
		UpdatedBy:  c.UpdatedBy,
		UpdatedAt:  c.UpdatedAt,
	}
}

func mapCreateCustomerDTOToDomain(dto CreateCustomerDTO) (domain.Customer, error) {
	return domain.NewCustomer(
		dto.Name,
		dto.Phone,
		dto.Email,
		dto.Document,
		dto.Address,
		dto.Notes,
		dto.CreatedBy,
	)
}

func mapUpdateCustomerDTOToDomain(dto UpdateCustomerDTO) domain.UpdateCustomer {
	return domain.UpdateCustomer{
		Name:      dto.Name,
		Phone:     dto.Phone,
		Email:     dto.Email,
		Document:  dto.Document,
		Address:   dto.Address,
		Notes:     dto.Notes,
		IsActive:  dto.IsActive,
		UpdatedBy: dto.UpdatedBy,
	}
}
