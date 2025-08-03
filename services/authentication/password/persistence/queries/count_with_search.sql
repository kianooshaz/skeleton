SELECT COUNT(`id`)
FROM `passwords`
WHERE `account_id` = $1
    AND `deleted_at` IS NULL