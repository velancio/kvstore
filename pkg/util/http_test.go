package util

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleGrpcError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode int
	}{
		{
			name:     "no error",
			err:      nil,
			wantCode: 0, // no error written
		},
		{
			name:     "grpc error with invalid argument",
			err:      status.Errorf(codes.InvalidArgument, "invalid argument"),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "grpc error key not found",
			err:      status.Errorf(codes.NotFound, "key not found"),
			wantCode: http.StatusNotFound,
		},
		{
			name:     "unknown grpc error",
			err:      status.Errorf(codes.Unknown, "unknown error"),
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			if HandleGrpcError(w, tt.err) {
				if w.Code != tt.wantCode {
					t.Errorf("HandleGrpcError() wrote code %d, want %d", w.Code, tt.wantCode)
				}
			} else {
				if w.Code != 200 && w.Code != 201 {
					t.Errorf("HandleGrpcError() wrote code %d, want 200 0r 201", w.Code)
				}
			}
		})
	}
}
