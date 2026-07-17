package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
)

type WhatsAppController struct {
	whatsappService application.WhatsAppService
	verifyToken     string
}

func NewWhatsAppController(whatsappService application.WhatsAppService, verifyToken string) WhatsAppController {
	return WhatsAppController{
		whatsappService: whatsappService,
		verifyToken:     verifyToken,
	}
}

func (ctrl *WhatsAppController) VerifyWebhook(ctx *gin.Context) {
	mode := ctx.Query("hub.mode")
	token := ctx.Query("hub.verify_token")
	challenge := ctx.Query("hub.challenge")

	if mode != "subscribe" || token != ctrl.verifyToken {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "verification failed"})
		return
	}

	ctx.String(http.StatusOK, challenge)
}

func (ctrl *WhatsAppController) ReceiveMessage(ctx *gin.Context) {
	var payload map[string]any
	if err := ctx.BindJSON(&payload); err != nil {
		ctx.Error(err)
		return
	}

	err := ctrl.whatsappService.ReceiveMessage(ctx.Request.Context(), payload)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (ctrl *WhatsAppController) SendMessage(ctx *gin.Context) {
	var dto struct {
		ConversationID string `json:"conversation_id" validate:"required"`
		Content        string `json:"content" validate:"required"`
	}

	if err := ctx.BindJSON(&dto); err != nil {
		ctx.Error(err)
		return
	}

	err := ctrl.whatsappService.SendMessage(ctx.Request.Context(), dto.ConversationID, dto.Content)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "sent"})
}
