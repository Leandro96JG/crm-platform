package database

import (
	"encoding/json"
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type WhatsAppMessageDTO struct {
	MessageID      string    `db:"message_id"`
	ConversationID string    `db:"conversation_id"`
	WAMessageID    string    `db:"wa_message_id"`
	Direction      string    `db:"direction"`
	MessageType    string    `db:"message_type"`
	Content        string    `db:"content"`
	MediaURL       string    `db:"media_url"`
	Metadata       string    `db:"metadata"`
	AiProcessed    bool      `db:"ai_processed"`
	AiResponse     string    `db:"ai_response"`
	CreatedAt      time.Time `db:"created_at"`
}

func mapWhatsAppMessageToDTO(m domain.WhatsAppMessage) (WhatsAppMessageDTO, error) {
	metadataBytes, err := jsonMarshal(m.Metadata)
	if err != nil {
		return WhatsAppMessageDTO{}, err
	}

	return WhatsAppMessageDTO{
		MessageID:      m.MessageID,
		ConversationID: m.ConversationID,
		WAMessageID:    m.WAMessageID,
		Direction:      m.Direction,
		MessageType:    m.MessageType,
		Content:        m.Content,
		MediaURL:       m.MediaURL,
		Metadata:       string(metadataBytes),
		AiProcessed:    m.AiProcessed,
		AiResponse:     m.AiResponse,
		CreatedAt:      m.CreatedAt,
	}, nil
}

func mapDTOToWhatsAppMessage(dto WhatsAppMessageDTO) (domain.WhatsAppMessage, error) {
	m := domain.WhatsAppMessage{
		MessageID:      dto.MessageID,
		ConversationID: dto.ConversationID,
		WAMessageID:    dto.WAMessageID,
		Direction:      dto.Direction,
		MessageType:    dto.MessageType,
		Content:        dto.Content,
		MediaURL:       dto.MediaURL,
		AiProcessed:    dto.AiProcessed,
		AiResponse:     dto.AiResponse,
		CreatedAt:      dto.CreatedAt,
	}

	if err := m.SetMetadata([]byte(dto.Metadata)); err != nil {
		return domain.WhatsAppMessage{}, err
	}

	return m, nil
}

func mapDTOsToWhatsAppMessages(dtos []WhatsAppMessageDTO) ([]domain.WhatsAppMessage, error) {
	messages := make([]domain.WhatsAppMessage, 0, len(dtos))
	for _, dto := range dtos {
		m, err := mapDTOToWhatsAppMessage(dto)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func jsonMarshal(v any) ([]byte, error) {
	return json.Marshal(v)
}
