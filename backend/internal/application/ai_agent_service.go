package application

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type aiAgentService struct {
	sessionRepo       domain.ChatbotSessionRepository
	conversationRepo  domain.WhatsAppConversationRepository
	openaiClient      OpenAIClient
	planchaSvc        PlanchaService
	orderSvc          OrderService
}

type OpenAIClient interface {
	ChatCompletion(ctx context.Context, systemPrompt string, messages []domain.ChatMessage) (string, error)
}

type AiAgentService interface {
	ProcessInboundMessage(ctx context.Context, conversationID string, messageContent string) (string, error)
}

const salesSystemPrompt = `Eres un vendedor de stickers por WhatsApp. Tu objetivo: saluda al cliente, pregunta que necesita, recomienda productos del catalogo, cotiza precios, y cierra la venta. Si el cliente pide descuento, esta enojado, o pide un diseno complejo, escala a humano. No inventes precios ni productos. Responde en espanol de Argentina, de forma amable y directa.`

func NewAiAgentService(
	sessionRepo domain.ChatbotSessionRepository,
	conversationRepo domain.WhatsAppConversationRepository,
	openaiClient OpenAIClient,
	planchaSvc PlanchaService,
	orderSvc OrderService,
) AiAgentService {
	return &aiAgentService{
		sessionRepo:      sessionRepo,
		conversationRepo: conversationRepo,
		openaiClient:     openaiClient,
		planchaSvc:       planchaSvc,
		orderSvc:         orderSvc,
	}
}

func (s *aiAgentService) ProcessInboundMessage(ctx context.Context, conversationID string, messageContent string) (string, error) {
	session, err := s.sessionRepo.GetActiveByConversationID(ctx, conversationID)
	if err != nil {
		return "", err
	}

	if session == nil {
		newSession, err := domain.NewChatbotSession(conversationID, salesSystemPrompt, "gpt-4o-mini")
		if err != nil {
			return "", err
		}

		sessionID, err := s.sessionRepo.Create(ctx, newSession)
		if err != nil {
			return "", err
		}

		session, err = s.sessionRepo.GetByID(ctx, sessionID)
		if err != nil {
			return "", err
		}
	}

	session.AddMessage("user", messageContent)

	if err := s.sessionRepo.Update(ctx, *session); err != nil {
		return "", err
	}

	aiResponse, err := s.openaiClient.ChatCompletion(ctx, session.SystemPrompt, session.Messages)
	if err != nil {
		return "", fmt.Errorf("AI chat completion failed: %w", err)
	}

	session.AddMessage("assistant", aiResponse)

	if err := s.sessionRepo.Update(ctx, *session); err != nil {
		return "", err
	}

	if s.shouldEscalate(aiResponse) {
		session.EscalateToHuman()
		_ = s.sessionRepo.Update(ctx, *session)

		conv, err := s.conversationRepo.GetByID(ctx, conversationID)
		if err == nil {
			conv.AiHandled = false
			_ = s.conversationRepo.Update(ctx, *conv)
		}

		return aiResponse + "\n\nUn operador humano te va a atender en breve.", nil
	}

	if s.isOrderConfirmation(aiResponse) {
		log.Printf("AI agent confirmed order for conversation %s", conversationID)
	}

	return aiResponse, nil
}

func (s *aiAgentService) shouldEscalate(response string) bool {
	escalationKeywords := []string{
		"hablar con un humano",
		"hablar con un operador",
		"reclamo",
		"queja",
		"descuento",
		"no me sirve",
		"no me gusta",
		"quiero hablar con",
	}

	responseBytes, _ := json.Marshal(response)
	responseStr := string(responseBytes)

	for _, keyword := range escalationKeywords {
		if contains(responseStr, keyword) {
			return true
		}
	}

	return false
}

func (s *aiAgentService) isOrderConfirmation(response string) bool {
	confirmKeywords := []string{
		"pedido confirmado",
		"orden creada",
		"te tomo el pedido",
		"confirmo tu pedido",
	}

	for _, keyword := range confirmKeywords {
		if contains(response, keyword) {
			return true
		}
	}

	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if toLower(s[i+j]) != toLower(substr[j]) {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func toLower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + 32
	}
	return c
}
