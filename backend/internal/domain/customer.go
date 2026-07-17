package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer Customer) (string, error)
	GetByID(ctx context.Context, customerID string) (*Customer, error)
	Search(ctx context.Context, filters CustomerFilters) (PagingResult[Customer], error)
	Update(ctx context.Context, customer Customer) error
	Delete(ctx context.Context, customerID string) error
}

type Customer struct {
	CustomerID string
	Name       string
	Phone      string
	Email      string
	Document   string
	Address    string
	Notes      string
	IsActive   bool
	CreatedBy  string
	CreatedAt  time.Time
	UpdatedBy  string
	UpdatedAt  time.Time
}

type CustomerFilters struct {
	CustomerID []string
	Search     *string
	IsActive   *bool
	PagingFilter
}

func NewCustomer(
	name string,
	phone string,
	email string,
	document string,
	address string,
	notes string,
	createdBy string,
) (Customer, error) {
	if name == "" {
		return Customer{}, NewValidationError("customer name is required", nil)
	}

	now := time.Now().UTC()
	customerID, err := uuid.NewUUID()
	if err != nil {
		return Customer{}, err
	}

	return Customer{
		CustomerID: customerID.String(),
		Name:       name,
		Phone:      phone,
		Email:      email,
		Document:   document,
		Address:    address,
		Notes:      notes,
		IsActive:   true,
		CreatedBy:  createdBy,
		CreatedAt:  now,
		UpdatedBy:  createdBy,
		UpdatedAt:  now,
	}, nil
}

type UpdateCustomer struct {
	Name      *string
	Phone     *string
	Email     *string
	Document  *string
	Address   *string
	Notes     *string
	IsActive  *bool
	UpdatedBy string
}

func (c *Customer) MergeUpdate(update UpdateCustomer) {
	c.UpdatedAt = time.Now().UTC()
	c.UpdatedBy = update.UpdatedBy

	if update.Name != nil {
		c.Name = *update.Name
	}
	if update.Phone != nil {
		c.Phone = *update.Phone
	}
	if update.Email != nil {
		c.Email = *update.Email
	}
	if update.Document != nil {
		c.Document = *update.Document
	}
	if update.Address != nil {
		c.Address = *update.Address
	}
	if update.Notes != nil {
		c.Notes = *update.Notes
	}
	if update.IsActive != nil {
		c.IsActive = *update.IsActive
	}
}
