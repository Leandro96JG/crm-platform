package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ChatbotSessionRepository interface {
	Create(ctx context.Context, session ChatbotSession) (string, error)
	GetByID(ctx context.Context, sessionID string) (*ChatbotSession, error)
	GetActiveByConversationID(ctx context.Context, conversationID string) (*ChatbotSession, error)
	Update(ctx context.Context, session ChatbotSession) error
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatbotSession struct {
	SessionID       string
	ConversationID  string
	Status          string
	AiModel         string
	SystemPrompt    string
	Messages        []ChatMessage
	EscalatedToHuman bool
	EscalatedAt     *time.Time
	ResolvedAt      *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewChatbotSession(
	conversationID string,
	systemPrompt string,
	aiModel string,
) (ChatbotSession, error) {
	sessionID, err := uuid.NewUUID()
	if err != nil {
		return ChatbotSession{}, err
	}

	now := time.Now().UTC()

	return ChatbotSession{
		SessionID:      sessionID.String(),
		ConversationID: conversationID,
		Status:         "active",
		AiModel:        aiModel,
		SystemPrompt:   systemPrompt,
		Messages:       []ChatMessage{},
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

func (s *ChatbotSession) AddMessage(role string, content string) {
	s.Messages = append(s.Messages, ChatMessage{Role: role, Content: content})
	s.UpdatedAt = time.Now().UTC()
}

func (s *ChatbotSession) EscalateToHuman() {
	now := time.Now().UTC()
	s.EscalatedToHuman = true
	s.EscalatedAt = &now
	s.Status = "escalated"
	s.UpdatedAt = now
}

func (s *ChatbotSession) Resolve() {
	now := time.Now().UTC()
	s.Status = "resolved"
	s.ResolvedAt = &now
	s.UpdatedAt = now
}

func (s *ChatbotSession) MarshalMessages() ([]byte, error) {
	return json.Marshal(s.Messages)
}

func (s *ChatbotSession) UnmarshalMessages(data []byte) error {
	return json.Unmarshal(data, &s.Messages)
}
