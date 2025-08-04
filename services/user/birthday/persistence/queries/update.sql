UPDATE birthdays
SET date_of_birth = $2,
    age = $3,
    updated_at = NOW()
WHERE id = $1