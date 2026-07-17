package application

import (
	"context"
	"fmt"
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type whatsappService struct {
	conversationRepo domain.WhatsAppConversationRepository
	messageRepo      domain.WhatsAppMessageRepository
	sessionRepo      domain.ChatbotSessionRepository
	whatsappClient   WhatsAppClient
	aiAgent          AiAgentService
}

type WhatsAppClient interface {
	SendText(ctx context.Context, to string, message string) error
	SendTemplate(ctx context.Context, to string, templateName string, parameters map[string]string) error
	VerifyWebhook(ctx context.Context, token string, challenge string) (string, error)
}

type WhatsAppService interface {
	ReceiveMessage(ctx context.Context, payload map[string]any) error
	SendMessage(ctx context.Context, conversationID string, content string) error
	SendOrderNotification(ctx context.Context, conversationID string, orderNumber string, status string) error
}

func NewWhatsAppService(
	conversationRepo domain.WhatsAppConversationRepository,
	messageRepo domain.WhatsAppMessageRepository,
	sessionRepo domain.ChatbotSessionRepository,
	whatsappClient WhatsAppClient,
	aiAgent AiAgentService,
) WhatsAppService {
	return &whatsappService{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
		sessionRepo:      sessionRepo,
		whatsappClient:   whatsappClient,
		aiAgent:          aiAgent,
	}
}

func (s *whatsappService) ReceiveMessage(ctx context.Context, payload map[string]any) error {
	entries, ok := payload["entry"].([]any)
	if !ok || len(entries) == 0 {
		return fmt.Errorf("invalid payload: no entries")
	}

	entry := entries[0].(map[string]any)
	changes, ok := entry["changes"].([]any)
	if !ok || len(changes) == 0 {
		return fmt.Errorf("invalid payload: no changes")
	}

	change := changes[0].(map[string]any)
	value, ok := change["value"].(map[string]any)
	if !ok {
		return fmt.Errorf("invalid payload: no value")
	}

	messages, ok := value["messages"].([]any)
	if !ok || len(messages) == 0 {
		return nil
	}

	msg := messages[0].(map[string]any)
	from := msg["from"].(string)
	msgType := msg["type"].(string)
	waMsgID := msg["id"].(string)

	var content string
	if msgType == "text" {
		textObj := msg["text"].(map[string]any)
		content = textObj["body"].(string)
	}

	conversation, err := s.conversationRepo.GetByPhone(ctx, from)
	if err != nil {
		return err
	}

	var conversationID string
	if conversation == nil {
		newConv, err := domain.NewWhatsAppConversation(from, "")
		if err != nil {
			return err
		}

		conversationID, err = s.conversationRepo.Create(ctx, newConv)
		if err != nil {
			return err
		}
	} else {
		conversationID = conversation.ConversationID
		now := time.Now().UTC()
		conversation.LastMessageAt = &now
		_ = s.conversationRepo.Update(ctx, *conversation)
	}

	inboundMsg, err := domain.NewInboundMessage(conversationID, waMsgID, msgType, content, "", nil)
	if err != nil {
		return err
	}

	_, err = s.messageRepo.Create(ctx, inboundMsg)
	if err != nil {
		return err
	}

	if conversation == nil || conversation.AiHandled {
		aiResponse, err := s.aiAgent.ProcessInboundMessage(ctx, conversationID, content)
		if err != nil {
			return err
		}

		if aiResponse != "" {
			outboundMsg, err := domain.NewOutboundMessage(conversationID, "text", aiResponse, "")
			if err != nil {
				return err
			}

			_, err = s.messageRepo.Create(ctx, outboundMsg)
			if err != nil {
				return err
			}

			return s.whatsappClient.SendText(ctx, from, aiResponse)
		}
	}

	return nil
}

func (s *whatsappService) SendMessage(ctx context.Context, conversationID string, content string) error {
	conversation, err := s.conversationRepo.GetByID(ctx, conversationID)
	if err != nil {
		return err
	}

	outboundMsg, err := domain.NewOutboundMessage(conversationID, "text", content, "")
	if err != nil {
		return err
	}

	_, err = s.messageRepo.Create(ctx, outboundMsg)
	if err != nil {
		return err
	}

	return s.whatsappClient.SendText(ctx, conversation.WhatsappPhone, content)
}

func (s *whatsappService) SendOrderNotification(ctx context.Context, conversationID string, orderNumber string, status string) error {
	var message string
	switch status {
	case "ready":
		message = fmt.Sprintf("¡Tu pedido %s está listo! Pasá a retirarlo cuando quieras 🎉", orderNumber)
	case "delivered":
		message = fmt.Sprintf("¡Gracias por tu compra %s! Esperamos verte pronto 😊", orderNumber)
	default:
		message = fmt.Sprintf("Tu pedido %s cambió a estado: %s", orderNumber, status)
	}

	return s.SendMessage(ctx, conversationID, message)
}
