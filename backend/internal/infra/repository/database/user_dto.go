package database

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type UserDTO struct {
	UserID       string    `db:"user_id"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	Role         string    `db:"role"`
	IsActive     bool      `db:"is_active"`
	CreatedBy    string    `db:"created_by"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedBy    string    `db:"updated_by"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func mapUserToDTO(u domain.User) UserDTO {
	return UserDTO{
		UserID:       u.UserID,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		Name:         u.Name,
		Email:        u.Email,
		Role:         u.Role,
		IsActive:     u.IsActive,
		CreatedBy:    u.CreatedBy,
		CreatedAt:    u.CreatedAt,
		UpdatedBy:    u.UpdatedBy,
		UpdatedAt:    u.UpdatedAt,
	}
}

func mapDTOToUser(dto UserDTO) domain.User {
	return domain.User{
		UserID:       dto.UserID,
		Username:     dto.Username,
		PasswordHash: dto.PasswordHash,
		Name:         dto.Name,
		Email:        dto.Email,
		Role:         dto.Role,
		IsActive:     dto.IsActive,
		CreatedBy:    dto.CreatedBy,
		CreatedAt:    dto.CreatedAt,
		UpdatedBy:    dto.UpdatedBy,
		UpdatedAt:    dto.UpdatedAt,
	}
}

func mapDTOsToUsers(dtos []UserDTO) []domain.User {
	users := make([]domain.User, 0, len(dtos))
	for _, dto := range dtos {
		users = append(users, mapDTOToUser(dto))
	}
	return users
}

var validUserSortColumns = map[string]bool{
	"created_at": true,
	"updated_at": true,
	"name":       true,
	"email":      true,
	"username":   true,
	"role":       true,
}
