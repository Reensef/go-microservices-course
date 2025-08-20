-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE orders (
    uuid UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    user_uuid UUID NOT NULL,
    part_ids VARCHAR(12)[] NOT NULL DEFAULT '{}',
    transaction_uuid UUID NULL,
    total_price NUMERIC NOT NULL,
    payment_method INT NOT NULL,
    order_status INT NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW(),
    deleted_at timestamp NULL
);

-- Триггер для обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_orders_updated_at
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Удаление триггера
DROP TRIGGER IF EXISTS update_orders_updated_at ON orders;
-- Удаление функции
DROP FUNCTION IF EXISTS update_updated_at_column();
-- Удаление таблицы
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
