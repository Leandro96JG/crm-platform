package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type StickerMaterialRepository interface {
	Create(ctx context.Context, material StickerMaterial) (string, error)
	GetByID(ctx context.Context, materialID string) (*StickerMaterial, error)
	Search(ctx context.Context, filters StickerMaterialFilters) (PagingResult[StickerMaterial], error)
	Update(ctx context.Context, material StickerMaterial) error
	Delete(ctx context.Context, materialID string) error
}

type StickerMaterial struct {
	MaterialID   string
	Name         string
	Description  string
	MaterialType string
	Finish       string
	IsCuttable   bool
	IsPrintable  bool
	BaseCost     float64
	IsActive     bool
	CreatedAt    time.Time
}

type StickerMaterialFilters struct {
	MaterialType []string
	Finish       []string
	IsActive     *bool
	PagingFilter
}

func NewStickerMaterial(
	name string,
	description string,
	materialType string,
	finish string,
	isCuttable bool,
	isPrintable bool,
	baseCost float64,
) (StickerMaterial, error) {
	materialID, err := uuid.NewUUID()
	if err != nil {
		return StickerMaterial{}, err
	}

	return StickerMaterial{
		MaterialID:   materialID.String(),
		Name:         name,
		Description:  description,
		MaterialType: materialType,
		Finish:       finish,
		IsCuttable:   isCuttable,
		IsPrintable:  isPrintable,
		BaseCost:     baseCost,
		IsActive:     true,
		CreatedAt:    time.Now().UTC(),
	}, nil
}
