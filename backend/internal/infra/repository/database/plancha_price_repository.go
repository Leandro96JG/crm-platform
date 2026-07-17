package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/jmoiron/sqlx"
)

type planchaPriceRepository struct {
	client *sqlx.DB
}

func NewPlanchaPriceRepository(client *sqlx.DB) domain.PlanchaPriceRepository {
	return &planchaPriceRepository{client: client}
}

func (r *planchaPriceRepository) Create(ctx context.Context, price domain.PlanchaPrice) (string, error) {
	dto, err := mapPlanchaPriceToDTO(price)
	if err != nil {
		return "", err
	}

	_, err = executor(ctx, r.client).NamedExecContext(
		ctx,
		`INSERT INTO plancha_prices
		(price_id, plancha_id, material_id, base_price, min_quantity, bulk_discount, is_active, created_at)
		VALUES
		(:price_id, :plancha_id, :material_id, :base_price, :min_quantity, :bulk_discount, :is_active, :created_at)`,
		dto,
	)
	if err != nil {
		return "", err
	}

	return price.PriceID, nil
}

func (r *planchaPriceRepository) GetByID(ctx context.Context, priceID string) (*domain.PlanchaPrice, error) {
	if priceID == "" {
		return nil, domain.NewValidationError("priceID is required", nil)
	}

	var dto PlanchaPriceDTO
	err := executor(ctx, r.client).GetContext(ctx, &dto, "SELECT * FROM plancha_prices WHERE price_id=$1", priceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("no price found with this id", map[string]any{"price_id": priceID})
		}
		return nil, err
	}

	price, err := mapDTOToPlanchaPrice(dto)
	if err != nil {
		return nil, err
	}

	return &price, nil
}

func (r *planchaPriceRepository) GetByPlanchaAndMaterial(ctx context.Context, planchaID string, materialID string) (*domain.PlanchaPrice, error) {
	var dto PlanchaPriceDTO
	err := executor(ctx, r.client).GetContext(ctx, &dto,
		"SELECT * FROM plancha_prices WHERE plancha_id=$1 AND material_id=$2 AND is_active=true",
		planchaID, materialID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("no price found", map[string]any{"plancha_id": planchaID, "material_id": materialID})
		}
		return nil, err
	}

	price, err := mapDTOToPlanchaPrice(dto)
	if err != nil {
		return nil, err
	}

	return &price, nil
}

func (r *planchaPriceRepository) SearchByPlancha(ctx context.Context, planchaID string) ([]domain.PlanchaPrice, error) {
	var dtos []PlanchaPriceDTO
	err := executor(ctx, r.client).SelectContext(ctx, &dtos, "SELECT * FROM plancha_prices WHERE plancha_id=$1", planchaID)
	if err != nil {
		return nil, err
	}

	return mapDTOsToPlanchaPrices(dtos)
}

func (r *planchaPriceRepository) Update(ctx context.Context, price domain.PlanchaPrice) error {
	dto, err := mapPlanchaPriceToDTO(price)
	if err != nil {
		return err
	}

	_, err = executor(ctx, r.client).NamedExecContext(
		ctx,
		`UPDATE plancha_prices SET
		base_price = :base_price, min_quantity = :min_quantity, bulk_discount = :bulk_discount, is_active = :is_active
		WHERE price_id = :price_id`,
		dto,
	)

	return err
}

func (r *planchaPriceRepository) Delete(ctx context.Context, priceID string) error {
	if priceID == "" {
		return domain.NewValidationError("priceID is required", nil)
	}

	_, err := executor(ctx, r.client).ExecContext(ctx, "DELETE FROM plancha_prices WHERE price_id=$1", priceID)
	return err
}
