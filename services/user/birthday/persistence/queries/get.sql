SELECT id,
    user_id,
    date_of_birth,
    age,
    created_at,
    updated_at
FROM birthdays
WHERE id = $1