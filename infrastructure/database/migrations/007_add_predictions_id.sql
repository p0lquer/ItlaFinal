ALTER TABLE predictions ADD COLUMN IF NOT EXISTS order_id UUID REFERENCES orders(id);
ALTER TABLE predictions
ALTER COLUMN order_id SET NOT NULL;
ALTER TABLE predictions
ADD CONSTRAINT predictions_order_id_unique UNIQUE (order_id);