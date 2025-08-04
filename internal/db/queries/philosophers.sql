-- name: ListPhilosophers :many
SELECT id, name, date_born, date_died, birthplace, interests, portrait_uri, bio, created_at
FROM philosophers
ORDER BY name;