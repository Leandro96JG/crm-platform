-- Datos de prueba para probar listados y paginado.
-- Todos los IDs generados usan el prefijo 'seed-' para poder revertirlos.

-- ── Materiales ──────────────────────────────────────────────
INSERT INTO sticker_materials (material_id, name, description, material_type, finish, is_cuttable, is_printable, base_cost)
VALUES
    ('seed-mat-01', 'Vinilo blanco brillante', 'Vinilo estándar', 'vinyl', 'brillante', true, true, 12.50),
    ('seed-mat-02', 'Vinilo transparente', 'Fondo transparente', 'vinyl', 'mate', true, true, 15.00),
    ('seed-mat-03', 'Papel adhesivo', 'Interior', 'paper', 'mate', true, true, 8.00),
    ('seed-mat-04', 'Holográfico', 'Efecto arcoíris', 'holographic', 'brillante', true, true, 22.00),
    ('seed-mat-05', 'Vinilo laminado', 'Resistente al agua', 'vinyl', 'satinado', true, true, 18.00),
    ('seed-mat-06', 'Poliéster metalizado', 'Acabado plata', 'polyester', 'metalizado', true, true, 25.00)
ON CONFLICT (material_id) DO NOTHING;

-- ── Planchas (32) ───────────────────────────────────────────
INSERT INTO planchas (plancha_id, name, description, category, subcategory, preview_image_url, is_active, created_by, updated_by)
SELECT
    'seed-plancha-' || lpad(g::text, 3, '0'),
    (ARRAY['Gatitos Kawaii','Frases Motivadoras','Logos Deportivos','Flores Vintage','Memes Populares',
           'Astronautas','Comida Rápida','Plantas Suculentas','Superhéroes','Mandalas'])[1 + (g % 10)] || ' #' || g,
    'Plancha de diseño para stickers troquelados lote ' || g,
    (ARRAY['single','pack','custom','premium'])[1 + (g % 4)],
    (ARRAY['infantil','adultos','empresas','eventos','navidad'])[1 + (g % 5)],
    'https://picsum.photos/seed/plancha' || g || '/200/200',
    (g % 7 <> 0),
    '00000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000001'
FROM generate_series(1, 32) AS g
ON CONFLICT (plancha_id) DO NOTHING;

-- ── Pedidos (40) ────────────────────────────────────────────
INSERT INTO orders (order_id, order_number, customer_id, status, source, ai_handled, total, urgency, created_at, created_by, updated_by)
SELECT
    'seed-order-' || lpad(g::text, 4, '0'),
    'PED-2026-' || lpad(g::text, 4, '0'),
    (ARRAY['María González','Juan Pérez','Ana Silva','Carlos Ruiz','Lucía Fernández',
           'Pedro Martín','Sofía López','Diego Torres','Valentina Díaz','Mateo Romero'])[1 + (g % 10)],
    (ARRAY['pending','approved','in_production','ready','delivered','cancelled'])[1 + (g % 6)],
    (ARRAY['whatsapp','manual'])[1 + (g % 2)],
    (g % 2 = 0),
    round((45 + (g * 13.7) + (g % 5) * 20)::numeric, 2),
    (ARRAY['normal','normal','normal','urgent'])[1 + (g % 4)],
    now() - (g || ' hours')::interval,
    '00000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000001'
FROM generate_series(1, 40) AS g
ON CONFLICT (order_id) DO NOTHING;

-- ── Items de pedido (1 por pedido) ──────────────────────────
INSERT INTO order_items (item_id, order_id, plancha_id, material_id, sheet_quantity, unit_price, subtotal, sort_order)
SELECT
    'seed-item-' || lpad(g::text, 4, '0'),
    'seed-order-' || lpad(g::text, 4, '0'),
    'seed-plancha-' || lpad((1 + (g % 32))::text, 3, '0'),
    'seed-mat-0' || (1 + (g % 6)),
    1 + (g % 5),
    round((10 + (g % 20))::numeric, 2),
    round(((1 + (g % 5)) * (10 + (g % 20)))::numeric, 2),
    0
FROM generate_series(1, 40) AS g
ON CONFLICT (item_id) DO NOTHING;

-- ── Trabajos de producción (40) ─────────────────────────────
INSERT INTO print_jobs (job_id, order_item_id, job_type, status, queue_position, notes, copies, created_by)
SELECT
    'seed-job-' || lpad(g::text, 4, '0'),
    'seed-item-' || lpad(g::text, 4, '0'),
    job_type,
    CASE WHEN job_type = 'cut'
         THEN (ARRAY['queued','cutting','cut','failed'])[1 + (g % 4)]
         ELSE (ARRAY['queued','printing','printed','failed'])[1 + (g % 4)]
    END,
    g,
    'Pedido PED-2026-' || lpad(g::text, 4, '0'),
    1 + (g % 3),
    '00000000-0000-0000-0000-000000000001'
FROM (
    SELECT g, (ARRAY['print','cut'])[1 + (g % 2)] AS job_type
    FROM generate_series(1, 40) AS g
) AS s
ON CONFLICT (job_id) DO NOTHING;

-- ── Usuarios (24) ───────────────────────────────────────────
-- password_hash corresponde a 'admin123' (mismo hash que el admin seed).
INSERT INTO users (user_id, username, password_hash, name, email, role, is_active, created_by, updated_by)
SELECT
    ('00000000-0000-0000-0000-0000000010' || lpad(g::text, 2, '0'))::uuid,
    'operador' || lpad(g::text, 2, '0'),
    '$2a$10$rmKmO2zXrigg/ImKKjjQ1uAjkAi3216rfuuM4c0xsurOiRO.mIo9u',
    (ARRAY['Laura','Martín','Camila','Andrés','Paula','Tomás','Elena','Nicolás','Julia','Bruno','Carla','Iván'])[1 + (g % 12)]
        || ' ' ||
    (ARRAY['García','Rodríguez','Martínez','López','Sánchez','Pérez','Gómez','Fernández'])[1 + (g % 8)],
    'operador' || lpad(g::text, 2, '0') || '@crmstickers.com',
    (ARRAY['operator','admin','admin_operator'])[1 + (g % 3)],
    (g % 6 <> 0),
    'system',
    'system'
FROM generate_series(1, 24) AS g
ON CONFLICT (user_id) DO NOTHING;
