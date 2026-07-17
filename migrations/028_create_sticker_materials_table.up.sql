CREATE TABLE IF NOT EXISTS sticker_materials (
    material_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT DEFAULT '',
    material_type TEXT NOT NULL DEFAULT 'vinyl',
    finish TEXT NOT NULL DEFAULT 'mate',
    is_cuttable BOOLEAN NOT NULL DEFAULT true,
    is_printable BOOLEAN NOT NULL DEFAULT true,
    base_cost DECIMAL(10,2) NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_sticker_materials_type ON sticker_materials (material_type);
CREATE INDEX IF NOT EXISTS idx_sticker_materials_active ON sticker_materials (is_active);
