package transport

import (
	"censys/proto/gen/proto"
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"reflect"
	"testing"
)

type mockKvStore struct {
	value string
}

func (m *mockKvStore) Get(ctx context.Context, key string) (string, bool) {
	if key == "" {
		return "", false
	}
	return m.value, true
}

func (m *mockKvStore) Set(ctx context.Context, key string, value string) error {
	if key == "" {
		return status.Errorf(codes.InvalidArgument, "key cannot be empty")
	}
	m.value = value
	return nil
}

func (m *mockKvStore) Delete(ctx context.Context, key string) error {
	if key == "" {
		return status.Errorf(codes.InvalidArgument, "key cannot be empty")
	}
	m.value = ""
	return nil
}

func TestKvStoreServer_Get(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		wantResp *proto.GetResponse
		wantErr  error
	}{
		{
			name:     "valid key",
			key:      "test-key",
			value:    "test-value",
			wantResp: &proto.GetResponse{Value: "test-value", Success: true},
		},
		{
			name:    "missing key",
			key:     "",
			wantErr: status.Errorf(codes.NotFound, "key not found"),
		},
		{
			name:    "key not found",
			key:     "non-existent-key",
			wantErr: status.Errorf(codes.NotFound, "key not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock store
			store := &mockKvStore{}
			if tt.key != "" {
				store.value = tt.value
			}

			// Create a KvStoreServer
			server := &KvStoreServer{Store: store}

			// Call the Get method
			resp, err := server.Get(context.Background(), &proto.GetRequest{Key: tt.key})

			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantResp != nil && !reflect.DeepEqual(resp, tt.wantResp) {
				t.Errorf("Get() response = %v, want %v", resp, tt.wantResp)
			}
		})
	}
}

func TestKvStoreServer_Set(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		wantResp *proto.SetResponse
		wantErr  error
	}{
		{
			name:     "valid key and value",
			key:      "test-key",
			value:    "test-value",
			wantResp: &proto.SetResponse{Success: true},
		},
		{
			name:    "missing key",
			key:     "",
			wantErr: status.Errorf(codes.InvalidArgument, "key cannot be empty"),
		},
		{
			name:     "missing value",
			key:      "test-key",
			value:    "",
			wantResp: &proto.SetResponse{Success: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock store
			store := &mockKvStore{}

			// Create a KvStoreServer
			server := &KvStoreServer{Store: store}

			// Call the Set method
			resp, err := server.Set(context.Background(), &proto.SetRequest{Key: tt.key, Value: tt.value})

			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantResp != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("Set() response = %v, want %v", resp, tt.wantResp)
			}
		})
	}
}

func TestKvStoreServer_Delete(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		wantResp *proto.DeleteResponse
		wantErr  error
	}{
		{
			name:     "valid key",
			key:      "test-key",
			wantResp: &proto.DeleteResponse{Success: true},
		},
		{
			name:    "missing key",
			key:     "",
			wantErr: status.Errorf(codes.NotFound, "key not found"),
		},
		{
			name:    "key not found",
			key:     "non-existent-key",
			wantErr: status.Errorf(codes.NotFound, "key not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock store
			store := &mockKvStore{}
			if tt.key != "" {
				store.value = tt.key
			}

			// Create a KvStoreServer
			server := &KvStoreServer{Store: store}

			// Call the Delete method
			resp, err := server.Delete(context.Background(), &proto.DeleteRequest{Key: tt.key})

			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantResp != nil && !reflect.DeepEqual(resp, tt.wantResp) {
				t.Errorf("Delete() response = %v, want %v", resp, tt.wantResp)
			}
		})
	}
}
