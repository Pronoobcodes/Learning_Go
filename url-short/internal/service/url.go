package service

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"url-shortener/internal/store"
)

var ErrInvalidURL = errors.New("invalid URL")

type URLRepository interface {
	Create(ctx context.Context, url store.URL) (store.URL, error)
	GetByCode(ctx context.Context, code string) (store.URL, error)
	Delete(ctx context.Context, code string) error
	IncrementClicks(ctx context.Context, code string) error
	GetURLAndIncrementClicks(ctx context.Context, code string) (store.URL, error)
}

type URLService struct {
	repo URLRepository
}

func NewURLService(repo URLRepository) *URLService {
	return &URLService{
		repo: repo,
	}
}

func (s *URLService) CreateURL(ctx context.Context, originalURL string) (store.URL, error) {
	originalURL = strings.TrimSpace(originalURL)

	parsedURL, err := url.ParseRequestURI(originalURL)
	if err != nil || parsedURL.Host == "" {
		return store.URL{}, ErrInvalidURL
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return store.URL{}, ErrInvalidURL
	}

	for i := 0; i < 5; i++ {
		code, err := GenerateCode(6)
		if err != nil {
			return store.URL{}, err
		}

		url := store.URL{
			Code:        code,
			OriginalURL: originalURL,
		}

		createdURL, err := s.repo.Create(ctx, url)
		if err == nil {
			return createdURL, nil
		}

		if !errors.Is(err, store.ErrCodeAlreadyExists) {
			return store.URL{}, err
		}
	}
	return store.URL{}, errors.New("could not generate a unique short code")
}

func (s *URLService) GetURL(ctx context.Context, code string) (store.URL, error) {
	code = strings.TrimSpace(code)

	if code == "" {
		return store.URL{}, ErrInvalidURL
	}

	return s.repo.GetByCode(ctx, code)
}

func (s *URLService) DeleteURL(ctx context.Context, code string) error {
	code = strings.TrimSpace(code)

	if code == "" {
		return ErrInvalidURL
	}

	return s.repo.Delete(ctx, code)
}

func (s *URLService) IncrementClicks(ctx context.Context, code string) error {
	code = strings.TrimSpace(code)

	if code == "" {
		return ErrInvalidURL
	}

	return s.repo.IncrementClicks(ctx, code)
}

func (s *URLService) GetURLAndIncrementClicks(ctx context.Context, code string) (store.URL, error) {
	code = strings.TrimSpace(code)

	if code == "" {
		return store.URL{}, ErrInvalidURL
	}

	return s.repo.GetURLAndIncrementClicks(ctx, code)
}
