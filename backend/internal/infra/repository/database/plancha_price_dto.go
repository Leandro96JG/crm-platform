package database

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type PlanchaPriceDTO struct {
	PriceID      string    `db:"price_id"`
	PlanchaID    string    `db:"plancha_id"`
	MaterialID   string    `db:"material_id"`
	BasePrice    float64   `db:"base_price"`
	MinQuantity  int       `db:"min_quantity"`
	BulkDiscount string    `db:"bulk_discount"`
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
}

func mapPlanchaPriceToDTO(p domain.PlanchaPrice) (PlanchaPriceDTO, error) {
	bulkJSON, err := p.MarshalBulkDiscount()
	if err != nil {
		return PlanchaPriceDTO{}, err
	}

	return PlanchaPriceDTO{
		PriceID:      p.PriceID,
		PlanchaID:    p.PlanchaID,
		MaterialID:   p.MaterialID,
		BasePrice:    p.BasePrice,
		MinQuantity:  p.MinQuantity,
		BulkDiscount: string(bulkJSON),
		IsActive:     p.IsActive,
		CreatedAt:    p.CreatedAt,
	}, nil
}

func mapDTOToPlanchaPrice(dto PlanchaPriceDTO) (domain.PlanchaPrice, error) {
	p := domain.PlanchaPrice{
		PriceID:     dto.PriceID,
		PlanchaID:   dto.PlanchaID,
		MaterialID:  dto.MaterialID,
		BasePrice:   dto.BasePrice,
		MinQuantity: dto.MinQuantity,
		IsActive:    dto.IsActive,
		CreatedAt:   dto.CreatedAt,
	}

	if err := p.UnmarshalBulkDiscount([]byte(dto.BulkDiscount)); err != nil {
		return domain.PlanchaPrice{}, err
	}

	return p, nil
}

func mapDTOsToPlanchaPrices(dtos []PlanchaPriceDTO) ([]domain.PlanchaPrice, error) {
	prices := make([]domain.PlanchaPrice, 0, len(dtos))
	for _, dto := range dtos {
		p, err := mapDTOToPlanchaPrice(dto)
		if err != nil {
			return nil, err
		}
		prices = append(prices, p)
	}
	return prices, nil
}
