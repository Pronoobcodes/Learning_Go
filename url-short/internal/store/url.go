package store

import "time"

type URL struct {
	ID          int64     `json:"id"`
	Code        string    `json:"code"`
	OriginalURL string    `json:"original_url"`
	Clicks      int64     `json:"clicks"`
	CreatedAt   time.Time `json:"created_at"`
}
