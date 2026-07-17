package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/jmoiron/sqlx"
)

type orderRepository struct {
	client *sqlx.DB
}

func NewOrderRepository(client *sqlx.DB) domain.OrderRepository {
	return &orderRepository{client: client}
}

func (r *orderRepository) Create(ctx context.Context, order domain.Order) (string, error) {
	dto := mapOrderToDTO(order)

	_, err := executor(ctx, r.client).NamedExecContext(
		ctx,
		`INSERT INTO orders
		(order_id, order_number, customer_id, status, source, ai_agent_id, ai_handled, assigned_to, notes, total, urgency, completed_at, created_by, created_at, updated_by, updated_at)
		VALUES
		(:order_id, :order_number, :customer_id, :status, :source, :ai_agent_id, :ai_handled, :assigned_to, :notes, :total, :urgency, :completed_at, :created_by, :created_at, :updated_by, :updated_at)`,
		dto,
	)
	if err != nil {
		return "", err
	}

	return order.OrderID, nil
}

func (r *orderRepository) GetByID(ctx context.Context, orderID string) (*domain.Order, error) {
	if orderID == "" {
		return nil, domain.NewValidationError("orderID is required", nil)
	}

	var dto OrderDTO
	err := executor(ctx, r.client).GetContext(ctx, &dto, "SELECT * FROM orders WHERE order_id=$1", orderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("no order found with this id", map[string]any{"order_id": orderID})
		}
		return nil, err
	}

	order := mapDTOToOrder(dto)
	return &order, nil
}

func (r *orderRepository) Search(ctx context.Context, filters domain.OrderFilters) (domain.PagingResult[domain.Order], error) {
	whereQuery := []string{"1=1"}
	whereArgs := make([]any, 0)

	whereQuery, whereArgs = prepareInQuery(filters.OrderID, whereQuery, whereArgs, "order_id")
	whereQuery, whereArgs = prepareInQuery(filters.CustomerID, whereQuery, whereArgs, "customer_id")
	whereQuery, whereArgs = prepareInQuery(filters.Status, whereQuery, whereArgs, "status")
	whereQuery, whereArgs = prepareInQuery(filters.Source, whereQuery, whereArgs, "source")
	whereQuery, whereArgs = prepareInQuery(filters.Urgency, whereQuery, whereArgs, "urgency")

	if filters.StartDate != nil {
		whereQuery = append(whereQuery, fmt.Sprintf("created_at >= $%d", len(whereArgs)+1))
		whereArgs = append(whereArgs, *filters.StartDate)
	}
	if filters.EndDate != nil {
		whereQuery = append(whereQuery, fmt.Sprintf("created_at <= $%d", len(whereArgs)+1))
		whereArgs = append(whereArgs, *filters.EndDate)
	}

	limitQuery := fmt.Sprintf("LIMIT $%d OFFSET $%d", len(whereArgs)+1, len(whereArgs)+2)
	limitArgs := append(whereArgs, filters.Limit, filters.Offset)

	orderBy := buildOrderBy(filters.SortBy, filters.SortOrder, validOrderSortColumns)

	query := fmt.Sprintf("SELECT * FROM orders WHERE %s ORDER BY %s %s", joinWhere(whereQuery), orderBy, limitQuery)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM orders WHERE %s", joinWhere(whereQuery))

	var dtos []OrderDTO
	err := executor(ctx, r.client).SelectContext(ctx, &dtos, query, limitArgs...)
	if err != nil {
		return domain.PagingResult[domain.Order]{}, err
	}

	var count int
	err = executor(ctx, r.client).GetContext(ctx, &count, countQuery, whereArgs...)
	if err != nil {
		return domain.PagingResult[domain.Order]{}, err
	}

	return domain.PagingResult[domain.Order]{
		Result: mapDTOsToOrders(dtos),
		Paging: domain.Paging{Total: count, Limit: filters.Limit, Offset: filters.Offset},
	}, nil
}

func (r *orderRepository) Update(ctx context.Context, order domain.Order) error {
	dto := mapOrderToDTO(order)

	_, err := executor(ctx, r.client).NamedExecContext(
		ctx,
		`UPDATE orders SET
		status = :status, assigned_to = :assigned_to, notes = :notes, total = :total, urgency = :urgency,
		completed_at = :completed_at, updated_by = :updated_by, updated_at = :updated_at
		WHERE order_id = :order_id`,
		dto,
	)

	return err
}

var validOrderSortColumns = map[string]bool{
	"created_at":  true,
	"updated_at":  true,
	"order_number": true,
	"status":      true,
	"urgency":     true,
	"total":       true,
}
