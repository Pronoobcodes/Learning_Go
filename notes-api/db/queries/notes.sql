-- -- name: CountNotes :one
-- SELECT COUNT(*)::BIGINT FROM notes;


-- name: CreateNote :one
INSERT INTO
    notes (user_id, title, content)
VALUES (sqlc.arg(user_id), sqlc.arg(title), sqlc.arg(content))
RETURNING *;

-- name: ListNotes :many
SELECT * FROM notes
WHERE user_id = sqlc.arg(user_id)
ORDER BY created_at DESC;

-- name: GetNote :one
SELECT * FROM notes
WHERE id = sqlc.arg(id) AND user_id = sqlc.arg(user_id);

-- name: UpdateNote :one
UPDATE notes
SET
    title = sqlc.arg(title),
    content = sqlc.arg(content),
    updated_at = NOW()
WHERE
    id = sqlc.arg(id)
    AND user_id = sqlc.arg(user_id)
RETURNING *;

-- name: DeleteNote :execrows
DELETE FROM notes
WHERE id = sqlc.arg(id) AND user_id = sqlc.arg(user_id);

-- name: AddTagToNote :exec
WITH inserted_tag AS (
    INSERT INTO tags (name)
    SELECT sqlc.arg(name)
    WHERE EXISTS (
        SELECT 1
        FROM notes
        WHERE notes.id = sqlc.arg(note_id)
          AND notes.user_id = sqlc.arg(user_id)
    )
    ON CONFLICT (name) DO NOTHING
    RETURNING id
)
INSERT INTO note_tags (note_id, tag_id)
SELECT
    sqlc.arg(note_id),
    COALESCE(
        (SELECT id FROM inserted_tag),
        (SELECT id FROM tags WHERE tags.name = sqlc.arg(name))
    )
WHERE EXISTS (
    SELECT 1
    FROM notes
    WHERE notes.id = sqlc.arg(note_id)
      AND notes.user_id = sqlc.arg(user_id)
)
ON CONFLICT (note_id, tag_id) DO NOTHING;

-- name: SearchNotes :many
SELECT *
FROM notes
WHERE
    user_id = sqlc.arg(user_id)
    AND (
        sqlc.arg(search)::text = ''
        OR title ILIKE '%' || sqlc.arg(search)::text || '%'
        OR content ILIKE '%' || sqlc.arg(search)::text || '%'
    )
    AND (
        sqlc.arg(tag)::text = ''
        OR EXISTS (
            SELECT 1
            FROM note_tags nt
            JOIN tags t ON t.id = nt.tag_id
            WHERE nt.note_id = notes.id
              AND t.name = sqlc.arg(tag)::text
        )
    )
ORDER BY created_at DESC;