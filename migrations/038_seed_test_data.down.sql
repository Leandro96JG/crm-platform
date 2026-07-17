DELETE FROM print_jobs WHERE job_id LIKE 'seed-%';
DELETE FROM order_items WHERE item_id LIKE 'seed-%';
DELETE FROM orders WHERE order_id LIKE 'seed-%';
DELETE FROM planchas WHERE plancha_id LIKE 'seed-%';
DELETE FROM sticker_materials WHERE material_id LIKE 'seed-%';
DELETE FROM users WHERE username LIKE 'operador%';
