package rest

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/domain"
)

type UserController struct {
	userService application.UserService
}

type CreateUserDTO struct {
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Name      string `json:"name" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Role      string `json:"role" validate:"required"`
	CreatedBy string `json:"created_by" validate:"required"`
}

type UpdateUserDTO struct {
	Name      *string `json:"name"`
	Email     *string `json:"email"`
	Role      *string `json:"role"`
	IsActive  *bool   `json:"is_active"`
	UpdatedBy string  `json:"updated_by" validate:"required"`
}

type ChangePasswordDTO struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type UserResponseDTO struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Active    bool      `json:"active"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedBy string    `json:"updated_by"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PagingResponseDTO struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func mapUserToResponseDTO(u domain.User) UserResponseDTO {
	firstName, lastName := splitName(u.Name)
	return UserResponseDTO{
		UserID:    u.UserID,
		Username:  u.Username,
		FirstName: firstName,
		LastName:  lastName,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		Active:    u.IsActive,
		CreatedBy: u.CreatedBy,
		CreatedAt: u.CreatedAt,
		UpdatedBy: u.UpdatedBy,
		UpdatedAt: u.UpdatedAt,
	}
}

func splitName(name string) (string, string) {
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return "", ""
	}
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], strings.Join(parts[1:], " ")
}

func NewUserController(userService application.UserService) UserController {
	return UserController{
		userService: userService,
	}
}

func (ctrl *UserController) CreateUser(ctx *gin.Context) {
	var dto CreateUserDTO
	if err := ctx.BindJSON(&dto); err != nil {
		ctx.Error(err)
		return
	}

	user, err := domain.NewUser(dto.Username, dto.Password, dto.Name, dto.Email, dto.Role, dto.CreatedBy)
	if err != nil {
		ctx.Error(err)
		return
	}

	userID, err := ctrl.userService.CreateUser(ctx.Request.Context(), user)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user_id": userID})
}

func (ctrl *UserController) SearchUser(ctx *gin.Context) {
	filters := ctrl.parseUserFilters(ctx)

	users, err := ctrl.userService.SearchUsers(ctx.Request.Context(), filters)
	if err != nil {
		ctx.Error(err)
		return
	}

	result := make([]UserResponseDTO, 0, len(users.Result))
	for _, u := range users.Result {
		result = append(result, mapUserToResponseDTO(u))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": result,
		"paging": PagingResponseDTO{
			Total:  users.Paging.Total,
			Limit:  users.Paging.Limit,
			Offset: users.Paging.Offset,
		},
	})
}

func (ctrl *UserController) GetUser(ctx *gin.Context) {
	userID := ctx.Param("userID")
	if userID == "" {
		ctx.Error(domain.NewValidationError("user_id is required", nil))
		return
	}

	user, err := ctrl.userService.GetUserByID(ctx.Request.Context(), userID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, mapUserToResponseDTO(*user))
}

func (ctrl *UserController) UpdateUser(ctx *gin.Context) {
	userID := ctx.Param("userID")
	if userID == "" {
		ctx.Error(domain.NewValidationError("user_id is required", nil))
		return
	}

	var dto UpdateUserDTO
	if err := ctx.BindJSON(&dto); err != nil {
		ctx.Error(err)
		return
	}

	update := domain.UpdateUser{
		Name:      dto.Name,
		Email:     dto.Email,
		Role:      dto.Role,
		IsActive:  dto.IsActive,
		UpdatedBy: dto.UpdatedBy,
	}

	err := ctrl.userService.UpdateUser(ctx.Request.Context(), userID, update)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (ctrl *UserController) DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("userID")
	if userID == "" {
		ctx.Error(domain.NewValidationError("user_id is required", nil))
		return
	}

	err := ctrl.userService.DeleteUser(ctx.Request.Context(), userID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (ctrl *UserController) ChangePassword(ctx *gin.Context) {
	userID := ctx.Param("userID")
	if userID == "" {
		ctx.Error(domain.NewValidationError("user_id is required", nil))
		return
	}

	var dto ChangePasswordDTO
	if err := ctx.BindJSON(&dto); err != nil {
		ctx.Error(err)
		return
	}

	err := ctrl.userService.ChangePassword(ctx.Request.Context(), userID, dto.OldPassword, dto.NewPassword)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (ctrl *UserController) parseUserFilters(ctx *gin.Context) domain.UserFilters {
	filters := domain.UserFilters{
		PagingFilter: domain.PagingFilter{
			Limit:     10,
			Offset:    0,
			SortBy:    "created_at",
			SortOrder: "DESC",
		},
	}

	if name := ctx.QueryArray("name"); len(name) > 0 {
		filters.Username = name
	}
	if role := ctx.QueryArray("role"); len(role) > 0 {
		filters.Role = role
	}
	if active := ctx.Query("is_active"); active != "" {
		isActive := active == "true"
		filters.IsActive = &isActive
	}
	if limit := ctx.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			filters.Limit = l
		}
	}
	if offset := ctx.Query("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			filters.Offset = o
		}
	}
	if sortBy := ctx.Query("sort_by"); sortBy != "" {
		filters.SortBy = sortBy
	}
	if sortOrder := strings.ToUpper(ctx.Query("sort_order")); sortOrder == "ASC" || sortOrder == "DESC" {
		filters.SortOrder = sortOrder
	}

	return filters
}
