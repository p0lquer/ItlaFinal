-- Clientes
CREATE TABLE IF NOT EXISTS customers (
    id          VARCHAR(36) PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    phone       VARCHAR(20),
    email       VARCHAR(100),
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Órdenes
CREATE TABLE IF NOT EXISTS orders (
    id              VARCHAR(36) PRIMARY KEY,
    customer_id     VARCHAR(36) NOT NULL REFERENCES customers(id),
    service_type    VARCHAR(50) NOT NULL,
    pieces_count    INT NOT NULL DEFAULT 1,
    notes           TEXT,
    status          VARCHAR(20) NOT NULL DEFAULT 'recibida',
    estimated_time  FLOAT NOT NULL DEFAULT 60,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    ready_at        TIMESTAMP
);

-- Predicciones / historial de tiempos
CREATE TABLE IF NOT EXISTS predictions (
    id              VARCHAR(36) PRIMARY KEY,
    service_type    VARCHAR(50) NOT NULL,
    pieces_count    INT NOT NULL,
    estimated       FLOAT NOT NULL,
    actual          FLOAT,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Datos de prueba
INSERT INTO customers (id, name, phone, email) VALUES
    ('c1', 'Juan Pérez',    '809-555-0001', 'juan@mail.com'),
    ('c2', 'María Gómez',   '809-555-0002', 'maria@mail.com');