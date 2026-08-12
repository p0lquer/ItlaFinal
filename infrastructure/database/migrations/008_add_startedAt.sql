-- Una orden recibida todavía no ha iniciado su procesamiento.
ALTER TABLE orders ADD COLUMN IF NOT EXISTS started_at TIMESTAMP;

CREATE UNIQUE INDEX IF NOT EXISTS customers_user_id_unique
    ON customers (user_id) WHERE user_id IS NOT NULL;
