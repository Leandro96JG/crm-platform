package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(ctx context.Context, order Order) (string, error)
	GetByID(ctx context.Context, orderID string) (*Order, error)
	Search(ctx context.Context, filters OrderFilters) (PagingResult[Order], error)
	Update(ctx context.Context, order Order) error
}

type OrderStatus string

const (
	OrderPending     OrderStatus = "pending"
	OrderApproved    OrderStatus = "approved"
	OrderInProduction OrderStatus = "in_production"
	OrderReady       OrderStatus = "ready"
	OrderDelivered   OrderStatus = "delivered"
	OrderCancelled   OrderStatus = "cancelled"
)

type Order struct {
	OrderID     string
	OrderNumber string
	CustomerID  string
	Status      OrderStatus
	Source      string
	AiAgentID   string
	AiHandled   bool
	AssignedTo  string
	Notes       string
	Total       float64
	Urgency     string
	CompletedAt *time.Time
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedBy   string
	UpdatedAt   time.Time
	Items       []OrderItem
}

type OrderFilters struct {
	OrderID    []string
	CustomerID []string
	Status     []string
	Source     []string
	Urgency    []string
	StartDate  *string
	EndDate    *string
	PagingFilter
}

func NewOrder(
	customerID string,
	source string,
	notes string,
	urgency string,
	createdBy string,
	items []OrderItem,
) (Order, error) {
	now := time.Now().UTC()
	orderID, err := uuid.NewUUID()
	if err != nil {
		return Order{}, err
	}

	orderNumber := fmt.Sprintf("ORD-%s-%04d", now.Format("20060102"), now.UnixMilli()%10000)

	total := 0.0
	for _, item := range items {
		total += item.Subtotal
	}

	return Order{
		OrderID:     orderID.String(),
		OrderNumber: orderNumber,
		CustomerID:  customerID,
		Status:      OrderPending,
		Source:      source,
		AiHandled:   false,
		Notes:       notes,
		Total:       total,
		Urgency:     urgency,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedBy:   createdBy,
		UpdatedAt:   now,
		Items:       items,
	}, nil
}

type UpdateOrder struct {
	Status     *OrderStatus
	AssignedTo *string
	Notes      *string
	Urgency    *string
	Total      *float64
	CompletedAt *time.Time
	UpdatedBy  string
}

func (o *Order) MergeUpdate(update UpdateOrder) {
	o.UpdatedAt = time.Now().UTC()
	o.UpdatedBy = update.UpdatedBy

	if update.Status != nil {
		o.Status = *update.Status
	}
	if update.AssignedTo != nil {
		o.AssignedTo = *update.AssignedTo
	}
	if update.Notes != nil {
		o.Notes = *update.Notes
	}
	if update.Urgency != nil {
		o.Urgency = *update.Urgency
	}
	if update.Total != nil {
		o.Total = *update.Total
	}
	if update.CompletedAt != nil {
		o.CompletedAt = update.CompletedAt
	}
}
