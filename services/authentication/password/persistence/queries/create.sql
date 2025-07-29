INSERT INTO passwords (
        id,
        account_id,
        password_hash,
        created_at,
        updated_at,
        deleted_at
    )
VALUES ($1, $2, $3, NOW(), NOW(), NULL)