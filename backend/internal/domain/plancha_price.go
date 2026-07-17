package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type PlanchaPriceRepository interface {
	Create(ctx context.Context, price PlanchaPrice) (string, error)
	GetByID(ctx context.Context, priceID string) (*PlanchaPrice, error)
	GetByPlanchaAndMaterial(ctx context.Context, planchaID string, materialID string) (*PlanchaPrice, error)
	SearchByPlancha(ctx context.Context, planchaID string) ([]PlanchaPrice, error)
	Update(ctx context.Context, price PlanchaPrice) error
	Delete(ctx context.Context, priceID string) error
}

type BulkDiscount struct {
	MinQuantity int     `json:"min_quantity"`
	UnitPrice   float64 `json:"unit_price"`
}

type PlanchaPrice struct {
	PriceID      string
	PlanchaID    string
	MaterialID   string
	BasePrice    float64
	MinQuantity  int
	BulkDiscount []BulkDiscount
	IsActive     bool
	CreatedAt    time.Time
}

func NewPlanchaPrice(
	planchaID string,
	materialID string,
	basePrice float64,
	minQuantity int,
	bulkDiscount []BulkDiscount,
) (PlanchaPrice, error) {
	priceID, err := uuid.NewUUID()
	if err != nil {
		return PlanchaPrice{}, err
	}

	if bulkDiscount == nil {
		bulkDiscount = []BulkDiscount{}
	}

	return PlanchaPrice{
		PriceID:      priceID.String(),
		PlanchaID:    planchaID,
		MaterialID:   materialID,
		BasePrice:    basePrice,
		MinQuantity:  minQuantity,
		BulkDiscount: bulkDiscount,
		IsActive:     true,
		CreatedAt:    time.Now().UTC(),
	}, nil
}

func (p PlanchaPrice) CalculatePrice(quantity int) float64 {
	if quantity <= p.MinQuantity {
		return p.BasePrice * float64(quantity)
	}

	bestPrice := p.BasePrice * float64(quantity)
	for _, discount := range p.BulkDiscount {
		if quantity >= discount.MinQuantity {
			total := discount.UnitPrice * float64(quantity)
			if total < bestPrice {
				bestPrice = total
			}
		}
	}
	return bestPrice
}

func (p PlanchaPrice) MarshalBulkDiscount() ([]byte, error) {
	return json.Marshal(p.BulkDiscount)
}

func (p *PlanchaPrice) UnmarshalBulkDiscount(data []byte) error {
	return json.Unmarshal(data, &p.BulkDiscount)
}
