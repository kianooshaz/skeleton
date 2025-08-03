UPDATE `passwords`
SET `deleted_at` = NOW()
WHERE `id` = $1
    AND `deleted_at` IS NULL