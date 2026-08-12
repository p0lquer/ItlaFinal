-- Each base_price is the unit tariff. The final amount is calculated by the
-- API as tariff × max(ceil(weight in lb), pieces), with a minimum of one.
INSERT INTO service_types (name, description, base_price, price_per_weight, price_per_piece) VALUES
 ('Lavado en Seco', 'Limpieza especializada para prendas delicadas.', 120, 0, 0),
 ('Lavado y Secado', 'Lavado, secado y doblado para prendas de uso diario.', 150, 0, 0),
 ('Planchado', 'Planchado profesional por pieza con manejo cuidadoso.', 115, 0, 0),
 ('Lavado, Secado y Planchado', 'Servicio completo: lavado, secado y planchado.', 215, 0, 0)
ON CONFLICT (name) DO UPDATE SET description=EXCLUDED.description, base_price=EXCLUDED.base_price, price_per_weight=0, price_per_piece=0;

CREATE TABLE IF NOT EXISTS payments (
 id VARCHAR(36) PRIMARY KEY,
 order_id VARCHAR(36) NOT NULL UNIQUE REFERENCES orders(id) ON DELETE RESTRICT,
 amount NUMERIC(12,2) NOT NULL CHECK (amount >= 0),
 currency VARCHAR(3) NOT NULL DEFAULT 'DOP',
 method VARCHAR(20) NOT NULL CHECK (method IN ('efectivo','tarjeta','transferencia')),
 status VARCHAR(20) NOT NULL CHECK (status = 'paid'),
 receipt_number VARCHAR(40) NOT NULL UNIQUE,
 paid_at TIMESTAMP NOT NULL DEFAULT NOW()
);
