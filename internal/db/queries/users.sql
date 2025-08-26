-- name: GetUser :one
SELECT id, username, email, password, created_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, username, email, password, created_at
FROM users
WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (
    username, email, password, created_at
) VALUES (
    $1, $2, $3, NOW()
)
RETURNING id;