package application

import (
	"context"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type customerService struct {
	customerRepo domain.CustomerRepository
}

type CustomerService interface {
	CreateCustomer(ctx context.Context, customer domain.Customer) (string, error)
	GetCustomerByID(ctx context.Context, customerID string) (*domain.Customer, error)
	SearchCustomers(ctx context.Context, filters domain.CustomerFilters) (domain.PagingResult[domain.Customer], error)
	UpdateCustomer(ctx context.Context, customerID string, update domain.UpdateCustomer) error
	DeleteCustomer(ctx context.Context, customerID string) error
}

func NewCustomerService(customerRepo domain.CustomerRepository) CustomerService {
	return &customerService{
		customerRepo: customerRepo,
	}
}

func (s *customerService) CreateCustomer(ctx context.Context, customer domain.Customer) (string, error) {
	return s.customerRepo.Create(ctx, customer)
}

func (s *customerService) GetCustomerByID(ctx context.Context, customerID string) (*domain.Customer, error) {
	return s.customerRepo.GetByID(ctx, customerID)
}

func (s *customerService) SearchCustomers(ctx context.Context, filters domain.CustomerFilters) (domain.PagingResult[domain.Customer], error) {
	return s.customerRepo.Search(ctx, filters)
}

func (s *customerService) UpdateCustomer(ctx context.Context, customerID string, update domain.UpdateCustomer) error {
	customer, err := s.customerRepo.GetByID(ctx, customerID)
	if err != nil {
		return err
	}

	customer.MergeUpdate(update)
	return s.customerRepo.Update(ctx, *customer)
}

func (s *customerService) DeleteCustomer(ctx context.Context, customerID string) error {
	return s.customerRepo.Delete(ctx, customerID)
}
