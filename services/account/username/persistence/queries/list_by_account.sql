SELECT `id`,
    `username`,
    `account_id`,
    `status`,
    `created_at`,
    `updated_at`
FROM `usernames`
WHERE `account_id` = $1
    AND `deleted_at` IS NULL