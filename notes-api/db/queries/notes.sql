-- name: CountNotes :one
SELECT COUNT(*)::BIGINT FROM notes;