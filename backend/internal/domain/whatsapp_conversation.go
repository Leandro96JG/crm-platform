package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type WhatsAppConversationRepository interface {
	Create(ctx context.Context, conversation WhatsAppConversation) (string, error)
	GetByID(ctx context.Context, conversationID string) (*WhatsAppConversation, error)
	GetByPhone(ctx context.Context, phone string) (*WhatsAppConversation, error)
	Update(ctx context.Context, conversation WhatsAppConversation) error
}

type WhatsAppConversation struct {
	ConversationID    string
	CustomerID        string
	WhatsappPhone     string
	WAConversationID  string
	Status            string
	AiHandled         bool
	AiAgentSessionID  string
	LastMessageAt     *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func NewWhatsAppConversation(
	whatsappPhone string,
	customerID string,
) (WhatsAppConversation, error) {
	conversationID, err := uuid.NewUUID()
	if err != nil {
		return WhatsAppConversation{}, err
	}

	now := time.Now().UTC()

	return WhatsAppConversation{
		ConversationID:   conversationID.String(),
		CustomerID:       customerID,
		WhatsappPhone:    whatsappPhone,
		Status:           "active",
		AiHandled:        true,
		LastMessageAt:    &now,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}
