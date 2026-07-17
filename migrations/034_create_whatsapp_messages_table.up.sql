CREATE TABLE IF NOT EXISTS whatsapp_messages (
    message_id TEXT PRIMARY KEY,
    conversation_id TEXT NOT NULL REFERENCES whatsapp_conversations(conversation_id),
    wa_message_id TEXT DEFAULT '',
    direction TEXT NOT NULL DEFAULT 'inbound',
    message_type TEXT NOT NULL DEFAULT 'text',
    content TEXT DEFAULT '',
    media_url TEXT DEFAULT '',
    metadata JSONB DEFAULT '{}',
    ai_processed BOOLEAN NOT NULL DEFAULT false,
    ai_response TEXT DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_whatsapp_msg_conversation ON whatsapp_messages (conversation_id);
CREATE INDEX IF NOT EXISTS idx_whatsapp_msg_direction ON whatsapp_messages (direction);
CREATE INDEX IF NOT EXISTS idx_whatsapp_msg_created ON whatsapp_messages (created_at);
