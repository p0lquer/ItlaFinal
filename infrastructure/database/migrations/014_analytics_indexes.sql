-- Analytics only reads a recent, seven-day time window for its two trend
-- charts. These indexes keep those range joins bounded as orders/payments grow.
-- The partial indexes deliberately exclude historical rows that the trend
-- queries cannot use, reducing write and cache overhead.

CREATE INDEX IF NOT EXISTS orders_created_at_idx
    ON orders (created_at);

CREATE INDEX IF NOT EXISTS orders_delivered_updated_at_idx
    ON orders (updated_at)
    WHERE status = 'entregada';

CREATE INDEX IF NOT EXISTS payments_paid_at_amount_idx
    ON payments (paid_at) INCLUDE (amount)
    WHERE status = 'paid';
