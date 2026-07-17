package database

import (
	"context"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/jmoiron/sqlx"
)

type orderItemRepository struct {
	client *sqlx.DB
}

func NewOrderItemRepository(client *sqlx.DB) *orderItemRepository {
	return &orderItemRepository{client: client}
}

func (r *orderItemRepository) CreateBatch(ctx context.Context, items []domain.OrderItem) error {
	dtos := make([]OrderItemDTO, 0, len(items))
	for _, item := range items {
		dtos = append(dtos, mapOrderItemToDTO(item))
	}

	_, err := executor(ctx, r.client).NamedExecContext(
		ctx,
		`INSERT INTO order_items
		(item_id, order_id, plancha_id, material_id, sheet_quantity, unit_price, subtotal, custom_design_file, custom_design_notes, sort_order)
		VALUES
		(:item_id, :order_id, :plancha_id, :material_id, :sheet_quantity, :unit_price, :subtotal, :custom_design_file, :custom_design_notes, :sort_order)`,
		dtos,
	)

	return err
}

func (r *orderItemRepository) GetByOrderID(ctx context.Context, orderID string) ([]domain.OrderItem, error) {
	var dtos []OrderItemDTO
	err := executor(ctx, r.client).SelectContext(ctx, &dtos, "SELECT * FROM order_items WHERE order_id=$1 ORDER BY sort_order", orderID)
	if err != nil {
		return nil, err
	}

	return mapDTOsToOrderItems(dtos), nil
}
