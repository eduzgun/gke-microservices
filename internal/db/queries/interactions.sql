-- name: CreateLikeInteraction :exec
INSERT INTO interactions (
    user_id,
    philosopher_id,
    username,
    type,
    content,
    created_at
) VALUES ($1, $2, $3, 'like', NULL, NOW());

-- name: CreateCommentInteraction :one
INSERT INTO interactions (
    user_id,
    philosopher_id,
    username,
    type,
    content,
    created_at
) VALUES ($1, $2, $3, 'comment', $4, NOW()) RETURNING id, user_id, philosopher_id, username, content, created_at;

-- name: GetInteractionsByPhilosopher :many
SELECT 
    i.id,
    i.user_id,
    i.philosopher_id,
    i.username,
    i.type,
    COALESCE(i.content, '') as content,
    i.created_at,
    u.username
FROM interactions i
JOIN users u ON i.user_id = u.id
WHERE i.philosopher_id = $1
ORDER BY i.created_at DESC;

-- name: GetUserInteraction :one
SELECT id, user_id, philosopher_id, username, type, COALESCE(content, '') as content, created_at
FROM interactions
WHERE user_id = $1 
  AND philosopher_id = $2 
  AND type = $3;

-- name: DeleteInteraction :exec
DELETE FROM interactions
WHERE id = $1;

-- name: CountLikesByPhilosopher :one
SELECT COUNT(*) FROM interactions
WHERE philosopher_id = $1 AND type = 'like';

-- name: HasUserLiked :one
SELECT EXISTS (
    SELECT 1 FROM interactions
    WHERE user_id = $1 
      AND philosopher_id = $2 
      AND type = 'like'
);