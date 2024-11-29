package main

import (
	"censys/pkg/transport"
	pb "censys/proto/gen/proto"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
)

// NewServer creates a new http server
func NewServer(server transport.Server) http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("GET /store", server.HandleGet)
	router.HandleFunc("POST /store", server.HandleSet)
	router.HandleFunc("DELETE /store/{key}", server.HandleDelete)
	return router
}

// LoadConfig loads config from .env
func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file \n")
	}
}

func main() {
	// Load config from .env
	LoadConfig()

	// Connect to kvstore service
	grpcConnection := fmt.Sprintf("%s:%s", os.Getenv("KVSTORE_HOST"), os.Getenv("KVSTORE_PORT"))
	conn, err := grpc.NewClient(grpcConnection,
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect to kvstore: %s", err)
	}

	// Create grpc client to communicate with kvstore
	store := pb.NewKvStoreServiceClient(conn)
	server := &transport.GrpcServer{
		Store: store,
	}

	// Start http server
	apiPort := fmt.Sprintf(":%s", os.Getenv("API_PORT"))
	err = http.ListenAndServe(apiPort, NewServer(server))
	if err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
