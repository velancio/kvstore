package transport

import (
	"censys/pkg/util"
	pb "censys/proto/gen/proto"
	"context"
	"encoding/json"
	"net/http"
)

// GrpcServer represents the gRPC server
type GrpcServer struct {
	Store pb.KvStoreServiceClient
}

// HandleGet handles GET requests to retrieve a value from the store
func (s *GrpcServer) HandleGet(w http.ResponseWriter, r *http.Request) {
	// Extract key from query parameter or request body
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing key", http.StatusBadRequest)
		return
	}

	// Make gRPC call to retrieve value
	resp, err := s.Store.Get(context.Background(), &pb.GetRequest{
		Key: key,
	})

	// Handle error and return appropriate http status code
	if util.HandleGrpcError(w, err) {
		return
	}

	// Successful response
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]string{
		"value": resp.Value,
	})
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// HandleSet handles POST requests to set a value in the store
func (s *GrpcServer) HandleSet(w http.ResponseWriter, r *http.Request) {
	// Decode request body
	var req KvPair
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Failed to decode request", http.StatusInternalServerError)
		return
	}

	// Validate key and value
	if util.ValidateKvPair(req.Key, req.Value) != nil {
		http.Error(w, "Invalid key/value pair", http.StatusBadRequest)
		return
	}

	// Make gRPC call to set value
	resp, err := s.Store.Set(context.Background(), &pb.SetRequest{
		Key:   req.Key,
		Value: req.Value,
	})

	// Handle error and return appropriate http status code
	if util.HandleGrpcError(w, err) {
		return
	}

	// Successful response
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]bool{
		"success": resp.Success, // Use the boolean from SetResponse
	})
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// HandleDelete handles DELETE requests to delete a value from the store
func (s *GrpcServer) HandleDelete(w http.ResponseWriter, r *http.Request) {
	// Extract key from query parameter or request body
	key := r.PathValue("key")

	if key == "" {
		http.Error(w, "Missing key", http.StatusBadRequest)
		return
	}

	// Make gRPC call to delete value
	resp, err := s.Store.Delete(context.Background(), &pb.DeleteRequest{
		Key: key,
	})

	// Handle error and return appropriate http status code
	if util.HandleGrpcError(w, err) {
		return
	}

	// Successful response
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]bool{
		"success": resp.Success,
	})
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
