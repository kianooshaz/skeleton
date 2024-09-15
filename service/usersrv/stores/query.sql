-- name: Create :one
INSERT INTO users (
    id, created_at
) VALUES (
             $1, NOW()
         )
RETURNING *;

-- name: Get :one
SELECT id, created_at FROM users
WHERE id = $1 LIMIT 1;

-- name: List :many
SELECT id, created_at FROM users
ORDER BY
    CASE WHEN sqlc.arg(order_by)::varchar = 'created_at_ASC' THEN created_at END,
    CASE WHEN sqlc.arg(order_by)::varchar = 'created_at_DESC' THEN created_at END DESC ;

-- name: Count :one
SELECT COUNT(id) FROM users;
