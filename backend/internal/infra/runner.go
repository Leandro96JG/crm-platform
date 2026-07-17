package infra

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/infra/ai"
	"github.com/icrxz/crm-api-core/internal/infra/config"
	"github.com/icrxz/crm-api-core/internal/infra/entrypoint"
	"github.com/icrxz/crm-api-core/internal/infra/entrypoint/middleware"
	"github.com/icrxz/crm-api-core/internal/infra/entrypoint/rest"
	"github.com/icrxz/crm-api-core/internal/infra/repository/database"
	"github.com/icrxz/crm-api-core/internal/infra/whatsapp"
)

func RunApp() error {
	_ = context.Background()

	appConfig, err := config.Load()
	if err != nil {
		return err
	}

	// database
	sqlDB, err := database.NewDatabase(appConfig.Database)
	if err != nil {
		return err
	}
	defer func() {
		err := sqlDB.Close()
		panic(err)
	}()

	// repositories
	userRepository := database.NewUserRepository(sqlDB)

	// sticker repositories
	planchaRepository := database.NewPlanchaRepository(sqlDB)
	materialRepository := database.NewStickerMaterialRepository(sqlDB)
	planchaPriceRepository := database.NewPlanchaPriceRepository(sqlDB)
	customerRepository := database.NewCustomerRepository(sqlDB)
	orderRepository := database.NewOrderRepository(sqlDB)
	printJobRepository := database.NewPrintJobRepository(sqlDB)
	whatsappConversationRepository := database.NewWhatsAppConversationRepository(sqlDB)
	whatsappMessageRepository := database.NewWhatsAppMessageRepository(sqlDB)
	chatbotSessionRepository := database.NewChatbotSessionRepository(sqlDB)

	orderItemRepository := database.NewOrderItemRepository(sqlDB)

	// services
	userService := application.NewUserService(userRepository)
	authService := application.NewAuthService(userRepository, appConfig.SecretKey())

	// sticker + whatsapp + ai services
	planchaService := application.NewPlanchaService(planchaRepository, materialRepository, planchaPriceRepository)
	customerService := application.NewCustomerService(customerRepository)
	orderService := application.NewOrderService(orderRepository, planchaRepository, planchaPriceRepository, printJobRepository, orderItemRepository)
	printingService := application.NewPrintingService(printJobRepository, orderRepository, orderService)
	openaiClient := ai.NewClient(appConfig.OpenAI.APIKeyValue(), appConfig.OpenAI.Model)
	aiAgentService := application.NewAiAgentService(chatbotSessionRepository, whatsappConversationRepository, openaiClient, planchaService, orderService)
	whatsappClient := whatsapp.NewClient(appConfig.WhatsApp.PhoneNumberID, appConfig.WhatsApp.AccessTokenValue(), appConfig.WhatsApp.APIVersion)
	whatsappService := application.NewWhatsAppService(whatsappConversationRepository, whatsappMessageRepository, chatbotSessionRepository, whatsappClient, aiAgentService)

	// controllers
	pingController := rest.NewPingController()
	userController := rest.NewUserController(userService)
	authController := rest.NewAuthController(authService)

	// sticker controllers
	planchaController := rest.NewPlanchaController(planchaService)
	customerController := rest.NewCustomerController(customerService)
	orderController := rest.NewOrderController(orderService)
	printingController := rest.NewPrintingController(printingService)
	whatsappController := rest.NewWhatsAppController(whatsappService, appConfig.WhatsApp.WebhookSecretValue())

	// middlewares
	authMiddleware := middleware.NewAuthenticationMiddleware(authService)

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(entrypoint.CustomErrorEncoder())

	err = router.SetTrustedProxies(nil)
	if err != nil {
		return err
	}

	entrypoint.LoadRoutes(
		router,
		pingController,
		userController,
		authController,
		authMiddleware,
	)

	entrypoint.LoadStickerRoutes(
		router,
		authMiddleware,
		planchaController,
		orderController,
		printingController,
		whatsappController,
		customerController,
	)

	return router.Run()
}
