DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'orders' AND column_name = 'estimaed_cost')
       AND NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'orders' AND column_name = 'estimated_cost') THEN
        ALTER TABLE orders RENAME COLUMN estimaed_cost TO estimated_cost;
    END IF;
END $$;
