-- name: ListPromotions :many
SELECT
	*
FROM
	promotions
WHERE
	deleted_at IS NULL;

