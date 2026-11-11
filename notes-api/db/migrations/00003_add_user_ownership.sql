-- +goose Up

ALTER TABLE notes
ADD COLUMN user_id BIGINT REFERENCES users (id) ON DELETE CASCADE;

ALTER TABLE bookmarks
ADD COLUMN user_id BIGINT REFERENCES users (id) ON DELETE CASCADE;

CREATE INDEX idx_notes_user_id ON notes (user_id);

CREATE INDEX idx_bookmarks_user_id ON bookmarks (user_id);

-- +goose Down

DROP INDEX idx_bookmarks_user_id;

DROP INDEX idx_notes_user_id;

ALTER TABLE bookmarks DROP COLUMN user_id;

ALTER TABLE notes DROP COLUMN user_id;