package service

import (
	"errors"
	"strings"

	"url-shortener/internal/store"
)

var ErrInvalidURL = errors.New("invalid URL")

type URLRepository interface {
	Create(url store.URL) error
	GetByCode(code string) (store.URL, error)
	Delete(code string) error
	IncrementClicks(code string) error
}

type URLService struct {
	repo URLRepository
}

func NewURLService(repo URLRepository) *URLService {
	return &URLService{
		repo: repo,
	}
}

func (s *URLService) CreateURL(originalURL string) (store.URL, error) {
	originalURL = strings.TrimSpace(originalURL)

	if originalURL == "" {
		return store.URL{}, ErrInvalidURL
	}

	code, err := GenerateCode(6)
	if err != nil {
		return store.URL{}, err
	}

	url := store.URL{
		Code:        code,
		OriginalURL: originalURL,
	}

	err = s.repo.Create(url)
	if err != nil {
		return store.URL{}, err
	}

	return url, nil
}

func (s *URLService) GetURL(code string) (store.URL, error) {
	code = strings.TrimSpace(code)

	if code == "" {
		return store.URL{}, ErrInvalidURL
	}

	return s.repo.GetByCode(code)
}

func (s *URLService) DeleteURL(code string) error {
	code = strings.TrimSpace(code)

	if code == "" {
		return ErrInvalidURL
	}

	return s.repo.Delete(code)
}

func (s *URLService) IncrementClicks(code string) error {
	return s.repo.IncrementClicks(code)
}
