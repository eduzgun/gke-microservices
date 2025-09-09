-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id  SERIAL PRIMARY KEY,
    username         TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd