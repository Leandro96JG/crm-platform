package database

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type PlanchaDTO struct {
	PlanchaID       string    `db:"plancha_id"`
	Name            string    `db:"name"`
	Description     string    `db:"description"`
	Category        string    `db:"category"`
	Subcategory     string    `db:"subcategory"`
	LayoutFileURL   string    `db:"layout_file_url"`
	PreviewImageURL string    `db:"preview_image_url"`
	Notes           string    `db:"notes"`
	IsActive        bool      `db:"is_active"`
	CreatedBy       string    `db:"created_by"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedBy       string    `db:"updated_by"`
	UpdatedAt       time.Time `db:"updated_at"`
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

func mapPlanchaDTOToPlancha(dto PlanchaDTO) domain.Plancha {
	return domain.Plancha{
		PlanchaID:       dto.PlanchaID,
		Name:            dto.Name,
		Description:     dto.Description,
		Category:        dto.Category,
		Subcategory:     dto.Subcategory,
		LayoutFileURL:   dto.LayoutFileURL,
		PreviewImageURL: dto.PreviewImageURL,
		Notes:           dto.Notes,
		IsActive:        dto.IsActive,
		CreatedBy:       dto.CreatedBy,
		CreatedAt:       dto.CreatedAt,
		UpdatedBy:       dto.UpdatedBy,
		UpdatedAt:       dto.UpdatedAt,
	}
}

func mapPlanchaDTOsToPlanchas(dtos []PlanchaDTO) []domain.Plancha {
	planchas := make([]domain.Plancha, 0, len(dtos))
	for _, dto := range dtos {
		planchas = append(planchas, mapPlanchaDTOToPlancha(dto))
	}
	return planchas
}
