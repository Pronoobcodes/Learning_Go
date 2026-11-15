package store

import (
	"errors"
	"sync"
)

var ErrURLNotFound = errors.New("URL not found")
var ErrCodeAlreadyExists = errors.New("short code already exists")

type MemoryRepository struct {
	mu   sync.RWMutex
	urls map[string]URL
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		urls: make(map[string]URL),
	}
}

func (m *MemoryRepository) Create(url URL) (URL, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.urls[url.Code]; exists {
		return URL{}, ErrCodeAlreadyExists
	}

	m.urls[url.Code] = url
	return url, nil
}

func (m *MemoryRepository) GetByCode(code string) (URL, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	url, exists := m.urls[code]
	if !exists {
		return URL{}, ErrURLNotFound
	}

	return url, nil
}

func (m *MemoryRepository) Delete(code string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.urls[code]; !exists {
		return ErrURLNotFound
	}

	delete(m.urls, code)
	return nil
}

func (m *MemoryRepository) IncrementClicks(code string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	url, exists := m.urls[code]
	if !exists {
		return ErrURLNotFound
	}

	url.Clicks++
	m.urls[code] = url

	return nil
}
