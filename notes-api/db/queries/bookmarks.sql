-- name: CreateBookmark :one
INSERT INTO
    bookmarks (url, title, description)
VALUES ($1, $2, $3) RETURNING id,
    url,
    title,
    description,
    created_at,
    updated_at;

-- name: ListBookmarks :many
SELECT
    id,
    url,
    title,
    description,
    created_at,
    updated_at
FROM bookmarks
ORDER BY id DESC;

-- name: GetBookmark :one
SELECT
    id,
    url,
    title,
    description,
    created_at,
    updated_at
FROM bookmarks
WHERE
    id = $1;

-- name: UpdateBookmark :one
UPDATE bookmarks
SET
    url = $2,
    title = $3,
    description = $4,
    updated_at = now()
WHERE
    id = $1 RETURNING id,
    url,
    title,
    description,
    created_at,
    updated_at;

-- name: DeleteBookmark :execrows
DELETE FROM bookmarks WHERE id = $1;