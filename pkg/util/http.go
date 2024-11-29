package util

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

// HandleGrpcError converts gRPC errors to HTTP errors
func HandleGrpcError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	// Extract gRPC status for detailed error handling
	st, ok := status.FromError(err)
	if !ok {
		// Generic error handling if not a gRPC status error
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return true
	}

	// Map gRPC status codes to appropriate HTTP status codes
	switch st.Code() {
	case codes.NotFound:
		http.Error(w, st.Message(), http.StatusNotFound)
	case codes.InvalidArgument:
		http.Error(w, st.Message(), http.StatusBadRequest)
	default:
		http.Error(w, "Unknown error", http.StatusInternalServerError)
	}
	return true
}
