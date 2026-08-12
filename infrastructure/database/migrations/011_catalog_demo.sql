-- Catalogo inicial realista para dejar el entorno de demostracion listo al
-- levantar PostgreSQL por primera vez. No modifica catalogos personalizados.
INSERT INTO service_types (name, description, base_price, price_per_weight, price_per_piece) VALUES
    ('Lavado y Secado', 'Lavado completo, secado y doblado para prendas de uso diario.', 50, 5, 0),
    ('Planchado', 'Planchado profesional por pieza con manejo cuidadoso.', 30, 0, 2),
    ('Lavado en Seco', 'Limpieza especializada para prendas delicadas.', 80, 8, 0)
ON CONFLICT (name) DO UPDATE SET
    description = EXCLUDED.description,
    base_price = CASE WHEN service_types.base_price = 0 THEN EXCLUDED.base_price ELSE service_types.base_price END,
    price_per_weight = CASE WHEN service_types.price_per_weight = 0 THEN EXCLUDED.price_per_weight ELSE service_types.price_per_weight END,
    price_per_piece = CASE WHEN service_types.price_per_piece = 0 THEN EXCLUDED.price_per_piece ELSE service_types.price_per_piece END;
