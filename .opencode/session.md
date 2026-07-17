# CRM Platform - Sesión

## Proyecto

CRM para gestión de siniestros, órdenes y producción de planchas (stickers).  
Stack: Go (backend) + Next.js (frontend) + PostgreSQL.

---

## Flujo del sistema

### 1. Login / Autenticación
- Endpoint: `POST /crm/core/api/v1/login` (backend)
- Frontend: formulario en `/login` con email + contraseña
- LoginDTO usa `email` como campo de login
- `authService.Login()` busca por email primero, luego por username
- Retorna JWT token + datos del usuario
- Las rutas protegidas usan middleware JWT

### 2. Gestión de Casos (Siniestros)
- CRUD de casos con estados: Borrador → Nuevo → Información del cliente → Esperando técnico → En curso → Informe → Pago → Comprobante → Cerrado / Cancelado
- Cada estado tiene un formulario específico (`form-details/`)
- Asignación de técnico, productos, transacciones
- Adjuntos (imágenes) por caso
- Comentarios por caso
- Laudo (informe) descargable

### 3. Clientes, Técnicos y Aseguradoras
- **Clientes**: CRUD básico con documento, dirección, contacto
- **Técnicos (Partners)**: CRUD con tipo (Montador/Tapicero), clave PIX, datos de pago
- **Aseguradoras (Contractors)**: CRUD con razón social, CNPJ, contacto

### 4. Pagos y Transacciones
- Transacciones de tipo MO (mano de obra), desplazamiento, piezas
- Roles: técnico (outgoing) vs aseguradora (incoming)
- Confirmación de pago con comprobante adjunto

### 5. Panel de Control Interno
- Tabla resumen con datos de casos, técnicos, valores
- Filtros por mes, año, estado, aseguradora, técnico

### 6. Gamificación / Dashboards
- KPIs: Casos finalizados, SLA a tiempo
- Gráfico de evolución de desempeño (TMA)
- Ranking de atendimiento
- Logros y recompensas
- Bono de asistencia (AttendanceBonus)
- Filtro de período (hoy/semana/mes)

### 7. Órdenes y Producción (Planchas/Stickers)
- Catálogo de planchas (materiales, precios)
- Pedidos con items
- Cola de impresión y corte

---

## Cambios realizados en esta sesión

### Login (Backend)
- LoginDTO cambió de `username` a `email`
- `authService.Login()` ahora busca por email primero, luego username
- `userService.Authenticate()` también busca por email primero
- Seed migration puesta como idempotente (`ON CONFLICT DO NOTHING`)
- Arreglado estado dirty de migración 036 en `schema_migrations`

### Traducción Portugués → Español (Frontend)
- **~130 archivos modificados**
- Locales: `date-fns` (es), `Intl.NumberFormat` (es-ES), HTML lang (es)
- Tipos: `caseStatusMap`, `casePriorityMap`, `paymentOptionMap`, `TransactionDescMap`
- Componentes: login, navegación, casos, clientes, técnicos, aseguradoras, pagos, panel, dashboards, usuarios, perfil
- Servicios (~45 archivos): mensajes de error/éxito en español
- Uppy: locale cambiado a `es_ES`
- Tests (~20): assertions actualizadas

---

## Pendientes / Notas
- El backend se comunica con la API en inglés/portugués; los valores de descripción de transacciones vienen de la API en portugués (MO, Peças, Deslocamento, etc.) y se mapean a español en el frontend
- Los datos de prueba (seed) usan el hash bcrypt de "admin123"
- Usuario admin por defecto: `admin@crmstickers.com` / `admin123`
