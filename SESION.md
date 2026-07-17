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

### Crear pedido desde la UI (P0)
- **Formulario de alta de pedidos** en `/orders/new` (botón "+ Nuevo pedido" en `/orders`).
- Servicio `createOrder` (`services/orders/create.ts`, `POST /orders`): llena `created_by`
  con la sesión, revalida `/orders` y `/home`.
- Servicios de planchas `fetchMaterials` (`GET /materials`) y `calculatePrice`
  (`GET /planchas/:id/calculate-price`) en `services/planchas/` + tipo `CalculatePriceResult`.
- Componente client `components/orders/order-form.tsx`: cliente (texto libre), urgencia,
  notas, ítems dinámicos (plancha/material/cantidad) con **cálculo de precio en vivo** por
  ítem y total en ARS; snackbar de éxito y redirect a `/orders`.

### Página de detalle de pedido (P0)
- **`/orders/[order_id]`**: servicio `fetchOrder` (`services/orders/fetch-one.ts`,
  `GET /orders/:orderID`). El N° de pedido en la tabla enlaza al detalle.
- Muestra ítems (plancha/material resueltos a nombre vía `fetchPlanchas`+`fetchMaterials`,
  cantidad, precio unit., subtotal), total, panel de info (estado, cliente, origen,
  urgencia, asignado), fechas (creado/actualizado/completado) y notas.
- Cambio de estado en línea reusando `OrderStatusActions`.
- Helper `formatDateTime` (locale es-AR) en `utils/status.ts`.

### Mutaciones de estado en la UI (P0 parcial)
- **Cambiar estado de pedido**: servicio `updateOrderStatus` (`services/orders/update-status.ts`,
  `PUT /orders/:id/status`) + dropdown de acciones en `components/orders/status-actions.tsx`
  (columna "Acciones" en la tabla de pedidos). Solo ofrece transiciones válidas.
- **Avanzar producción**: servicio `updatePrintJobStatus` (`services/printing/update-status.ts`,
  `PUT /print-jobs/:id/status`) + botones en `components/printing/job-actions.tsx`
  (integrados en `queue.tsx`). Transiciones según `isValidTransition` del backend.
- Mapas de transiciones en `utils/status.ts`: `getNextOrderStatuses`,
  `getNextPrintJobStatuses` (reflejan las reglas del backend).

### Skeletons de carga + microinteracciones (solo Tailwind)
- Keyframe `fadeIn` (opacity + translateY) en `tailwind.config.ts` → clase `animate-fadeIn`.
- Skeletons con `animate-pulse`: `kpi-card/skeleton.tsx`, y en `components/common/skeletons/`
  (`orders-table-skeleton`, `production-queue-skeleton`, `table-skeleton` genérico,
  `page-header-skeleton`).
- Home usa `<Suspense>` por sección (streaming con fetch real): componentes async
  `components/home/` (`kpi-section`, `recent-orders`, `production-queue`).
- `loading.tsx` por ruta en orders / printing / stickers / users (App Router los
  muestra durante la navegación). Se quitaron los `<Suspense>` con "Cargando..." inútiles.
- Fade-in escalonado por fila/tarjeta (`animationDelay: i*40ms`) en todas las tablas
  y colas. Pulso sutil en badge para estados `in_production` / `printing` / `cutting`
  (`StatusBadge` acepta prop `pulse`).

### Bug fix: paginación devolvía 0 en KPIs
- **Causa raíz:** `domain.Paging` (Go) no tenía tags JSON → serializaba `Total`/`Limit`/`Offset`
  (mayúscula) pero el tipo TS `SearchResponse.paging` espera minúscula → `paging.total`
  quedaba `undefined` → KPIs mostraban 0.
- **Fix:** tags JSON `json:"total"` etc. en `backend/internal/domain/paging.go`.
  Corrige orders, planchas y printing (users ya usaba su propio `PagingResponseDTO`).

