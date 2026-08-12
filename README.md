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

## AdministraciÃ³n

Define `ADMIN_EMAIL` y `ADMIN_PASSWORD` (y opcionalmente `ADMIN_NAME`) antes de iniciar la API para aprovisionar la Ãºnica cuenta administrativa inicial. El registro pÃºblico no puede crear administradores. Desde `/api/admin/users`, un administrador puede buscar, paginar, bloquear, desbloquear o eliminar cuentas sin Ã³rdenes; el bloqueo invalida las sesiones ya emitidas.

## Cobros y factura de demostraciÃ³n

El cobro es una **simulaciÃ³n acadÃ©mica**: no contacta pasarelas, no procesa tarjetas y no almacena informaciÃ³n financiera. Solo registra el mÃ©todo elegido y crea un comprobante interno para la demostraciÃ³n. Una orden puede pagarse una sola vez y exclusivamente cuando su estado es `lista`.

La factura PDF es un comprobante visual generado por el navegador; no es una factura fiscal y no incluye NCF, RNC ni validez tributaria. Para operar comercialmente se requerirÃ­an una pasarela real con webhooks y un proveedor de facturaciÃ³n fiscal. Los recibos del cliente estÃ¡n disponibles en `GET /api/payments/mine`, y una factura pagada en `GET /api/orders/:id/invoice`.

## Despliegue reproducible

Configura secretos reales en un archivo `.env` no versionado y ejecuta:

```powershell
docker compose -f docker-compose.production.yml up --build -d
```

El servicio `migrate` aplica los scripts SQL antes de iniciar la API. En producciÃ³n utiliza HTTPS y un proxy inverso; no expongas PostgreSQL ni uses valores de ejemplo.

## Variables requeridas

`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `JWT_SECRET`, `OPERATOR_KEY`, `PORT` y `WS_ALLOWED_ORIGINS`.
