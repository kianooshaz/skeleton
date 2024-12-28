-- name: Create :one
INSERT INTO passwords (
    id, user_id, password_hash, created_at, deleted_at
) VALUES (
             $1, $2, $3, NOW(), NULL
         )
RETURNING id, user_id, password_hash, created_at, deleted_at;

-- name: Delete :exec
UPDATE
    passwords
SET 
    deleted_at = NOW()
WHERE
    id = $1 AND deleted_at IS NULL;

-- name: GetByUserID :one
SELECT 
    id, user_id, password_hash, created_at, deleted_at
FROM
    passwords
WHERE
    user_id = $1 AND deleted_at IS NULL;

-- name: History :many
SELECT 
    id, user_id, password_hash, created_at, deleted_at
FROM
    passwords
WHERE
    user_id = $1;