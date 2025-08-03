SELECT `id`,
    `username`,
    `account_id`,
    `status`,
    `created_at`,
    `updated_at`
FROM `usernames`
WHERE `id` = $1
    AND `deleted_at` IS NULL