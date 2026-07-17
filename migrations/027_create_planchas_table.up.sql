CREATE TABLE IF NOT EXISTS planchas (
    plancha_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT DEFAULT '',
    category TEXT NOT NULL DEFAULT 'single',
    subcategory TEXT DEFAULT '',
    layout_file_url TEXT DEFAULT '',
    preview_image_url TEXT DEFAULT '',
    notes TEXT DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_by TEXT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_planchas_category ON planchas (category);
CREATE INDEX IF NOT EXISTS idx_planchas_active ON planchas (is_active);
