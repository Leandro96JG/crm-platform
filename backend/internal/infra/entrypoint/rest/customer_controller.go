package rest

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/domain"
)

type CustomerController struct {
	customerService application.CustomerService
}

func NewCustomerController(customerService application.CustomerService) CustomerController {
	return CustomerController{
		customerService: customerService,
	}
}

func (ctrl *CustomerController) CreateCustomer(ctx *gin.Context) {
	var dto CreateCustomerDTO
	if err := ctx.BindJSON(&dto); err != nil {
		ctx.Error(err)
		return
	}

	customer, err := mapCreateCustomerDTOToDomain(dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	customerID, err := ctrl.customerService.CreateCustomer(ctx.Request.Context(), customer)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"customer_id": customerID})
}

func (ctrl *CustomerController) GetCustomer(ctx *gin.Context) {
	customerID := ctx.Param("customerID")
	if customerID == "" {
		ctx.Error(domain.NewValidationError("customer_id is required", nil))
		return
	}

	customer, err := ctrl.customerService.GetCustomerByID(ctx.Request.Context(), customerID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, mapCustomerToCustomerDTO(*customer))
}

func (ctrl *CustomerController) SearchCustomers(ctx *gin.Context) {
	filters := ctrl.parseFilters(ctx)

	customers, err := ctrl.customerService.SearchCustomers(ctx.Request.Context(), filters)
	if err != nil {
		ctx.Error(err)
		return
	}

	dtos := make([]CustomerDTO, len(customers.Result))
	for i, c := range customers.Result {
		dtos[i] = mapCustomerToCustomerDTO(c)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": dtos,
		"paging": customers.Paging,
	})
}

func (ctrl *CustomerController) UpdateCustomer(ctx *gin.Context) {
	customerID := ctx.Param("customerID")
	if customerID == "" {
		ctx.Error(domain.NewValidationError("customer_id is required", nil))
		return
	}

	var dto UpdateCustomerDTO
	if err := ctx.BindJSON(&dto); err != nil {
		ctx.Error(err)
		return
	}

	update := mapUpdateCustomerDTOToDomain(dto)
	err := ctrl.customerService.UpdateCustomer(ctx.Request.Context(), customerID, update)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (ctrl *CustomerController) DeleteCustomer(ctx *gin.Context) {
	customerID := ctx.Param("customerID")
	if customerID == "" {
		ctx.Error(domain.NewValidationError("customer_id is required", nil))
		return
	}

	err := ctrl.customerService.DeleteCustomer(ctx.Request.Context(), customerID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (ctrl *CustomerController) parseFilters(ctx *gin.Context) domain.CustomerFilters {
	filters := domain.CustomerFilters{
		PagingFilter: domain.PagingFilter{
			Limit:     10,
			Offset:    0,
			SortBy:    "created_at",
			SortOrder: "DESC",
		},
	}

	if search := ctx.Query("search"); search != "" {
		filters.Search = &search
	}
	if active := ctx.Query("is_active"); active != "" {
		isActive := active == "true"
		filters.IsActive = &isActive
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
