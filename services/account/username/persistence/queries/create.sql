INSERT INTO usernames (
        id,
        username,
        account_id,
        status,
        created_at,
        updated_at,
        deleted_at
    )
VALUES ($1, $2, $3, $4, NOW(), NOW(), NULL)