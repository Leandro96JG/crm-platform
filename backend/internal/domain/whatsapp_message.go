package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type WhatsAppMessageRepository interface {
	Create(ctx context.Context, message WhatsAppMessage) (string, error)
	GetByConversationID(ctx context.Context, conversationID string) ([]WhatsAppMessage, error)
}

type WhatsAppMessage struct {
	MessageID      string
	ConversationID string
	WAMessageID    string
	Direction      string
	MessageType    string
	Content        string
	MediaURL       string
	Metadata       map[string]any
	AiProcessed    bool
	AiResponse     string
	CreatedAt      time.Time
}

func NewInboundMessage(
	conversationID string,
	waMessageID string,
	messageType string,
	content string,
	mediaURL string,
	metadata map[string]any,
) (WhatsAppMessage, error) {
	messageID, err := uuid.NewUUID()
	if err != nil {
		return WhatsAppMessage{}, err
	}

	return WhatsAppMessage{
		MessageID:      messageID.String(),
		ConversationID: conversationID,
		WAMessageID:    waMessageID,
		Direction:      "inbound",
		MessageType:    messageType,
		Content:        content,
		MediaURL:       mediaURL,
		Metadata:       metadata,
		CreatedAt:      time.Now().UTC(),
	}, nil
}

func NewOutboundMessage(
	conversationID string,
	messageType string,
	content string,
	mediaURL string,
) (WhatsAppMessage, error) {
	messageID, err := uuid.NewUUID()
	if err != nil {
		return WhatsAppMessage{}, err
	}

	return WhatsAppMessage{
		MessageID:      messageID.String(),
		ConversationID: conversationID,
		Direction:      "outbound",
		MessageType:    messageType,
		Content:        content,
		MediaURL:       mediaURL,
		CreatedAt:      time.Now().UTC(),
	}, nil
}

func (m *WhatsAppMessage) SetMetadata(data []byte) error {
	return json.Unmarshal(data, &m.Metadata)
}
