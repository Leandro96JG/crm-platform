package entrypoint

import (
	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/infra/entrypoint/middleware"
	"github.com/icrxz/crm-api-core/internal/infra/entrypoint/rest"
)

func LoadStickerRoutes(
	app *gin.Engine,
	authMiddleware middleware.AuthenticationMiddleware,
	planchaController rest.PlanchaController,
	orderController rest.OrderController,
	printingController rest.PrintingController,
	whatsappController rest.WhatsAppController,
) {
	authGroup := app.Group("/crm/core/api/v1")
	authGroup.Use(authMiddleware.Authenticate())

	publicGroup := app.Group("/crm/core/api/v1")

	// WhatsApp webhook (public)
	publicGroup.GET("/webhook/whatsapp", whatsappController.VerifyWebhook)
	publicGroup.POST("/webhook/whatsapp", whatsappController.ReceiveMessage)

	// WhatsApp messaging (authenticated)
	authGroup.POST("/messages/send", whatsappController.SendMessage)

	// Planchas (catalog)
	authGroup.POST("/planchas", planchaController.CreatePlancha)
	authGroup.GET("/planchas", planchaController.SearchPlanchas)
	authGroup.GET("/planchas/:planchaID", planchaController.GetPlancha)
	authGroup.PUT("/planchas/:planchaID", planchaController.UpdatePlancha)
	authGroup.DELETE("/planchas/:planchaID", planchaController.DeletePlancha)
	authGroup.GET("/planchas/:planchaID/prices", planchaController.GetPlanchaPrices)
	authGroup.GET("/planchas/:planchaID/calculate-price", planchaController.CalculatePrice)

	// Sticker materials
	authGroup.POST("/materials", planchaController.CreateMaterial)
	authGroup.GET("/materials", planchaController.SearchMaterials)

	// Plancha prices
	authGroup.POST("/prices", planchaController.CreatePrice)

	// Orders
	authGroup.POST("/orders", orderController.CreateOrder)
	authGroup.GET("/orders", orderController.SearchOrders)
	authGroup.GET("/orders/:orderID", orderController.GetOrder)
	authGroup.PUT("/orders/:orderID/status", orderController.UpdateOrderStatus)

	// Print jobs (production queue)
	authGroup.GET("/print-jobs", printingController.GetPrintQueue)
	authGroup.PUT("/print-jobs/:jobID/status", printingController.UpdatePrintJobStatus)
}
