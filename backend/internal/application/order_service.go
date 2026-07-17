package application

import (
	"context"
	"fmt"
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type orderItemLoader interface {
	GetByOrderID(ctx context.Context, orderID string) ([]domain.OrderItem, error)
	CreateBatch(ctx context.Context, items []domain.OrderItem) error
}

type orderService struct {
	orderRepo        domain.OrderRepository
	planchaRepo      domain.PlanchaRepository
	priceRepo        domain.PlanchaPriceRepository
	printRepo        domain.PrintJobRepository
	orderItemLoader  orderItemLoader
}

type OrderService interface {
	CreateOrder(ctx context.Context, order domain.Order, items []domain.OrderItem) (string, error)
	GetOrderByID(ctx context.Context, orderID string) (*domain.Order, error)
	SearchOrders(ctx context.Context, filters domain.OrderFilters) (domain.PagingResult[domain.Order], error)
	UpdateOrderStatus(ctx context.Context, orderID string, status domain.OrderStatus, updatedBy string) error
}

func NewOrderService(
	orderRepo domain.OrderRepository,
	planchaRepo domain.PlanchaRepository,
	priceRepo domain.PlanchaPriceRepository,
	printRepo domain.PrintJobRepository,
	orderItemLoader orderItemLoader,
) OrderService {
	return &orderService{
		orderRepo:       orderRepo,
		planchaRepo:     planchaRepo,
		priceRepo:       priceRepo,
		printRepo:       printRepo,
		orderItemLoader: orderItemLoader,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, order domain.Order, items []domain.OrderItem) (string, error) {
	orderID, err := s.orderRepo.Create(ctx, order)
	if err != nil {
		return "", err
	}

	for i := range items {
		items[i].OrderID = orderID
		items[i].SortOrder = i
	}

	if err := s.orderItemLoader.CreateBatch(ctx, items); err != nil {
		return "", err
	}

	return orderID, nil
}

func (s *orderService) GetOrderByID(ctx context.Context, orderID string) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	items, err := s.orderItemLoader.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	order.Items = items

	return order, nil
}

func (s *orderService) SearchOrders(ctx context.Context, filters domain.OrderFilters) (domain.PagingResult[domain.Order], error) {
	return s.orderRepo.Search(ctx, filters)
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, orderID string, status domain.OrderStatus, updatedBy string) error {
	order, err := s.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	update := domain.UpdateOrder{
		Status:    &status,
		UpdatedBy: updatedBy,
	}

	if status == domain.OrderDelivered || status == domain.OrderCancelled {
		now := time.Now().UTC()
		update.CompletedAt = &now
	}

	order.MergeUpdate(update)

	if err := s.orderRepo.Update(ctx, *order); err != nil {
		return err
	}

	if status == domain.OrderApproved {
		if err := s.createPrintJobs(ctx, *order); err != nil {
			return fmt.Errorf("failed to create print jobs: %w", err)
		}
	}

	return nil
}

func (s *orderService) createPrintJobs(ctx context.Context, order domain.Order) error {
	for i, item := range order.Items {
		printJob, err := domain.NewPrintJob(
			item.ItemID,
			domain.PrintJobTypePrint,
			"",
			fmt.Sprintf("Print job for order %s item %d", order.OrderNumber, i+1),
			item.SheetQuantity,
			i,
			order.UpdatedBy,
		)
		if err != nil {
			return err
		}

		if _, err := s.printRepo.Create(ctx, printJob); err != nil {
			return err
		}

		cutJob, err := domain.NewPrintJob(
			item.ItemID,
			domain.PrintJobTypeCut,
			"",
			fmt.Sprintf("Cut job for order %s item %d", order.OrderNumber, i+1),
			item.SheetQuantity,
			i,
			order.UpdatedBy,
		)
		if err != nil {
			return err
		}

		if _, err := s.printRepo.Create(ctx, cutJob); err != nil {
			return err
		}
	}

	return nil
}
