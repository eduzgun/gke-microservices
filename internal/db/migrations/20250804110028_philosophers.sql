-- +goose Up
-- +goose StatementBegin
CREATE TABLE philosophers (
    id           SERIAL PRIMARY KEY,
    name         TEXT NOT NULL,
    date_born    TEXT,
    date_died    TEXT,
    birthplace   TEXT,
    interests    TEXT[],
    portrait_uri TEXT,
    bio          TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS philosophers;
-- +goose StatementEnd
