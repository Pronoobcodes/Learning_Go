package store

import (
	"context"
	"errors"

	db "url-shortener/internal/store/db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PostgresRepository struct {
	queries *db.Queries
}

func NewPostgresRepository(queries *db.Queries) *PostgresRepository {
	return &PostgresRepository{
		queries: queries,
	}
}

func (p *PostgresRepository) Create(ctx context.Context, url URL) (URL, error) {
	result, err := p.queries.CreateURL(ctx, db.CreateURLParams{
		Code:        url.Code,
		OriginalUrl: url.OriginalURL,
	})
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return URL{}, ErrCodeAlreadyExists
		}

		return URL{}, err
	}

	return URL{
		ID:          result.ID,
		Code:        result.Code,
		OriginalURL: result.OriginalUrl,
		Clicks:      result.Clicks,
		CreatedAt:   result.CreatedAt.Time,
	}, nil
}

func (p *PostgresRepository) GetByCode(ctx context.Context, code string) (URL, error) {
	result, err := p.queries.GetURLByCode(ctx, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return URL{}, ErrURLNotFound
		}

		return URL{}, err
	}

	return URL{
		ID:          result.ID,
		Code:        result.Code,
		OriginalURL: result.OriginalUrl,
		Clicks:      result.Clicks,
		CreatedAt:   result.CreatedAt.Time,
	}, nil
}

func (p *PostgresRepository) Delete(ctx context.Context, code string) error {
	rows, err := p.queries.DeleteURL(ctx, code)
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrURLNotFound
	}
	return nil
}

func (p *PostgresRepository) IncrementClicks(ctx context.Context, code string) error {
	rows, err := p.queries.IncrementClicks(ctx, code)
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrURLNotFound
	}
	return nil
}

func (p *PostgresRepository) GetURLAndIncrementClicks(ctx context.Context, code string) (URL, error) {
	result, err := p.queries.GetURLAndIncrementClicks(ctx, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return URL{}, ErrURLNotFound
		}

		return URL{}, err
	}

	return URL{
		ID:          result.ID,
		Code:        result.Code,
		OriginalURL: result.OriginalUrl,
		Clicks:      result.Clicks,
		CreatedAt:   result.CreatedAt.Time,
	}, nil
}
