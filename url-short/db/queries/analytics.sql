-- name: RecordClick :exec
INSERT INTO url_clicks (code) VALUES ($1);

-- name: GetURLClickCount :one
SELECT COUNT(*)::BIGINT FROM url_clicks WHERE code = $1;
