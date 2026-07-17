package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/jmoiron/sqlx"
)

type whatsappConversationRepository struct {
	client *sqlx.DB
}

func NewWhatsAppConversationRepository(client *sqlx.DB) domain.WhatsAppConversationRepository {
	return &whatsappConversationRepository{client: client}
}

func (r *whatsappConversationRepository) Create(ctx context.Context, conversation domain.WhatsAppConversation) (string, error) {
	dto := mapWhatsAppConversationToDTO(conversation)

	_, err := executor(ctx, r.client).NamedExecContext(
		ctx,
		`INSERT INTO whatsapp_conversations
		(conversation_id, customer_id, whatsapp_phone, wa_conversation_id, status, ai_handled, ai_agent_session_id, last_message_at, created_at, updated_at)
		VALUES
		(:conversation_id, :customer_id, :whatsapp_phone, :wa_conversation_id, :status, :ai_handled, :ai_agent_session_id, :last_message_at, :created_at, :updated_at)`,
		dto,
	)
	if err != nil {
		return "", err
	}

	return conversation.ConversationID, nil
}

func (r *whatsappConversationRepository) GetByID(ctx context.Context, conversationID string) (*domain.WhatsAppConversation, error) {
	if conversationID == "" {
		return nil, domain.NewValidationError("conversationID is required", nil)
	}

	var dto WhatsAppConversationDTO
	err := executor(ctx, r.client).GetContext(ctx, &dto, "SELECT * FROM whatsapp_conversations WHERE conversation_id=$1", conversationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("no conversation found", map[string]any{"conversation_id": conversationID})
		}
		return nil, err
	}

	c := mapDTOToWhatsAppConversation(dto)
	return &c, nil
}

func (r *whatsappConversationRepository) GetByPhone(ctx context.Context, phone string) (*domain.WhatsAppConversation, error) {
	if phone == "" {
		return nil, domain.NewValidationError("phone is required", nil)
	}

	var dto WhatsAppConversationDTO
	err := executor(ctx, r.client).GetContext(ctx, &dto, "SELECT * FROM whatsapp_conversations WHERE whatsapp_phone=$1 AND status='active'", phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	c := mapDTOToWhatsAppConversation(dto)
	return &c, nil
}

func (r *whatsappConversationRepository) Update(ctx context.Context, conversation domain.WhatsAppConversation) error {
	dto := mapWhatsAppConversationToDTO(conversation)

	_, err := executor(ctx, r.client).NamedExecContext(
		ctx,
		`UPDATE whatsapp_conversations SET
		customer_id = :customer_id, wa_conversation_id = :wa_conversation_id, status = :status,
		ai_handled = :ai_handled, ai_agent_session_id = :ai_agent_session_id, last_message_at = :last_message_at,
		updated_at = :updated_at
		WHERE conversation_id = :conversation_id`,
		dto,
	)

	return err
}
