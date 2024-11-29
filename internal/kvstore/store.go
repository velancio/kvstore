package kvstore

import (
	"context"
)

// KeyValueStore is an interface for a key-value store
type KeyValueStore interface {
	Set(ctx context.Context, key string, value string) error
	Get(ctx context.Context, key string) (string, bool)
	Delete(ctx context.Context, key string) error
}
