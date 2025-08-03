SELECT `id`,
    `account_id`,
    `password_hash`,
    `created_at`,
    `updated_at`
FROM `passwords`
WHERE `account_id` = $1
    AND `deleted_at` IS NULL
ORDER BY `created_at` DESC
LIMIT $2