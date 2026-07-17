CREATE TABLE IF NOT EXISTS whatsapp_conversations (
    conversation_id TEXT PRIMARY KEY,
    customer_id TEXT,
    whatsapp_phone TEXT NOT NULL,
    wa_conversation_id TEXT DEFAULT '',
    status TEXT NOT NULL DEFAULT 'active',
    ai_handled BOOLEAN NOT NULL DEFAULT true,
    ai_agent_session_id TEXT DEFAULT '',
    last_message_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_whatsapp_conv_phone ON whatsapp_conversations (whatsapp_phone);
CREATE INDEX IF NOT EXISTS idx_whatsapp_conv_status ON whatsapp_conversations (status);
CREATE INDEX IF NOT EXISTS idx_whatsapp_conv_customer ON whatsapp_conversations (customer_id);
