package entrypoint

import (
	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/infra/entrypoint/middleware"
	"github.com/icrxz/crm-api-core/internal/infra/entrypoint/rest"
)

func LoadRoutes(
	app *gin.Engine,
	pingController rest.PingController,
	userController rest.UserController,
	authController rest.AuthController,
	authMiddleware middleware.AuthenticationMiddleware,
) {
	authGroup := app.Group("/crm/core/api/v1")
	authGroup.Use(authMiddleware.Authenticate())

	publicGroup := app.Group("/crm/core/api/v1")

	// miscellaneous
	app.GET("/ping", pingController.Pong)

	// user
	publicGroup.POST("/users", userController.CreateUser)
	authGroup.GET("/users", userController.SearchUser)
	authGroup.GET("/users/:userID", userController.GetUser)
	authGroup.PUT("/users/:userID", userController.UpdateUser)
	authGroup.DELETE("/users/:userID", userController.DeleteUser)
	authGroup.PUT("/users/:userID/password", userController.ChangePassword)

	// auth
	publicGroup.POST("/login", authController.Login)
	authGroup.POST("/logout", authController.Logout)
}
