# Sistema de Gestión de Órdenes — ITLA Proyecto Final

## Requisitos
- Go 1.22+
- Node.js 18+
- Docker Desktop (recomendado) o PostgreSQL instalado

## Cómo correr el proyecto

### 1. Clonar el repositorio
git clone <url-del-repo>

### 2. Base de datos
Con Docker (recomendado):
docker compose up -d

Sin Docker: crear la base de datos manualmente
y correr infrastructure/database/migrations/001_init.sql

### 3. Backend
cd ITLAFINAL
cp .env.example .env   ← llenar con tus credenciales
go run cmd/server/main.go

## URLs
- Backend:  http://localhost:8080
- API Docs: http://localhost:8080/swagger
