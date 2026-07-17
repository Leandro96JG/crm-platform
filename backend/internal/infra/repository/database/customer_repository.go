package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/jmoiron/sqlx"
)

type customerRepository struct {
	client *sqlx.DB
}

func NewCustomerRepository(client *sqlx.DB) domain.CustomerRepository {
	return &customerRepository{client: client}
}

func (r *customerRepository) Create(ctx context.Context, customer domain.Customer) (string, error) {
	customerDTO := mapCustomerToCustomerDTO(customer)

	_, err := executor(ctx, r.client).NamedExecContext(
		ctx,
		`INSERT INTO customers
		(customer_id, name, phone, email, document, address, notes, is_active, created_by, created_at, updated_by, updated_at)
		VALUES
		(:customer_id, :name, :phone, :email, :document, :address, :notes, :is_active, :created_by, :created_at, :updated_by, :updated_at)`,
		customerDTO,
	)
	if err != nil {
		return "", err
	}

	return customer.CustomerID, nil
}

func (r *customerRepository) GetByID(ctx context.Context, customerID string) (*domain.Customer, error) {
	if customerID == "" {
		return nil, domain.NewValidationError("customerID is required", nil)
	}

	var dto CustomerDTO
	err := executor(ctx, r.client).GetContext(ctx, &dto, "SELECT * FROM customers WHERE customer_id=$1", customerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("no customer found with this id", map[string]any{"customer_id": customerID})
		}
		return nil, err
	}

	customer := mapCustomerDTOToCustomer(dto)
	return &customer, nil
}

func (r *customerRepository) Search(ctx context.Context, filters domain.CustomerFilters) (domain.PagingResult[domain.Customer], error) {
	whereQuery := []string{"1=1"}
	whereArgs := make([]any, 0)

	whereQuery, whereArgs = prepareInQuery(filters.CustomerID, whereQuery, whereArgs, "customer_id")

	if filters.Search != nil && *filters.Search != "" {
		idx := len(whereArgs) + 1
		whereQuery = append(whereQuery, fmt.Sprintf(
			"(name ILIKE '%%' || $%d::text || '%%' OR phone ILIKE '%%' || $%d::text || '%%' OR email ILIKE '%%' || $%d::text || '%%')",
			idx, idx, idx,
		))
		whereArgs = append(whereArgs, *filters.Search)
	}

	if filters.IsActive != nil {
		whereQuery = append(whereQuery, fmt.Sprintf("is_active = $%d", len(whereArgs)+1))
		whereArgs = append(whereArgs, *filters.IsActive)
	}

	limitQuery := fmt.Sprintf("LIMIT $%d OFFSET $%d", len(whereArgs)+1, len(whereArgs)+2)
	limitArgs := append(whereArgs, filters.Limit, filters.Offset)

	orderBy := buildOrderBy(filters.SortBy, filters.SortOrder, validCustomerSortColumns)

	query := fmt.Sprintf("SELECT * FROM customers WHERE %s ORDER BY %s %s", joinWhere(whereQuery), orderBy, limitQuery)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM customers WHERE %s", joinWhere(whereQuery))

	var customersDTO []CustomerDTO
	err := executor(ctx, r.client).SelectContext(ctx, &customersDTO, query, limitArgs...)
	if err != nil {
		return domain.PagingResult[domain.Customer]{}, err
	}

	var count int
	err = executor(ctx, r.client).GetContext(ctx, &count, countQuery, whereArgs...)
	if err != nil {
		return domain.PagingResult[domain.Customer]{}, err
	}

	return domain.PagingResult[domain.Customer]{
		Result: mapCustomerDTOsToCustomers(customersDTO),
		Paging: domain.Paging{
			Total:  count,
			Limit:  filters.Limit,
			Offset: filters.Offset,
		},
	}, nil
}

func (r *customerRepository) Update(ctx context.Context, customer domain.Customer) error {
	dto := mapCustomerToCustomerDTO(customer)

	_, err := executor(ctx, r.client).NamedExecContext(
		ctx,
		`UPDATE customers SET
		name = :name, phone = :phone, email = :email, document = :document,
		address = :address, notes = :notes, is_active = :is_active,
		updated_by = :updated_by, updated_at = :updated_at
		WHERE customer_id = :customer_id`,
		dto,
	)

	return err
}

func (r *customerRepository) Delete(ctx context.Context, customerID string) error {
	if customerID == "" {
		return domain.NewValidationError("customerID is required", nil)
	}

	_, err := executor(ctx, r.client).ExecContext(ctx, "DELETE FROM customers WHERE customer_id=$1", customerID)
	return err
}

var validCustomerSortColumns = map[string]bool{
	"created_at": true,
	"updated_at": true,
	"name":       true,
}
