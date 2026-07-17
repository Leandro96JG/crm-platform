---
name: stickers
description: "Módulo de stickers y planchas para impresión y corte"
---

## Módulo Stickers

### Estructura
- `backend/internal/domain/plancha.go` - Modelo Plancha
- `backend/internal/domain/order.go` - Modelo Order
- `backend/internal/domain/print_job.go` - Modelo PrintJob
- `backend/internal/application/plancha_service.go` - Lógica de planchas
- `backend/internal/application/order_service.go` - Lógica de órdenes
- `backend/internal/application/printing_service.go` - Lógica de cola
- `backend/internal/infra/repository/database/plancha_repository.go` - Persistencia
- `backend/internal/infra/repository/database/order_repository.go` - Persistencia
- `backend/internal/infra/repository/database/print_job_repository.go` - Persistencia
- `backend/internal/infra/entrypoint/rest/plancha_controller.go` - Endpoints
- `backend/internal/infra/entrypoint/rest/order_controller.go` - Endpoints
- `backend/internal/infra/entrypoint/rest/printing_controller.go` - Endpoints

### Tablas DB
- `planchas` - Catálogo de productos (planchas/stickers)
- `sticker_materials` - Materiales disponibles
- `plancha_prices` - Precios por plancha + material
- `orders` - Órdenes de producción
- `order_items` - Items de cada orden
- `print_jobs` - Cola de trabajos manuales

### Flujo
1. Cliente pide plancha del catálogo vía WhatsApp/Web
2. AI o admin crea orden (status: pending)
3. Admin aprueba orden (status: approved)
4. Sistema crea print_jobs automáticamente
5. Admin imprime con Epson → marca "printed"
6. Admin corta con plotter → marca "cut"
7. Admin empaqueta → marca "ready"
8. Sistema notifica al cliente
9. Admin marca "delivered"

### Estados de orden
pending → approved → in_production → ready → delivered
                                                      → cancelled (en cualquier momento)

### Notas
- No se guardan posiciones XY de stickers
- El diseño se maneja como archivo PDF completo (layout_file_url)
- La impresión y corte son 100% manuales, el sistema solo organiza la cola
