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

-- name: IsPromotionCodAactive :one
SELECT EXISTS (
    SELECT 1 FROM promotions
    WHERE
        promotion_code = ?
        AND CURRENT_TIMESTAMP BETWEEN start_date AND end_date
        AND deleted_at IS NULL
) AS is_exists;