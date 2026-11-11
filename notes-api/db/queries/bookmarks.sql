-- name: CreateBookmark :one
INSERT INTO
    bookmarks (user_id, url, title, description)
VALUES (sqlc.arg(user_id), sqlc.arg(url), sqlc.arg(title), sqlc.arg(description))
RETURNING *;

-- name: ListBookmarks :many
SELECT * FROM bookmarks
WHERE user_id = sqlc.arg(user_id)
ORDER BY id DESC;

-- name: GetBookmark :one
SELECT * FROM bookmarks
WHERE id = sqlc.arg(id) AND user_id = sqlc.arg(user_id);

-- name: UpdateBookmark :one
UPDATE bookmarks
SET
    url = sqlc.arg(url),
    title = sqlc.arg(title),
    description = sqlc.arg(description),
    updated_at = now()
WHERE
    id = sqlc.arg(id)
    AND user_id = sqlc.arg(user_id)
RETURNING *;

-- name: DeleteBookmark :execrows
DELETE FROM bookmarks
WHERE id = sqlc.arg(id) AND user_id = sqlc.arg(user_id);

-- name: AddTagToBookmark :exec
WITH inserted_tag AS (
    INSERT INTO tags (name)
    SELECT sqlc.arg(name)
    WHERE EXISTS (
        SELECT 1
        FROM bookmarks
        WHERE bookmarks.id = sqlc.arg(bookmark_id)
          AND bookmarks.user_id = sqlc.arg(user_id)
    )
    ON CONFLICT (name) DO NOTHING
    RETURNING id
)
INSERT INTO bookmark_tags (bookmark_id, tag_id)
SELECT
        sqlc.arg(bookmark_id),
        COALESCE(
            (SELECT id FROM inserted_tag),
            (SELECT id FROM tags WHERE tags.name = sqlc.arg(name))
        )
WHERE EXISTS (
    SELECT 1
    FROM bookmarks
    WHERE bookmarks.id = sqlc.arg(bookmark_id)
      AND bookmarks.user_id = sqlc.arg(user_id)
)
ON CONFLICT (bookmark_id, tag_id) DO NOTHING;

-- name: SearchBookmarks :many
SELECT b.*
FROM bookmarks b
WHERE
    b.user_id = sqlc.arg(user_id)
    AND
    (
        sqlc.arg(search)::text = ''
        OR b.title ILIKE '%' || sqlc.arg(search)::text || '%'
        OR b.description ILIKE '%' || sqlc.arg(search)::text || '%'
        OR b.url ILIKE '%' || sqlc.arg(search)::text || '%'
    )
    AND
    (
        sqlc.arg(tag)::text = ''
        OR EXISTS (
            SELECT 1
            FROM bookmark_tags bt
            JOIN tags t ON t.id = bt.tag_id
            WHERE bt.bookmark_id = b.id
              AND t.name = sqlc.arg(tag)::text
        )
    )
ORDER BY b.id DESC;