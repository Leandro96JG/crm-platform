package rest

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type CreateOrderDTO struct {
	CustomerID string          `json:"customer_id" validate:"required"`
	Source     string          `json:"source"`
	Notes      string          `json:"notes"`
	Urgency    string          `json:"urgency"`
	CreatedBy  string          `json:"created_by" validate:"required"`
	Items      []CreateOrderItemDTO `json:"items" validate:"required,min=1"`
}

type CreateOrderItemDTO struct {
	PlanchaID         string  `json:"plancha_id" validate:"required"`
	MaterialID        string  `json:"material_id" validate:"required"`
	SheetQuantity     int     `json:"sheet_quantity" validate:"required,min=1"`
	UnitPrice         float64 `json:"unit_price" validate:"required"`
	CustomDesignFile  string  `json:"custom_design_file"`
	CustomDesignNotes string  `json:"custom_design_notes"`
}

type OrderDTO struct {
	OrderID     string       `json:"order_id"`
	OrderNumber string       `json:"order_number"`
	CustomerID  string       `json:"customer_id"`
	Status      string       `json:"status"`
	Source      string       `json:"source"`
	AiHandled   bool         `json:"ai_handled"`
	AssignedTo  string       `json:"assigned_to"`
	Notes       string       `json:"notes"`
	Total       float64      `json:"total"`
	Urgency     string       `json:"urgency"`
	CompletedAt *time.Time   `json:"completed_at"`
	CreatedBy   string       `json:"created_by"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedBy   string       `json:"updated_by"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Items       []OrderItemDTO `json:"items"`
}

type OrderItemDTO struct {
	ItemID            string  `json:"item_id"`
	PlanchaID         string  `json:"plancha_id"`
	MaterialID        string  `json:"material_id"`
	SheetQuantity     int     `json:"sheet_quantity"`
	UnitPrice         float64 `json:"unit_price"`
	Subtotal          float64 `json:"subtotal"`
	CustomDesignFile  string  `json:"custom_design_file"`
	CustomDesignNotes string  `json:"custom_design_notes"`
}

type UpdateOrderStatusDTO struct {
	Status    string `json:"status" validate:"required"`
	UpdatedBy string `json:"updated_by" validate:"required"`
}

func mapOrderToOrderDTO(o domain.Order) OrderDTO {
	items := make([]OrderItemDTO, len(o.Items))
	for i, item := range o.Items {
		items[i] = OrderItemDTO{
			ItemID:            item.ItemID,
			PlanchaID:         item.PlanchaID,
			MaterialID:        item.MaterialID,
			SheetQuantity:     item.SheetQuantity,
			UnitPrice:         item.UnitPrice,
			Subtotal:          item.Subtotal,
			CustomDesignFile:  item.CustomDesignFile,
			CustomDesignNotes: item.CustomDesignNotes,
		}
	}

	return OrderDTO{
		OrderID:     o.OrderID,
		OrderNumber: o.OrderNumber,
		CustomerID:  o.CustomerID,
		Status:      string(o.Status),
		Source:      o.Source,
		AiHandled:   o.AiHandled,
		AssignedTo:  o.AssignedTo,
		Notes:       o.Notes,
		Total:       o.Total,
		Urgency:     o.Urgency,
		CompletedAt: o.CompletedAt,
		CreatedBy:   o.CreatedBy,
		CreatedAt:   o.CreatedAt,
		UpdatedBy:   o.UpdatedBy,
		UpdatedAt:   o.UpdatedAt,
		Items:       items,
	}
}

func mapSearchResultToOrderDTOs(result domain.PagingResult[domain.Order]) ([]OrderDTO, domain.Paging) {
	dtos := make([]OrderDTO, len(result.Result))
	for i, o := range result.Result {
		dtos[i] = mapOrderToOrderDTO(o)
	}
	return dtos, result.Paging
}
