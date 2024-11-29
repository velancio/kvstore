package kvstore

import (
	"context"
	"fmt"
	"sync"
)

// InMemoryStore is an in-memory store
type InMemoryStore struct {
	Data sync.Map
}

// Set sets a value for a key
func (s *InMemoryStore) Set(ctx context.Context, key string, value string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if key == "" {
			return fmt.Errorf("key cannot be empty")
		}
		s.Data.Store(key, value)
		return nil
	}
}

// Get gets a value for a key
func (s *InMemoryStore) Get(ctx context.Context, key string) (string, bool) {
	select {
	case <-ctx.Done():
		return "", false
	default:
		value, ok := s.Data.Load(key)
		if !ok {
			return "", false
		}
		return value.(string), true
	}
}

// Delete deletes a value for a key
func (s *InMemoryStore) Delete(ctx context.Context, key string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if key == "" {
			return fmt.Errorf("key cannot be empty")
		}
		s.Data.Delete(key)
		return nil
	}
}
