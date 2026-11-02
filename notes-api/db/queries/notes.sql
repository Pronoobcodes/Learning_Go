-- -- name: CountNotes :one
-- SELECT COUNT(*)::BIGINT FROM notes;


-- name: CreateNote :one
INSERT INTO
    notes (title, content)
VALUES ($1, $2) RETURNING id,
    title,
    content,
    created_at,
    updated_at;

-- name: ListNotes :many
SELECT
    id,
    title,
    content,
    created_at,
    updated_at
FROM notes
ORDER BY id DESC;

-- name: GetNote :one
SELECT
    id,
    title,
    content,
    created_at,
    updated_at
FROM notes
WHERE
    id = $1;

-- name: UpdateNote :one
UPDATE notes
SET
    title = $2,
    content = $3,
    updated_at = now()
WHERE
    id = $1 RETURNING id,
    title,
    content,
    created_at,
    updated_at;

-- name: DeleteNote :execrows
DELETE FROM notes WHERE id = $1;

-- name: AddTagToNote :exec
WITH inserted_tag AS (
    INSERT INTO tags (name)
    VALUES ($2)
    ON CONFLICT (name) DO NOTHING
    RETURNING id
)
INSERT INTO note_tags (note_id, tag_id)
VALUES (
    $1,
    COALESCE(
        (SELECT id FROM inserted_tag),
        (SELECT id FROM tags WHERE tags.name = $2)
    )
)
ON CONFLICT (note_id, tag_id) DO NOTHING;

-- name: SearchNotes :many
SELECT
    n.id,
    n.title,
    n.content,
    n.created_at,
    n.updated_at
FROM notes n
WHERE
    (
        sqlc.arg(search)::text = ''
        OR n.title ILIKE '%' || sqlc.arg(search)::text || '%'
        OR n.content ILIKE '%' || sqlc.arg(search)::text || '%'
    )
    AND
    (
        sqlc.arg(tag)::text = ''
        OR EXISTS (
            SELECT 1
            FROM note_tags nt
            JOIN tags t ON t.id = nt.tag_id
            WHERE nt.note_id = n.id
              AND t.name = sqlc.arg(tag)::text
        )
    )
ORDER BY n.id DESC;