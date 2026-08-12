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

La API escucha en `http://localhost:8080`, Swagger en `http://localhost:8080/swagger/index.html`, healthcheck en `http://localhost:8080/health` y WebSocket en `ws://localhost:8080/ws`. El WebSocket requiere el JWT del cliente mediante `?access_token=<token>`.

Para detener la BD sin borrar datos: `docker compose down`. Para reiniciar datos de desarrollo: `docker compose down -v` (destructivo).

## Datos de demostraciÃ³n

En una base de datos nueva, las migraciones cargan el catÃ¡logo base de lavado y secado, planchado y lavado en seco, con precios y descripciones para que el sistema pueda demostrarse de inmediato. Para aplicar migraciones nuevas sobre una base creada antes de esta versiÃ³n, ejecuta los scripts de `infrastructure/database/migrations` en orden o recrea Ãºnicamente el volumen de desarrollo si no contiene datos que quieras conservar.

## OperaciÃ³n y trazabilidad

Las Ã³rdenes siguen el flujo `recibida -> en_proceso -> lista -> entregada`; no se permiten saltos ni retrocesos. Cada cambio queda registrado en el historial de estados con su autor, rol y fecha. El endpoint de health devuelve `503` si PostgreSQL no responde, para que un monitor pueda detectar una dependencia caÃ­da.

## Variables requeridas

`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `JWT_SECRET`, `OPERATOR_KEY`, `PORT` y `WS_ALLOWED_ORIGINS`.
