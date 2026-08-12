# TimeGoBetter API

API Go para la gestión de órdenes de lavandería. Requiere Go 1.26+, Docker Desktop y PostgreSQL 16.

## Inicio local

1. Copia `.env.example` a `.env` y reemplaza los secretos de ejemplo. `JWT_SECRET` debe tener al menos 32 caracteres.
2. Inicia PostgreSQL:

```powershell
docker compose up -d
docker compose ps
```

3. Ejecuta las validaciones y servidor:

```powershell
go test ./...
go vet ./...
go run ./cmd/server
```

La API escucha en `http://localhost:8080`, Swagger en `http://localhost:8080/swagger/index.html` y WebSocket en `ws://localhost:8080/ws`. El WebSocket requiere el JWT del cliente mediante `?access_token=<token>`.

Para detener la BD sin borrar datos: `docker compose down`. Para reiniciar datos de desarrollo: `docker compose down -v` (destructivo).

## Variables requeridas

`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `JWT_SECRET`, `OPERATOR_KEY`, `PORT` y `WS_ALLOWED_ORIGINS`.
