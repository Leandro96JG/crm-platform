package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/jmoiron/sqlx"
)

type planchaRepository struct {
	client *sqlx.DB
}

func NewPlanchaRepository(client *sqlx.DB) domain.PlanchaRepository {
	return &planchaRepository{client: client}
}

func (r *planchaRepository) Create(ctx context.Context, plancha domain.Plancha) (string, error) {
	planchaDTO := mapPlanchaToPlanchaDTO(plancha)

	_, err := executor(ctx, r.client).NamedExecContext(
		ctx,
		`INSERT INTO planchas
		(plancha_id, name, description, category, subcategory, layout_file_url, preview_image_url, notes, is_active, created_by, created_at, updated_by, updated_at)
		VALUES
		(:plancha_id, :name, :description, :category, :subcategory, :layout_file_url, :preview_image_url, :notes, :is_active, :created_by, :created_at, :updated_by, :updated_at)`,
		planchaDTO,
	)
	if err != nil {
		return "", err
	}

	return plancha.PlanchaID, nil
}

func (r *planchaRepository) GetByID(ctx context.Context, planchaID string) (*domain.Plancha, error) {
	if planchaID == "" {
		return nil, domain.NewValidationError("planchaID is required", nil)
	}

	var dto PlanchaDTO
	err := executor(ctx, r.client).GetContext(ctx, &dto, "SELECT * FROM planchas WHERE plancha_id=$1", planchaID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("no plancha found with this id", map[string]any{"plancha_id": planchaID})
		}
		return nil, err
	}

	plancha := mapPlanchaDTOToPlancha(dto)
	return &plancha, nil
}

func (r *planchaRepository) Search(ctx context.Context, filters domain.PlanchaFilters) (domain.PagingResult[domain.Plancha], error) {
	whereQuery := []string{"1=1"}
	whereArgs := make([]any, 0)

	whereQuery, whereArgs = prepareInQuery(filters.PlanchaID, whereQuery, whereArgs, "plancha_id")
	whereQuery, whereArgs = prepareInQuery(filters.Category, whereQuery, whereArgs, "category")
	whereQuery, whereArgs = prepareInQuery(filters.Subcategory, whereQuery, whereArgs, "subcategory")

	if filters.IsActive != nil {
		whereQuery = append(whereQuery, fmt.Sprintf("is_active = $%d", len(whereArgs)+1))
		whereArgs = append(whereArgs, *filters.IsActive)
	}

	limitQuery := fmt.Sprintf("LIMIT $%d OFFSET $%d", len(whereArgs)+1, len(whereArgs)+2)
	limitArgs := append(whereArgs, filters.Limit, filters.Offset)

	orderBy := buildOrderBy(filters.SortBy, filters.SortOrder, validPlanchaSortColumns)

	query := fmt.Sprintf("SELECT * FROM planchas WHERE %s ORDER BY %s %s", joinWhere(whereQuery), orderBy, limitQuery)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM planchas WHERE %s", joinWhere(whereQuery))

	var planchasDTO []PlanchaDTO
	err := executor(ctx, r.client).SelectContext(ctx, &planchasDTO, query, limitArgs...)
	if err != nil {
		return domain.PagingResult[domain.Plancha]{}, err
	}

	var count int
	err = executor(ctx, r.client).GetContext(ctx, &count, countQuery, whereArgs...)
	if err != nil {
		return domain.PagingResult[domain.Plancha]{}, err
	}

	return domain.PagingResult[domain.Plancha]{
		Result: mapPlanchaDTOsToPlanchas(planchasDTO),
		Paging: domain.Paging{
			Total:  count,
			Limit:  filters.Limit,
			Offset: filters.Offset,
		},
	}, nil
}

func (r *planchaRepository) Update(ctx context.Context, plancha domain.Plancha) error {
	dto := mapPlanchaToPlanchaDTO(plancha)

	_, err := executor(ctx, r.client).NamedExecContext(
		ctx,
		`UPDATE planchas SET
		name = :name, description = :description, category = :category, subcategory = :subcategory,
		layout_file_url = :layout_file_url, preview_image_url = :preview_image_url, notes = :notes,
		is_active = :is_active, updated_by = :updated_by, updated_at = :updated_at
		WHERE plancha_id = :plancha_id`,
		dto,
	)

	return err
}

func (r *planchaRepository) Delete(ctx context.Context, planchaID string) error {
	if planchaID == "" {
		return domain.NewValidationError("planchaID is required", nil)
	}

	_, err := executor(ctx, r.client).ExecContext(ctx, "DELETE FROM planchas WHERE plancha_id=$1", planchaID)
	return err
}

var validPlanchaSortColumns = map[string]bool{
	"created_at": true,
	"updated_at": true,
	"name":       true,
	"category":   true,
}

func joinWhere(conditions []string) string {
	result := conditions[0]
	for i := 1; i < len(conditions); i++ {
		result += " AND " + conditions[i]
	}
	return result
}
