UPDATE `usernames`
SET `status` = $2,
    `updated_at` = NOW()
WHERE `id` = $1
    AND `deleted_at` IS NULL