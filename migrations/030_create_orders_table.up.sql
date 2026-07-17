CREATE TABLE IF NOT EXISTS orders (
    order_id TEXT PRIMARY KEY,
    order_number TEXT NOT NULL UNIQUE,
    customer_id TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    source TEXT NOT NULL DEFAULT 'whatsapp',
    ai_agent_id TEXT DEFAULT '',
    ai_handled BOOLEAN NOT NULL DEFAULT false,
    assigned_to TEXT DEFAULT '',
    notes TEXT DEFAULT '',
    total DECIMAL(12,2) NOT NULL DEFAULT 0,
    urgency TEXT NOT NULL DEFAULT 'normal',
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    created_by TEXT NOT NULL DEFAULT '',
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_by TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_orders_status ON orders (status);
CREATE INDEX IF NOT EXISTS idx_orders_customer ON orders (customer_id);
CREATE INDEX IF NOT EXISTS idx_orders_urgency ON orders (urgency);
CREATE INDEX IF NOT EXISTS idx_orders_created ON orders (created_at);
