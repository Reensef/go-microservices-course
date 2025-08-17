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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Удаление таблицы
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
