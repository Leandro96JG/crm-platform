package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
)

type AuthController struct {
	authService application.AuthService
}

type LoginDTO struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func NewAuthController(authService application.AuthService) AuthController {
	return AuthController{
		authService: authService,
	}
}

func (ctrl *AuthController) Login(ctx *gin.Context) {
	var dto LoginDTO
	if err := ctx.BindJSON(&dto); err != nil {
		ctx.Error(err)
		return
	}

	token, user, err := ctrl.authService.Login(ctx.Request.Context(), dto.Email, dto.Password)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

func (ctrl *AuthController) Logout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
