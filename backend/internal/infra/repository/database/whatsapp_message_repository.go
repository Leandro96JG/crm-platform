package database

import (
	"context"
	"encoding/json"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/jmoiron/sqlx"
)

type whatsappMessageRepository struct {
	client *sqlx.DB
}

func NewWhatsAppMessageRepository(client *sqlx.DB) domain.WhatsAppMessageRepository {
	return &whatsappMessageRepository{client: client}
}

func (r *whatsappMessageRepository) Create(ctx context.Context, message domain.WhatsAppMessage) (string, error) {
	metadataBytes, err := json.Marshal(message.Metadata)
	if err != nil {
		return "", err
	}

	_, err = executor(ctx, r.client).ExecContext(ctx,
		`INSERT INTO whatsapp_messages
		(message_id, conversation_id, wa_message_id, direction, message_type, content, media_url, metadata, ai_processed, ai_response, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		message.MessageID, message.ConversationID, message.WAMessageID, message.Direction,
		message.MessageType, message.Content, message.MediaURL, string(metadataBytes),
		message.AiProcessed, message.AiResponse, message.CreatedAt,
	)
	if err != nil {
		return "", err
	}

	return message.MessageID, nil
}

func (r *whatsappMessageRepository) GetByConversationID(ctx context.Context, conversationID string) ([]domain.WhatsAppMessage, error) {
	type dto struct {
		MessageID      string `db:"message_id"`
		ConversationID string `db:"conversation_id"`
		WAMessageID    string `db:"wa_message_id"`
		Direction      string `db:"direction"`
		MessageType    string `db:"message_type"`
		Content        string `db:"content"`
		MediaURL       string `db:"media_url"`
		Metadata       string `db:"metadata"`
		AiProcessed    bool   `db:"ai_processed"`
		AiResponse     string `db:"ai_response"`
		CreatedAt      string `db:"created_at"`
	}

	var raw []dto
	err := executor(ctx, r.client).SelectContext(ctx, &raw,
		"SELECT * FROM whatsapp_messages WHERE conversation_id=$1 ORDER BY created_at ASC", conversationID)
	if err != nil {
		return nil, err
	}

	messages := make([]domain.WhatsAppMessage, 0, len(raw))
	for _, r := range raw {
		var metadata map[string]any
		_ = json.Unmarshal([]byte(r.Metadata), &metadata)

		messages = append(messages, domain.WhatsAppMessage{
			MessageID:      r.MessageID,
			ConversationID: r.ConversationID,
			WAMessageID:    r.WAMessageID,
			Direction:      r.Direction,
			MessageType:    r.MessageType,
			Content:        r.Content,
			MediaURL:       r.MediaURL,
			Metadata:       metadata,
			AiProcessed:    r.AiProcessed,
			AiResponse:     r.AiResponse,
		})
	}

	return messages, nil
}
