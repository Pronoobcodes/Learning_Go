-- +goose Up

CREATE TABLE url_clicks (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(20) NOT NULL REFERENCES urls (code) ON DELETE CASCADE,
    clicked_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_url_clicks_code ON url_clicks (code);

-- +goose Down

DROP TABLE url_clicks;