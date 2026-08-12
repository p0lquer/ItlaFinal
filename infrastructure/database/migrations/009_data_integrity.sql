-- Las predicciones pertenecen a una orden y no deben quedar huérfanas al
-- eliminarla. El nombre de la FK puede variar entre instalaciones, por eso se
-- localiza dinámicamente antes de recrearla.
DO $$
DECLARE constraint_name text;
BEGIN
    SELECT tc.constraint_name INTO constraint_name
    FROM information_schema.table_constraints tc
    JOIN information_schema.key_column_usage kcu
      ON tc.constraint_name = kcu.constraint_name
     AND tc.table_schema = kcu.table_schema
    WHERE tc.table_name = 'predictions'
      AND tc.constraint_type = 'FOREIGN KEY'
      AND kcu.column_name = 'order_id'
    LIMIT 1;

    IF constraint_name IS NOT NULL THEN
        EXECUTE format('ALTER TABLE predictions DROP CONSTRAINT %I', constraint_name);
    END IF;
END $$;

ALTER TABLE predictions
    ADD CONSTRAINT predictions_order_id_fkey
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE;
