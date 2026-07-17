package rest

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/domain"
)

type PlanchaController struct {
	planchaService application.PlanchaService
}

func NewPlanchaController(planchaService application.PlanchaService) PlanchaController {
	return PlanchaController{
		planchaService: planchaService,
	}
}

func (ctrl *PlanchaController) CreatePlancha(ctx *gin.Context) {
	var dto CreatePlanchaDTO
	if err := ctx.BindJSON(&dto); err != nil {
		ctx.Error(err)
		return
	}

	plancha, err := mapCreatePlanchaDTOToDomain(dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	planchaID, err := ctrl.planchaService.CreatePlancha(ctx.Request.Context(), plancha)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"plancha_id": planchaID})
}

func (ctrl *PlanchaController) GetPlancha(ctx *gin.Context) {
	planchaID := ctx.Param("planchaID")
	if planchaID == "" {
		ctx.Error(domain.NewValidationError("plancha_id is required", nil))
		return
	}

	plancha, err := ctrl.planchaService.GetPlanchaByID(ctx.Request.Context(), planchaID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, mapPlanchaToPlanchaDTO(*plancha))
}

func (ctrl *PlanchaController) SearchPlanchas(ctx *gin.Context) {
	filters := ctrl.parseFilters(ctx)

	planchas, err := ctrl.planchaService.SearchPlanchas(ctx.Request.Context(), filters)
	if err != nil {
		ctx.Error(err)
		return
	}

	dtos := make([]PlanchaDTO, len(planchas.Result))
	for i, p := range planchas.Result {
		dtos[i] = mapPlanchaToPlanchaDTO(p)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": dtos,
		"paging": planchas.Paging,
	})
}

func (ctrl *PlanchaController) UpdatePlancha(ctx *gin.Context) {
	planchaID := ctx.Param("planchaID")
	if planchaID == "" {
		ctx.Error(domain.NewValidationError("plancha_id is required", nil))
		return
	}

	var dto UpdatePlanchaDTO
	if err := ctx.BindJSON(&dto); err != nil {
		ctx.Error(err)
		return
	}

	update := mapUpdatePlanchaDTOToDomain(dto)
	err := ctrl.planchaService.UpdatePlancha(ctx.Request.Context(), planchaID, update)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (ctrl *PlanchaController) DeletePlancha(ctx *gin.Context) {
	planchaID := ctx.Param("planchaID")
	if planchaID == "" {
		ctx.Error(domain.NewValidationError("plancha_id is required", nil))
		return
	}

	err := ctrl.planchaService.DeletePlancha(ctx.Request.Context(), planchaID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (ctrl *PlanchaController) CreateMaterial(ctx *gin.Context) {
	var dto CreateMaterialDTO
	if err := ctx.BindJSON(&dto); err != nil {
		ctx.Error(err)
		return
	}

	material, err := domain.NewStickerMaterial(
		dto.Name, dto.Description, dto.MaterialType, dto.Finish,
		dto.IsCuttable, dto.IsPrintable, dto.BaseCost,
	)
	if err != nil {
		ctx.Error(err)
		return
	}

	materialID, err := ctrl.planchaService.CreateMaterial(ctx.Request.Context(), material)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"material_id": materialID})
}

func (ctrl *PlanchaController) SearchMaterials(ctx *gin.Context) {
	materials, err := ctrl.planchaService.SearchMaterials(ctx.Request.Context(), domain.StickerMaterialFilters{
		PagingFilter: domain.PagingFilter{Limit: 100, Offset: 0},
	})
	if err != nil {
		ctx.Error(err)
		return
	}

	dtos := make([]StickerMaterialDTO, len(materials.Result))
	for i, m := range materials.Result {
		dtos[i] = mapMaterialToDTO(m)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": dtos,
		"paging": materials.Paging,
	})
}

func (ctrl *PlanchaController) CreatePrice(ctx *gin.Context) {
	var dto CreatePlanchaPriceDTO
	if err := ctx.BindJSON(&dto); err != nil {
		ctx.Error(err)
		return
	}

	price, err := domain.NewPlanchaPrice(
		dto.PlanchaID, dto.MaterialID, dto.BasePrice, dto.MinQuantity, dto.BulkDiscount,
	)
	if err != nil {
		ctx.Error(err)
		return
	}

	priceID, err := ctrl.planchaService.CreatePrice(ctx.Request.Context(), price)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"price_id": priceID})
}

func (ctrl *PlanchaController) GetPlanchaPrices(ctx *gin.Context) {
	planchaID := ctx.Param("planchaID")
	if planchaID == "" {
		ctx.Error(domain.NewValidationError("plancha_id is required", nil))
		return
	}

	prices, err := ctrl.planchaService.GetPlanchaPrices(ctx.Request.Context(), planchaID)
	if err != nil {
		ctx.Error(err)
		return
	}

	dtos := make([]PlanchaPriceDTO, len(prices))
	for i, p := range prices {
		dtos[i] = mapPriceToDTO(p)
	}

	ctx.JSON(http.StatusOK, gin.H{"prices": dtos})
}

func (ctrl *PlanchaController) CalculatePrice(ctx *gin.Context) {
	planchaID := ctx.Param("planchaID")
	materialID := ctx.Query("material_id")
	quantityStr := ctx.Query("quantity")

	if planchaID == "" || materialID == "" || quantityStr == "" {
		ctx.Error(domain.NewValidationError("plancha_id, material_id, and quantity are required", nil))
		return
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		ctx.Error(domain.NewValidationError("quantity must be a number", nil))
		return
	}

	total, err := ctrl.planchaService.CalculatePrice(ctx.Request.Context(), planchaID, materialID, quantity)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"plancha_id":  planchaID,
		"material_id": materialID,
		"quantity":    quantity,
		"total":       total,
	})
}

func (ctrl *PlanchaController) parseFilters(ctx *gin.Context) domain.PlanchaFilters {
	filters := domain.PlanchaFilters{
		PagingFilter: domain.PagingFilter{
			Limit:     10,
			Offset:    0,
			SortBy:    "created_at",
			SortOrder: "DESC",
		},
	}

	if category := ctx.QueryArray("category"); len(category) > 0 {
		filters.Category = category
	}
	if subcategory := ctx.QueryArray("subcategory"); len(subcategory) > 0 {
		filters.Subcategory = subcategory
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
