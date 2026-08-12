CREATE TABLE IF NOT EXISTS service_types (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Algunos tipos iniciales de ejemplo
INSERT INTO service_types (name, description) VALUES
    ('Lavado y Secado',  'Lavado completo con secado incluido'),
    ('Planchado',        'Servicio de planchado de prendas'),
    ('Lavado en Seco',   'Limpieza en seco para prendas delicadas')
ON CONFLICT (name) DO NOTHING;
