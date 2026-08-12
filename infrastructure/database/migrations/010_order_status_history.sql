-- Auditoria inmutable de las transiciones de una orden. El estado vigente
-- sigue viviendo en `orders`; esta tabla permite explicar como llego alli.
CREATE TABLE IF NOT EXISTS order_status_history (
    id               BIGSERIAL PRIMARY KEY,
    order_id         VARCHAR(36) NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    from_status      VARCHAR(20),
    to_status        VARCHAR(20) NOT NULL,
    changed_by       VARCHAR(36),
    changed_by_role  VARCHAR(20),
    description      TEXT,
    changed_at       TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS order_status_history_order_changed_at_idx
    ON order_status_history (order_id, changed_at DESC);
