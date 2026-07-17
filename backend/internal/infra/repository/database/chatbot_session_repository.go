package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/jmoiron/sqlx"
)

type chatbotSessionRepository struct {
	client *sqlx.DB
}

func NewChatbotSessionRepository(client *sqlx.DB) domain.ChatbotSessionRepository {
	return &chatbotSessionRepository{client: client}
}

func (r *chatbotSessionRepository) Create(ctx context.Context, session domain.ChatbotSession) (string, error) {
	messagesBytes, err := session.MarshalMessages()
	if err != nil {
		return "", err
	}

	_, err = executor(ctx, r.client).ExecContext(ctx,
		`INSERT INTO chatbot_sessions
		(session_id, conversation_id, status, ai_model, system_prompt, messages, escalated_to_human, escalated_at, resolved_at, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		session.SessionID, session.ConversationID, session.Status, session.AiModel,
		session.SystemPrompt, string(messagesBytes), session.EscalatedToHuman,
		session.EscalatedAt, session.ResolvedAt, session.CreatedAt, session.UpdatedAt,
	)
	if err != nil {
		return "", err
	}

	return session.SessionID, nil
}

func (r *chatbotSessionRepository) GetByID(ctx context.Context, sessionID string) (*domain.ChatbotSession, error) {
	if sessionID == "" {
		return nil, domain.NewValidationError("sessionID is required", nil)
	}

	type row struct {
		SessionID        string     `db:"session_id"`
		ConversationID   string     `db:"conversation_id"`
		Status           string     `db:"status"`
		AiModel          string     `db:"ai_model"`
		SystemPrompt     string     `db:"system_prompt"`
		Messages         string     `db:"messages"`
		EscalatedToHuman bool       `db:"escalated_to_human"`
		EscalatedAt      *time.Time `db:"escalated_at"`
		ResolvedAt       *time.Time `db:"resolved_at"`
		CreatedAt        time.Time  `db:"created_at"`
		UpdatedAt        time.Time  `db:"updated_at"`
	}

	var rowData row
	err := executor(ctx, r.client).GetContext(ctx, &rowData, "SELECT * FROM chatbot_sessions WHERE session_id=$1", sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("session not found", map[string]any{"session_id": sessionID})
		}
		return nil, err
	}

	session := domain.ChatbotSession{
		SessionID:        rowData.SessionID,
		ConversationID:   rowData.ConversationID,
		Status:           rowData.Status,
		AiModel:          rowData.AiModel,
		SystemPrompt:     rowData.SystemPrompt,
		EscalatedToHuman: rowData.EscalatedToHuman,
		EscalatedAt:      rowData.EscalatedAt,
		ResolvedAt:       rowData.ResolvedAt,
		CreatedAt:        rowData.CreatedAt,
		UpdatedAt:        rowData.UpdatedAt,
	}

	if err := session.UnmarshalMessages([]byte(rowData.Messages)); err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *chatbotSessionRepository) GetActiveByConversationID(ctx context.Context, conversationID string) (*domain.ChatbotSession, error) {
	type row struct {
		SessionID        string     `db:"session_id"`
		ConversationID   string     `db:"conversation_id"`
		Status           string     `db:"status"`
		AiModel          string     `db:"ai_model"`
		SystemPrompt     string     `db:"system_prompt"`
		Messages         string     `db:"messages"`
		EscalatedToHuman bool       `db:"escalated_to_human"`
		EscalatedAt      *time.Time `db:"escalated_at"`
		ResolvedAt       *time.Time `db:"resolved_at"`
		CreatedAt        time.Time  `db:"created_at"`
		UpdatedAt        time.Time  `db:"updated_at"`
	}

	var rowData row
	err := executor(ctx, r.client).GetContext(ctx, &rowData,
		"SELECT * FROM chatbot_sessions WHERE conversation_id=$1 AND status='active' ORDER BY created_at DESC LIMIT 1",
		conversationID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	session := domain.ChatbotSession{
		SessionID:        rowData.SessionID,
		ConversationID:   rowData.ConversationID,
		Status:           rowData.Status,
		AiModel:          rowData.AiModel,
		SystemPrompt:     rowData.SystemPrompt,
		EscalatedToHuman: rowData.EscalatedToHuman,
		EscalatedAt:      rowData.EscalatedAt,
		ResolvedAt:       rowData.ResolvedAt,
		CreatedAt:        rowData.CreatedAt,
		UpdatedAt:        rowData.UpdatedAt,
	}

	var messages []domain.ChatMessage
	if err := json.Unmarshal([]byte(rowData.Messages), &messages); err != nil {
		return nil, err
	}
	session.Messages = messages

	return &session, nil
}

func (r *chatbotSessionRepository) Update(ctx context.Context, session domain.ChatbotSession) error {
	messagesBytes, err := session.MarshalMessages()
	if err != nil {
		return err
	}

	_, err = executor(ctx, r.client).ExecContext(ctx,
		`UPDATE chatbot_sessions SET
		status=$1, messages=$2, escalated_to_human=$3, escalated_at=$4, resolved_at=$5, updated_at=$6
		WHERE session_id=$7`,
		session.Status, string(messagesBytes), session.EscalatedToHuman,
		session.EscalatedAt, session.ResolvedAt, session.UpdatedAt,
		session.SessionID,
	)

	return err
}
