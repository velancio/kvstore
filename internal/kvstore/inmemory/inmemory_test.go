package kvstore

import (
	"context"

	"sync"
	"testing"
)

func TestInMemoryStore_Get(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		wantResp string
		wantErr  bool
	}{
		{
			name:     "valid key",
			key:      "test-key",
			wantResp: "test-value",
		},
		{
			name:    "missing key",
			key:     "",
			wantErr: true,
		},
		{
			name:    "key not found",
			key:     "non-existent-key",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create an in-memory store
			store := &InMemoryStore{
				Data: sync.Map{},
			}

			// Set a value for the key
			if tt.key != "" {
				store.Set(context.Background(), tt.key, tt.wantResp)
			}

			// Call the Get method
			resp, ok := store.Get(context.Background(), tt.key)

			if !ok && !tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", ok, tt.wantErr)
				return
			}

			if resp != tt.wantResp {
				t.Errorf("Get() response = %v, want %v", resp, tt.wantResp)
			}
		})
	}
}

func TestInMemoryStore_Set(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		value   string
		wantErr bool
	}{
		{
			name:  "valid key and value",
			key:   "test-key",
			value: "test-value",
		},
		{
			name:    "missing key",
			key:     "",
			wantErr: true,
		},
		{
			name:  "missing value",
			key:   "test-key",
			value: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create an in-memory store
			store := &InMemoryStore{
				Data: sync.Map{},
			}
			// Call the Set method
			err := store.Set(context.Background(), tt.key, tt.value)

			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestInMemoryStore_Delete(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{
			name: "valid key",
			key:  "test-key",
		},
		{
			name:    "missing key",
			key:     "",
			wantErr: true,
		},
		{
			name: "key not found",
			key:  "non-existent-key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create an in-memory store
			store := &InMemoryStore{
				Data: sync.Map{},
			}

			// Set a value for the key
			if tt.key != "" {
				store.Set(context.Background(), tt.key, "test-value")
			}

			// Call the Delete method
			err := store.Delete(context.Background(), tt.key)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
