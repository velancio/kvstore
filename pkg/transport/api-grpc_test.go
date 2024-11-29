package transport

import (
	"bytes"
	pb "censys/proto/gen/proto"
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Create a mock store
type mockStore struct {
	err     error
	value   string
	success bool
}

func (m *mockStore) Set(ctx context.Context, in *pb.SetRequest, opts ...grpc.CallOption) (*pb.SetResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &pb.SetResponse{Success: m.success}, nil
}

func (m *mockStore) Delete(ctx context.Context, in *pb.DeleteRequest, opts ...grpc.CallOption) (*pb.DeleteResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &pb.DeleteResponse{Success: m.success}, nil
}

func (m *mockStore) Get(ctx context.Context, in *pb.GetRequest, opts ...grpc.CallOption) (*pb.GetResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &pb.GetResponse{Value: m.value}, nil
}

// Test the HandleGet function
func TestHandleGet(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		wantCode       int
		wantResp       string
		grpcStoreError error
	}{
		{
			name:     "valid key",
			key:      "test-key",
			wantCode: http.StatusOK,
			wantResp: "{\"value\":\"test-value\"}\n",
		},
		{
			name:     "missing key text",
			key:      "",
			wantCode: http.StatusBadRequest,
		},
		{
			name:           "key not found in store",
			key:            "test-key",
			grpcStoreError: status.Errorf(codes.NotFound, "key not found"),
			wantCode:       http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/store?key="+tt.key, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Mock the gRPC store
			store := &mockStore{}
			if tt.grpcStoreError != nil {
				store.err = tt.grpcStoreError
			} else {
				store.value = "test-value"
			}

			s := &GrpcServer{Store: store}
			s.HandleGet(w, req)

			if w.Code != tt.wantCode {
				t.Errorf("HandleGet() wrote code %d, want %d", w.Code, tt.wantCode)
			}

			if tt.wantResp != "" {
				resp := w.Body.String()
				if resp != tt.wantResp {
					t.Errorf("HandleGet() response = %s, want %s", resp, tt.wantResp)
				}
			}
		})
	}
}

// Test the HandleSet function
func TestHandleSet(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		value          string
		wantCode       int
		wantResp       string
		grpcStoreError error
	}{
		{
			name:     "valid key",
			key:      "test-key",
			value:    "test-value",
			wantCode: http.StatusOK,
			wantResp: "{\"success\":true}\n",
		},
		{
			name:     "Missing value",
			key:      "test-key",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "missing value text",
			value:    "test-value",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "missing key and value text",
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			// Create JSON payload
			payload := make(map[string]string)
			if tt.value != "" {
				payload["value"] = tt.value
			}
			if tt.key != "" {
				payload["key"] = tt.key
			}
			// Convert payload to JSON
			jsonPayload, err := json.Marshal(payload)
			if err != nil {
				t.Fatalf("Failed to create JSON payload: %v", err)
			}

			// Create request with JSON body
			req, err := http.NewRequest("POST", "/store", bytes.NewBuffer(jsonPayload))
			if err != nil {
				t.Fatal(err)
			}

			// Mock the gRPC store
			store := &mockStore{
				success: true,
			}
			if tt.grpcStoreError != nil {
				store.err = tt.grpcStoreError
			} else {
				store.value = "test-value"
			}

			s := &GrpcServer{Store: store}
			s.HandleSet(w, req)

			if w.Code != tt.wantCode {
				t.Errorf("HandleSet() wrote code %d, want %d", w.Code, tt.wantCode)
			}

			if tt.wantResp != "" {
				resp := w.Body.String()
				if resp != tt.wantResp {
					t.Errorf("HandleSet() response = %s, want %s", resp, tt.wantResp)
				}
			}
		})
	}
}

// Test the HandleDelete function
func TestHandleDelete(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		wantCode       int
		wantResp       string
		grpcStoreError error
	}{
		{
			name:     "valid key",
			key:      "test-key",
			wantCode: http.StatusOK,
			wantResp: "{\"success\":true}\n",
		},
		{
			name:     "missing key",
			key:      "",
			wantCode: http.StatusBadRequest,
		},
		{
			name:           "grpc store error",
			key:            "test-key",
			grpcStoreError: status.Errorf(codes.Internal, "internal error"),
			wantCode:       http.StatusInternalServerError,
		},
		{
			name:           "key not found",
			key:            "test-key",
			grpcStoreError: status.Errorf(codes.NotFound, "not found"),
			wantCode:       http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest("DELETE", "/store/{key}", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.SetPathValue("key", tt.key)

			// Mock the gRPC store
			store := &mockStore{}
			if tt.grpcStoreError != nil {
				store.err = tt.grpcStoreError
			} else {
				store.success = true
			}

			s := &GrpcServer{Store: store}
			s.HandleDelete(w, req)

			if w.Code != tt.wantCode {
				t.Errorf("HandleDelete() wrote code %d, want %d", w.Code, tt.wantCode)
			}

			if tt.wantResp != "" {
				resp := w.Body.String()
				if resp != tt.wantResp {
					t.Errorf("HandleDelete() response = %s, want %s", resp, tt.wantResp)
				}
			}
		})
	}
}
