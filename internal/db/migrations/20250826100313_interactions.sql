-- db/migrations/002_add_interactions.up.sql
-- +goose Up
CREATE TABLE interactions (
    id               SERIAL PRIMARY KEY,
    user_id          INT NOT NULL REFERENCES users(id),
    philosopher_id   INT NOT NULL REFERENCES philosophers(id),
    username         TEXT NOT NULL,
    type             TEXT NOT NULL CHECK (type IN ('like', 'comment')),
    content          TEXT,
    created_at       TIMESTAMPTZ DEFAULT NOW(),

    --Ensure content is NULL for likes
    CONSTRAINT content_required_for_comment CHECK (
        (type = 'comment' AND content IS NOT NULL) OR
        (type = 'like' AND content IS NULL)
    )
);

--Index for fast lookup
CREATE INDEX idx_interactions_philosopher ON interactions(philosopher_id);
CREATE INDEX idx_interactions_user ON interactions(user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_interactions_philosopher;
DROP INDEX IF EXISTS idx_interactions_user;
DROP TABLE IF EXISTS interactions;