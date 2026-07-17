CREATE TABLE IF NOT EXISTS chatbot_sessions (
    session_id TEXT PRIMARY KEY,
    conversation_id TEXT NOT NULL REFERENCES whatsapp_conversations(conversation_id),
    status TEXT NOT NULL DEFAULT 'active',
    ai_model TEXT NOT NULL DEFAULT 'gpt-4o-mini',
    system_prompt TEXT DEFAULT '',
    messages JSONB NOT NULL DEFAULT '[]',
    escalated_to_human BOOLEAN NOT NULL DEFAULT false,
    escalated_at TIMESTAMP,
    resolved_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_chatbot_sessions_conversation ON chatbot_sessions (conversation_id);
CREATE INDEX IF NOT EXISTS idx_chatbot_sessions_status ON chatbot_sessions (status);
