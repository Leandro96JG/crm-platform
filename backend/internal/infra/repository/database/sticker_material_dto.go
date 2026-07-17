package database

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type StickerMaterialDTO struct {
	MaterialID   string    `db:"material_id"`
	Name         string    `db:"name"`
	Description  string    `db:"description"`
	MaterialType string    `db:"material_type"`
	Finish       string    `db:"finish"`
	IsCuttable   bool      `db:"is_cuttable"`
	IsPrintable  bool      `db:"is_printable"`
	BaseCost     float64   `db:"base_cost"`
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
}

func mapStickerMaterialToDTO(m domain.StickerMaterial) StickerMaterialDTO {
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

func mapDTOToStickerMaterial(dto StickerMaterialDTO) domain.StickerMaterial {
	return domain.StickerMaterial{
		MaterialID:   dto.MaterialID,
		Name:         dto.Name,
		Description:  dto.Description,
		MaterialType: dto.MaterialType,
		Finish:       dto.Finish,
		IsCuttable:   dto.IsCuttable,
		IsPrintable:  dto.IsPrintable,
		BaseCost:     dto.BaseCost,
		IsActive:     dto.IsActive,
		CreatedAt:    dto.CreatedAt,
	}
}

func mapDTOsToStickerMaterials(dtos []StickerMaterialDTO) []domain.StickerMaterial {
	result := make([]domain.StickerMaterial, 0, len(dtos))
	for _, dto := range dtos {
		result = append(result, mapDTOToStickerMaterial(dto))
	}
	return result
}
