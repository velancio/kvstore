package transport

import "net/http"

// Server is an interface for handling HTTP requests
type Server interface {
	HandleGet(w http.ResponseWriter, r *http.Request)
	HandleSet(w http.ResponseWriter, r *http.Request)
	HandleDelete(w http.ResponseWriter, r *http.Request)
}

// KvPair represents a key-value pair
type KvPair struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}
