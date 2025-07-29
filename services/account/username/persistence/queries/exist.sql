SELECT EXISTS(
        SELECT 1
        FROM usernames
        WHERE username = $1
            AND deleted_at IS NULL
    )