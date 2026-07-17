# Sesión CRM "Viva" — Taller de Stickers

## Objetivo
CRM para un taller de venta de stickers por plancha (A4), con captación por WhatsApp
+ agente IA (GPT-4o-mini) y flujo de impresión/corte manual (Epson L3210 + plotter).
Marca: **Viva**. Moneda: **peso argentino (ARS, locale es-AR)**.

## Stack
- **Backend:** Go 1.26, Gin, PostgreSQL 16, SQLX
- **Frontend:** Next.js 16 (App Router), Node 20, NextAuth, Tailwind 3.4, HeroUI, Heroicons
- **APIs externas:** Meta WhatsApp Cloud API, OpenAI
- **Infra:** Docker Compose (postgres + backend + frontend), AWS S3 (attachments)

## Convenciones de trabajo
- Tras editar backend/frontend: `docker compose up -d --build <servicio>` (¡el
  `build` solo reconstruye la imagen, `--build` con `up` recrea el contenedor!).
- Recargar navegador con Ctrl+Shift+R (caché).
- Login: `admin@crmstickers.com` / `admin123`.
- Migraciones: se copian a la imagen del backend y corren solas al arrancar
  (golang-migrate, tabla `schema_migrations`).

---

## ✅ Completado

### Migración de dominio (seguros → stickers)
- Eliminado TODO el código de la app vieja de seguros (casos, aseguradoras,
  clientes, técnicos, pagos, colas, TV). Migración `037_drop_insurance_tables`.
- Traducción completa de portugués a español.

### Rediseño visual "Viva"
- Sistema de tema claro/oscuro (default oscuro) con tokens CSS en `global.css`
  (`paper`, `ink`, `cut`, `teal`, `st-*`), `ThemeProvider`, toggle y anti-flash.
- Logo Viva, sidebar oscuro, Topbar, KpiCard, PanelCard, StatusBadge, iconos.
- Login rediseñado (pantalla dividida, marca Viva).
- Páginas rediseñadas: Home (4 KPIs + pedidos recientes + cola), Pedidos,
  Planchas, Producción, Usuarios, Mi Perfil.
- Eliminado código muerto: `acme-logo`, `card`, `badge`, `dropdown`, `tooltip`,
  `search`, `parser.ts` (moneda BRL + docs CPF/CNPJ brasileños).

### Backend actual (funcional)
- **Auth:** `POST /login` (JWT HS256, 24h), `POST /logout` (stub).
- **Usuarios:** CRUD completo + cambio de contraseña. DTO de respuesta con tags
  JSON snake_case (`user_id`, `first_name`, `last_name`, `active`);
  `password_hash` ya NO se expone.
- **Planchas:** CRUD + precios + cálculo de precio.
- **Materiales:** solo `POST` y `GET` (sin editar/borrar).
- **Pedidos:** crear, buscar, ver, cambiar estado (auto-genera print jobs).
- **Producción:** `GET /print-jobs`, cambiar estado (transiciones validadas).
- **WhatsApp/IA:** webhook + envío + agente OpenAI (requiere credenciales).

### Datos de prueba
- Migración `038_seed_test_data`: 6 materiales, 32 planchas, 40 pedidos + items,
  40 trabajos de producción, 24 usuarios. Idempotente y reversible.

### Moneda
- `ARS` / `es-AR` en `utils/status.ts`. Eliminada toda referencia a Brasil/BRL.

---

## 🚧 Pendiente / Próximas funcionalidades

> El dashboard hoy es **casi 100% de solo lectura** (única mutación real: el
> propio perfil). El mayor valor está en habilitar las acciones CRUD y cerrar
> los gaps de datos.

### P0 — Crítico (hace el CRM realmente usable)
- [ ] **Mutaciones de pedidos en la UI**
  - [ ] Crear pedido (formulario: cliente, items con plancha+material+cantidad).
  - [ ] Cambiar estado del pedido (dropdown/acciones: pending→approved→…→delivered).
  - [ ] Página de detalle `/orders/[order_id]` con items y timeline de estado.
  - [ ] Servicios frontend de mutación (POST/PUT) en `services/orders`.
- [ ] **Avanzar producción desde la UI**
  - [ ] Botones en la cola: imprimir→impreso, cortar→cortado, marcar fallido.
  - [ ] Servicio `PUT /print-jobs/:id/status` desde el frontend.
- [ ] **Módulo de Clientes** (gap de datos grande)
  - [ ] Backend: tabla `customers`, dominio, repo, service, CRUD REST.
  - [ ] FK `orders.customer_id → customers` (hoy es texto libre).
  - [ ] Frontend: página de clientes + selección al crear pedido.
- [ ] **Arreglar enlace roto** `/users/[user_id]`: crear página de detalle/edición
  de usuario o quitar el botón "ver".

### P1 — Alta (completar el catálogo y la navegación)
- [ ] **CRUD de planchas y materiales en la UI**
  - [ ] Crear/editar/borrar planchas; subir preview a S3.
  - [ ] Backend: `GET/PUT/DELETE /materials/:id` (faltan; service sin Update/Delete).
  - [ ] Editar/borrar precios (`PUT/DELETE`).
- [ ] **Paginación en todas las tablas**: `/orders` y `/stickers` reciben `paging`
  del backend pero no la renderizan (solo `/users` la muestra).
- [ ] **Filtros por estado/categoría con UI** (hoy solo vía `?status=`/`?category=`
  manual en la URL): dropdowns/chips en Pedidos, Planchas.
- [ ] **Búsqueda funcional en el Topbar** (hoy el input es decorativo, sin lógica).

### P2 — Media (valor de negocio)
- [ ] **Dashboard / KPIs reales**: endpoint backend de métricas (ingresos, pedidos
  por estado, producción del día) en vez de calcular en el frontend.
- [ ] **Notificaciones WhatsApp de estado**: `SendOrderNotification` existe pero no
  se invoca al cambiar el estado del pedido; conectarlo.
- [ ] **Vista de conversaciones WhatsApp / chatbot**: no hay módulo de chat ni
  historial (el modelo ya distingue origen WhatsApp/IA).
- [ ] **Vista de TV / pantalla de producción**: existe `components/tv/auto-refresh.tsx`
  huérfano; crear página `tv/` con la cola en tiempo real para el taller.
- [ ] **Flujo de primer login**: `FirstLoginModal` está creado pero no se monta;
  el backend `UpdateUserDTO` no soporta `password`/`first_login_completed`.

### P3 — Baja / técnica
- [ ] **Logout real**: invalidar/expirar el JWT (hoy es stub sin revocación).
- [ ] **Paginación consistente en backend**: `GET /materials` ignora query params
  (limit hardcodeado 100); precios no paginan.
- [ ] **Cálculo de total del pedido**: `CreateOrder` ignora `bulk_discount` de
  `plancha_prices` (inconsistencia con `CalculatePrice`).
- [ ] **Asignación de pedidos**: `assigned_to` existe en el modelo pero no hay
  endpoint ni UI para asignar operador.
- [ ] Revisar componentes huérfanos restantes (`Modal` solo lo usa el modal de
  primer login no montado).

---

## 🎯 Próximo paso sugerido
Empezar por **P0 — cambiar estado de pedidos y avanzar producción desde la UI**:
es el flujo diario del taller, ya existe el backend (`PUT /orders/:id/status` y
`PUT /print-jobs/:id/status`), solo falta cablear servicios de mutación + botones.
