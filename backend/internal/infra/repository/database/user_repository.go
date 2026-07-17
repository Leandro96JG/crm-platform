package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	client *sqlx.DB
}

func NewUserRepository(client *sqlx.DB) domain.UserRepository {
	return &userRepository{client: client}
}

func (r *userRepository) Create(ctx context.Context, user domain.User) (string, error) {
	dto := mapUserToDTO(user)

	_, err := executor(ctx, r.client).NamedExecContext(
		ctx,
		`INSERT INTO users
		(user_id, username, password_hash, name, email, role, is_active, created_by, created_at, updated_by, updated_at)
		VALUES
		(:user_id, :username, :password_hash, :name, :email, :role, :is_active, :created_by, :created_at, :updated_by, :updated_at)`,
		dto,
	)
	if err != nil {
		return "", err
	}

	return user.UserID, nil
}

func (r *userRepository) GetByID(ctx context.Context, userID string) (*domain.User, error) {
	if userID == "" {
		return nil, domain.NewValidationError("userID is required", nil)
	}

	var dto UserDTO
	err := executor(ctx, r.client).GetContext(ctx, &dto, "SELECT * FROM users WHERE user_id=$1", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("no user found with this id", map[string]any{"user_id": userID})
		}
		return nil, err
	}

	user := mapDTOToUser(dto)
	return &user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	if username == "" {
		return nil, domain.NewValidationError("username is required", nil)
	}

	var dto UserDTO
	err := executor(ctx, r.client).GetContext(ctx, &dto, "SELECT * FROM users WHERE username=$1", username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	user := mapDTOToUser(dto)
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	if email == "" {
		return nil, domain.NewValidationError("email is required", nil)
	}

	var dto UserDTO
	err := executor(ctx, r.client).GetContext(ctx, &dto, "SELECT * FROM users WHERE email=$1", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	user := mapDTOToUser(dto)
	return &user, nil
}

func (r *userRepository) Search(ctx context.Context, filters domain.UserFilters) (domain.PagingResult[domain.User], error) {
	whereQuery := []string{"1=1"}
	whereArgs := make([]any, 0)

	whereQuery, whereArgs = prepareInQuery(filters.UserID, whereQuery, whereArgs, "user_id")
	whereQuery, whereArgs = prepareInQuery(filters.Username, whereQuery, whereArgs, "username")
	whereQuery, whereArgs = prepareInQuery(filters.Role, whereQuery, whereArgs, "role")

	if filters.IsActive != nil {
		whereQuery = append(whereQuery, fmt.Sprintf("is_active = $%d", len(whereArgs)+1))
		whereArgs = append(whereArgs, *filters.IsActive)
	}

	limitQuery := fmt.Sprintf("LIMIT $%d OFFSET $%d", len(whereArgs)+1, len(whereArgs)+2)
	limitArgs := append(whereArgs, filters.Limit, filters.Offset)

	orderBy := buildOrderBy(filters.SortBy, filters.SortOrder, validUserSortColumns)

	query := fmt.Sprintf("SELECT * FROM users WHERE %s ORDER BY %s %s", joinWhere(whereQuery), orderBy, limitQuery)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users WHERE %s", joinWhere(whereQuery))

	var dtos []UserDTO
	err := executor(ctx, r.client).SelectContext(ctx, &dtos, query, limitArgs...)
	if err != nil {
		return domain.PagingResult[domain.User]{}, err
	}

	var count int
	err = executor(ctx, r.client).GetContext(ctx, &count, countQuery, whereArgs...)
	if err != nil {
		return domain.PagingResult[domain.User]{}, err
	}

	return domain.PagingResult[domain.User]{
		Result: mapDTOsToUsers(dtos),
		Paging: domain.Paging{Total: count, Limit: filters.Limit, Offset: filters.Offset},
	}, nil
}

func (r *userRepository) Update(ctx context.Context, user domain.User) error {
	dto := mapUserToDTO(user)

	_, err := executor(ctx, r.client).NamedExecContext(
		ctx,
		`UPDATE users SET
		username = :username, name = :name, email = :email, role = :role,
		is_active = :is_active, updated_by = :updated_by, updated_at = :updated_at
		WHERE user_id = :user_id`,
		dto,
	)

	return err
}

func (r *userRepository) Delete(ctx context.Context, userID string) error {
	if userID == "" {
		return domain.NewValidationError("userID is required", nil)
	}

	_, err := executor(ctx, r.client).ExecContext(ctx, "UPDATE users SET is_active=false, updated_at=NOW() WHERE user_id=$1", userID)
	return err
}
