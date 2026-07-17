package rest

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/domain"
)

type OrderController struct {
	orderService application.OrderService
}

func NewOrderController(orderService application.OrderService) OrderController {
	return OrderController{
		orderService: orderService,
	}
}

func (ctrl *OrderController) CreateOrder(ctx *gin.Context) {
	var dto CreateOrderDTO
	if err := ctx.BindJSON(&dto); err != nil {
		ctx.Error(err)
		return
	}

	items := make([]domain.OrderItem, 0, len(dto.Items))
	total := 0.0
	for i, itemDTO := range dto.Items {
		subtotal := itemDTO.UnitPrice * float64(itemDTO.SheetQuantity)
		total += subtotal

		item, err := domain.NewOrderItem(
			"",
			itemDTO.PlanchaID,
			itemDTO.MaterialID,
			itemDTO.SheetQuantity,
			itemDTO.UnitPrice,
			itemDTO.CustomDesignFile,
			itemDTO.CustomDesignNotes,
			i,
		)
		if err != nil {
			ctx.Error(err)
			return
		}
		items = append(items, item)
	}

	order, err := domain.NewOrder(
		dto.CustomerID,
		dto.Source,
		dto.Notes,
		dto.Urgency,
		dto.CreatedBy,
		items,
	)
	if err != nil {
		ctx.Error(err)
		return
	}

	order.Total = total

	orderID, err := ctrl.orderService.CreateOrder(ctx.Request.Context(), order, items)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"order_id": orderID, "order_number": order.OrderNumber})
}

func (ctrl *OrderController) GetOrder(ctx *gin.Context) {
	orderID := ctx.Param("orderID")
	if orderID == "" {
		ctx.Error(domain.NewValidationError("order_id is required", nil))
		return
	}

	order, err := ctrl.orderService.GetOrderByID(ctx.Request.Context(), orderID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, mapOrderToOrderDTO(*order))
}

func (ctrl *OrderController) SearchOrders(ctx *gin.Context) {
	filters := ctrl.parseFilters(ctx)

	orders, err := ctrl.orderService.SearchOrders(ctx.Request.Context(), filters)
	if err != nil {
		ctx.Error(err)
		return
	}

	dtos, paging := mapSearchResultToOrderDTOs(orders)
	ctx.JSON(http.StatusOK, gin.H{
		"result": dtos,
		"paging": paging,
	})
}

func (ctrl *OrderController) UpdateOrderStatus(ctx *gin.Context) {
	orderID := ctx.Param("orderID")
	if orderID == "" {
		ctx.Error(domain.NewValidationError("order_id is required", nil))
		return
	}

	var dto UpdateOrderStatusDTO
	if err := ctx.BindJSON(&dto); err != nil {
		ctx.Error(err)
		return
	}

	err := ctrl.orderService.UpdateOrderStatus(
		ctx.Request.Context(),
		orderID,
		domain.OrderStatus(dto.Status),
		dto.UpdatedBy,
	)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (ctrl *OrderController) parseFilters(ctx *gin.Context) domain.OrderFilters {
	filters := domain.OrderFilters{
		PagingFilter: domain.PagingFilter{
			Limit:     10,
			Offset:    0,
			SortBy:    "created_at",
			SortOrder: "DESC",
		},
	}

	if status := ctx.QueryArray("status"); len(status) > 0 {
		filters.Status = status
	}
	if source := ctx.QueryArray("source"); len(source) > 0 {
		filters.Source = source
	}
	if urgency := ctx.QueryArray("urgency"); len(urgency) > 0 {
		filters.Urgency = urgency
	}
	if startDate := ctx.Query("start_date"); startDate != "" {
		filters.StartDate = &startDate
	}
	if endDate := ctx.Query("end_date"); endDate != "" {
		filters.EndDate = &endDate
	}
	if limit := ctx.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			filters.Limit = l
		}
	}
	if offset := ctx.Query("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			filters.Offset = o
		}
	}
	if sortBy := ctx.Query("sort_by"); sortBy != "" {
		filters.SortBy = sortBy
	}
	if sortOrder := strings.ToUpper(ctx.Query("sort_order")); sortOrder == "ASC" || sortOrder == "DESC" {
		filters.SortOrder = sortOrder
	}

	return filters
}
