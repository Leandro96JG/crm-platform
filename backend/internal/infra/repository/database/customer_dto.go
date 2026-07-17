package database

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type CustomerDTO struct {
	CustomerID string    `db:"customer_id"`
	Name       string    `db:"name"`
	Phone      string    `db:"phone"`
	Email      string    `db:"email"`
	Document   string    `db:"document"`
	Address    string    `db:"address"`
	Notes      string    `db:"notes"`
	IsActive   bool      `db:"is_active"`
	CreatedBy  string    `db:"created_by"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedBy  string    `db:"updated_by"`
	UpdatedAt  time.Time `db:"updated_at"`
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

func mapCustomerDTOToCustomer(dto CustomerDTO) domain.Customer {
	return domain.Customer{
		CustomerID: dto.CustomerID,
		Name:       dto.Name,
		Phone:      dto.Phone,
		Email:      dto.Email,
		Document:   dto.Document,
		Address:    dto.Address,
		Notes:      dto.Notes,
		IsActive:   dto.IsActive,
		CreatedBy:  dto.CreatedBy,
		CreatedAt:  dto.CreatedAt,
		UpdatedBy:  dto.UpdatedBy,
		UpdatedAt:  dto.UpdatedAt,
	}
}

func mapCustomerDTOsToCustomers(dtos []CustomerDTO) []domain.Customer {
	customers := make([]domain.Customer, 0, len(dtos))
	for _, dto := range dtos {
		customers = append(customers, mapCustomerDTOToCustomer(dto))
	}
	return customers
}
