# CRM Viva

CRM para taller de stickers (Viva). Gestión de pedidos, catálogo de planchas,
cola de producción (impresión/corte), usuarios e integración con WhatsApp/IA.

Interfaz en español, moneda en pesos argentinos (ARS).

## Stack

| Capa | Tecnología |
|------|-----------|
| Backend | Go 1.26 · Gin · SQLX · PostgreSQL 16 |
| Frontend | Next.js 16 (App Router) · Node 20 · NextAuth · Tailwind 3.4 · HeroUI · Heroicons |
| Infra | Docker Compose (postgres · backend · frontend) · AWS S3 |

## Estructura del repositorio

```
backend/       API en Go (Gin), dominio y repositorios
frontend/      App Next.js (App Router)
migrations/    Migraciones SQL (numeradas)
agents/        Definiciones de agentes de IA (ventas / producción)
hardware/      Notas de setup de impresora y plotter
diseño/        Mockups HTML/CSS del dashboard
docker-compose.yml
```

## Requisitos

- [Docker Desktop](https://www.docker.com/) (con Docker Compose)
- [Node.js 20+](https://nodejs.org/) (solo si vas a correr el frontend fuera de Docker)
- [Go 1.26+](https://go.dev/) (solo si vas a correr el backend fuera de Docker)

## Puesta en marcha (Docker — recomendado)

El `docker-compose.yml` ya trae la configuración necesaria para desarrollo local,
así que **no hace falta crear archivos `.env` para levantar con Docker**.

```bash
git clone https://github.com/Leandro96JG/crm-platform.git
cd crm-platform
docker compose up -d --build
```

Servicios expuestos:

| Servicio | URL |
|----------|-----|
| Frontend | http://localhost:3000 |
| Backend (API) | http://localhost:8080 |
| PostgreSQL | localhost:5432 |

Las migraciones se aplican automáticamente al iniciar el backend.

### Credenciales de acceso

- **Usuario:** `admin@crmstickers.com`
- **Contraseña:** `admin123`

> Datos de prueba: la migración `038_seed_test_data` carga materiales, planchas,
> pedidos y trabajos de producción de ejemplo.

## Flujo de trabajo con Docker

Reconstruir un servicio tras cambiar código:

```bash
docker compose up -d --build frontend   # o backend
```

> **Importante:** `docker compose build` solo reconstruye la imagen; para ver los
> cambios hay que `up -d --build <servicio>` (recrea el contenedor). En el navegador,
> recargar con `Ctrl+Shift+R`.

Otros comandos útiles:

```bash
docker compose ps                # estado de contenedores
docker compose logs -f frontend  # logs en vivo
docker compose down              # detener todo
docker compose down -v           # detener y borrar la base de datos
```

## Desarrollo fuera de Docker (opcional)

Si preferís correr el frontend o el backend directamente, necesitás los archivos
de entorno. Copiá las plantillas y completá los valores:

```bash
# Frontend
cp frontend/.env.example frontend/.env.local

# Backend
cp backend/.env.example backend/.env
```

Frontend:

```bash
cd frontend
npm install
npm run dev      # http://localhost:3000
```

Backend:

```bash
cd backend
go run ./cmd/api
```

## Scripts del frontend

| Comando | Descripción |
|---------|-------------|
| `npm run dev` | Servidor de desarrollo |
| `npm run build` | Build de producción |
| `npm run start` | Servir el build |
| `npm run lint` | ESLint |
| `npm run test` | Tests (Jest) |

## Trabajar desde varias PCs

```bash
git pull                          # antes de empezar a editar
# ...cambios...
git add -A
git commit -m "descripción"
git push
```

> Los archivos `.env` **no** están versionados (contienen secretos). Recreálos en
> cada máquina a partir de los `.env.example`. Para desarrollo con Docker no son
> necesarios; solo si corrés los servicios fuera de Docker.
