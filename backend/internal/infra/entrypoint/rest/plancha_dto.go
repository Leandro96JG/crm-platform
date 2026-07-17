package rest

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type CreatePlanchaDTO struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Category    string `json:"category" validate:"required"`
	Subcategory string `json:"subcategory"`
	LayoutURL   string `json:"layout_file_url"`
	PreviewURL  string `json:"preview_image_url"`
	Notes       string `json:"notes"`
	CreatedBy   string `json:"created_by" validate:"required"`
}

type UpdatePlanchaDTO struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	Subcategory *string `json:"subcategory"`
	LayoutURL   *string `json:"layout_file_url"`
	PreviewURL  *string `json:"preview_image_url"`
	Notes       *string `json:"notes"`
	IsActive    *bool   `json:"is_active"`
	UpdatedBy   string  `json:"updated_by" validate:"required"`
}

type PlanchaDTO struct {
	PlanchaID       string    `json:"plancha_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Category        string    `json:"category"`
	Subcategory     string    `json:"subcategory"`
	LayoutFileURL   string    `json:"layout_file_url"`
	PreviewImageURL string    `json:"preview_image_url"`
	Notes           string    `json:"notes"`
	IsActive        bool      `json:"is_active"`
	CreatedBy       string    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedBy       string    `json:"updated_by"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateMaterialDTO struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	MaterialType string `json:"material_type" validate:"required"`
	Finish      string `json:"finish"`
	IsCuttable  bool   `json:"is_cuttable"`
	IsPrintable bool   `json:"is_printable"`
	BaseCost    float64 `json:"base_cost"`
}

type StickerMaterialDTO struct {
	MaterialID   string    `json:"material_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	MaterialType string    `json:"material_type"`
	Finish       string    `json:"finish"`
	IsCuttable   bool      `json:"is_cuttable"`
	IsPrintable  bool      `json:"is_printable"`
	BaseCost     float64   `json:"base_cost"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreatePlanchaPriceDTO struct {
	PlanchaID    string             `json:"plancha_id" validate:"required"`
	MaterialID   string             `json:"material_id" validate:"required"`
	BasePrice    float64            `json:"base_price" validate:"required"`
	MinQuantity  int                `json:"min_quantity"`
	BulkDiscount []domain.BulkDiscount `json:"bulk_discount"`
}

type PlanchaPriceDTO struct {
	PriceID      string               `json:"price_id"`
	PlanchaID    string               `json:"plancha_id"`
	MaterialID   string               `json:"material_id"`
	BasePrice    float64              `json:"base_price"`
	MinQuantity  int                  `json:"min_quantity"`
	BulkDiscount []domain.BulkDiscount `json:"bulk_discount"`
	IsActive     bool                 `json:"is_active"`
}

func mapPlanchaToPlanchaDTO(p domain.Plancha) PlanchaDTO {
	return PlanchaDTO{
		PlanchaID:       p.PlanchaID,
		Name:            p.Name,
		Description:     p.Description,
		Category:        p.Category,
		Subcategory:     p.Subcategory,
		LayoutFileURL:   p.LayoutFileURL,
		PreviewImageURL: p.PreviewImageURL,
		Notes:           p.Notes,
		IsActive:        p.IsActive,
		CreatedBy:       p.CreatedBy,
		CreatedAt:       p.CreatedAt,
		UpdatedBy:       p.UpdatedBy,
		UpdatedAt:       p.UpdatedAt,
	}
}

func mapCreatePlanchaDTOToDomain(dto CreatePlanchaDTO) (domain.Plancha, error) {
	return domain.NewPlancha(
		dto.Name,
		dto.Description,
		dto.Category,
		dto.Subcategory,
		dto.LayoutURL,
		dto.Notes,
		dto.CreatedBy,
	)
}

func mapUpdatePlanchaDTOToDomain(dto UpdatePlanchaDTO) domain.UpdatePlancha {
	return domain.UpdatePlancha{
		Name:        dto.Name,
		Description: dto.Description,
		Category:    dto.Category,
		Subcategory: dto.Subcategory,
		LayoutFileURL:   dto.LayoutURL,
		PreviewImageURL: dto.PreviewURL,
		Notes:       dto.Notes,
		IsActive:    dto.IsActive,
		UpdatedBy:   dto.UpdatedBy,
	}
}

func mapMaterialToDTO(m domain.StickerMaterial) StickerMaterialDTO {
	return StickerMaterialDTO{
		MaterialID:   m.MaterialID,
		Name:         m.Name,
		Description:  m.Description,
		MaterialType: m.MaterialType,
		Finish:       m.Finish,
		IsCuttable:   m.IsCuttable,
		IsPrintable:  m.IsPrintable,
		BaseCost:     m.BaseCost,
		IsActive:     m.IsActive,
		CreatedAt:    m.CreatedAt,
	}
}

func mapPriceToDTO(p domain.PlanchaPrice) PlanchaPriceDTO {
	return PlanchaPriceDTO{
		PriceID:      p.PriceID,
		PlanchaID:    p.PlanchaID,
		MaterialID:   p.MaterialID,
		BasePrice:    p.BasePrice,
		MinQuantity:  p.MinQuantity,
		BulkDiscount: p.BulkDiscount,
		IsActive:     p.IsActive,
	}
}
