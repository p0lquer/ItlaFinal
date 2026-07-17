CREATE TABLE IF NOT EXISTS service_types (
    id          VARCHAR(36) PRIMARY KEY,
    name        VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Algunos tipos iniciales de ejemplo
INSERT INTO service_types (id, name, description) VALUES
    ('st1', 'Lavado y Secado',  'Lavado completo con secado incluido'),
    ('st2', 'Planchado',        'Servicio de planchado de prendas'),
    ('st3', 'Lavado en Seco',   'Limpieza en seco para prendas delicadas')
ON CONFLICT (name) DO NOTHING;