-- name: ListPhilosophers :many
SELECT id, name, date_born, date_died, birthplace, interests, portrait_uri, bio, created_at
FROM philosophers
ORDER BY name;

-- name: GetPhilosopher :one
SELECT id, name, date_born, date_died, birthplace, interests, portrait_uri, bio, created_at
FROM philosophers
WHERE id = $1;

-- name: CreatePhilosopher :one
INSERT INTO philosophers (
    name, date_born, date_died, birthplace, interests, portrait_uri, bio, created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, NOW()
)
RETURNING id;