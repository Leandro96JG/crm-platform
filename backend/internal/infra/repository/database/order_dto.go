package database

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type OrderDTO struct {
	OrderID     string     `db:"order_id"`
	OrderNumber string     `db:"order_number"`
	CustomerID  string     `db:"customer_id"`
	Status      string     `db:"status"`
	Source      string     `db:"source"`
	AiAgentID   string     `db:"ai_agent_id"`
	AiHandled   bool       `db:"ai_handled"`
	AssignedTo  string     `db:"assigned_to"`
	Notes       string     `db:"notes"`
	Total       float64    `db:"total"`
	Urgency     string     `db:"urgency"`
	CompletedAt *time.Time `db:"completed_at"`
	CreatedBy   string     `db:"created_by"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedBy   string     `db:"updated_by"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

func mapOrderToDTO(o domain.Order) OrderDTO {
	return OrderDTO{
		OrderID:     o.OrderID,
		OrderNumber: o.OrderNumber,
		CustomerID:  o.CustomerID,
		Status:      string(o.Status),
		Source:      o.Source,
		AiAgentID:   o.AiAgentID,
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
	}
}

func mapDTOToOrder(dto OrderDTO) domain.Order {
	return domain.Order{
		OrderID:     dto.OrderID,
		OrderNumber: dto.OrderNumber,
		CustomerID:  dto.CustomerID,
		Status:      domain.OrderStatus(dto.Status),
		Source:      dto.Source,
		AiAgentID:   dto.AiAgentID,
		AiHandled:   dto.AiHandled,
		AssignedTo:  dto.AssignedTo,
		Notes:       dto.Notes,
		Total:       dto.Total,
		Urgency:     dto.Urgency,
		CompletedAt: dto.CompletedAt,
		CreatedBy:   dto.CreatedBy,
		CreatedAt:   dto.CreatedAt,
		UpdatedBy:   dto.UpdatedBy,
		UpdatedAt:   dto.UpdatedAt,
	}
}

func mapDTOsToOrders(dtos []OrderDTO) []domain.Order {
	orders := make([]domain.Order, 0, len(dtos))
	for _, dto := range dtos {
		orders = append(orders, mapDTOToOrder(dto))
	}
	return orders
}