### Módulo de Clientes (P0)
- **Backend completo (CRUD REST)**: migración `039_create_customers_table` (tabla `customers`
  con name/phone/email/document/address/notes/is_active + timestamps, índices por
  name/phone/is_active). `domain/customer.go`, `application/customer_service.go`,
  `infra/repository/database/customer_repository.go` + `customer_dto.go`,
  `rest/customer_controller.go` + `customer_dto.go`. Wiring en `runner.go`, rutas en
  `routes_stickers.go`: `POST/GET /customers`, `GET/PUT/DELETE /customers/:customerID`.
  Búsqueda por texto (`?search=` sobre name OR phone OR email). Probado end-to-end
  (create/read/search/update/delete OK).
- **Frontend**: tipo `Customer` + servicios `services/customers/` (search, fetch-one,
  create, update, delete). Página `/customers` con `CustomersManager` (tabla + búsqueda +
  modal de alta/edición + borrado con confirmación) y link "Clientes" en el sidebar.
- **Integración con pedidos**: `OrderForm` recibe `customers` y el campo cliente usa
  `<datalist>` (autocompletado). `customer_id` en `orders` SIGUE siendo texto libre (se
  guarda el nombre elegido) — **FK `orders.customer_id → customers` queda pendiente**
  (requiere migración de datos de los pedidos existentes).

### Repositorio Git
- Subido a **https://github.com/Leandro96JG/crm-platform** (rama `main`).
- `.gitignore` raíz (excluye `node_modules`, `.env`, builds Go/Next, logs).
- Plantillas `backend/.env.example` y `frontend/.env.example` (sin secretos).
- `README.md` raíz con setup Docker, credenciales, scripts y flujo multi-PC.
- **Los `.env` reales NO están versionados** — recrear desde `.env.example` en cada PC.
  Para Docker no hacen falta (el `docker-compose.yml` trae las env vars inline).

---

## 🚧 Pendiente / Próximas funcionalidades

> El dashboard hoy es **casi 100% de solo lectura** (única mutación real: el
> propio perfil). El mayor valor está en habilitar las acciones CRUD y cerrar
> los gaps de datos.

### P0 — Crítico (hace el CRM realmente usable)
- [ ] **Mutaciones de pedidos en la UI**
  - [x] Crear pedido (formulario: cliente, items con plancha+material+cantidad). ✅
  - [x] Cambiar estado del pedido (dropdown de acciones con transiciones válidas). ✅
  - [x] Página de detalle `/orders/[order_id]` con items (timeline de estado pendiente). ✅
  - [x] Servicio frontend de mutación `updateOrderStatus` (`PUT /orders/:id/status`). ✅
- [x] **Avanzar producción desde la UI** ✅
  - [x] Botones en la cola según transiciones válidas (imprimir/cortar/fallido/etc.). ✅
  - [x] Servicio `updatePrintJobStatus` (`PUT /print-jobs/:id/status`) desde el frontend. ✅
- [ ] **Módulo de Clientes** (gap de datos grande)
  - [x] Backend: tabla `customers`, dominio, repo, service, CRUD REST. ✅
  - [x] Frontend: página de clientes + selección al crear pedido (datalist). ✅
  - [ ] FK `orders.customer_id → customers` (hoy sigue texto libre; requiere migrar datos).
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
- [ ] Revisar componentes huérfanos restantes (`Modal` ahora también lo usa el módulo
  de Clientes).
- [ ] **FK `orders.customer_id → customers`**: migrar los pedidos existentes (customer_id
  texto libre → id real) y refactorizar el alta de pedido a `customer_id`.

---

## 🎯 Próximo paso sugerido
Pedidos (crear/detalle/estado), producción y Clientes (CRUD + UI) ya están hechos.
Los siguientes pasos de mayor valor:
1. **Arreglar enlace roto** `/users/[user_id]` (detalle/edición o quitar el botón "ver").
2. **FK `orders.customer_id → customers`**: migrar datos de pedidos existentes y usar
   `customer_id` real en el alta de pedidos (hoy se guarda el nombre como texto).
3. **Timeline de estado** en el detalle de pedido (historial de transiciones; requiere
   auditoría en el backend — tabla `order_status_history`).
