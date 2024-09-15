-- name: GetPromotion :one
SELECT * FROM promotions
WHERE id = ? LIMIT 1;

-- name: ListPromotions :many
SELECT * FROM promotions
WHERE deleted_at IS NULL;

-- name: CreatePromotion :execresult
INSERT INTO promotions (
    promotion_code,
    start_date,
    end_date
) VALUES (
    ?, ?, ?
);

-- name: IsPromotionCodeTaken :one
SELECT COUNT(1) > 0 AS is_exists
FROM promotions
WHERE
    promotion_code = ?
    AND deleted_at IS NULL
    AND sqlc.arg(start_date) < end_date
    AND sqlc.arg(end_date) > start_date;