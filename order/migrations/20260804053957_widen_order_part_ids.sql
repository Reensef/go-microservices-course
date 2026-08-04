-- +goose Up
-- +goose StatementBegin
-- part_ids хранит MongoDB ObjectID деталей в hex-виде (24 символа), а не 12 байт
-- сырого ObjectID — VARCHAR(12) был занижен и обрезал бы реальные ID.
ALTER TABLE orders ALTER COLUMN part_ids TYPE VARCHAR(24)[];
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders ALTER COLUMN part_ids TYPE VARCHAR(12)[];
-- +goose StatementEnd
