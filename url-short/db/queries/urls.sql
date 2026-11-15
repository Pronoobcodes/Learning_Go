-- name: CreateURL :one
INSERT INTO
    urls (code, original_url)
VALUES ($1, $2) RETURNING id,
    code,
    original_url,
    clicks,
    created_at;

-- name: GetURLByCode :one
SELECT
    id,
    code,
    original_url,
    clicks,
    created_at
FROM urls
WHERE
    code = $1
LIMIT 1;

-- name: DeleteURL :exec
DELETE FROM urls WHERE code = $1;

-- name: IncrementClicks :exec
UPDATE urls SET clicks = clicks + 1 WHERE code = $1;

-- name: ListURLs :many
SELECT
    id,
    code,
    original_url,
    clicks,
    created_at
FROM urls
ORDER BY created_at DESC;