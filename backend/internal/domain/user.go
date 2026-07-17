package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user User) (string, error)
	GetByID(ctx context.Context, userID string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Search(ctx context.Context, filters UserFilters) (PagingResult[User], error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, userID string) error
}

type User struct {
	UserID       string
	Username     string
	PasswordHash string
	Name         string
	Email        string
	Role         string
	IsActive     bool
	CreatedBy    string
	CreatedAt    time.Time
	UpdatedBy    string
	UpdatedAt    time.Time
}

type UserFilters struct {
	UserID   []string
	Username []string
	Role     []string
	IsActive *bool
	PagingFilter
}

type UpdateUser struct {
	Name      *string
	Email     *string
	Role      *string
	IsActive  *bool
	UpdatedBy string
}

func NewUser(username, passwordHash, name, email, role, createdBy string) (User, error) {
	now := time.Now().UTC()
	id, err := uuid.NewUUID()
	if err != nil {
		return User{}, err
	}
	return User{
		UserID:       id.String(),
		Username:     username,
		PasswordHash: passwordHash,
		Name:         name,
		Email:        email,
		Role:         role,
		IsActive:     true,
		CreatedBy:    createdBy,
		CreatedAt:    now,
		UpdatedBy:    createdBy,
		UpdatedAt:    now,
	}, nil
}

func (u *User) MergeUpdate(update UpdateUser) {
	u.UpdatedAt = time.Now().UTC()
	u.UpdatedBy = update.UpdatedBy
	if update.Name != nil {
		u.Name = *update.Name
	}
	if update.Email != nil {
		u.Email = *update.Email
	}
	if update.Role != nil {
		u.Role = *update.Role
	}
	if update.IsActive != nil {
		u.IsActive = *update.IsActive
	}
}
