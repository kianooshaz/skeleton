-- name: Create :one
INSERT INTO usernames (
    id, user_id, organization_id, status, created_at, updated_at, deleted_at
) VALUES (
             $1, $2, $3, $4, NOW(), NOW(), NULL
         ) RETURNING id, user_id, organization_id, status, created_at, updated_at, deleted_at;


-- name: Update :exec
UPDATE usernames
SET 
    user_id = $2,
    organization_id = $3,
    status = $4,
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: Get :one
SELECT id, user_id, organization_id, status, created_at, updated_at FROM usernames
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListByUser :many
SELECT id, user_id, organization_id, status, created_at, updated_at FROM usernames
WHERE user_id = $1 AND deleted_at IS NULL;

-- name: ListByUserAndOrganization :many
SELECT id, user_id, organization_id, status, created_at, updated_at FROM usernames
WHERE user_id = $1 AND organization_id = $2 AND deleted_at IS NULL;

-- name: CountByUser :one
SELECT COUNT(id) FROM usernames
WHERE user_id = $1 AND deleted_at IS NULL;


-- name: CountByUserAndOrganization :one
SELECT COUNT(id) FROM usernames
WHERE user_id = $1 AND deleted_at IS NULL;

-- name: Count :one
SELECT COUNT(id) FROM usernames
WHERE id = $1 AND deleted_at IS NULL;

-- name: Delete :exec
UPDATE usernames
SET 
    deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;
