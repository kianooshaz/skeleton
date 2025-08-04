SELECT EXISTS(
        SELECT 1
        FROM birthdays
        WHERE user_id = $1
    )