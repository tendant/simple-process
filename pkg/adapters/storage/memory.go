package storage

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"sync"

	"github.com/tendant/simple-process/pkg/adapters"
)

// InMemoryStorage is an in-memory implementation of the Storage interface.
// It is useful for testing and local development.
type InMemoryStorage struct {
	mu    sync.RWMutex
	blobs map[string][]byte
}

// NewInMemoryStorage creates a new InMemoryStorage.
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		blobs: make(map[string][]byte),
	}
}

// Get returns a reader for the given blob location.
func (s *InMemoryStorage) Get(ctx context.Context, location string) (io.ReadCloser, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, ok := s.blobs[location]
	if !ok {
		return nil, adapters.ErrNotFound
	}

	return io.NopCloser(bytes.NewReader(data)), nil
}

// Put uploads a blob from a reader to the given location.
func (s *InMemoryStorage) Put(ctx context.Context, location string, reader io.Reader) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.blobs[location] = data
	return nil
}

// PresignGet generates a presigned URL for getting a blob.
// For in-memory storage, it returns a data URI.
func (s *InMemoryStorage) PresignGet(ctx context.Context, location string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, ok := s.blobs[location]
	if !ok {
		return "", adapters.ErrNotFound
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	return "data:application/octet-stream;base64," + encoded, nil
}
