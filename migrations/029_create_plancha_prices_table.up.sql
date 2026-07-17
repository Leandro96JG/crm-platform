CREATE TABLE IF NOT EXISTS plancha_prices (
    price_id TEXT PRIMARY KEY,
    plancha_id TEXT NOT NULL REFERENCES planchas(plancha_id),
    material_id TEXT NOT NULL REFERENCES sticker_materials(material_id),
    base_price DECIMAL(10,2) NOT NULL DEFAULT 0,
    min_quantity INT NOT NULL DEFAULT 1,
    bulk_discount JSONB DEFAULT '[]',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE(plancha_id, material_id)
);

CREATE INDEX IF NOT EXISTS idx_plancha_prices_plancha ON plancha_prices (plancha_id);
CREATE INDEX IF NOT EXISTS idx_plancha_prices_material ON plancha_prices (material_id);
