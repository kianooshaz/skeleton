-- name: Create :one
INSERT INTO usernames (
    id, username_value, user_id, is_primary, status, created_at, updated_at, deleted_at
) VALUES (
             $1, $2, $3, $4, $5, NOW(), NOW(), NULL
         ) RETURNING id, username_value, user_id, status, is_primary, created_at, updated_at, deleted_at;


-- name: Update :exec
UPDATE usernames
SET 
    is_primary = $1,
    status = $2,
    updated_at = NOW()
WHERE id = $3 AND deleted_at IS NULL;

-- name: List :many
SELECT id, username_value, user_id, is_primary, status, created_at, updated_at FROM usernames
WHERE user_id = $1 AND deleted_at IS NULL;

-- name: GetByUsername :one
SELECT id, username_value, user_id, is_primary, status, created_at, updated_at FROM usernames
WHERE username_value = $1 AND deleted_at IS NULL;

-- name: CountByUserID :one
SELECT COUNT(id) FROM usernames
WHERE user_id = $1 AND deleted_at IS NULL;

-- name: CountByUsername :one
SELECT COUNT(id) FROM usernames
WHERE username_value = $1 AND deleted_at IS NULL;

-- name: Delete :exec
UPDATE usernames
SET 
    deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;
