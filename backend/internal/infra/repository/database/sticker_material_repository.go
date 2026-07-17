package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/jmoiron/sqlx"
)

type stickerMaterialRepository struct {
	client *sqlx.DB
}

func NewStickerMaterialRepository(client *sqlx.DB) domain.StickerMaterialRepository {
	return &stickerMaterialRepository{client: client}
}

func (r *stickerMaterialRepository) Create(ctx context.Context, material domain.StickerMaterial) (string, error) {
	dto := mapStickerMaterialToDTO(material)

	_, err := executor(ctx, r.client).NamedExecContext(
		ctx,
		`INSERT INTO sticker_materials
		(material_id, name, description, material_type, finish, is_cuttable, is_printable, base_cost, is_active, created_at)
		VALUES
		(:material_id, :name, :description, :material_type, :finish, :is_cuttable, :is_printable, :base_cost, :is_active, :created_at)`,
		dto,
	)
	if err != nil {
		return "", err
	}

	return material.MaterialID, nil
}

func (r *stickerMaterialRepository) GetByID(ctx context.Context, materialID string) (*domain.StickerMaterial, error) {
	if materialID == "" {
		return nil, domain.NewValidationError("materialID is required", nil)
	}

	var dto StickerMaterialDTO
	err := executor(ctx, r.client).GetContext(ctx, &dto, "SELECT * FROM sticker_materials WHERE material_id=$1", materialID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("no material found with this id", map[string]any{"material_id": materialID})
		}
		return nil, err
	}

	material := mapDTOToStickerMaterial(dto)
	return &material, nil
}

func (r *stickerMaterialRepository) Search(ctx context.Context, filters domain.StickerMaterialFilters) (domain.PagingResult[domain.StickerMaterial], error) {
	whereQuery := []string{"1=1"}
	whereArgs := make([]any, 0)

	whereQuery, whereArgs = prepareInQuery(filters.MaterialType, whereQuery, whereArgs, "material_type")
	whereQuery, whereArgs = prepareInQuery(filters.Finish, whereQuery, whereArgs, "finish")

	if filters.IsActive != nil {
		whereQuery = append(whereQuery, fmt.Sprintf("is_active = $%d", len(whereArgs)+1))
		whereArgs = append(whereArgs, *filters.IsActive)
	}

	limitQuery := fmt.Sprintf("LIMIT $%d OFFSET $%d", len(whereArgs)+1, len(whereArgs)+2)
	limitArgs := append(whereArgs, filters.Limit, filters.Offset)

	orderBy := buildOrderBy(filters.SortBy, filters.SortOrder, validMaterialSortColumns)

	query := fmt.Sprintf("SELECT * FROM sticker_materials WHERE %s ORDER BY %s %s", joinWhere(whereQuery), orderBy, limitQuery)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM sticker_materials WHERE %s", joinWhere(whereQuery))

	var dtos []StickerMaterialDTO
	err := executor(ctx, r.client).SelectContext(ctx, &dtos, query, limitArgs...)
	if err != nil {
		return domain.PagingResult[domain.StickerMaterial]{}, err
	}

	var count int
	err = executor(ctx, r.client).GetContext(ctx, &count, countQuery, whereArgs...)
	if err != nil {
		return domain.PagingResult[domain.StickerMaterial]{}, err
	}

	return domain.PagingResult[domain.StickerMaterial]{
		Result: mapDTOsToStickerMaterials(dtos),
		Paging: domain.Paging{Total: count, Limit: filters.Limit, Offset: filters.Offset},
	}, nil
}

func (r *stickerMaterialRepository) Update(ctx context.Context, material domain.StickerMaterial) error {
	dto := mapStickerMaterialToDTO(material)

	_, err := executor(ctx, r.client).NamedExecContext(
		ctx,
		`UPDATE sticker_materials SET
		name = :name, description = :description, material_type = :material_type, finish = :finish,
		is_cuttable = :is_cuttable, is_printable = :is_printable, base_cost = :base_cost, is_active = :is_active
		WHERE material_id = :material_id`,
		dto,
	)

	return err
}

func (r *stickerMaterialRepository) Delete(ctx context.Context, materialID string) error {
	if materialID == "" {
		return domain.NewValidationError("materialID is required", nil)
	}

	_, err := executor(ctx, r.client).ExecContext(ctx, "DELETE FROM sticker_materials WHERE material_id=$1", materialID)
	return err
}

var validMaterialSortColumns = map[string]bool{
	"created_at":    true,
	"name":          true,
	"material_type": true,
	"base_cost":     true,
}
