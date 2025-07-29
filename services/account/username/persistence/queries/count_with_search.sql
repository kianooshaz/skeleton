SELECT COUNT(id)
FROM usernames
WHERE account_id = $1
    AND deleted_at IS NULL