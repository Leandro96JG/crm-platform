package database

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type WhatsAppConversationDTO struct {
	ConversationID   string     `db:"conversation_id"`
	CustomerID       string     `db:"customer_id"`
	WhatsappPhone    string     `db:"whatsapp_phone"`
	WAConversationID string     `db:"wa_conversation_id"`
	Status           string     `db:"status"`
	AiHandled        bool       `db:"ai_handled"`
	AiAgentSessionID string     `db:"ai_agent_session_id"`
	LastMessageAt    *time.Time `db:"last_message_at"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
}

func mapWhatsAppConversationToDTO(c domain.WhatsAppConversation) WhatsAppConversationDTO {
	return WhatsAppConversationDTO{
		ConversationID:   c.ConversationID,
		CustomerID:       c.CustomerID,
		WhatsappPhone:    c.WhatsappPhone,
		WAConversationID: c.WAConversationID,
		Status:           c.Status,
		AiHandled:        c.AiHandled,
		AiAgentSessionID: c.AiAgentSessionID,
		LastMessageAt:    c.LastMessageAt,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
	}
}

func mapDTOToWhatsAppConversation(dto WhatsAppConversationDTO) domain.WhatsAppConversation {
	return domain.WhatsAppConversation{
		ConversationID:   dto.ConversationID,
		CustomerID:       dto.CustomerID,
		WhatsappPhone:    dto.WhatsappPhone,
		WAConversationID: dto.WAConversationID,
		Status:           dto.Status,
		AiHandled:        dto.AiHandled,
		AiAgentSessionID: dto.AiAgentSessionID,
		LastMessageAt:    dto.LastMessageAt,
		CreatedAt:        dto.CreatedAt,
		UpdatedAt:        dto.UpdatedAt,
	}
}
