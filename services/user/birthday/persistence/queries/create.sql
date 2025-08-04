INSERT INTO birthdays (
        id,
        user_id,
        date_of_birth,
        age,
        created_at,
        updated_at
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        NOW(),
        NOW()
    )