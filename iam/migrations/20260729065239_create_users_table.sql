-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    uuid UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    login VARCHAR NOT NULL UNIQUE,
    email VARCHAR NOT NULL,
    password_hash TEXT NOT NULL,
    notification_methods JSONB NOT NULL DEFAULT '[]',
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
